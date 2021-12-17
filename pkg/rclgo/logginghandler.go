package rclgo

// #include <rcl/logging.h>
import "C"

// LoggingOutputHandler is the function signature of logging output handling.
// Backwards compatibility is not guaranteed for this type alias. Use it only if
// necessary.
type LoggingOutputHandler = func(
	location *C.rcutils_log_location_t,
	severity C.int,
	name *C.char,
	timestamp C.rcutils_time_point_value_t,
	format *C.char,
	args *C.va_list,
)

// GetLoggingOutputHandler returns the current logging output handler.
func GetLoggingOutputHandler() LoggingOutputHandler {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	return currentLoggingOutputHandler
}

// SetLoggingOutputHandler sets the current logging output handler to h. If h ==
// nil, DefaultLoggingOutputHandler is used.
func SetLoggingOutputHandler(h LoggingOutputHandler) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	if h == nil {
		currentLoggingOutputHandler = DefaultLoggingOutputHandler
	} else {
		currentLoggingOutputHandler = h
	}
}

// DefaultLoggingOutputHandler is the logging output handler used by default,
// which logs messages based on ROS parameters used to initialize the logging
// system.
func DefaultLoggingOutputHandler(
	location *C.rcutils_log_location_t,
	severity C.int,
	name *C.char,
	timestamp C.rcutils_time_point_value_t,
	format *C.char,
	args *C.va_list,
) {
	C.rcl_logging_multiple_output_handler(location, severity, name, timestamp, format, args)
}

//export loggingOutputHandler
func loggingOutputHandler(
	location *C.rcutils_log_location_t,
	severity C.int,
	name *C.char,
	timestamp C.rcutils_time_point_value_t,
	format *C.char,
	args *C.va_list,
) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	currentLoggingOutputHandler(location, severity, name, timestamp, format, args)
}
