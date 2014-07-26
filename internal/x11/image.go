// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/xcb_image.h>

#cgo LDFLAGS: -lxcb -lxcb-image
*/
import "C"

import (
	"unsafe"
)

type Image C.xcb_image_t

func ImageCreate(width, height uint16, format int, xpad, depth, bpp, unit uint8, byteOrder, bitOrder int, base unsafe.Pointer, bytes uint32, data *uint8) *Image {
	return (*Image)(unsafe.Pointer(C.xcb_image_create(
		C.uint16_t(width),
		C.uint16_t(height),
		C.xcb_image_format_t(format),
		C.uint8_t(xpad),
		C.uint8_t(depth),
		C.uint8_t(bpp),
		C.uint8_t(unit),
		C.xcb_image_order_t(byteOrder),
		C.xcb_image_order_t(bitOrder),
		base,
		C.uint32_t(bytes),
		(*C.uint8_t)(unsafe.Pointer(data)),
	)))
}

func (c *Connection) ImageNative(i *Image, convert bool) *Image {
	cConvert := C.int(0)
	if convert {
		cConvert = 1
	}
	ret := C.xcb_image_native(
		c.c(),
		(*C.xcb_image_t)(unsafe.Pointer(i)),
		cConvert,
	)
	return (*Image)(unsafe.Pointer(ret))
}

func (c *Connection) ImagePut(draw Drawable, gc GContext, image *Image, x, y int16, leftPad uint8) VoidCookie {
	return VoidCookie(C.xcb_image_put(
		c.c(),
		C.xcb_drawable_t(draw),
		C.xcb_gcontext_t(gc),
		(*C.xcb_image_t)(unsafe.Pointer(image)),
		C.int16_t(x),
		C.int16_t(y),
		C.uint8_t(leftPad),
	))
}
