/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/
package rclgo

/*
#cgo LDFLAGS: -ldl

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <dlfcn.h>

#include <rosidl_runtime_c/message_type_support_struct.h>

char* formatString(const char* format, ...) {
	va_list args1, args2;
	va_start(args1, format);
	va_copy(args2, args1);
	size_t len = 1 + vsnprintf(NULL, 0, format, args1);
	va_end(args1);
	char* buf = malloc(len);
	if (buf != NULL) {
		vsprintf(buf, format, args2);
	}
	va_end(args2);
	return buf;
}

typedef rosidl_message_type_support_t * (*GetTypeSupportFunc)();

const char* loadTypeSupport(
	const char* pkgName,
	const char* ifaceName,
	void** lib,
	void** typeSupport
) {
	char* libName = formatString(
		"lib%s__rosidl_typesupport_c.so",
		pkgName
	);
	if (libName == NULL) {
		return "allocation failed";
	}
	*lib = dlopen(libName, RTLD_LAZY);
	if (*lib == NULL) {
		free(libName);
		return dlerror();
	}
	char* tsName = formatString(
		"rosidl_typesupport_c__get_message_type_support_handle__%s__msg__%s",
		pkgName, ifaceName
	);
	if (tsName == NULL) {
		free(libName);
		dlclose(*lib);
		return "allocation failed";
	}
	void* tsSym = dlsym(*lib, tsName);
	if (tsSym == NULL) {
		free(libName);
		dlclose(*lib);
		free(tsName);
		return dlerror();
	}
	*typeSupport = ((GetTypeSupportFunc)tsSym)();
	free(libName);
	free(tsName);
	return NULL;
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

type dynamicMessageTypeSupport struct {
	lib         unsafe.Pointer // void*
	typeSupport unsafe.Pointer // rosidl_message_type_support_t*
}

// LoadDynamicMessageTypeSupport loads a message type support implementation
// dynamically.
//
// MessageTypeSupport instances returned by LoadDynamicMessageTypeSupport
// support use cases related to handling only serialized messages. Methods New,
// PrepareMemory, ReleaseMemory, AsCStruct and AsGoStruct will panic.
//
// Backwards compatibility is not guaranteed for this API. Use it only if
// necessary.
func LoadDynamicMessageTypeSupport(pkgName, msgName string) (types.MessageTypeSupport, error) {
	cPkgName := C.CString(pkgName)
	defer C.free(unsafe.Pointer(cPkgName))
	cIfaceName := C.CString(msgName)
	defer C.free(unsafe.Pointer(cIfaceName))
	ts := new(dynamicMessageTypeSupport)
	err := C.loadTypeSupport(cPkgName, cIfaceName, &ts.lib, &ts.typeSupport)
	if err != nil {
		return nil, fmt.Errorf("failed to load type support: %v", C.GoString(err))
	}
	runtime.SetFinalizer(ts, func(g *dynamicMessageTypeSupport) {
		C.dlclose(g.lib)
	})
	return ts, nil
}

func (g *dynamicMessageTypeSupport) New() types.Message {
	panic("not supported")
}

func (g *dynamicMessageTypeSupport) PrepareMemory() unsafe.Pointer {
	panic("not supported")
}

func (g *dynamicMessageTypeSupport) ReleaseMemory(unsafe.Pointer) {
	panic("not supported")
}

func (g *dynamicMessageTypeSupport) AsCStruct(unsafe.Pointer, types.Message) {
	panic("not supported")
}

func (g *dynamicMessageTypeSupport) AsGoStruct(types.Message, unsafe.Pointer) {
	panic("not supported")
}

func (g *dynamicMessageTypeSupport) TypeSupport() unsafe.Pointer { // *C.rosidl_message_type_support_t
	return g.typeSupport
}
