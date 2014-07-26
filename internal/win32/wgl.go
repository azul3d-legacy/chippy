// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package win32

/*
#define UNICODE
#include <windows.h>
#include <GL/gl.h>

#cgo LDFLAGS: -lopengl32

extern HGLRC chippy_wglCreateContextAttribsARB(void* p, HDC hDC, HGLRC hshareContext, const int* attribList);
extern char* chippy_wglGetExtensionsStringARB(void* p, HDC hdc);
extern BOOL chippy_wglSwapIntervalEXT(void* p, int interval);

extern BOOL chippy_wglGetPixelFormatAttribivARB(void* p, HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, int* piValues);
extern BOOL chippy_wglGetPixelFormatAttribfvARB(void* p, HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, FLOAT* pfValues);
*/
import "C"

import (
	"unsafe"
)

// OpenGL functions
type HGLRC C.HGLRC

func WglCreateContext(hdc HDC) HGLRC {
	return HGLRC(C.wglCreateContext(C.HDC(hdc)))
}

func WglDeleteContext(hglrc HGLRC) bool {
	return C.wglDeleteContext(C.HGLRC(hglrc)) != 0
}

func WglGetProcAddress(lpszProc string) uintptr {
	cstr := (*C.CHAR)(unsafe.Pointer(C.CString(lpszProc)))
	defer C.free(unsafe.Pointer(cstr))

	r := uintptr(unsafe.Pointer(C.wglGetProcAddress(cstr)))
	negativeOne := -1
	if r == 0 || r == 1 || r == 2 || r == 3 || r == uintptr(negativeOne) {
		return 0
	}
	return uintptr(r)
}

func WglMakeCurrent(hdc HDC, hglrc HGLRC) bool {
	return C.wglMakeCurrent(C.HDC(hdc), C.HGLRC(hglrc)) != 0
}

func WglShareLists(hglrc1, hglrc2 HGLRC) bool {
	return C.wglShareLists(C.HGLRC(hglrc1), C.HGLRC(hglrc2)) != 0
}

const (
	WGL_CONTEXT_MAJOR_VERSION_ARB = 0x2091
	WGL_CONTEXT_MINOR_VERSION_ARB = 0x2092
	WGL_CONTEXT_LAYER_PLANE_ARB   = 0x2093
	WGL_CONTEXT_FLAGS_ARB         = 0x2094
	WGL_CONTEXT_PROFILE_MASK_ARB  = 0x9126

	//Accepted as bits in the attribute value for WGL_CONTEXT_FLAGS in <*attribList>:

	WGL_CONTEXT_DEBUG_BIT_ARB              = 0x0001
	WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB = 0x0002

	// Accepted as bits in the attribute value for WGL_CONTEXT_PROFILE_MASK_ARB in <*attribList>:

	WGL_CONTEXT_CORE_PROFILE_BIT_ARB          = 0x00000001
	WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB = 0x00000002

	// New errors returned by GetLastError:

	ERROR_INVALID_VERSION_ARB = 0x2095
	ERROR_INVALID_PROFILE_ARB = 0x2096

	// Accepted in the <piAttributes> parameter array of
	// wglGetPixelFormatAttribivARB, and wglGetPixelFormatAttribfvARB, and
	// as a type in the <piAttribIList> and <pfAttribFList> parameter
	// arrays of wglChoosePixelFormatARB:
	WGL_NUMBER_PIXEL_FORMATS_ARB    = 0x2000
	WGL_DRAW_TO_WINDOW_ARB          = 0x2001
	WGL_DRAW_TO_BITMAP_ARB          = 0x2002
	WGL_ACCELERATION_ARB            = 0x2003
	WGL_NEED_PALETTE_ARB            = 0x2004
	WGL_NEED_SYSTEM_PALETTE_ARB     = 0x2005
	WGL_SWAP_LAYER_BUFFERS_ARB      = 0x2006
	WGL_SWAP_METHOD_ARB             = 0x2007
	WGL_NUMBER_OVERLAYS_ARB         = 0x2008
	WGL_NUMBER_UNDERLAYS_ARB        = 0x2009
	WGL_TRANSPARENT_ARB             = 0x200A
	WGL_TRANSPARENT_RED_VALUE_ARB   = 0x2037
	WGL_TRANSPARENT_GREEN_VALUE_ARB = 0x2038
	WGL_TRANSPARENT_BLUE_VALUE_ARB  = 0x2039
	WGL_TRANSPARENT_ALPHA_VALUE_ARB = 0x203A
	WGL_TRANSPARENT_INDEX_VALUE_ARB = 0x203B
	WGL_SHARE_DEPTH_ARB             = 0x200C
	WGL_SHARE_STENCIL_ARB           = 0x200D
	WGL_SHARE_ACCUM_ARB             = 0x200E
	WGL_SUPPORT_GDI_ARB             = 0x200F
	WGL_SUPPORT_OPENGL_ARB          = 0x2010
	WGL_DOUBLE_BUFFER_ARB           = 0x2011
	WGL_STEREO_ARB                  = 0x2012
	WGL_PIXEL_TYPE_ARB              = 0x2013
	WGL_COLOR_BITS_ARB              = 0x2014
	WGL_RED_BITS_ARB                = 0x2015
	WGL_RED_SHIFT_ARB               = 0x2016
	WGL_GREEN_BITS_ARB              = 0x2017
	WGL_GREEN_SHIFT_ARB             = 0x2018
	WGL_BLUE_BITS_ARB               = 0x2019
	WGL_BLUE_SHIFT_ARB              = 0x201A
	WGL_ALPHA_BITS_ARB              = 0x201B
	WGL_ALPHA_SHIFT_ARB             = 0x201C
	WGL_ACCUM_BITS_ARB              = 0x201D
	WGL_ACCUM_RED_BITS_ARB          = 0x201E
	WGL_ACCUM_GREEN_BITS_ARB        = 0x201F
	WGL_ACCUM_BLUE_BITS_ARB         = 0x2020
	WGL_ACCUM_ALPHA_BITS_ARB        = 0x2021
	WGL_DEPTH_BITS_ARB              = 0x2022
	WGL_STENCIL_BITS_ARB            = 0x2023
	WGL_AUX_BUFFERS_ARB             = 0x2024

	// Accepted as a value in the <piAttribIList> and <pfAttribFList>
	// parameter arrays of wglChoosePixelFormatARB, and returned in the
	// <piValues> parameter array of wglGetPixelFormatAttribivARB, and the
	// <pfValues> parameter array of wglGetPixelFormatAttribfvARB:
	WGL_NO_ACCELERATION_ARB      = 0x2025
	WGL_GENERIC_ACCELERATION_ARB = 0x2026
	WGL_FULL_ACCELERATION_ARB    = 0x2027

	WGL_SWAP_EXCHANGE_ARB  = 0x2028
	WGL_SWAP_COPY_ARB      = 0x2029
	WGL_SWAP_UNDEFINED_ARB = 0x202A

	WGL_TYPE_RGBA_ARB       = 0x202B
	WGL_TYPE_COLORINDEX_ARB = 0x202C

	WGL_SAMPLE_BUFFERS_ARB = 0x2041
	WGL_SAMPLES_ARB        = 0x2042
)

func WglCreateContextAttribsARB(hdc HDC, hglrc HGLRC, attribs []Int) (HGLRC, bool) {
	ptr := WglGetProcAddress("wglCreateContextAttribsARB")
	if ptr == 0 {
		return nil, false
	}
	ret := C.chippy_wglCreateContextAttribsARB(unsafe.Pointer(ptr), C.HDC(hdc), C.HGLRC(hglrc), (*C.int)(unsafe.Pointer(&attribs[0])))
	return HGLRC(ret), true
}

func WglGetExtensionsStringARB(hdc HDC) (string, bool) {
	ptr := WglGetProcAddress("wglGetExtensionsStringARB")
	if ptr == 0 {
		return "", false
	}
	ret := C.chippy_wglGetExtensionsStringARB(unsafe.Pointer(ptr), C.HDC(hdc))
	return C.GoString(ret), true
}

func WglSwapIntervalEXT(interval int) bool {
	ptr := WglGetProcAddress("wglSwapIntervalEXT")
	if ptr == 0 {
		return false
	}
	return C.chippy_wglSwapIntervalEXT(unsafe.Pointer(ptr), C.int(interval)) != 0
}

type Float C.FLOAT

var wglGetPixelFormatAttribivARBPtr unsafe.Pointer

func WglGetPixelFormatAttribivARB(hdc HDC, iPixelFormat, iLayerPlane Int, piAttributes []Int, piValues *Int) bool {
	if wglGetPixelFormatAttribivARBPtr == nil {
		wglGetPixelFormatAttribivARBPtr = unsafe.Pointer(WglGetProcAddress("wglGetPixelFormatAttribivARB"))
	}

	return C.chippy_wglGetPixelFormatAttribivARB(
		wglGetPixelFormatAttribivARBPtr,
		C.HDC(hdc),
		C.int(iPixelFormat),
		C.int(iLayerPlane),
		C.UINT(len(piAttributes)),
		(*C.int)(unsafe.Pointer(&piAttributes[0])),
		(*C.int)(unsafe.Pointer(piValues)),
	) == 1
}

var wglGetPixelFormatAttribfvARBPtr unsafe.Pointer

func WglGetPixelFormatAttribfvARB(hdc HDC, iPixelFormat, iLayerPlane Int, piAttributes []Int, pfValues *Float) bool {
	if wglGetPixelFormatAttribfvARBPtr == nil {
		wglGetPixelFormatAttribfvARBPtr = unsafe.Pointer(WglGetProcAddress("wglGetPixelFormatAttribfvARB"))
	}

	return C.chippy_wglGetPixelFormatAttribfvARB(
		wglGetPixelFormatAttribfvARBPtr,
		C.HDC(hdc),
		C.int(iPixelFormat),
		C.int(iLayerPlane),
		C.UINT(len(piAttributes)),
		(*C.int)(unsafe.Pointer(&piAttributes[0])),
		(*C.FLOAT)(unsafe.Pointer(pfValues)),
	) == 1
}

const (
	GL_VERSION = 0x1F02
)

func GlGetString(name uint32) string {
	return C.GoString((*C.char)(unsafe.Pointer(C.glGetString(C.GLenum(name)))))
}
