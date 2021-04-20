/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

/// The severity levels of log messages / loggers. // Copypaste from /opt/ros/foxy/include/rcutils/logging.h
type RCUTILS_LOG_SEVERITY uint32

var RCUTILS_LOG_SEVERITY_UNSET = 0  ///< The unset log level
var RCUTILS_LOG_SEVERITY_DEBUG = 10 ///< The debug log level
var RCUTILS_LOG_SEVERITY_INFO = 20  ///< The info log level
var RCUTILS_LOG_SEVERITY_WARN = 30  ///< The warn log level
var RCUTILS_LOG_SEVERITY_ERROR = 40 ///< The error log level
var RCUTILS_LOG_SEVERITY_FATAL = 50 ///< The fatal log level
