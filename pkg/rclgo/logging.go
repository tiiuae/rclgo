/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#cgo CFLAGS: -Wno-format-security

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
    const char* format
) {
    rcutils_log(location, severity, name, format);
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
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// The severity levels of log messages / loggers.
type LogSeverity uint32

// Copypaste from /opt/ros/galactic/include/rcutils/logging.h

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
// If the logging system has not yet been initialized on the first call of
// NewContext, logging is initialized by NewContext using the arguments passed
// to it. Unlike InitLogging, NewContext will not update logging configuration
// if logging has already been initialized.
func InitLogging(args *Args) error {
	return rclInitLogging(args, true)
}

var (
	loggingMutex                sync.Mutex
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
	rc := C.rcl_logging_configure_with_output_handler(
		&rclArgs.parsed,
		loggingAllocator,
		(*[0]byte)(C.loggingOutputHandler),
	)
	runtime.KeepAlive(rclArgs)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "rclInitLogging -> rcl_logging_configure_with_output_handler()")
	}
	return nil
}

// logNamed logs msg with severity level to logger named name. logNamed should
// not be used for logging directly. Use one of the exported logging functions
// instead.
func logNamed(level LogSeverity, name, msg string) error {
	// rcutils_log takes a C-style format string as an argument. We want to do
	// string formatting in Go to avoid the pitfalls of C string formatting
	// functions, so the message must be escaped before calling rcutils_log.
	buf := make([]byte, 0, len(msg)+1)
	for i := range msg {
		buf = append(buf, msg[i])
		if msg[i] == '%' {
			buf = append(buf, '%')
		}
	}
	buf = append(buf, 0) // Null terminator

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
		loc.line_number = C.ulong(line)
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	bufHeader := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	C.rcutils_log_wrapper(
		&loc,
		C.int(level),
		cname,
		(*C.char)(unsafe.Pointer(bufHeader.Data)),
	)
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
	level := C.rcutils_logging_get_logger_level(l.cname)
	runtime.KeepAlive(l)
	if level < 0 {
		return 0, errors.New("failed to get log level")
	}
	return LogSeverity(level), nil
}

// SetLevel sets the logging level of l.
func (l *Logger) SetLevel(level LogSeverity) error {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	ret := C.rcutils_logging_set_logger_level(l.cname, C.int(level))
	runtime.KeepAlive(l)
	if ret != C.RCL_RET_OK {
		return errorsCastC(ret, "SetLoggerLevel")
	}
	return nil
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
	level := C.rcutils_logging_get_logger_effective_level(l.cname)
	runtime.KeepAlive(l)
	if level < 0 {
		return 0, errors.New("failed to get effective log level")
	}
	return LogSeverity(level), nil
}

func (l *Logger) Log(level LogSeverity, a ...interface{}) error {
	return logNamed(level, l.name, fmt.Sprint(a...))
}

func (l *Logger) Debug(a ...interface{}) error {
	return logNamed(LogSeverityDebug, l.name, fmt.Sprint(a...))
}
func (l *Logger) Info(a ...interface{}) error {
	return logNamed(LogSeverityInfo, l.name, fmt.Sprint(a...))
}
func (l *Logger) Warn(a ...interface{}) error {
	return logNamed(LogSeverityWarn, l.name, fmt.Sprint(a...))
}
func (l *Logger) Error(a ...interface{}) error {
	return logNamed(LogSeverityError, l.name, fmt.Sprint(a...))
}
func (l *Logger) Fatal(a ...interface{}) error {
	return logNamed(LogSeverityFatal, l.name, fmt.Sprint(a...))
}

func sprintln(a ...interface{}) string {
	b := fmt.Sprintln(a...)
	return b[:len(b)-1]
}

func (l *Logger) Logln(level LogSeverity, a ...interface{}) error {
	return logNamed(level, l.name, sprintln(a...))
}

func (l *Logger) Debugln(a ...interface{}) error {
	return logNamed(LogSeverityDebug, l.name, sprintln(a...))
}
func (l *Logger) Infoln(a ...interface{}) error {
	return logNamed(LogSeverityInfo, l.name, sprintln(a...))
}
func (l *Logger) Warnln(a ...interface{}) error {
	return logNamed(LogSeverityWarn, l.name, sprintln(a...))
}
func (l *Logger) Errorln(a ...interface{}) error {
	return logNamed(LogSeverityError, l.name, sprintln(a...))
}
func (l *Logger) Fatalln(a ...interface{}) error {
	return logNamed(LogSeverityFatal, l.name, sprintln(a...))
}

func (l *Logger) Logf(level LogSeverity, format string, a ...interface{}) error {
	return logNamed(level, l.name, fmt.Sprintf(format, a...))
}

func (l *Logger) Debugf(format string, a ...interface{}) error {
	return logNamed(LogSeverityDebug, l.name, fmt.Sprintf(format, a...))
}
func (l *Logger) Infof(format string, a ...interface{}) error {
	return logNamed(LogSeverityInfo, l.name, fmt.Sprintf(format, a...))
}
func (l *Logger) Warnf(format string, a ...interface{}) error {
	return logNamed(LogSeverityWarn, l.name, fmt.Sprintf(format, a...))
}
func (l *Logger) Errorf(format string, a ...interface{}) error {
	return logNamed(LogSeverityError, l.name, fmt.Sprintf(format, a...))
}
func (l *Logger) Fatalf(format string, a ...interface{}) error {
	return logNamed(LogSeverityFatal, l.name, fmt.Sprintf(format, a...))
}
