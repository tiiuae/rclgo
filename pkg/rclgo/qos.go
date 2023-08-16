/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

// #include "rmw/rmw.h"
import "C"

import (
	"time"
)

const (
	DurationInfinite    = 9223372036*time.Second + 854775807*time.Nanosecond
	DurationUnspecified = time.Duration(0)
)

type HistoryPolicy int

const (
	HistorySystemDefault HistoryPolicy = iota
	HistoryKeepLast
	HistoryKeepAll
	HistoryUnknown
)

type ReliabilityPolicy int

const (
	ReliabilitySystemDefault ReliabilityPolicy = iota
	ReliabilityReliable
	ReliabilityBestEffort
	ReliabilityUnknown
)

type DurabilityPolicy int

const (
	DurabilitySystemDefault DurabilityPolicy = iota
	DurabilityTransientLocal
	DurabilityVolatile
	DurabilityUnknown
)
const DeadlineDefault = DurationUnspecified

const LifespanDefault = DurationUnspecified

type LivelinessPolicy int

const (
	LivelinessSystemDefault LivelinessPolicy = iota
	LivelinessAutomatic
	_
	LivelinessManualByTopic
	LivelinessUnknown
)

const LivelinessLeaseDurationDefault = DurationUnspecified

type QosProfile struct {
	History                      HistoryPolicy
	Depth                        int
	Reliability                  ReliabilityPolicy
	Durability                   DurabilityPolicy
	Deadline                     time.Duration
	Lifespan                     time.Duration
	Liveliness                   LivelinessPolicy
	LivelinessLeaseDuration      time.Duration
	AvoidRosNamespaceConventions bool
}

func NewDefaultQosProfile() QosProfile {
	return QosProfile{
		History:                      HistoryKeepLast,
		Depth:                        10,
		Reliability:                  ReliabilityReliable,
		Durability:                   DurabilityVolatile,
		Deadline:                     DeadlineDefault,
		Lifespan:                     LifespanDefault,
		Liveliness:                   LivelinessSystemDefault,
		LivelinessLeaseDuration:      LivelinessLeaseDurationDefault,
		AvoidRosNamespaceConventions: false,
	}
}

func NewDefaultServiceQosProfile() QosProfile {
	return NewDefaultQosProfile()
}

func (p *QosProfile) asCStruct(dst *C.rmw_qos_profile_t) {
	dst.history = uint32(p.History)
	dst.depth = C.ulong(p.Depth)
	dst.reliability = uint32(p.Reliability)
	dst.durability = uint32(p.Durability)
	dst.deadline = C.rmw_time_t{nsec: C.ulong(p.Deadline)}
	dst.lifespan = C.rmw_time_t{nsec: C.ulong(p.Lifespan)}
	dst.liveliness = uint32(p.Liveliness)
	dst.liveliness_lease_duration = C.rmw_time_t{nsec: C.ulong(p.LivelinessLeaseDuration)}
	dst.avoid_ros_namespace_conventions = C.bool(p.AvoidRosNamespaceConventions)
}
