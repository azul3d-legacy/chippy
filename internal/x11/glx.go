// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build linux

package x11

/*
#include <stdlib.h>
#include <GL/glx.h>
#include <GL/gl.h>
#include <X11/Xlib-xcb.h>

#cgo LDFLAGS: -lX11 -lGL

GLXContext chippy_glXCreateNewContext(void* p, void* dpy, void* config, int render_type, GLXContext share_list, Bool direct);
Bool chippy_glXMakeContextCurrent(void* p, void* dpy, GLXDrawable draw, GLXDrawable read, GLXContext ctx);
GLXWindow chippy_glXCreateWindow(void* p, void* dpy, void* config, Window win, const int *attrib_list);
void chippy_glXDestroyWindow(void* p, void* dpy, GLXWindow win);
void chippy_glXDestroyContext(void* p, void* dpy, GLXContext ctx);
Bool chippy_glXQueryVersion(void* p, void* dpy, int *maj, int *min);
void chippy_glXSwapBuffers(void* p, void* dpy, GLXDrawable drawable);
void* chippy_glXGetFBConfigs(void* p, void* dpy, int screen, int *nelements);
const char* chippy_glXQueryExtensionsString(void* p, void* dpy, int screen);
int chippy_glXGetFBConfigAttrib(void* p, void* dpy, void* config, int attribute, int *value);
GLXContext chippy_glXGetCurrentContext(void* p);
XVisualInfo* chippy_glXGetVisualFromFBConfig(void* p, void* dpy, void* config);
GLubyte* chippy_glGetString(void* p, GLenum v);

// Extensions below here.

GLXContext chippy_glXCreateContextAttribsARB(void* p, void* dpy, void* config, GLXContext share, Bool direct, const int* attribs);
void chippy_glXSwapIntervalEXT(void* p, void* dpy, GLXDrawable d, int interval);
int chippy_glXSwapIntervalMESA(void* p, int interval);
int chippy_glXSwapIntervalSGI(void* p, int interval);
*/
import "C"

import (
	"reflect"
	"unsafe"
)

const (
	GLX_DOUBLEBUFFER     = C.GLX_DOUBLEBUFFER
	GLX_STEREO           = C.GLX_STEREO
	GLX_AUX_BUFFERS      = C.GLX_AUX_BUFFERS
	GLX_RED_SIZE         = C.GLX_RED_SIZE
	GLX_GREEN_SIZE       = C.GLX_GREEN_SIZE
	GLX_BLUE_SIZE        = C.GLX_BLUE_SIZE
	GLX_ALPHA_SIZE       = C.GLX_ALPHA_SIZE
	GLX_DEPTH_SIZE       = C.GLX_DEPTH_SIZE
	GLX_STENCIL_SIZE     = C.GLX_STENCIL_SIZE
	GLX_ACCUM_RED_SIZE   = C.GLX_ACCUM_RED_SIZE
	GLX_ACCUM_GREEN_SIZE = C.GLX_ACCUM_GREEN_SIZE
	GLX_ACCUM_BLUE_SIZE  = C.GLX_ACCUM_BLUE_SIZE
	GLX_ACCUM_ALPHA_SIZE = C.GLX_ACCUM_ALPHA_SIZE

	GLX_SAMPLE_BUFFERS = C.GLX_SAMPLE_BUFFERS
	GLX_SAMPLES        = C.GLX_SAMPLES

	GLX_TRANSPARENT_TYPE  = C.GLX_TRANSPARENT_TYPE
	GLX_NONE              = C.GLX_NONE
	GLX_TRANSPARENT_RGB   = C.GLX_TRANSPARENT_RGB
	GLX_TRANSPARENT_INDEX = C.GLX_TRANSPARENT_INDEX

	GLX_TRANSPARENT_INDEX_VALUE = C.GLX_TRANSPARENT_INDEX_VALUE
	GLX_TRANSPARENT_RED_VALUE   = C.GLX_TRANSPARENT_RED_VALUE
	GLX_TRANSPARENT_GREEN_VALUE = C.GLX_TRANSPARENT_GREEN_VALUE
	GLX_TRANSPARENT_BLUE_VALUE  = C.GLX_TRANSPARENT_BLUE_VALUE
	GLX_TRANSPARENT_ALPHA_VALUE = C.GLX_TRANSPARENT_ALPHA_VALUE

	GLX_RENDER_TYPE      = C.GLX_RENDER_TYPE
	GLX_RGBA_TYPE        = C.GLX_RGBA_TYPE
	GLX_COLOR_INDEX_TYPE = C.GLX_COLOR_INDEX_TYPE

	GLX_X_VISUAL_TYPE = C.GLX_X_VISUAL_TYPE
	GLX_TRUE_COLOR    = C.GLX_TRUE_COLOR

	GLX_CONFIG_CAVEAT         = C.GLX_CONFIG_CAVEAT
	GLX_SLOW_CONFIG           = C.GLX_SLOW_CONFIG
	GLX_NON_CONFORMANT_CONFIG = C.GLX_NON_CONFORMANT_CONFIG

	GLX_CONTEXT_MAJOR_VERSION_ARB             = C.GLX_CONTEXT_MAJOR_VERSION_ARB
	GLX_CONTEXT_MINOR_VERSION_ARB             = C.GLX_CONTEXT_MINOR_VERSION_ARB
	GLX_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB    = C.GLX_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB
	GLX_CONTEXT_DEBUG_BIT_ARB                 = C.GLX_CONTEXT_DEBUG_BIT_ARB
	GLX_CONTEXT_FLAGS_ARB                     = C.GLX_CONTEXT_FLAGS_ARB
	GLX_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB = C.GLX_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB
	GLX_CONTEXT_CORE_PROFILE_BIT_ARB          = C.GLX_CONTEXT_CORE_PROFILE_BIT_ARB
	GLX_CONTEXT_PROFILE_MASK_ARB              = C.GLX_CONTEXT_PROFILE_MASK_ARB
)

type (
	GLXContext  C.GLXContext
	GLXDrawable C.GLXDrawable
	GLXFBConfig uintptr
	GLXWindow   C.GLXWindow
)

var glXCreateNewContextPtr unsafe.Pointer

func (d *Display) GLXCreateNewContext(config GLXFBConfig, renderType int, shareList GLXContext, direct bool) GLXContext {
	if glXCreateNewContextPtr == nil {
		glXCreateNewContextPtr = GLXGetProcAddressARB("glXCreateNewContext")
	}

	d.Lock()
	defer d.Unlock()
	cDirect := C.Bool(0)
	if direct {
		cDirect = 1
	}
	return GLXContext(C.chippy_glXCreateNewContext(
		glXCreateNewContextPtr,
		d.ptr(),
		unsafe.Pointer(config),
		C.int(renderType),
		C.GLXContext(shareList),
		cDirect,
	))
}

var glXDestroyContextPtr unsafe.Pointer

func (d *Display) GLXDestroyContext(ctx GLXContext) {
	if glXDestroyContextPtr == nil {
		glXDestroyContextPtr = GLXGetProcAddressARB("glXDestroyContext")
	}

	d.Lock()
	defer d.Unlock()
	C.chippy_glXDestroyContext(
		glXDestroyContextPtr,
		d.ptr(),
		C.GLXContext(ctx),
	)
}

var glXMakeContextCurrentPtr unsafe.Pointer

func (d *Display) GLXMakeContextCurrent(draw, read GLXDrawable, ctx GLXContext) int {
	if glXMakeContextCurrentPtr == nil {
		glXMakeContextCurrentPtr = GLXGetProcAddressARB("glXMakeContextCurrent")
	}

	d.Lock()
	defer d.Unlock()
	return int(C.chippy_glXMakeContextCurrent(
		glXMakeContextCurrentPtr,
		d.ptr(),
		C.GLXDrawable(draw),
		C.GLXDrawable(read),
		C.GLXContext(ctx),
	))
}

var glXSwapBuffersPtr unsafe.Pointer

func (d *Display) GLXSwapBuffers(drawable GLXDrawable) {
	if glXSwapBuffersPtr == nil {
		glXSwapBuffersPtr = GLXGetProcAddressARB("glXSwapBuffers")
	}

	d.Lock()
	defer d.Unlock()
	C.chippy_glXSwapBuffers(
		glXSwapBuffersPtr,
		d.ptr(),
		C.GLXDrawable(drawable),
	)
}

type Int C.int

var glXQueryVersionPtr unsafe.Pointer

func (d *Display) GLXQueryVersion(maj, min *Int) bool {
	if glXQueryVersionPtr == nil {
		glXQueryVersionPtr = GLXGetProcAddressARB("glXQueryVersion")
	}

	d.Lock()
	defer d.Unlock()
	return C.chippy_glXQueryVersion(
		glXQueryVersionPtr,
		d.ptr(),
		(*C.int)(unsafe.Pointer(maj)),
		(*C.int)(unsafe.Pointer(min)),
	) != 0
}

var glXQueryExtensionsStringPtr unsafe.Pointer

func (d *Display) GLXQueryExtensionsString(screen int) string {
	if glXQueryExtensionsStringPtr == nil {
		glXQueryExtensionsStringPtr = GLXGetProcAddressARB("glXQueryExtensionsString")
	}

	d.Lock()
	defer d.Unlock()
	data := C.chippy_glXQueryExtensionsString(
		glXQueryExtensionsStringPtr,
		d.ptr(),
		C.int(screen),
	)
	return C.GoString(data)
}

var glXGetFBConfigAttribPtr unsafe.Pointer

func (d *Display) GLXGetFBConfigAttrib(config GLXFBConfig, attrib int) (value Int, ret int) {
	if glXGetFBConfigAttribPtr == nil {
		glXGetFBConfigAttribPtr = GLXGetProcAddressARB("glXGetFBConfigAttrib")
	}

	d.Lock()
	defer d.Unlock()
	ret = int(C.chippy_glXGetFBConfigAttrib(
		glXGetFBConfigAttribPtr,
		d.ptr(),
		unsafe.Pointer(config),
		C.int(attrib),
		(*C.int)(unsafe.Pointer(&value)),
	))
	return
}

var glXGetFBConfigsPtr unsafe.Pointer

func (d *Display) GLXGetFBConfigs(screen int) (configs []GLXFBConfig) {
	if glXGetFBConfigsPtr == nil {
		glXGetFBConfigsPtr = GLXGetProcAddressARB("glXGetFBConfigs")
	}

	d.Lock()
	defer d.Unlock()
	var nConfigs C.int
	cConfigs := C.chippy_glXGetFBConfigs(
		glXGetFBConfigsPtr,
		d.ptr(),
		C.int(screen),
		&nConfigs,
	)
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&configs))
	sliceHeader.Len = int(nConfigs)
	sliceHeader.Cap = int(nConfigs)
	sliceHeader.Data = uintptr(unsafe.Pointer(cConfigs))
	return
}

var glXCreateWindowPtr unsafe.Pointer

func (d *Display) GLXCreateWindow(config GLXFBConfig, win Window) GLXWindow {
	if glXCreateWindowPtr == nil {
		glXCreateWindowPtr = GLXGetProcAddressARB("glXCreateWindow")
	}

	d.Lock()
	defer d.Unlock()
	return GLXWindow(C.chippy_glXCreateWindow(
		glXCreateWindowPtr,
		d.ptr(),
		unsafe.Pointer(config),
		C.Window(win),
		nil,
	))
}

var glXDestroyWindowPtr unsafe.Pointer

func (d *Display) GLXDestroyWindow(win GLXWindow) {
	if glXDestroyWindowPtr == nil {
		glXDestroyWindowPtr = GLXGetProcAddressARB("glXDestroyWindow")
	}

	d.Lock()
	defer d.Unlock()
	C.chippy_glXDestroyWindow(
		glXDestroyWindowPtr,
		d.ptr(),
		C.GLXWindow(win),
	)
}

var glXGetVisualFromFBConfigPtr unsafe.Pointer

func (d *Display) GLXGetVisualFromFBConfig(config GLXFBConfig) *XVisualInfo {
	if glXGetVisualFromFBConfigPtr == nil {
		glXGetVisualFromFBConfigPtr = GLXGetProcAddressARB("glXGetVisualFromFBConfig")
	}

	d.Lock()
	defer d.Unlock()
	return (*XVisualInfo)(unsafe.Pointer(C.chippy_glXGetVisualFromFBConfig(
		glXGetVisualFromFBConfigPtr,
		d.ptr(),
		unsafe.Pointer(config),
	)))
}

var glXGetCurrentContextPtr unsafe.Pointer

func GLXGetCurrentContext() GLXContext {
	if glXGetCurrentContextPtr == nil {
		glXGetCurrentContextPtr = GLXGetProcAddressARB("glXGetCurrentContext")
	}

	return GLXContext(C.chippy_glXGetCurrentContext(
		glXGetCurrentContextPtr,
	))
}

func GLXGetProcAddressARB(p string) unsafe.Pointer {
	cstr := C.CString(p)
	defer C.free(unsafe.Pointer(cstr))
	return unsafe.Pointer(C.glXGetProcAddressARB(
		(*C.GLubyte)(unsafe.Pointer(cstr)),
	))
}

var glXCreateContextAttribsARBPtr unsafe.Pointer

func (d *Display) GLXCreateContextAttribsARB(config GLXFBConfig, share GLXContext, direct bool, attribs *Int) GLXContext {
	if glXCreateContextAttribsARBPtr == nil {
		glXCreateContextAttribsARBPtr = GLXGetProcAddressARB("glXCreateContextAttribsARB")
	}

	d.Lock()
	defer d.Unlock()
	cDirect := C.Bool(0)
	if direct {
		cDirect = 1
	}
	return GLXContext(C.chippy_glXCreateContextAttribsARB(
		glXCreateContextAttribsARBPtr,
		d.ptr(),
		unsafe.Pointer(config),
		C.GLXContext(share),
		cDirect,
		(*C.int)(unsafe.Pointer(attribs)),
	))
}

var glXSwapIntervalEXTPtr unsafe.Pointer

func (d *Display) GLXSwapIntervalEXT(drawable GLXDrawable, interval int) {
	if glXSwapIntervalEXTPtr == nil {
		glXSwapIntervalEXTPtr = GLXGetProcAddressARB("glXSwapIntervalEXT")
	}

	d.Lock()
	defer d.Unlock()
	C.chippy_glXSwapIntervalEXT(
		glXSwapIntervalEXTPtr,
		d.ptr(),
		C.GLXDrawable(drawable),
		C.int(interval),
	)
}

var glXSwapIntervalMESAPtr unsafe.Pointer

func (d *Display) GLXSwapIntervalMESA(interval int) int {
	if glXSwapIntervalMESAPtr == nil {
		glXSwapIntervalMESAPtr = GLXGetProcAddressARB("glXSwapIntervalMESA")
	}

	d.Lock()
	defer d.Unlock()
	// Should be OK to use glXSwapIntervalEXT but with the MESA pointer because
	// they have the same typedef
	return int(C.chippy_glXSwapIntervalMESA(
		glXSwapIntervalMESAPtr,
		C.int(interval),
	))
}

var glXSwapIntervalSGIPtr unsafe.Pointer

func (d *Display) GLXSwapIntervalSGI(interval int) int {
	if glXSwapIntervalSGIPtr == nil {
		glXSwapIntervalSGIPtr = GLXGetProcAddressARB("glXSwapIntervalSGI")
	}

	d.Lock()
	defer d.Unlock()
	// Should be OK to use glXSwapIntervalEXT but with the SGI pointer because
	// they have the same typedef
	return int(C.chippy_glXSwapIntervalSGI(
		glXSwapIntervalSGIPtr,
		C.int(interval),
	))
}

func XFree(ptr unsafe.Pointer) {
	C.XFree(ptr)
}

const (
	GL_VERSION = 0x1F02
)

var glGetStringPtr unsafe.Pointer

func GlGetString(name uint32) string {
	if glGetStringPtr == nil {
		glGetStringPtr = GLXGetProcAddressARB("glGetString")
	}

	ret := C.chippy_glGetString(
		glGetStringPtr,
		C.GLenum(name),
	)
	return C.GoString((*C.char)(unsafe.Pointer(ret)))
}
