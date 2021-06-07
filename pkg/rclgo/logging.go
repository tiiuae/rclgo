/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/// The severity levels of log messages / loggers. // Copypaste from /opt/ros/foxy/include/rcutils/logging.h
type LogSeverity uint32

const (
	LogSeverityUnset LogSeverity = 0  ///< The unset log level
	LogSeverityDebug LogSeverity = 10 ///< The debug log level
	LogSeverityInfo  LogSeverity = 20 ///< The info log level
	LogSeverityWarn  LogSeverity = 30 ///< The warn log level
	LogSeverityError LogSeverity = 40 ///< The error log level
	LogSeverityFatal LogSeverity = 50 ///< The fatal log level
)
