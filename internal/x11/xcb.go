// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

// Yes, I know it's ugly to use both Xlib and XCB; but suprisingly it is not
// really possible to just use XCB for what we are doing. For one thing GLX
// does not work with XCB (without mixing in Xlib, of course). For another,
// XCB's keyboard support is rudimentary at best, e.g. Only Xlib: XLookupKeysym
// supports UCS encoding, etc..
//
// It looks like this is *the way to do it*, and it looks like other
// open source projects do it this way too; so, there is that as well.

/*
#include <stdlib.h>
#include <string.h>
#include <X11/Xlib-xcb.h>

#cgo LDFLAGS: -lxcb

void chippy_send_client_message(xcb_connection_t *c, xcb_window_t window, xcb_window_t dest, xcb_atom_t atom, uint32_t data_len, uint32_t *data) {
	xcb_client_message_event_t ev;
	memset(&ev, 0, sizeof(xcb_client_message_event_t));

	ev.response_type = XCB_CLIENT_MESSAGE;
	ev.window = window;
	ev.format = 32;
	ev.type = atom;

	for(; data_len != 0; data_len--) {
		ev.data.data32[0] = data[1];
	}

	xcb_send_event(c, 0, dest, XCB_EVENT_MASK_SUBSTRUCTURE_NOTIFY|XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT, (char *) &ev);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

var (
	ConnError                 = errors.New("XCB_CONN_ERROR: socket, pipe, or other stream error occurred.")
	ConnClosedExtNotSupported = errors.New("XCB_CONN_CLOSED_EXT_NOTSUPPORTED: extension is not supported.")
	ConnClosedMemInsufficient = errors.New("XCB_CONN_CLOSED_MEM_INSUFFICIENT: memory not available.")
	ConnClosedReqLenExceed    = errors.New("XCB_CONN_CLOSED_REQ_LEN_EXCEED: exceeded request length that server accepts.")
	ConnClosedParseErr        = errors.New("XCB_CONN_CLOSED_PARSE_ERR: error during parsing display string.")
	ConnClosedInvalidScreen   = errors.New("XCB_CONN_CLOSED_INVALID_SCREEN: the server does not have a screen matching the display.")
)

const (
	TIME_CURRENT_TIME         = C.XCB_TIME_CURRENT_TIME
	COPY_FROM_PARENT          = C.XCB_COPY_FROM_PARENT
	WINDOW_CLASS_INPUT_OUTPUT = C.XCB_WINDOW_CLASS_INPUT_OUTPUT
	PROP_MODE_REPLACE         = C.XCB_PROP_MODE_REPLACE

	ATOM_ANY      = C.XCB_ATOM_ANY
	ATOM_CARDINAL = C.XCB_ATOM_CARDINAL
	ATOM_WM_NAME  = C.XCB_ATOM_WM_NAME
	ATOM_STRING   = C.XCB_ATOM_STRING
	ATOM_ATOM     = C.XCB_ATOM_ATOM

	CONFIG_WINDOW_X          = C.XCB_CONFIG_WINDOW_X
	CONFIG_WINDOW_Y          = C.XCB_CONFIG_WINDOW_Y
	CONFIG_WINDOW_WIDTH      = C.XCB_CONFIG_WINDOW_WIDTH
	CONFIG_WINDOW_HEIGHT     = C.XCB_CONFIG_WINDOW_HEIGHT
	CONFIG_WINDOW_STACK_MODE = C.XCB_CONFIG_WINDOW_STACK_MODE

	STACK_MODE_ABOVE = C.XCB_STACK_MODE_ABOVE
	CW_EVENT_MASK    = C.XCB_CW_EVENT_MASK
	CW_BACK_PIXEL    = C.XCB_CW_BACK_PIXEL
	CW_BORDER_PIXEL  = C.XCB_CW_BORDER_PIXEL
	CW_COLORMAP      = C.XCB_CW_COLORMAP
	CW_CURSOR        = C.XCB_CW_CURSOR

	GC_FOREGROUND = C.XCB_GC_FOREGROUND
	GC_BACKGROUND = C.XCB_GC_BACKGROUND

	IMAGE_FORMAT_Z_PIXMAP  = C.XCB_IMAGE_FORMAT_Z_PIXMAP
	IMAGE_FORMAT_XY_PIXMAP = C.XCB_IMAGE_FORMAT_XY_PIXMAP

	IMAGE_ORDER_LSB_FIRST = C.XCB_IMAGE_ORDER_LSB_FIRST
	IMAGE_ORDER_MSB_FIRST = C.XCB_IMAGE_ORDER_MSB_FIRST
)

type RequestError struct {
	ResponseType C.uint8_t
	ErrorCode    C.uint8_t
	Sequence     C.uint16_t
	BadValue     C.uint32_t
	MinorOpcode  C.uint16_t
	MajorOpcode  C.uint8_t
	Pad0         C.uint8_t
}

type MotifWMHints struct {
	Flags       C.ulong
	Functions   C.ulong
	Decorations C.ulong
	InputMode   C.long
	Status      C.ulong
}

func xcbError(e *C.xcb_generic_error_t) string {
	return fmt.Sprintf("XCB error=%d, sequence=%d, resource=%d, minor=%d, major=%d", e.error_code, e.sequence, e.resource_id, e.minor_code, e.major_code)
}

func findError(code C.int) error {
	switch code {
	case C.XCB_CONN_ERROR:
		return ConnError
	case C.XCB_CONN_CLOSED_EXT_NOTSUPPORTED:
		return ConnClosedExtNotSupported
	case C.XCB_CONN_CLOSED_MEM_INSUFFICIENT:
		return ConnClosedMemInsufficient
	case C.XCB_CONN_CLOSED_REQ_LEN_EXCEED:
		return ConnClosedReqLenExceed
	case C.XCB_CONN_CLOSED_PARSE_ERR:
		return ConnClosedParseErr
	case 6: //Only in some newer versions: C.XCB_CONN_CLOSED_INVALID_SCREEN:
		return ConnClosedInvalidScreen
	}
	return nil
}

type (
	Window           C.xcb_window_t
	Connection       C.xcb_connection_t
	Colormap         C.xcb_colormap_t
	VisualId         C.xcb_visualid_t
	Timestamp        C.xcb_timestamp_t
	VoidCookie       C.xcb_void_cookie_t
	Atom             C.xcb_atom_t
	InternAtomCookie C.xcb_intern_atom_cookie_t
	Keycode          C.xcb_keycode_t
	Keysym           C.xcb_keysym_t
	Button           C.xcb_button_t
	Cursor           C.xcb_cursor_t

	Drawable C.xcb_drawable_t
	Long     C.long
)

func (c *InternAtomCookie) c() *C.xcb_intern_atom_cookie_t {
	return (*C.xcb_intern_atom_cookie_t)(unsafe.Pointer(c))
}

type Screen struct {
	Root                Window
	DefaultColormap     Colormap
	WhitePixel          C.uint32_t
	BlackPixel          C.uint32_t
	CurrentInputMasks   C.uint32_t
	WidthInPixels       C.uint16_t
	HeightInPixels      C.uint16_t
	WidthInMillimeters  C.uint16_t
	HeightInMillimeters C.uint16_t
	MinInstalledMaps    C.uint16_t
	MaxInstalledMaps    C.uint16_t
	RootVisual          VisualId
	BackingStores       C.uint8_t
	SaveUnders          C.uint8_t
	RootDepth           C.uint8_t
	AllowedDepthsLen    C.uint8_t
}

func (s *Screen) c() *C.xcb_screen_t {
	return (*C.xcb_screen_t)(unsafe.Pointer(s))
}

func (c *Connection) c() *C.xcb_connection_t {
	return (*C.xcb_connection_t)(unsafe.Pointer(c))
}

func (c *Connection) GetFileDescriptor() int {
	return int(C.xcb_get_file_descriptor(c.c()))
}

func (c *Connection) HasError() error {
	return findError(C.xcb_connection_has_error(c.c()))
}

func (c *Connection) Disconnect() {
	C.xcb_disconnect(c.c())
}

func (c *Connection) ScreenOfDisplay(screen int) *Screen {
	iter := C.xcb_setup_roots_iterator(C.xcb_get_setup(c.c()))
	for iter.rem != 0 {
		if screen == 0 {
			return (*Screen)(unsafe.Pointer(iter.data))
		}

		C.xcb_screen_next(&iter)
		screen--
	}
	return nil
}

type Setup struct {
	Status                   C.uint8_t
	Pad0                     C.uint8_t
	ProtocolMajorVersion     C.uint16_t
	ProtocolMinorVersion     C.uint16_t
	Length                   C.uint16_t
	ReleaseNumber            C.uint32_t
	ResourceIdBase           C.uint32_t
	ResourceIdMask           C.uint32_t
	MotionBufferSize         C.uint32_t
	VenderLen                C.uint16_t
	MaximumRequestLength     C.uint16_t
	RootsLen                 C.uint8_t
	PixmapFormatsLen         C.uint8_t
	ImageByteOrder           C.uint8_t
	BitmapFormatBitOrder     C.uint8_t
	BitmapFormatScanlineUnit C.uint8_t
	BitmapFormatScanlinePad  C.uint8_t
	MinKeycode               Keycode
	MaxKeycode               Keycode
	Pad1                     [4]C.uint8_t
}

func (s *Setup) c() *C.xcb_setup_t {
	return (*C.xcb_setup_t)(unsafe.Pointer(s))
}
func (c *Connection) GetSetup() *Setup {
	csetup := C.xcb_get_setup(c.c())
	return (*Setup)(unsafe.Pointer(csetup))
}

func (c *Connection) SetupRootsLength(r *Setup) int {
	return int(C.xcb_setup_roots_length(r.c()))
}

func (c *Connection) RootsIterator(callback func(s *Screen) (again bool)) {
	iter := C.xcb_setup_roots_iterator(C.xcb_get_setup(c.c()))
	for iter.rem != 0 {
		s := (*Screen)(unsafe.Pointer(iter.data))
		if !callback(s) {
			break
		}
		C.xcb_screen_next(&iter)
	}
}

type Depth struct {
	Depth      C.uint8_t
	Pad0       C.uint8_t
	VisualsLen C.uint16_t
	Pad1       [4]C.uint8_t
}

func (c *Connection) ScreenAllowedDepthsIterator(s *Screen, callback func(d *Depth) (again bool)) {
	iter := C.xcb_screen_allowed_depths_iterator((*C.xcb_screen_t)(unsafe.Pointer(s)))
	for iter.rem != 0 {
		d := (*Depth)(unsafe.Pointer(iter.data))
		if !callback(d) {
			return
		}
		C.xcb_depth_next(&iter)
	}
}

type Format struct {
	Depth        C.uint8_t
	BitsPerPixel C.uint8_t
	ScanlinePad  C.uint8_t
	Pad0         [5]C.uint8_t
}

func (c *Connection) PixmapFormatsIterator(callback func(*Format) (again bool)) {
	iter := C.xcb_setup_pixmap_formats_iterator(C.xcb_get_setup(c.c()))
	for iter.rem != 0 {
		f := (*Format)(unsafe.Pointer(iter.data))
		if !callback(f) {
			return
		}
		C.xcb_format_next(&iter)
	}
}

func (c *Connection) FindPixmapFormat(depth, bpp uint8) (foundFormat *Format) {
	c.PixmapFormatsIterator(func(f *Format) bool {
		if uint8(f.Depth) == depth && uint8(f.BitsPerPixel) == bpp {
			foundFormat = f
			return false
		}
		return true
	})
	return
}

type VisualType struct {
	VisualId        VisualId
	Class           C.uint8_t
	BitsPerRgbValue C.uint8_t
	ColormapEntries C.uint16_t
	RedMask         C.uint32_t
	GreenMask       C.uint32_t
	BlueMask        C.uint32_t
	Pad0            [4]C.uint8_t
}

func (d *Depth) Iterate(callback func(vis *VisualType) (again bool)) {
	iter := C.xcb_depth_visuals_iterator((*C.xcb_depth_t)(unsafe.Pointer(d)))
	for iter.rem != 0 {
		vis := (*VisualType)(unsafe.Pointer(iter.data))
		if !callback(vis) {
			return
		}
		C.xcb_visualtype_next(&iter)
	}
}

func (c *Connection) CreateColormapChecked(allocAll bool, mid Colormap, w Window, v VisualId) VoidCookie {
	calloc := C.XCB_COLORMAP_ALLOC_NONE
	if allocAll {
		calloc = C.XCB_COLORMAP_ALLOC_ALL
	}
	return VoidCookie(C.xcb_create_colormap_checked(
		c.c(),
		C.uint8_t(calloc),
		C.xcb_colormap_t(mid),
		C.xcb_window_t(w),
		C.xcb_visualid_t(v),
	))
}

type (
	Pixmap   C.xcb_pixmap_t
	GContext C.xcb_gcontext_t
)

func (c *Connection) CreatePixmapChecked(depth uint8, pid Pixmap, d Drawable, width, height uint16) VoidCookie {
	return VoidCookie(C.xcb_create_pixmap_checked(
		c.c(),
		C.uint8_t(depth),
		C.xcb_pixmap_t(pid),
		C.xcb_drawable_t(d),
		C.uint16_t(width),
		C.uint16_t(height),
	))
}

func (c *Connection) CopyArea(src, dst Drawable, gc GContext, srcX, srcY, dstX, dstY int16, width, height uint16) VoidCookie {
	return VoidCookie(C.xcb_copy_area(
		c.c(),
		C.xcb_drawable_t(src),
		C.xcb_drawable_t(dst),
		C.xcb_gcontext_t(gc),
		C.int16_t(srcX),
		C.int16_t(srcY),
		C.int16_t(dstX),
		C.int16_t(dstY),
		C.uint16_t(width),
		C.uint16_t(height),
	))
}

func (c *Connection) PutImageChecked(format uint8, d Drawable, gc GContext, width, height uint16, dstX, dstY int16, leftPad, depth uint8, dataLen uint32, data unsafe.Pointer) VoidCookie {
	return VoidCookie(C.xcb_put_image_checked(
		c.c(),
		C.uint8_t(format),
		C.xcb_drawable_t(d),
		C.xcb_gcontext_t(gc),
		C.uint16_t(width),
		C.uint16_t(height),
		C.int16_t(dstX),
		C.int16_t(dstY),
		C.uint8_t(leftPad),
		C.uint8_t(depth),
		C.uint32_t(dataLen),
		(*C.uint8_t)(data),
	))
}

func (c *Connection) CreateGCChecked(cid GContext, d Drawable, valueMask uint32, values *uint32) VoidCookie {
	return VoidCookie(C.xcb_create_gc_checked(
		c.c(),
		C.xcb_gcontext_t(cid),
		C.xcb_drawable_t(d),
		C.uint32_t(valueMask),
		(*C.uint32_t)(unsafe.Pointer(values)),
	))
}

func (c *Connection) ChangeGCChecked(cid GContext, valueMask uint32, values *uint32) VoidCookie {
	return VoidCookie(C.xcb_change_gc_checked(
		c.c(),
		C.xcb_gcontext_t(cid),
		C.uint32_t(valueMask),
		(*C.uint32_t)(unsafe.Pointer(values)),
	))
}

func (c *Connection) FreeGC(gc GContext) VoidCookie {
	return VoidCookie(C.xcb_free_gc(
		c.c(),
		C.xcb_gcontext_t(gc),
	))
}

func (c *Connection) FreeCursor(cursor Cursor) VoidCookie {
	return VoidCookie(C.xcb_free_cursor(
		c.c(),
		C.xcb_cursor_t(cursor),
	))
}

func (c *Connection) FreePixmap(pixmap Pixmap) VoidCookie {
	return VoidCookie(C.xcb_free_pixmap(
		c.c(),
		C.xcb_pixmap_t(pixmap),
	))
}

func (c *Connection) GenerateId() Window {
	return Window(C.xcb_generate_id(c.c()))
}

func (c *Connection) CreateWindowChecked(depth uint8, wid, parent Window, x, y int16, width, height, border, class uint16, visual VisualId, valueMask uint32, valueList *uint32) VoidCookie {
	return VoidCookie(C.xcb_create_window_checked(
		c.c(), C.uint8_t(depth),
		C.xcb_window_t(wid),
		C.xcb_window_t(parent),
		C.int16_t(x), C.int16_t(y),
		C.uint16_t(width), C.uint16_t(height),
		C.uint16_t(border), C.uint16_t(class),
		C.xcb_visualid_t(visual),
		C.uint32_t(valueMask),
		(*C.uint32_t)(unsafe.Pointer(valueList)),
	))
}

func (c *Connection) DestroyWindow(w Window) VoidCookie {
	return VoidCookie(C.xcb_destroy_window(
		c.c(),
		C.xcb_window_t(w),
	))
}

func (c *Connection) MapWindow(w Window) VoidCookie {
	return VoidCookie(C.xcb_map_window(c.c(), C.xcb_window_t(w)))
}

func (c *Connection) UnmapWindow(w Window) VoidCookie {
	return VoidCookie(C.xcb_unmap_window(c.c(), C.xcb_window_t(w)))
}

func (c *Connection) RequestCheck(cookie VoidCookie) (err error) {
	e := C.xcb_request_check(c.c(), C.xcb_void_cookie_t(cookie))
	if e != nil {
		err = errors.New(xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) ConfigureWindow(w Window, flags uint16, values []uint32) {
	C.xcb_configure_window(c.c(), C.xcb_window_t(w), C.uint16_t(flags), (*C.uint32_t)(unsafe.Pointer(&values[0])))
}

func (c *Connection) InternAtom(onlyIfExists bool, name string) Atom {
	cOnlyIfExists := C.uint8_t(0)
	if onlyIfExists {
		cOnlyIfExists = 1
	}
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	cookie := C.xcb_intern_atom(c.c(), cOnlyIfExists, C.uint16_t(len(name)), cstr)

	var e *C.xcb_generic_error_t
	reply := C.xcb_intern_atom_reply(c.c(), cookie, &e)
	if e != nil || reply == nil {
		return 0
	}
	return Atom(reply.atom)
}

type EGenericEvent struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Pad          [7]C.uint32_t
	FullSequence C.uint32_t
}
type GenericEvent struct {
	*EGenericEvent
	Free func()
}

func (c *Connection) WaitForEvent() (ev *GenericEvent) {
	cev := C.xcb_wait_for_event(c.c())
	if cev != nil {
		ev = new(GenericEvent)
		ev.EGenericEvent = (*EGenericEvent)(unsafe.Pointer(cev))
		ev.Free = func() {
			C.free(unsafe.Pointer(ev.EGenericEvent))
		}
	}
	return
}

func (c *Connection) ChangeWindowAttributes(w Window, mask uint32, values *uint32) VoidCookie {
	return VoidCookie(C.xcb_change_window_attributes(
		c.c(),
		C.xcb_window_t(w),
		C.uint32_t(mask),
		(*C.uint32_t)(unsafe.Pointer(values)),
	))
}

func (c *Connection) ChangeProperty(mode uint8, w Window, prop Atom, Type Atom, format uint8, dataLen uint32, data unsafe.Pointer) VoidCookie {
	return VoidCookie(C.xcb_change_property(
		c.c(), C.uint8_t(mode),
		C.xcb_window_t(w),
		C.xcb_atom_t(prop),
		C.xcb_atom_t(Type),
		C.uint8_t(format),
		C.uint32_t(dataLen), data,
	))
}

func (c *Connection) ChangePropertyChecked(mode uint8, w Window, prop Atom, Type Atom, format uint8, dataLen uint32, data unsafe.Pointer) VoidCookie {
	return VoidCookie(C.xcb_change_property_checked(
		c.c(), C.uint8_t(mode),
		C.xcb_window_t(w),
		C.xcb_atom_t(prop),
		C.xcb_atom_t(Type),
		C.uint8_t(format),
		C.uint32_t(dataLen), data,
	))
}

func (c *Connection) WarpPointer(src, dst Window, srcX, srcY int16, srcWidth, srcHeight uint16, dstX, dstY int16) VoidCookie {
	return VoidCookie(C.xcb_warp_pointer(
		c.c(),
		C.xcb_window_t(src),
		C.xcb_window_t(dst),
		C.int16_t(srcX),
		C.int16_t(srcY),
		C.uint16_t(srcWidth),
		C.uint16_t(srcHeight),
		C.int16_t(dstX),
		C.int16_t(dstY),
	))
}

func (c *Connection) DeleteProperty(w Window, prop Atom) {
	C.xcb_delete_property(c.c(), C.xcb_window_t(w), C.xcb_atom_t(prop))
}

type EGetPropertyReply struct {
	ResponseType C.uint8_t
	Format       C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
	Type         Atom
	BytesAfter   C.uint32_t
	ValueLen     C.uint32_t
	Pad0         [12]C.uint8_t
}
type GetPropertyReply struct {
	*EGetPropertyReply
}

func (c *Connection) GetProperty(delete bool, w Window, prop Atom, Type Atom, offset uint32, length uint32) (reply *GetPropertyReply, ptr unsafe.Pointer, ptrLen int, err error) {
	cDelete := C.uint8_t(0)
	if delete {
		cDelete = 1
	}
	cookie := C.xcb_get_property(
		c.c(), cDelete,
		C.xcb_window_t(w),
		C.xcb_atom_t(prop),
		C.xcb_atom_t(Type),
		C.uint32_t(offset),
		C.uint32_t(length),
	)
	var e *C.xcb_generic_error_t
	cReply := C.xcb_get_property_reply(c.c(), cookie, &e)
	if e == nil {
		ptr = C.xcb_get_property_value(cReply)
		ptrLen = int(C.xcb_get_property_value_length(cReply))

		reply = new(GetPropertyReply)
		reply.EGetPropertyReply = (*EGetPropertyReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *GetPropertyReply) {
			C.free(unsafe.Pointer(f.EGetPropertyReply))
		})

	} else {
		err = errors.New("XcbGetPropertyReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) SendEvent(propogate bool, dest Window, eventMask uint32, event unsafe.Pointer) {
	cPropogate := C.uint8_t(0)
	if propogate {
		cPropogate = 1
	}

	C.xcb_send_event(
		c.c(), cPropogate,
		C.xcb_window_t(dest),
		C.uint32_t(eventMask),
		(*C.char)(event),
	)
}

func (c *Connection) SendClientMessage(w, dest Window, atom Atom, dataLen uint32, data *uint32) {
	C.chippy_send_client_message(
		c.c(),
		C.xcb_window_t(w),
		C.xcb_window_t(dest),
		C.xcb_atom_t(atom),
		C.uint32_t(dataLen),
		(*C.uint32_t)(unsafe.Pointer(data)),
	)
}

func (c *Connection) Flush() {
	C.xcb_flush(c.c())
}

func (c *Connection) QueryExtension(name string) bool {
	cstr := []byte(name)
	cookie := C.xcb_query_extension(c.c(), C.uint16_t(len(cstr)), (*C.char)(unsafe.Pointer(&cstr[0])))
	var e *C.xcb_generic_error_t
	reply := C.xcb_query_extension_reply(c.c(), cookie, &e)
	if e != nil {
		return false
	}
	return reply.present == 1
}

func Connect(displayName string, screen *int) *Connection {
	var cstr *C.char
	if len(displayName) > 0 {
		cstr = C.CString(displayName)
		defer C.free(unsafe.Pointer(cstr))
	}
	ret := C.xcb_connect(cstr, (*C.int)(unsafe.Pointer(screen)))
	if ret == nil {
		return nil
	}
	return (*Connection)(unsafe.Pointer(ret))
}

type ETranslateCoordinatesReply struct {
	ResponseType C.uint8_t
	SameScreen   C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
	Child        Window
	DstX         C.uint16_t
	DstY         C.uint16_t
}
type TranslateCoordinatesReply struct {
	*ETranslateCoordinatesReply
}

func (c *TranslateCoordinatesReply) c() *C.xcb_translate_coordinates_reply_t {
	ptr := c.ETranslateCoordinatesReply
	return (*C.xcb_translate_coordinates_reply_t)(unsafe.Pointer(ptr))
}

type TranslateCoordinatesCookie C.xcb_translate_coordinates_cookie_t

func (c TranslateCoordinatesCookie) c() C.xcb_translate_coordinates_cookie_t {
	return C.xcb_translate_coordinates_cookie_t(c)
}

func (c *Connection) TranslateCoordinates(srcWindow, dstWindow Window, srcX, srcY int16) TranslateCoordinatesCookie {
	cookie := C.xcb_translate_coordinates(
		c.c(),
		C.xcb_window_t(srcWindow),
		C.xcb_window_t(dstWindow),
		C.int16_t(srcX), C.int16_t(srcY),
	)
	return TranslateCoordinatesCookie(cookie)
}

func (c *Connection) TranslateCoordinatesReply(cookie TranslateCoordinatesCookie) (reply *TranslateCoordinatesReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_translate_coordinates_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(TranslateCoordinatesReply)
		reply.ETranslateCoordinatesReply = (*ETranslateCoordinatesReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *TranslateCoordinatesReply) {
			C.free(unsafe.Pointer(f.ETranslateCoordinatesReply))
		})
	}
	if e != nil {
		err = errors.New("TranslateCoordinatesReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type EGetKeyboardMappingReply struct {
	ResponseType      C.uint8_t
	KeysymsPerKeycode C.uint8_t
	Sequence          C.uint16_t
	Length            C.uint32_t
	Pad0              [24]C.uint8_t
}
type GetKeyboardMappingReply struct {
	*EGetKeyboardMappingReply
}

func (c *GetKeyboardMappingReply) c() *C.xcb_get_keyboard_mapping_reply_t {
	ptr := c.EGetKeyboardMappingReply
	return (*C.xcb_get_keyboard_mapping_reply_t)(unsafe.Pointer(ptr))
}

type GetKeyboardMappingCookie C.xcb_get_keyboard_mapping_cookie_t

func (c GetKeyboardMappingCookie) c() C.xcb_get_keyboard_mapping_cookie_t {
	return C.xcb_get_keyboard_mapping_cookie_t(c)
}

func (c *Connection) GetKeyboardMapping(firstKeycode Keycode, count uint8) GetKeyboardMappingCookie {
	cookie := C.xcb_get_keyboard_mapping(
		c.c(),
		C.xcb_keycode_t(firstKeycode),
		C.uint8_t(count),
	)
	return GetKeyboardMappingCookie(cookie)
}

func (c *Connection) GetKeyboardMappingReply(cookie GetKeyboardMappingCookie) (reply *GetKeyboardMappingReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_get_keyboard_mapping_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(GetKeyboardMappingReply)
		reply.EGetKeyboardMappingReply = (*EGetKeyboardMappingReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *GetKeyboardMappingReply) {
			C.free(unsafe.Pointer(f.EGetKeyboardMappingReply))
		})
	}
	if e != nil {
		err = errors.New("GetKeyboardMappingReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

type Keysyms struct {
	Slice []Keysym
}

func (c *Connection) GetKeyboardMappingKeysyms(r *GetKeyboardMappingReply) (keysyms *Keysyms) {
	cKeysyms := C.xcb_get_keyboard_mapping_keysyms(r.c())
	numKeysyms := C.xcb_get_keyboard_mapping_keysyms_length(r.c())

	keysyms = new(Keysyms)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&keysyms.Slice))
	sliceHeader.Len = int(numKeysyms)
	sliceHeader.Cap = int(numKeysyms)
	sliceHeader.Data = uintptr(unsafe.Pointer(cKeysyms))
	return
}

const (
	GRAB_MODE_ASYNC = C.XCB_GRAB_MODE_ASYNC
)

type EGrabPointerReply struct {
	ResponseType C.uint8_t
	Status       C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
}
type GrabPointerReply struct {
	*EGrabPointerReply
}

func (c *GrabPointerReply) c() *C.xcb_grab_pointer_reply_t {
	ptr := c.EGrabPointerReply
	return (*C.xcb_grab_pointer_reply_t)(unsafe.Pointer(ptr))
}

type GrabPointerCookie C.xcb_grab_pointer_cookie_t

func (c GrabPointerCookie) c() C.xcb_grab_pointer_cookie_t {
	return C.xcb_grab_pointer_cookie_t(c)
}

func (c *Connection) GrabPointer(ownerEvents uint8, grabEvents Window, eventMask uint16, pointerMode, keyboardMode uint8, confineTo Window, cursor Cursor, time Timestamp) GrabPointerCookie {
	cookie := C.xcb_grab_pointer(
		c.c(),
		C.uint8_t(ownerEvents),
		C.xcb_window_t(grabEvents),
		C.uint16_t(eventMask),
		C.uint8_t(pointerMode),
		C.uint8_t(keyboardMode),
		C.xcb_window_t(confineTo),
		C.xcb_cursor_t(cursor),
		C.xcb_timestamp_t(time),
	)
	return GrabPointerCookie(cookie)
}

func (c *Connection) GrabPointerReply(cookie GrabPointerCookie) (reply *GrabPointerReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_grab_pointer_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(GrabPointerReply)
		reply.EGrabPointerReply = (*EGrabPointerReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *GrabPointerReply) {
			C.free(unsafe.Pointer(f.EGrabPointerReply))
		})
	}
	if e != nil {
		err = errors.New("GrabPointerReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) UngrabPointer(time Timestamp) VoidCookie {
	return VoidCookie(C.xcb_ungrab_pointer(
		c.c(),
		C.xcb_timestamp_t(time),
	))
}

type EGrabKeyboardReply struct {
	ResponseType C.uint8_t
	Status       C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
}
type GrabKeyboardReply struct {
	*EGrabKeyboardReply
}

func (c *GrabKeyboardReply) c() *C.xcb_grab_keyboard_reply_t {
	ptr := c.EGrabKeyboardReply
	return (*C.xcb_grab_keyboard_reply_t)(unsafe.Pointer(ptr))
}

type GrabKeyboardCookie C.xcb_grab_keyboard_cookie_t

func (c GrabKeyboardCookie) c() C.xcb_grab_keyboard_cookie_t {
	return C.xcb_grab_keyboard_cookie_t(c)
}

func (c *Connection) GrabKeyboard(ownerEvents uint8, grabWindow Window, time Timestamp, pointerMode, keyboardMode uint8) GrabKeyboardCookie {
	cookie := C.xcb_grab_keyboard(
		c.c(),
		C.uint8_t(ownerEvents),
		C.xcb_window_t(grabWindow),
		C.xcb_timestamp_t(time),
		C.uint8_t(pointerMode),
		C.uint8_t(keyboardMode),
	)
	return GrabKeyboardCookie(cookie)
}

func (c *Connection) GrabKeyboardReply(cookie GrabKeyboardCookie) (reply *GrabKeyboardReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_grab_keyboard_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(GrabKeyboardReply)
		reply.EGrabKeyboardReply = (*EGrabKeyboardReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *GrabKeyboardReply) {
			C.free(unsafe.Pointer(f.EGrabKeyboardReply))
		})
	}
	if e != nil {
		err = errors.New("GrabKeyboardReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

func (c *Connection) UngrabKeyboard(time Timestamp) VoidCookie {
	return VoidCookie(C.xcb_ungrab_keyboard(
		c.c(),
		C.xcb_timestamp_t(time),
	))
}

type EAllocColorReply struct {
	ResponseType C.uint8_t
	Pad0         C.uint8_t
	Sequence     C.uint16_t
	Length       C.uint32_t
	Red          C.uint16_t
	Green        C.uint16_t
	Blue         C.uint16_t
	Pad1         [2]C.uint8_t
	Pixel        C.uint32_t
}
type AllocColorReply struct {
	*EAllocColorReply
}

func (c *AllocColorReply) c() *C.xcb_alloc_color_reply_t {
	ptr := c.EAllocColorReply
	return (*C.xcb_alloc_color_reply_t)(unsafe.Pointer(ptr))
}

type AllocColorCookie C.xcb_alloc_color_cookie_t

func (c AllocColorCookie) c() C.xcb_alloc_color_cookie_t {
	return C.xcb_alloc_color_cookie_t(c)
}

func (c *Connection) AllocColor(cmap Colormap, red, green, blue uint16) AllocColorCookie {
	cookie := C.xcb_alloc_color(
		c.c(),
		C.xcb_colormap_t(cmap),
		C.uint16_t(red),
		C.uint16_t(green),
		C.uint16_t(blue),
	)
	return AllocColorCookie(cookie)
}

func (c *Connection) AllocColorReply(cookie AllocColorCookie) (reply *AllocColorReply, err error) {
	var e *C.xcb_generic_error_t
	cReply := C.xcb_alloc_color_reply(c.c(), cookie.c(), &e)
	if e == nil {
		reply = new(AllocColorReply)
		reply.EAllocColorReply = (*EAllocColorReply)(unsafe.Pointer(cReply))
		runtime.SetFinalizer(reply, func(f *AllocColorReply) {
			C.free(unsafe.Pointer(f.EAllocColorReply))
		})
	}
	if e != nil {
		err = errors.New("AllocColorReply(): " + xcbError(e))
		C.free(unsafe.Pointer(e))
	}
	return
}

const (
	COLORMAP_ALLOC_NONE = C.XCB_COLORMAP_ALLOC_NONE
)

func (c *Connection) CreateColormap(alloc uint8, mid Colormap, window Window, visual VisualId) VoidCookie {
	return VoidCookie(C.xcb_create_colormap(
		c.c(),
		C.uint8_t(alloc),
		C.xcb_colormap_t(mid),
		C.xcb_window_t(window),
		C.xcb_visualid_t(visual),
	))
}
