/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#cgo LDFLAGS: -lrmw_implementation
#cgo CFLAGS: -I/opt/ros/galactic/include

#include "rmw/rmw.h"
*/
import "C"

import (
	"time"
)

type RmwQosHistoryPolicy int

const (
	RmwQosHistoryPolicySystemDefault RmwQosHistoryPolicy = iota
	RmwQosHistoryPolicyKeepLast
	RmwQosHistoryPolicyKeepAll
	RmwQosHistoryPolicyUnknown
)

type RmwQosReliabilityPolicy int

const (
	RmwQosReliabilityPolicySystemDefault RmwQosReliabilityPolicy = iota
	RmwQosReliabilityPolicyReliable
	RmwQosReliabilityPolicyBestEffort
	RmwQosReliabilityPolicyUnknown
)

type RmwQosDurabilityPolicy int

const (
	RmwQosDurabilityPolicySystemDefault RmwQosDurabilityPolicy = iota
	RmwQosDurabilityPolicyTransientLocal
	RmwQosDurabilityPolicyVolatile
	RmwQosDurabilityPolicyUnknown
)

type RmwQosLivelinessPolicy int

const (
	RmwQosLivelinessPolicySystemDefault RmwQosLivelinessPolicy = iota
	RmwQosLivelinessPolicyAutomatic
	_
	RmwQosLivelinessPolicyManualByTopic
	RmwQosLivelinessPolicyUnknown
)

type RmwQosProfile struct {
	History                      RmwQosHistoryPolicy
	Depth                        int
	Reliability                  RmwQosReliabilityPolicy
	Durability                   RmwQosDurabilityPolicy
	Deadline                     time.Duration
	Lifespan                     time.Duration
	Liveliness                   RmwQosLivelinessPolicy
	LivelinessLeaseDuration      time.Duration
	AvoidRosNamespaceConventions bool
}

func NewRmwQosProfileDefault() RmwQosProfile {
	return RmwQosProfile{
		History:                      RmwQosHistoryPolicyKeepLast,
		Depth:                        10,
		Reliability:                  RmwQosReliabilityPolicyReliable,
		Durability:                   RmwQosDurabilityPolicyVolatile,
		Deadline:                     0,
		Lifespan:                     0,
		Liveliness:                   RmwQosLivelinessPolicySystemDefault,
		LivelinessLeaseDuration:      0,
		AvoidRosNamespaceConventions: false,
	}
}

func NewRmwQosProfileServicesDefault() RmwQosProfile {
	return NewRmwQosProfileDefault()
}

func (p *RmwQosProfile) asCStruct(dst *C.rmw_qos_profile_t) {
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
