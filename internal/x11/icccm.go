// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/xcb_icccm.h>

#cgo LDFLAGS: -lxcb -lxcb-icccm

// xcb_icccm 3.8+ prefixes things with 'icccm', we can support both of we
// are careful.
#ifdef XCB_ICCCM_NUM_WM_SIZE_HINTS_ELEMENTS

#define xcb_set_wm_size_hints xcb_icccm_set_wm_size_hints
#define xcb_wm_hints_t xcb_icccm_wm_hints_t

#define XCB_SIZE_HINT_US_POSITION   XCB_ICCCM_SIZE_HINT_US_POSITION
#define XCB_SIZE_HINT_US_SIZE       XCB_ICCCM_SIZE_HINT_US_SIZE
#define XCB_SIZE_HINT_P_POSITION    XCB_ICCCM_SIZE_HINT_P_POSITION
#define XCB_SIZE_HINT_P_SIZE        XCB_ICCCM_SIZE_HINT_P_SIZE
#define XCB_SIZE_HINT_P_MIN_SIZE    XCB_ICCCM_SIZE_HINT_P_MIN_SIZE
#define XCB_SIZE_HINT_P_MAX_SIZE    XCB_ICCCM_SIZE_HINT_P_MAX_SIZE
#define XCB_SIZE_HINT_P_RESIZE_INC  XCB_ICCCM_SIZE_HINT_P_RESIZE_INC
#define XCB_SIZE_HINT_P_ASPECT      XCB_ICCCM_SIZE_HINT_P_ASPECT
#define XCB_SIZE_HINT_BASE_SIZE     XCB_ICCCM_SIZE_HINT_BASE_SIZE
#define XCB_SIZE_HINT_P_WIN_GRAVITY XCB_ICCCM_SIZE_HINT_P_WIN_GRAVITY

#define XCB_WM_STATE_ICONIC XCB_ICCCM_WM_STATE_ICONIC
#define XCB_WM_STATE_NORMAL XCB_ICCCM_WM_STATE_NORMAL

#endif
*/
import "C"

import (
	"unsafe"
)

const (
	SIZE_HINT_US_POSITION   = C.XCB_SIZE_HINT_US_POSITION
	SIZE_HINT_US_SIZE       = C.XCB_SIZE_HINT_US_SIZE
	SIZE_HINT_P_POSITION    = C.XCB_SIZE_HINT_P_POSITION
	SIZE_HINT_P_SIZE        = C.XCB_SIZE_HINT_P_SIZE
	SIZE_HINT_P_MIN_SIZE    = C.XCB_SIZE_HINT_P_MIN_SIZE
	SIZE_HINT_P_MAX_SIZE    = C.XCB_SIZE_HINT_P_MAX_SIZE
	SIZE_HINT_P_RESIZE_INC  = C.XCB_SIZE_HINT_P_RESIZE_INC
	SIZE_HINT_P_ASPECT      = C.XCB_SIZE_HINT_P_ASPECT
	SIZE_HINT_BASE_SIZE     = C.XCB_SIZE_HINT_BASE_SIZE
	SIZE_HINT_P_WIN_GRAVITY = C.XCB_SIZE_HINT_P_WIN_GRAVITY

	WM_STATE_ICONIC = C.XCB_WM_STATE_ICONIC
	WM_STATE_NORMAL = C.XCB_WM_STATE_NORMAL
)

type Uint32 C.uint32_t
type Int32 C.int32_t

type SizeHints struct {
	Flags                      Uint32
	X, Y                       Int32
	Width, Height              Int32
	MinWidth, MinHeight        Int32
	MaxWidth, MaxHeight        Int32
	WidthInc, HeightInc        Int32
	MinAspectNum, MinAspectDen Int32
	MaxAspectNum, MaxAspectDen Int32
	BaseWidth, BaseHeight      Int32
	WinGravity                 Uint32
}

func (c *Connection) SetWmSizeHints(w Window, hints *SizeHints) {
	C.xcb_set_wm_size_hints(
		c.c(),
		C.xcb_window_t(w),
		C.XCB_ATOM_WM_NORMAL_HINTS,
		(*C.xcb_size_hints_t)(unsafe.Pointer(hints)),
	)
}
