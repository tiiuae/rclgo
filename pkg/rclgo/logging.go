/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#include <rcl/logging.h>

const rcutils_log_location_t zero_location = {
    .function_name = "",
    .file_name = "",
    .line_number = 0
};

// Variable argument functions can't be called from Go so a wrapper is required.
void rcutils_log_wrapper(
    const rcutils_log_location_t* location,
    int severity,
    const char* name,
    const char* format,
    const char* msg
) {
    rcutils_log(location, severity, name, format, msg);
}

void loggingOutputHandler(
    const rcutils_log_location_t* location,
    int severity,
    const char* name,
    const char* format,
    va_list* args
);
*/
import "C"

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

// The severity levels of log messages / loggers.
type LogSeverity uint32

// Copypaste from /opt/ros/jazzy/include/rcutils/logging.h

const (
	LogSeverityUnset LogSeverity = 0  ///< The unset log level
	LogSeverityDebug LogSeverity = 10 ///< The debug log level
	LogSeverityInfo  LogSeverity = 20 ///< The info log level
	LogSeverityWarn  LogSeverity = 30 ///< The warn log level
	LogSeverityError LogSeverity = 40 ///< The error log level
	LogSeverityFatal LogSeverity = 50 ///< The fatal log level
)

func (s LogSeverity) String() string {
	switch s {
	case LogSeverityUnset:
		return "UNSET"
	case LogSeverityDebug:
		return "DEBUG"
	case LogSeverityInfo:
		return "INFO"
	case LogSeverityWarn:
		return "WARN"
	case LogSeverityError:
		return "ERROR"
	case LogSeverityFatal:
		return "FATAL"
	default:
		return ""
	}
}

// InitLogging initializes the logging system, which is required for using
// logging functionality.
//
// Logging configuration can be updated by calling InitLogging again with the
// desired args.
//
// If the logging system has not yet been initialized on the first call of Init
// or NewContext, logging is initialized by that call using the passed
// arguments. Unlike InitLogging, Init and NewContext will not update logging
// configuration if logging has already been initialized.
func InitLogging(args *Args) error {
	return rclInitLogging(args, true)
}

var (
	loggingMutex sync.Mutex
	// OS thread that holds loggingMutex across rcl_logging_configure; rcl logs
	// from inside it, and those callbacks run on this thread and must not relock.
	loggingMutexOwner           atomic.Int64
	currentLoggingOutputHandler = DefaultLoggingOutputHandler
	loggingInitialized          = false
	loggingAllocator            = func() *C.rcl_allocator_t {
		alloc := (*C.rcl_allocator_t)(C.malloc(C.sizeof_rcl_allocator_t))
		*alloc = C.rcl_get_default_allocator()
		return alloc
	}()
)

func rclInitLogging(rclArgs *Args, update bool) error {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	if loggingInitialized && !update {
		return nil
	}
	loggingInitialized = true
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	loggingMutexOwner.Store(int64(syscall.Gettid()))
	defer loggingMutexOwner.Store(0)
	err := rclCall("rclInitLogging -> rcl_logging_configure_with_output_handler()", func() C.rcl_ret_t {
		return C.rcl_logging_configure_with_output_handler(
			&rclArgs.parsed,
			loggingAllocator,
			(*[0]byte)(C.loggingOutputHandler),
		)
	})
	runtime.KeepAlive(rclArgs)
	return err
}

var rcutilsLogFormat = C.CString("%s")

// logNamed logs msg with severity level to logger named name. logNamed should
// not be used for logging directly. Use one of the exported logging functions
// instead.
func logNamed(level LogSeverity, name, msg string) error {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	loc := C.zero_location
	// Because convenience wrappers are provided for logging, two stack frames
	// must be skipped; one for logNamed and one for the wrapper.
	if pc, file, line, ok := runtime.Caller(2); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			loc.function_name = C.CString(fn.Name())
			defer C.free(unsafe.Pointer(loc.function_name))
		}
		loc.file_name = C.CString(file)
		defer C.free(unsafe.Pointer(loc.file_name))
		loc.line_number = C.size_t(line)
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.rcutils_log_wrapper(&loc, C.int(level), cname, rcutilsLogFormat, cmsg)
	return nil
}

var defaultLogger = &Logger{
	name:  "",
	cname: C.CString(""),
}

// Logger can be used to log messages using the ROS 2 logging system.
//
// Loggers are usable only after logging has been initialized. See InitLogging.
//
// Logging methods prefixed with "Log" take the logging level as the first
// parameter. Methods prefixed with the name of a logging level are shorthands
// to "Log" methods, and log using the prefixed logging level.
//
// Logging methods suffixed with "", "f" or "ln" format their arguments in the
// same way as fmt.Print, fmt.Printf and fmt.Println, respectively.
type Logger struct {
	name  string
	cname *C.char
}

var invalidLoggerNameRegex = regexp.MustCompile(`(^\.)|(\.$)|(\.\.)`)

// GetLogger returns the logger named name. If name is empty, the default logger
// is returned. Returns nil if name is invalid.
func GetLogger(name string) *Logger {
	if name == "" {
		return defaultLogger
	}
	if invalidLoggerNameRegex.MatchString(name) {
		return nil
	}
	l := &Logger{
		name:  name,
		cname: C.CString(name),
	}
	runtime.SetFinalizer(l, func(l *Logger) {
		C.free(unsafe.Pointer(l.cname))
		l.cname = nil
	})
	return l
}

func (l *Logger) Name() string {
	return l.name
}

// Parent returns the parent logger of l. If l has no parent, the default logger
// is returned.
func (l *Logger) Parent() *Logger {
	i := strings.LastIndexByte(l.name, '.')
	if i == -1 {
		return defaultLogger
	}
	return GetLogger(l.name[:i])
}

// Child returns the child logger of l named name. Returns nil if name is
// invalid.
func (l *Logger) Child(name string) *Logger {
	if l.name == "" {
		return GetLogger(name)
	}
	return GetLogger(l.name + "." + name)
}

// Level returns the logging level of l. Note that this is not necessarily the
// same as EffectiveLevel.
func (l *Logger) Level() (LogSeverity, error) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	level := C.rcutils_logging_get_logger_level(l.cname)
	runtime.KeepAlive(l)
	if level < 0 {
		return 0, errors.New("failed to get log level: " + errorString())
	}
	return LogSeverity(level), nil
}

// SetLevel sets the logging level of l.
func (l *Logger) SetLevel(level LogSeverity) error {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	err := rclCall("SetLoggerLevel", func() C.rcl_ret_t {
		return C.rcutils_logging_set_logger_level(l.cname, C.int(level))
	})
	runtime.KeepAlive(l)
	return err
}

// IsEnabledFor returns true if l can log messages whose severity is at least
// level and false if not.
func (l *Logger) IsEnabledFor(level LogSeverity) bool {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	b := C.rcutils_logging_logger_is_enabled_for(l.cname, C.int(level))
	runtime.KeepAlive(l)
	return bool(b)
}

// EffectiveLevel returns the effective logging level of l, which considers the
// logging levels of l's ancestors as well as the logging level of l itself.
// Note that this is not necessarily the same as Level.
func (l *Logger) EffectiveLevel() (LogSeverity, error) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	level := C.rcutils_logging_get_logger_effective_level(l.cname)
	runtime.KeepAlive(l)
	if level < 0 {
		return 0, errors.New("failed to get effective log level: " + errorString())
	}
	return LogSeverity(level), nil
}

func (l *Logger) Log(level LogSeverity, a ...any) error {
	return logNamed(level, l.name, fmt.Sprint(a...))
}

func (l *Logger) Debug(a ...any) error {
	return logNamed(LogSeverityDebug, l.name, fmt.Sprint(a...))
}

func (l *Logger) Info(a ...any) error {
	return logNamed(LogSeverityInfo, l.name, fmt.Sprint(a...))
}

func (l *Logger) Warn(a ...any) error {
	return logNamed(LogSeverityWarn, l.name, fmt.Sprint(a...))
}

func (l *Logger) Error(a ...any) error {
	return logNamed(LogSeverityError, l.name, fmt.Sprint(a...))
}

func (l *Logger) Fatal(a ...any) error {
	return logNamed(LogSeverityFatal, l.name, fmt.Sprint(a...))
}

func sprintln(a ...any) string {
	b := fmt.Sprintln(a...)
	return b[:len(b)-1]
}

func (l *Logger) Logln(level LogSeverity, a ...any) error {
	return logNamed(level, l.name, sprintln(a...))
}

func (l *Logger) Debugln(a ...any) error {
	return logNamed(LogSeverityDebug, l.name, sprintln(a...))
}

func (l *Logger) Infoln(a ...any) error {
	return logNamed(LogSeverityInfo, l.name, sprintln(a...))
}

func (l *Logger) Warnln(a ...any) error {
	return logNamed(LogSeverityWarn, l.name, sprintln(a...))
}

func (l *Logger) Errorln(a ...any) error {
	return logNamed(LogSeverityError, l.name, sprintln(a...))
}

func (l *Logger) Fatalln(a ...any) error {
	return logNamed(LogSeverityFatal, l.name, sprintln(a...))
}

func (l *Logger) Logf(level LogSeverity, format string, a ...any) error {
	return logNamed(level, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Debugf(format string, a ...any) error {
	return logNamed(LogSeverityDebug, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Infof(format string, a ...any) error {
	return logNamed(LogSeverityInfo, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Warnf(format string, a ...any) error {
	return logNamed(LogSeverityWarn, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Errorf(format string, a ...any) error {
	return logNamed(LogSeverityError, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Fatalf(format string, a ...any) error {
	return logNamed(LogSeverityFatal, l.name, fmt.Sprintf(format, a...))
}
