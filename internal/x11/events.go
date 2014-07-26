// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <xcb/xcb.h>

#cgo LDFLAGS: -lxcb
*/
import "C"

const (
	EVENT_MASK_NO_EVENT              = C.XCB_EVENT_MASK_NO_EVENT
	EVENT_MASK_KEY_PRESS             = C.XCB_EVENT_MASK_KEY_PRESS
	EVENT_MASK_KEY_RELEASE           = C.XCB_EVENT_MASK_KEY_RELEASE
	EVENT_MASK_BUTTON_PRESS          = C.XCB_EVENT_MASK_BUTTON_PRESS
	EVENT_MASK_BUTTON_RELEASE        = C.XCB_EVENT_MASK_BUTTON_RELEASE
	EVENT_MASK_ENTER_WINDOW          = C.XCB_EVENT_MASK_ENTER_WINDOW
	EVENT_MASK_LEAVE_WINDOW          = C.XCB_EVENT_MASK_LEAVE_WINDOW
	EVENT_MASK_POINTER_MOTION        = C.XCB_EVENT_MASK_POINTER_MOTION
	EVENT_MASK_POINTER_MOTION_HINT   = C.XCB_EVENT_MASK_POINTER_MOTION_HINT
	EVENT_MASK_BUTTON_1_MOTION       = C.XCB_EVENT_MASK_BUTTON_1_MOTION
	EVENT_MASK_BUTTON_2_MOTION       = C.XCB_EVENT_MASK_BUTTON_2_MOTION
	EVENT_MASK_BUTTON_3_MOTION       = C.XCB_EVENT_MASK_BUTTON_3_MOTION
	EVENT_MASK_BUTTON_4_MOTION       = C.XCB_EVENT_MASK_BUTTON_4_MOTION
	EVENT_MASK_BUTTON_5_MOTION       = C.XCB_EVENT_MASK_BUTTON_5_MOTION
	EVENT_MASK_BUTTON_MOTION         = C.XCB_EVENT_MASK_BUTTON_MOTION
	EVENT_MASK_KEYMAP_STATE          = C.XCB_EVENT_MASK_KEYMAP_STATE
	EVENT_MASK_EXPOSURE              = C.XCB_EVENT_MASK_EXPOSURE
	EVENT_MASK_VISIBILITY_CHANGE     = C.XCB_EVENT_MASK_VISIBILITY_CHANGE
	EVENT_MASK_STRUCTURE_NOTIFY      = C.XCB_EVENT_MASK_STRUCTURE_NOTIFY
	EVENT_MASK_RESIZE_REDIRECT       = C.XCB_EVENT_MASK_RESIZE_REDIRECT
	EVENT_MASK_SUBSTRUCTURE_NOTIFY   = C.XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY
	EVENT_MASK_SUBSTRUCTURE_REDIRECT = C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT
	EVENT_MASK_FOCUS_CHANGE          = C.XCB_EVENT_MASK_FOCUS_CHANGE
	EVENT_MASK_PROPERTY_CHANGE       = C.XCB_EVENT_MASK_PROPERTY_CHANGE
	EVENT_MASK_OWNER_GRAB_BUTTON     = C.XCB_EVENT_MASK_OWNER_GRAB_BUTTON

	KEY_PRESS    = C.XCB_KEY_PRESS
	KEY_RELEASE  = C.XCB_KEY_RELEASE
	BUTTON_PRESS = C.XCB_BUTTON_PRESS

	BUTTON_MASK_1     = C.XCB_BUTTON_MASK_1
	BUTTON_MASK_2     = C.XCB_BUTTON_MASK_2
	BUTTON_MASK_3     = C.XCB_BUTTON_MASK_3
	BUTTON_MASK_4     = C.XCB_BUTTON_MASK_4
	BUTTON_MASK_5     = C.XCB_BUTTON_MASK_5
	BUTTON_MASK_ANY   = C.XCB_BUTTON_MASK_ANY
	BUTTON_RELEASE    = C.XCB_BUTTON_RELEASE
	MOTION_NOTIFY     = C.XCB_MOTION_NOTIFY
	ENTER_NOTIFY      = C.XCB_ENTER_NOTIFY
	LEAVE_NOTIFY      = C.XCB_LEAVE_NOTIFY
	FOCUS_IN          = C.XCB_FOCUS_IN
	FOCUS_OUT         = C.XCB_FOCUS_OUT
	KEYMAP_NOTIFY     = C.XCB_KEYMAP_NOTIFY
	EXPOSE            = C.XCB_EXPOSE
	GRAPHICS_EXPOSURE = C.XCB_GRAPHICS_EXPOSURE
	NO_EXPOSURE       = C.XCB_NO_EXPOSURE
	VISIBILITY_NOTIFY = C.XCB_VISIBILITY_NOTIFY
	CREATE_NOTIFY     = C.XCB_CREATE_NOTIFY
	DESTROY_NOTIFY    = C.XCB_DESTROY_NOTIFY
	UNMAP_NOTIFY      = C.XCB_UNMAP_NOTIFY
	MAP_NOTIFY        = C.XCB_MAP_NOTIFY
	MAP_REQUEST       = C.XCB_MAP_REQUEST
	REPARENT_NOTIFY   = C.XCB_REPARENT_NOTIFY
	CONFIGURE_NOTIFY  = C.XCB_CONFIGURE_NOTIFY
	CONFIGURE_REQUEST = C.XCB_CONFIGURE_REQUEST
	GRAVITY_NOTIFY    = C.XCB_GRAVITY_NOTIFY
	RESIZE_REQUEST    = C.XCB_RESIZE_REQUEST
	CIRCULATE_NOTIFY  = C.XCB_CIRCULATE_NOTIFY
	CIRCULATE_REQUEST = C.XCB_CIRCULATE_REQUEST
	PROPERTY_NOTIFY   = C.XCB_PROPERTY_NOTIFY
	SELECTION_CLEAR   = C.XCB_SELECTION_CLEAR
	SELECTION_REQUEST = C.XCB_SELECTION_REQUEST
	COLORMAP_NOTIFY   = C.XCB_COLORMAP_NOTIFY
	CLIENT_MESSAGE    = C.XCB_CLIENT_MESSAGE
	SELECTION_NOTIFY  = C.XCB_SELECTION_NOTIFY
	MAPPING_NOTIFY    = C.XCB_MAPPING_NOTIFY

	COLORMAP_STATE_UNINSTALLED    = C.XCB_COLORMAP_STATE_UNINSTALLED
	COLORMAP_STATE_INSTALLED      = C.XCB_COLORMAP_STATE_INSTALLED
	PLACE_ON_TOP                  = C.XCB_PLACE_ON_TOP
	PLACE_ON_BOTTOM               = C.XCB_PLACE_ON_BOTTOM
	VISIBILITY_UNOBSCURED         = C.XCB_VISIBILITY_UNOBSCURED
	VISIBILITY_PARTIALLY_OBSCURED = C.XCB_VISIBILITY_PARTIALLY_OBSCURED
	VISIBILITY_FULLY_OBSCURED     = C.XCB_VISIBILITY_FULLY_OBSCURED
	MOTION_NORMAL                 = C.XCB_MOTION_NORMAL
	MOTION_HINT                   = C.XCB_MOTION_HINT

	NOTIFY_DETAIL_ANCESTOR          = C.XCB_NOTIFY_DETAIL_ANCESTOR
	NOTIFY_DETAIL_VIRTUAL           = C.XCB_NOTIFY_DETAIL_VIRTUAL
	NOTIFY_DETAIL_INFERIOR          = C.XCB_NOTIFY_DETAIL_INFERIOR
	NOTIFY_DETAIL_NONLINEAR         = C.XCB_NOTIFY_DETAIL_NONLINEAR
	NOTIFY_DETAIL_NONLINEAR_VIRTUAL = C.XCB_NOTIFY_DETAIL_NONLINEAR_VIRTUAL
	NOTIFY_DETAIL_POINTER           = C.XCB_NOTIFY_DETAIL_POINTER
	NOTIFY_DETAIL_POINTER_ROOT      = C.XCB_NOTIFY_DETAIL_POINTER_ROOT
	NOTIFY_DETAIL_NONE              = C.XCB_NOTIFY_DETAIL_NONE

	NOTIFY_MODE_NORMAL        = C.XCB_NOTIFY_MODE_NORMAL
	NOTIFY_MODE_GRAB          = C.XCB_NOTIFY_MODE_GRAB
	NOTIFY_MODE_UNGRAB        = C.XCB_NOTIFY_MODE_UNGRAB
	NOTIFY_MODE_WHILE_GRABBED = C.XCB_NOTIFY_MODE_WHILE_GRABBED

	MAPPING_MODIFIER = C.XCB_MAPPING_MODIFIER
	MAPPING_KEYBOARD = C.XCB_MAPPING_KEYBOARD
	MAPPING_POINTER  = C.XCB_MAPPING_POINTER

	MOD_MASK_SHIFT   = C.XCB_MOD_MASK_SHIFT
	MOD_MASK_LOCK    = C.XCB_MOD_MASK_LOCK
	MOD_MASK_CONTROL = C.XCB_MOD_MASK_CONTROL
	MOD_MASK_1       = C.XCB_MOD_MASK_1
	MOD_MASK_2       = C.XCB_MOD_MASK_2
	MOD_MASK_3       = C.XCB_MOD_MASK_3
	MOD_MASK_4       = C.XCB_MOD_MASK_4
	MOD_MASK_5       = C.XCB_MOD_MASK_5
	MOD_MASK_ANY     = C.XCB_MOD_MASK_ANY
)

type KeyPressEvent struct {
	ResponseType C.uint8_t
	Detail       Keycode
	Sequence     C.uint16_t
	Time         Timestamp
	Root         Window
	Event        Window
	Child        Window
	RootX        C.int16_t
	RootY        C.int16_t
	EventX       C.int16_t
	EventY       C.int16_t
	State        C.uint16_t
	SameScreen   C.uint8_t
	Pad0         C.uint8_t
}

type KeyReleaseEvent KeyPressEvent

type ButtonPressEvent struct {
	ResponseType C.uint8_t
	Detail       Button
	Sequence     C.uint16_t
	Time         Timestamp
	Root         Window
	Event        Window
	Child        Window
	RootX        C.int16_t
	RootY        C.int16_t
	EventX       C.int16_t
	EventY       C.int16_t
	State        C.uint16_t
	SameScreen   C.uint8_t
	Pad0         C.uint8_t
}

type ButtonReleaseEvent ButtonPressEvent

type MotionNotifyEvent struct {
	ResponseType C.uint8_t
	Detail       C.uint8_t
	Sequence     C.uint16_t
	Time         Timestamp
	Root         Window
	Event        Window
	Child        Window
	RootX        C.int16_t
	RootY        C.int16_t
	EventX       C.int16_t
	EventY       C.int16_t
	State        C.uint16_t
	SameScreen   C.uint8_t
	Pad0         C.uint8_t
}

type EnterNotifyEvent struct {
	ResponseType    C.uint8_t
	Detail          C.uint8_t
	Sequence        C.uint16_t
	Time            Timestamp
	Root            Window
	Event           Window
	Child           Window
	RootX           C.int16_t
	RootY           C.int16_t
	EventX          C.int16_t
	EventY          C.int16_t
	State           C.uint16_t
	Mode            C.uint8_t
	SameScreenFocus C.uint8_t
}

type LeaveNotifyEvent EnterNotifyEvent

type FocusInEvent struct {
	ResponseType C.uint8_t
	Detail       C.uint8_t
	Sequence     C.uint16_t
	Event        Window
	Mode         C.uint8_t
	Pad0         [3]C.uint8_t
}

type FocusOutEvent FocusInEvent

type KeymapNotifyEvent struct {
	ResponseType C.uint8_t
	Keys         [31]C.uint8_t
}

type ExposeEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	X            C.uint16_t
	Y            C.uint16_t
	Width        C.uint16_t
	Height       C.uint16_t
	Count        C.uint16_t
	Pad1         [2]C.uint8_t
}

type GraphicsExposureEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Drawable     Drawable
	X            C.uint16_t
	Y            C.uint16_t
	Width        C.uint16_t
	Height       C.uint16_t
	MinorOpcode  C.uint16_t
	Count        C.uint16_t
	MajorOpcode  C.uint16_t
	Pad1         [3]C.uint8_t
}

type NoExposureEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Drawable     Drawable
	MajorOpcode  C.uint16_t
	MinorOpcode  C.uint8_t
	Pad1         C.uint8_t
}

type VisibilityNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	State        C.uint8_t
	Pad1         [3]C.uint8_t
}

type CreateNotifyEvent struct {
	ResponseType     C.uint8_t
	Pad0             C.uint8_t
	Sequence         C.uint16_t
	Parent           Window
	Window           Window
	X                C.int16_t
	Y                C.int16_t
	Width            C.uint16_t
	Height           C.uint16_t
	BorderWidth      C.uint16_t
	OverrideRedirect C.uint8_t
	Pad1             C.uint8_t
}

type DestroyNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Event        Window
	Window       Window
}

type UnmapNotifyEvent struct {
	ResponseType  C.uint8_t
	Pad0          C.uint8_t
	Sequence      C.uint16_t
	Event         Window
	Window        Window
	FromConfigure C.uint8_t
	Pad1          [3]C.uint8_t
}

type MapNotifyEvent struct {
	ResponseType     C.uint8_t
	Pad0             C.uint8_t
	Sequence         C.uint16_t
	Event            Window
	Window           Window
	OverrideRedirect C.uint8_t
	Pad1             [3]C.uint8_t
}

type MapRequestEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Parent       Window
	Window       Window
}

type ReparentNotifyEvent struct {
	ResponseType     C.uint8_t
	Pad0             C.uint8_t
	Sequence         C.uint16_t
	Event            Window
	Window           Window
	Parent           Window
	X                C.int16_t
	Y                C.int16_t
	OverrideRedirect C.uint8_t
	Pad1             [3]C.uint8_t
}

type ConfigureNotifyEvent struct {
	ResponseType     C.uint8_t
	Pad0             C.uint8_t
	Sequence         C.uint16_t
	Event            Window
	Window           Window
	AboveSibling     Window
	X                C.int16_t
	Y                C.int16_t
	Width            C.uint16_t
	Height           C.uint16_t
	BorderWidth      C.uint16_t
	OverrideRedirect C.uint8_t
	Pad1             C.uint8_t
}

type ConfigureRequestEvent struct {
	ResponseType C.uint8_t
	StackMode    C.uint8_t
	Sequence     C.uint16_t
	Parent       Window
	Window       Window
	Sibling      Window
	X            C.int16_t
	Y            C.int16_t
	Width        C.uint16_t
	Height       C.uint16_t
	BorderWidth  C.uint16_t
	ValueMask    C.uint16_t
}

type GravityNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Event        Window
	Window       Window
	X            C.int16_t
	Y            C.int16_t
}

type ResizeRequestEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	Width        C.uint16_t
	Height       C.uint16_t
}

type CirculateNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Event        Window
	Window       Window
	Pad1         [4]C.uint8_t
	Place        C.uint8_t
	Pad2         [3]C.uint8_t
}

type CirculateRequestEvent CirculateNotifyEvent

type PropertyNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	Atom         Atom
	Time         Timestamp
	State        C.uint8_t
	Pad1         [3]C.uint8_t
}

type SelectionClearEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Time         Timestamp
	Owner        Window
	Selection    Atom
}

type SelectionRequestEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Time         Timestamp
	Owner        Window
	Requestor    Window
	Selection    Atom
	Target       Atom
	Property     Atom
}

type ColormapNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	Colormap     Colormap
	New          C.uint8_t
	State        C.uint8_t
	Pad1         [2]C.uint8_t
}

type SelectionNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Time         Timestamp
	Requestor    Window
	Selection    Atom
	Target       Atom
	Property     Atom
}

type ClientMessageEvent struct {
	ResponseType C.uint8_t
	Format       C.uint8_t
	Sequence     C.uint16_t
	Window       Window
	Type         Atom
	Data         [20]byte
}

type MappingNotifyEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Request      C.uint8_t
	FirstKeycode Keycode
	Count        C.uint8_t
	Pad1         C.uint8_t
}
