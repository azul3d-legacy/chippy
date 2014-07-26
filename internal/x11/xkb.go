// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <xcb/xkb.h>

#cgo LDFLAGS: -lxcb-xkb
*/
import "C"

import "unsafe"

const (
	XKB_STATE_NOTIFY        = C.XCB_XKB_STATE_NOTIFY
	XKB_MAP_NOTIFY          = C.XCB_XKB_MAP_NOTIFY
	XKB_NEW_KEYBOARD_NOTIFY = C.XCB_XKB_NEW_KEYBOARD_NOTIFY

	XKB_EVENT_TYPE_NEW_KEYBOARD_NOTIFY = C.XCB_XKB_EVENT_TYPE_NEW_KEYBOARD_NOTIFY
	XKB_EVENT_TYPE_MAP_NOTIFY          = C.XCB_XKB_EVENT_TYPE_MAP_NOTIFY
	XKB_EVENT_TYPE_STATE_NOTIFY        = C.XCB_XKB_EVENT_TYPE_STATE_NOTIFY
	XKB_NKN_DETAIL_KEYCODES            = C.XCB_XKB_NKN_DETAIL_KEYCODES
	XKB_MAP_PART_KEY_TYPES             = C.XCB_XKB_MAP_PART_KEY_TYPES
	XKB_MAP_PART_KEY_SYMS              = C.XCB_XKB_MAP_PART_KEY_SYMS
	XKB_MAP_PART_MODIFIER_MAP          = C.XCB_XKB_MAP_PART_MODIFIER_MAP
	XKB_MAP_PART_EXPLICIT_COMPONENTS   = C.XCB_XKB_MAP_PART_EXPLICIT_COMPONENTS
	XKB_MAP_PART_KEY_ACTIONS           = C.XCB_XKB_MAP_PART_KEY_ACTIONS
	XKB_MAP_PART_VIRTUAL_MODS          = C.XCB_XKB_MAP_PART_VIRTUAL_MODS
	XKB_MAP_PART_VIRTUAL_MOD_MAP       = C.XCB_XKB_MAP_PART_VIRTUAL_MOD_MAP
	XKB_STATE_PART_MODIFIER_BASE       = C.XCB_XKB_STATE_PART_MODIFIER_BASE
	XKB_STATE_PART_MODIFIER_LATCH      = C.XCB_XKB_STATE_PART_MODIFIER_LATCH
	XKB_STATE_PART_MODIFIER_LOCK       = C.XCB_XKB_STATE_PART_MODIFIER_LOCK
	XKB_STATE_PART_GROUP_BASE          = C.XCB_XKB_STATE_PART_GROUP_BASE
	XKB_STATE_PART_GROUP_LATCH         = C.XCB_XKB_STATE_PART_GROUP_LATCH
	XKB_STATE_PART_GROUP_LOCK          = C.XCB_XKB_STATE_PART_GROUP_LOCK
)

type (
	XkbDeviceSpec C.xcb_xkb_device_spec_t
)

type XkbSelectEventDetails struct {
	AffectNewKeyboard  uint16
	NewKeyboardDetails uint16
	AffectState        uint16
	StateDetails       uint16

	AffectCtrls           uint32
	CtrlDetails           uint32
	AffectIndicatorState  uint32
	IndicatorStateDetails uint32
	AffectIndicatorMap    uint32
	IndicatorMapDetails   uint32

	AffectNames  uint16
	NamesDetails uint16

	AffectCompat     uint8
	CompatDetails    uint8
	AffectBell       uint8
	BellDetails      uint8
	AffectMsgDetails uint8
	MsgDetails       uint8

	AffectAccessX  uint16
	AccessXDetails uint16
	AffectExtDev   uint16
	ExtdevDetails  uint16
}

func (c *Connection) XkbSelectEventsAuxChecked(deviceSpec XkbDeviceSpec, affectWhich, clear, selectAll, affectMap, pmap uint16, details *XkbSelectEventDetails) VoidCookie {
	return VoidCookie(C.xcb_xkb_select_events_aux_checked(
		c.c(),
		C.xcb_xkb_device_spec_t(deviceSpec),
		C.uint16_t(affectWhich),
		C.uint16_t(clear),
		C.uint16_t(selectAll),
		C.uint16_t(affectMap),
		C.uint16_t(pmap),
		(*C.xcb_xkb_select_events_details_t)(unsafe.Pointer(details)),
	))
}

type XkbAnyEvent struct {
	ResponseType C.uint8_t
	XkbType      C.uint8_t
	Sequence     C.uint16_t
	Time         Timestamp
	DeviceID     C.uint8_t
}

type XkbStateNotifyEvent C.xcb_xkb_state_notify_event_t

func (e *XkbStateNotifyEvent) BaseMods() uint8 {
	return uint8(e.baseMods)
}
func (e *XkbStateNotifyEvent) LatchedMods() uint8 {
	return uint8(e.latchedMods)
}
func (e *XkbStateNotifyEvent) LockedMods() uint8 {
	return uint8(e.lockedMods)
}
func (e *XkbStateNotifyEvent) BaseGroup() int16 {
	return int16(e.baseGroup)
}
func (e *XkbStateNotifyEvent) LatchedGroup() int16 {
	return int16(e.latchedGroup)
}
func (e *XkbStateNotifyEvent) LockedGroup() uint8 {
	return uint8(e.lockedGroup)
}
func (e *XkbStateNotifyEvent) Keycode() Keycode {
	return Keycode(e.keycode)
}
