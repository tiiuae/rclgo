/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package types

import "unsafe"

type Message interface {
	Clone() Message
	SetDefaults()
}

type MessageTypeSupport interface {
	New() Message
	PrepareMemory() unsafe.Pointer
	ReleaseMemory(unsafe.Pointer)
	AsCStruct(unsafe.Pointer, Message)
	AsGoStruct(Message, unsafe.Pointer)
	TypeSupport() unsafe.Pointer // *C.rosidl_message_type_support_t
}

type ServiceTypeSupport interface {
	Request() MessageTypeSupport
	Response() MessageTypeSupport
	TypeSupport() unsafe.Pointer // *C.rosidl_service_type_support_t
}
