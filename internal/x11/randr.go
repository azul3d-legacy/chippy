// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/randr.h>

#cgo LDFLAGS: -lxcb -lxcb-randr
*/
import "C"

import (
	"errors"
	"reflect"
	"runtime"
	"unsafe"
)

const (
	RANDR_MODE_FLAG_DOUBLE_SCAN = C.XCB_RANDR_MODE_FLAG_DOUBLE_SCAN
	RANDR_MODE_FLAG_INTERLACE   = C.XCB_RANDR_MODE_FLAG_INTERLACE
)

var (
	RandrSetConfigInvalidConfigTime = errors.New("RandrSetConfig(): Invalid config timestamp.")
	RandrSetConfigInvalidTime       = errors.New("RandrSetConfig(): Invalid timestamp.")
	RandrSetConfigFailed            = errors.New("RandrSetConfig(): Failed to set configuration.")
)

type RandrModeInfo struct {
	Id         C.uint32_t
	Width      C.uint16_t
	Height     C.uint16_t
	DotClock   C.uint32_t
	HSyncStart C.uint16_t
	HSyncEnd   C.uint16_t
	HTotal     C.uint16_t
	HSkew      C.uint16_t
	VSyncStart C.uint16_t
	VSyncEnd   C.uint16_t
	VTotal     C.uint16_t
	NameLen    C.uint16_t
	ModeFlags  C.uint32_t
}

type ERandrQueryVersionReply struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
	MajorVersion C.uint32_t
	MinorVersion C.uint32_t
	Pad1         [16]C.uint8_t
}
type RandrQueryVersionReply struct {
	*ERandrQueryVersionReply
}

type RandrQueryVersionCookie C.xcb_randr_query_version_cookie_t

func (c RandrQueryVersionCookie) c() C.xcb_randr_query_version_cookie_t {
	return C.xcb_randr_query_version_cookie_t(c)
}

func (c *Connection) RandrQueryVersion(majorVersion, minorVersion uint32) RandrQueryVersionCookie {
	cookie := C.xcb_randr_query_version(c.c(), C.uint32_t(minorVersion), C.uint32_t(majorVersion))
	return RandrQueryVersionCookie(cookie)
}

func (c *Connection) RandrQueryVersionReply(cookie RandrQueryVersionCookie) (reply *RandrQueryVersionReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_query_version_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrQueryVersionReply)
		reply.ERandrQueryVersionReply = (*ERandrQueryVersionReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrQueryVersionReply) {
			C.free(unsafe.Pointer(f.ERandrQueryVersionReply))
		})

	} else {
		err = errors.New("RandrQueryVersionReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type ERandrGetScreenResourcesCurrentReply struct {
	ResponseType    C.uint8_t
	Pad0            C.uint8_t
	Sequence        C.uint16_t
	Length          C.uint32_t
	Timestamp       Timestamp
	ConfigTimestamp Timestamp
	NumCrtcs        C.uint16_t
	NumOutputs      C.uint16_t
	NumModes        C.uint16_t
	NamesLen        C.uint16_t
	Pad1            [8]C.uint8_t
}
type RandrGetScreenResourcesCurrentReply struct {
	*ERandrGetScreenResourcesCurrentReply
}

func (c *RandrGetScreenResourcesCurrentReply) c() *C.xcb_randr_get_screen_resources_current_reply_t {
	ptr := c.ERandrGetScreenResourcesCurrentReply
	return (*C.xcb_randr_get_screen_resources_current_reply_t)(unsafe.Pointer(ptr))
}

type RandrGetScreenResourcesCurrentCookie C.xcb_randr_get_screen_resources_current_cookie_t

func (c RandrGetScreenResourcesCurrentCookie) c() C.xcb_randr_get_screen_resources_current_cookie_t {
	return C.xcb_randr_get_screen_resources_current_cookie_t(c)
}

func (c *Connection) RandrGetScreenResourcesCurrent(window Window) RandrGetScreenResourcesCurrentCookie {
	cookie := C.xcb_randr_get_screen_resources_current(c.c(), C.xcb_window_t(window))
	return RandrGetScreenResourcesCurrentCookie(cookie)
}

func (c *Connection) RandrGetScreenResourcesCurrentReply(cookie RandrGetScreenResourcesCurrentCookie) (reply *RandrGetScreenResourcesCurrentReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_get_screen_resources_current_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrGetScreenResourcesCurrentReply)
		reply.ERandrGetScreenResourcesCurrentReply = (*ERandrGetScreenResourcesCurrentReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrGetScreenResourcesCurrentReply) {
			C.free(unsafe.Pointer(f.ERandrGetScreenResourcesCurrentReply))
		})

	} else {
		err = errors.New("RandrGetScreenResourcesCurrentReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) RandrGetScreenResourcesCurrentModes(r *RandrGetScreenResourcesCurrentReply) (modes *RandrModeInfos) {
	cModes := C.xcb_randr_get_screen_resources_current_modes(r.c())
	numModes := C.xcb_randr_get_screen_resources_current_modes_length(r.c())

	modes = new(RandrModeInfos)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&modes.Slice))
	sliceHeader.Len = int(numModes)
	sliceHeader.Cap = int(numModes)
	sliceHeader.Data = uintptr(unsafe.Pointer(cModes))
	return
}

type RandrCrtcs struct {
	Slice []RandrCrtc
}

func (c *Connection) RandrGetScreenResourcesCurrentCrtcs(r *RandrGetScreenResourcesCurrentReply) (crtcs *RandrCrtcs) {
	cCrtcs := C.xcb_randr_get_screen_resources_current_crtcs(r.c())
	numCrtcs := C.xcb_randr_get_screen_resources_current_crtcs_length(r.c())

	crtcs = new(RandrCrtcs)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&crtcs.Slice))
	sliceHeader.Len = int(numCrtcs)
	sliceHeader.Cap = int(numCrtcs)
	sliceHeader.Data = uintptr(unsafe.Pointer(cCrtcs))
	return
}

type ERandrGetScreenResourcesReply struct {
	ResponseType    C.uint8_t
	Pad0            C.uint8_t
	Sequence        C.uint16_t
	Length          C.uint32_t
	Timestamp       Timestamp
	ConfigTimestamp Timestamp
	NumCrtcs        C.uint16_t
	NumOutputs      C.uint16_t
	NumModes        C.uint16_t
	NamesLen        C.uint16_t
	Pad1            [8]C.uint8_t
}
type RandrGetScreenResourcesReply struct {
	*ERandrGetScreenResourcesReply
}

func (c *RandrGetScreenResourcesReply) c() *C.xcb_randr_get_screen_resources_reply_t {
	ptr := c.ERandrGetScreenResourcesReply
	return (*C.xcb_randr_get_screen_resources_reply_t)(unsafe.Pointer(ptr))
}

type RandrGetScreenResourcesCookie C.xcb_randr_get_screen_resources_cookie_t

func (c RandrGetScreenResourcesCookie) c() C.xcb_randr_get_screen_resources_cookie_t {
	return C.xcb_randr_get_screen_resources_cookie_t(c)
}

func (c *Connection) RandrGetScreenResources(window Window) RandrGetScreenResourcesCookie {
	cookie := C.xcb_randr_get_screen_resources(c.c(), C.xcb_window_t(window))
	return RandrGetScreenResourcesCookie(cookie)
}

func (c *Connection) RandrGetScreenResourcesReply(cookie RandrGetScreenResourcesCookie) (reply *RandrGetScreenResourcesReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_get_screen_resources_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrGetScreenResourcesReply)
		reply.ERandrGetScreenResourcesReply = (*ERandrGetScreenResourcesReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrGetScreenResourcesReply) {
			C.free(unsafe.Pointer(f.ERandrGetScreenResourcesReply))
		})

	} else {
		err = errors.New("RandrGetScreenResourcesReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type RandrModeInfos struct {
	Slice []RandrModeInfo
}

func (c *Connection) RandrGetScreenResourcesModes(r *RandrGetScreenResourcesReply) (modes *RandrModeInfos) {
	cModes := C.xcb_randr_get_screen_resources_modes(r.c())
	numModes := C.xcb_randr_get_screen_resources_modes_length(r.c())

	modes = new(RandrModeInfos)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&modes.Slice))
	sliceHeader.Len = int(numModes)
	sliceHeader.Cap = int(numModes)
	sliceHeader.Data = uintptr(unsafe.Pointer(cModes))
	return
}

func (c *Connection) RandrGetScreenResourcesCrtcs(r *RandrGetScreenResourcesReply) (crtcs *RandrCrtcs) {
	cCrtcs := C.xcb_randr_get_screen_resources_crtcs(r.c())
	numCrtcs := C.xcb_randr_get_screen_resources_crtcs_length(r.c())

	crtcs = new(RandrCrtcs)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&crtcs.Slice))
	sliceHeader.Len = int(numCrtcs)
	sliceHeader.Cap = int(numCrtcs)
	sliceHeader.Data = uintptr(unsafe.Pointer(cCrtcs))
	return
}

type ERandrSetCrtcConfigReply struct {
	ResponseType C.uint8_t
	Status       C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
	Timestamp    Timestamp
	Pad0         [20]C.uint8_t
}
type RandrSetCrtcConfigReply struct {
	*ERandrSetCrtcConfigReply
}

func (c *RandrSetCrtcConfigReply) c() *C.xcb_randr_set_crtc_config_reply_t {
	ptr := c.ERandrSetCrtcConfigReply
	return (*C.xcb_randr_set_crtc_config_reply_t)(unsafe.Pointer(ptr))
}

type RandrSetCrtcConfigCookie C.xcb_randr_set_crtc_config_cookie_t

func (c RandrSetCrtcConfigCookie) c() C.xcb_randr_set_crtc_config_cookie_t {
	return C.xcb_randr_set_crtc_config_cookie_t(c)
}

type (
	RandrCrtc   C.xcb_randr_crtc_t
	RandrMode   C.xcb_randr_mode_t
	RandrOutput C.xcb_randr_output_t
)

func (c *Connection) RandrSetCrtcConfig(crtc RandrCrtc, timestamp Timestamp, configTimestamp Timestamp, x, y int16, mode RandrMode, rotation uint16, outputs []RandrOutput) RandrSetCrtcConfigCookie {
	var coutputs *C.xcb_randr_output_t
	var coutputsLen int
	if len(outputs) > 0 {
		coutputs = (*C.xcb_randr_output_t)(unsafe.Pointer(&outputs[0]))
		coutputsLen = len(outputs)
	}
	cookie := C.xcb_randr_set_crtc_config(
		c.c(), C.xcb_randr_crtc_t(crtc),
		C.xcb_timestamp_t(timestamp), C.xcb_timestamp_t(configTimestamp),
		C.int16_t(x), C.int16_t(y),
		C.xcb_randr_mode_t(mode), C.uint16_t(rotation),
		C.uint32_t(coutputsLen),
		coutputs,
	)
	return RandrSetCrtcConfigCookie(cookie)
}

func (c *Connection) RandrSetCrtcConfigReply(cookie RandrSetCrtcConfigCookie) (reply *RandrSetCrtcConfigReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_set_crtc_config_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrSetCrtcConfigReply)
		reply.ERandrSetCrtcConfigReply = (*ERandrSetCrtcConfigReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrSetCrtcConfigReply) {
			C.free(unsafe.Pointer(f.ERandrSetCrtcConfigReply))
		})

		switch reply.Status {
		case C.XCB_RANDR_SET_CONFIG_SUCCESS:
			break
		case C.XCB_RANDR_SET_CONFIG_INVALID_CONFIG_TIME:
			err = RandrSetConfigInvalidConfigTime
		case C.XCB_RANDR_SET_CONFIG_INVALID_TIME:
			err = RandrSetConfigInvalidTime
		case C.XCB_RANDR_SET_CONFIG_FAILED:
			err = RandrSetConfigFailed
		}

	} else {
		err = errors.New("RandrSetCrtcConfigReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type ERandrGetCrtcInfoReply struct {
	ResponseType       C.uint8_t
	Status             C.uint8_t
	Sequence           C.uint16_t
	Length             C.uint32_t
	Timestamp          Timestamp
	X                  C.int16_t
	Y                  C.int16_t
	Width              C.uint16_t
	Height             C.uint16_t
	Mode               RandrMode
	Rotation           C.uint16_t
	Rotations          C.uint16_t
	NumOutputs         C.uint16_t
	NumPossibleOutputs C.uint16_t
}
type RandrGetCrtcInfoReply struct {
	*ERandrGetCrtcInfoReply
}

func (c *RandrGetCrtcInfoReply) c() *C.xcb_randr_get_crtc_info_reply_t {
	ptr := c.ERandrGetCrtcInfoReply
	return (*C.xcb_randr_get_crtc_info_reply_t)(unsafe.Pointer(ptr))
}

type RandrGetCrtcInfoCookie C.xcb_randr_get_crtc_info_cookie_t

func (c RandrGetCrtcInfoCookie) c() C.xcb_randr_get_crtc_info_cookie_t {
	return C.xcb_randr_get_crtc_info_cookie_t(c)
}

func (c *Connection) RandrGetCrtcInfo(crtc RandrCrtc, configTimestamp Timestamp) RandrGetCrtcInfoCookie {
	cookie := C.xcb_randr_get_crtc_info(c.c(), C.xcb_randr_crtc_t(crtc), C.xcb_timestamp_t(configTimestamp))
	return RandrGetCrtcInfoCookie(cookie)
}

func (c *Connection) RandrGetCrtcInfoReply(cookie RandrGetCrtcInfoCookie) (reply *RandrGetCrtcInfoReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_get_crtc_info_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrGetCrtcInfoReply)
		reply.ERandrGetCrtcInfoReply = (*ERandrGetCrtcInfoReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrGetCrtcInfoReply) {
			C.free(unsafe.Pointer(f.ERandrGetCrtcInfoReply))
		})
	} else {
		err = errors.New("RandrGetCrtcInfoReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type RandrOutputs struct {
	Slice []RandrOutput
}

func (c *Connection) RandrGetCrtcInfoReplyOutputs(r *RandrGetCrtcInfoReply) (outputs *RandrOutputs) {
	cOutputs := C.xcb_randr_get_crtc_info_outputs(r.c())
	numOutputs := C.xcb_randr_get_crtc_info_outputs_length(r.c())

	outputs = new(RandrOutputs)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&outputs.Slice))
	sliceHeader.Len = int(numOutputs)
	sliceHeader.Cap = int(numOutputs)
	sliceHeader.Data = uintptr(unsafe.Pointer(cOutputs))
	return
}

type ERandrGetOutputInfoReply struct {
	ResponseType  C.uint8_t
	Status        C.uint8_t
	Sequence      C.uint16_t
	Length        C.uint32_t
	Timestamp     Timestamp
	Crtc          RandrCrtc
	MMWidth       C.uint32_t
	MMHeight      C.uint32_t
	Connection    C.uint8_t
	SubpixelOrder C.uint8_t
	NumCrtcs      C.uint16_t
	NumModes      C.uint16_t
	NumPreferred  C.uint16_t
	NumClones     C.uint16_t
	NameLen       C.uint16_t
}
type RandrGetOutputInfoReply struct {
	*ERandrGetOutputInfoReply
}

func (c *RandrGetOutputInfoReply) c() *C.xcb_randr_get_output_info_reply_t {
	ptr := c.ERandrGetOutputInfoReply
	return (*C.xcb_randr_get_output_info_reply_t)(unsafe.Pointer(ptr))
}

type RandrGetOutputInfoCookie C.xcb_randr_get_output_info_cookie_t

func (c RandrGetOutputInfoCookie) c() C.xcb_randr_get_output_info_cookie_t {
	return C.xcb_randr_get_output_info_cookie_t(c)
}

func (c *Connection) RandrGetOutputInfo(output RandrOutput, configTimestamp Timestamp) RandrGetOutputInfoCookie {
	cookie := C.xcb_randr_get_output_info(c.c(), C.xcb_randr_output_t(output), C.xcb_timestamp_t(configTimestamp))
	return RandrGetOutputInfoCookie(cookie)
}

func (c *Connection) RandrGetOutputInfoReply(cookie RandrGetOutputInfoCookie) (reply *RandrGetOutputInfoReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_randr_get_output_info_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(RandrGetOutputInfoReply)
		reply.ERandrGetOutputInfoReply = (*ERandrGetOutputInfoReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *RandrGetOutputInfoReply) {
			C.free(unsafe.Pointer(f.ERandrGetOutputInfoReply))
		})
	}
	if e != nil {
		err = errors.New("RandrGetOutputInfoReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) RandrGetOutputInfoName(r *RandrGetOutputInfoReply) (name string) {
	cstr := C.xcb_randr_get_output_info_name(r.c())
	len := C.xcb_randr_get_output_info_name_length(r.c())

	var ba []byte
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&ba))
	sliceHeader.Len = int(len)
	sliceHeader.Cap = int(len)
	sliceHeader.Data = uintptr(unsafe.Pointer(cstr))
	name = string(ba)
	return
}

type RandrModes struct {
	Slice []RandrMode
}

func (c *Connection) RandrGetOutputInfoModes(r *RandrGetOutputInfoReply) (modes *RandrModes) {
	cModes := C.xcb_randr_get_output_info_modes(r.c())
	numModes := C.xcb_randr_get_output_info_modes_length(r.c())

	modes = new(RandrModes)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&modes.Slice))
	sliceHeader.Len = int(numModes)
	sliceHeader.Cap = int(numModes)
	sliceHeader.Data = uintptr(unsafe.Pointer(cModes))
	return
}
