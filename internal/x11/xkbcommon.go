// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include "xkbcommon/xkbcommon-x11.h"

#cgo LDFLAGS: -lxcb -lXau -lXdmcp -lxcb-xkb
*/
import "C"

import (
	"reflect"
	"runtime"
	"unicode/utf8"
	"unsafe"
)

const (
	XKB_X11_MIN_MAJOR_XKB_VERSION = C.XKB_X11_MIN_MAJOR_XKB_VERSION
	XKB_X11_MIN_MINOR_XKB_VERSION = C.XKB_X11_MIN_MINOR_XKB_VERSION
	XKB_KEY_DOWN                  = C.XKB_KEY_DOWN
	XKB_KEY_UP                    = C.XKB_KEY_UP
)

type (
	XkbStateComponent            C.enum_xkb_state_component
	XkbContextFlags              C.enum_xkb_context_flags
	XkbKeyDirection              C.enum_xkb_key_direction
	XkbX11SetupXkbExtensionFlags C.enum_xkb_x11_setup_xkb_extension_flags
	XkbKeymapCompileFlags        C.enum_xkb_keymap_compile_flags

	XkbKeycode     C.xkb_keycode_t
	XkbKeysym      C.xkb_keysym_t
	XkbModMask     C.xkb_mod_mask_t
	XkbLayoutIndex C.xkb_layout_index_t
)

func (c *Connection) XkbX11SetupXkbExtension(majorXkbVersion, minorXkbVersion uint16, flags XkbX11SetupXkbExtensionFlags) (ret int, majorVersion, minorVersion uint16, baseEvent, baseError uint8) {
	ret = int(C.xkb_x11_setup_xkb_extension(
		c.c(),
		C.uint16_t(majorXkbVersion),
		C.uint16_t(minorXkbVersion),
		C.enum_xkb_x11_setup_xkb_extension_flags(flags),
		(*C.uint16_t)(unsafe.Pointer(&majorVersion)),
		(*C.uint16_t)(unsafe.Pointer(&minorVersion)),
		(*C.uint8_t)(unsafe.Pointer(&baseEvent)),
		(*C.uint8_t)(unsafe.Pointer(&baseError)),
	))
	return
}

// Get core/default keyboard device id (-1 on error)
func (c *Connection) XkbX11GetCoreKeyboardDeviceId() int32 {
	return int32(C.xkb_x11_get_core_keyboard_device_id(c.c()))
}

type XkbKeymap struct {
	c *C.struct_xkb_keymap
}

// Get keymap from current device state, nil on failure.
func (c *Connection) XkbX11KeymapNewFromDevice(context *XkbContext, deviceId int32, flags XkbKeymapCompileFlags) *XkbKeymap {
	km := new(XkbKeymap)
	km.c = C.xkb_x11_keymap_new_from_device(
		context.c,
		c.c(),
		C.int32_t(deviceId),
		C.enum_xkb_keymap_compile_flags(flags),
	)
	if km.c == nil {
		return nil
	}
	runtime.SetFinalizer(km, func(tmp *XkbKeymap) {
		C.free(unsafe.Pointer(tmp.c))
	})
	return km
}

type XkbContext struct {
	c *C.struct_xkb_context
}

// Create new context, nil on failure.
func XkbContextNew(flags XkbContextFlags) *XkbContext {
	c := new(XkbContext)
	c.c = C.xkb_context_new(
		C.enum_xkb_context_flags(flags),
	)
	if c.c == nil {
		return nil
	}
	runtime.SetFinalizer(c, func(tmp *XkbContext) {
		C.free(unsafe.Pointer(tmp.c))
	})
	return c
}

type XkbState struct {
	c *C.struct_xkb_state
}

func (s *XkbState) UpdateMask(depressedMods, latchedMods, lockedMods XkbModMask, depressedLayout, latchedLayout, lockedLayout XkbLayoutIndex) XkbStateComponent {
	return XkbStateComponent(C.xkb_state_update_mask(
		s.c,
		C.xkb_mod_mask_t(depressedMods),
		C.xkb_mod_mask_t(latchedMods),
		C.xkb_mod_mask_t(lockedMods),
		C.xkb_layout_index_t(depressedLayout),
		C.xkb_layout_index_t(latchedLayout),
		C.xkb_layout_index_t(lockedLayout),
	))
}

func (s *XkbState) LedNameIsActive(name string) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.xkb_state_led_name_is_active(
		s.c,
		cname,
	))
}

func (s *XkbState) KeyGetSyms(key XkbKeycode) (out []XkbKeysym) {
	var syms *C.xkb_keysym_t
	length := int(C.xkb_state_key_get_syms(
		s.c,
		C.xkb_keycode_t(key),
		&syms,
	))
	if length > 0 {
		h := (*reflect.SliceHeader)(unsafe.Pointer(&out))
		h.Len = length
		h.Cap = length
		h.Data = uintptr(unsafe.Pointer(syms))
	}
	return
}

func (s *XkbState) KeyGetOneSym(key XkbKeycode) XkbKeysym {
	return XkbKeysym(C.xkb_state_key_get_one_sym(
		s.c,
		C.xkb_keycode_t(key),
	))
}

func (s *XkbState) KeyGetLayout(key XkbKeycode) XkbLayoutIndex {
	return XkbLayoutIndex(C.xkb_state_key_get_layout(
		s.c,
		C.xkb_keycode_t(key),
	))
}

// Create keyboard state manager from a given keymap and device id, returns nil
// on error.
func (c *Connection) XkbX11StateNewFromDevice(keymap *XkbKeymap, deviceId int32) *XkbState {
	s := new(XkbState)
	s.c = C.xkb_x11_state_new_from_device(
		keymap.c,
		c.c(),
		C.int32_t(deviceId),
	)
	if s.c == nil {
		return nil
	}
	runtime.SetFinalizer(s, func(tmp *XkbState) {
		C.free(unsafe.Pointer(tmp.c))
	})
	return s
}

func XkbKeysymToUTF8(keysym XkbKeysym, buffer unsafe.Pointer, size int) int {
	return int(C.xkb_keysym_to_utf8(
		C.xkb_keysym_t(keysym),
		(*C.char)(buffer),
		C.size_t(size),
	))
}

func (k XkbKeysym) Rune() rune {
	var (
		buf    = make([]byte, 64) // 64*2 for first run.
		status = -1               // Buffer too small
	)
	for status == -1 {
		buf = make([]byte, len(buf)*2)
		status = XkbKeysymToUTF8(k, unsafe.Pointer(&buf[0]), len(buf))
	}
	r, _ := utf8.DecodeRune(buf)
	return r
}
