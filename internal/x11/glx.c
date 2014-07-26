// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <X11/Xlib-xcb.h>
#include <GL/glx.h>

#include "_cgo_export.h"

typedef GLXContext (*chippy_p_glXCreateNewContext) (Display *dpy, GLXFBConfig config, int render_type, GLXContext share_list, Bool direct);
GLXContext chippy_glXCreateNewContext(void* p, void* dpy, void* config, int render_type, GLXContext share_list, Bool direct) {
	chippy_p_glXCreateNewContext fn = (chippy_p_glXCreateNewContext)p;
	return fn((Display*)dpy, config, render_type, share_list, direct);
}

typedef Bool (*chippy_p_glXMakeContextCurrent) (Display *dpy, GLXDrawable draw, GLXDrawable read, GLXContext ctx);
Bool chippy_glXMakeContextCurrent(void* p, void* dpy, GLXDrawable draw, GLXDrawable read, GLXContext ctx) {
	chippy_p_glXMakeContextCurrent fn = (chippy_p_glXMakeContextCurrent)p;
	return fn((Display*)dpy, draw, read, ctx);
}

typedef GLXWindow (*chippy_p_glXCreateWindow) (Display *dpy, GLXFBConfig config, Window win, const int *attrib_list);
GLXWindow chippy_glXCreateWindow(void* p, void* dpy, void* config, Window win, const int *attrib_list) {
	chippy_p_glXCreateWindow fn = (chippy_p_glXCreateWindow)p;
	return fn((Display*)dpy, config, win, attrib_list);
}

typedef void (*chippy_p_glXDestroyWindow) (Display *dpy, GLXWindow win);
void chippy_glXDestroyWindow(void* p, void* dpy, GLXWindow win) {
	chippy_p_glXDestroyWindow fn = (chippy_p_glXDestroyWindow)p;
	fn((Display*)dpy, win);
}

typedef void (*chippy_p_glXDestroyContext) (Display *dpy, GLXContext ctx);
void chippy_glXDestroyContext(void* p, void* dpy, GLXContext ctx) {
	chippy_p_glXDestroyContext fn = (chippy_p_glXDestroyContext)p;
	fn((Display*)dpy, ctx);
}

typedef Bool (*chippy_p_glXQueryVersion) (Display *dpy, int *maj, int *min);
Bool chippy_glXQueryVersion(void* p, void* dpy, int *maj, int *min) {
	chippy_p_glXQueryVersion fn = (chippy_p_glXQueryVersion)p;
	return fn((Display*)dpy, maj, min);
}

typedef void (*chippy_p_glXSwapBuffers) (Display *dpy, GLXDrawable drawable);
void chippy_glXSwapBuffers(void* p, void* dpy, GLXDrawable drawable) {
	chippy_p_glXSwapBuffers fn = (chippy_p_glXSwapBuffers)p;
	fn((Display*)dpy, drawable);
}

typedef GLXFBConfig* (*chippy_p_glXGetFBConfigs) (Display *dpy, int screen, int *nelements);
void* chippy_glXGetFBConfigs(void* p, void* dpy, int screen, int *nelements) {
	chippy_p_glXGetFBConfigs fn = (chippy_p_glXGetFBConfigs)p;
	return fn((Display*)dpy, screen, nelements);
}

typedef const char* (*chippy_p_glXQueryExtensionsString) (Display *dpy, int screen);
const char* chippy_glXQueryExtensionsString(void* p, void* dpy, int screen) {
	chippy_p_glXQueryExtensionsString fn = (chippy_p_glXQueryExtensionsString)p;
	return fn((Display*)dpy, screen);
}

typedef int (*chippy_p_glXGetFBConfigAttrib) (Display *dpy, GLXFBConfig config, int attribute, int *value);
int chippy_glXGetFBConfigAttrib(void* p, void* dpy, void* config, int attribute, int *value) {
	chippy_p_glXGetFBConfigAttrib fn = (chippy_p_glXGetFBConfigAttrib)p;
	return fn((Display*)dpy, config, attribute, value);
}

typedef GLXContext (*chippy_p_glXGetCurrentContext) (void);
GLXContext chippy_glXGetCurrentContext(void* p) {
	chippy_p_glXGetCurrentContext fn = (chippy_p_glXGetCurrentContext)p;
	return fn();
}

typedef XVisualInfo* (*chippy_p_glXGetVisualFromFBConfig) (Display *dpy, GLXFBConfig config);
XVisualInfo* chippy_glXGetVisualFromFBConfig(void* p, void* dpy, void* config) {
	chippy_p_glXGetVisualFromFBConfig fn = (chippy_p_glXGetVisualFromFBConfig)p;
	return fn((Display*)dpy, config);
}

typedef GLubyte* (*chippy_p_glGetString) (GLenum v);
GLubyte* chippy_glGetString(void* p, GLenum v) {
	chippy_p_glGetString fn = (chippy_p_glGetString)p;
	return fn(v);
}


// Extensions below here.

typedef GLXContext (*chippy_p_glXCreateContextAttribsARB) (Display* dpy, GLXFBConfig config, GLXContext share, Bool direct, const int* attribs);
GLXContext chippy_glXCreateContextAttribsARB(void* p, Display* dpy, void* config, GLXContext share, Bool direct, const int* attribs) {
	chippy_p_glXCreateContextAttribsARB fn = (chippy_p_glXCreateContextAttribsARB)p;
	return fn((Display*)dpy, config, share, direct, attribs);
}

typedef void (*chippy_p_glXSwapIntervalEXT) (Display* dpy, GLXDrawable d, int interval);
void chippy_glXSwapIntervalEXT(void* p, Display* dpy, GLXDrawable d, int interval) {
	chippy_p_glXSwapIntervalEXT fn = (chippy_p_glXSwapIntervalEXT)p;
	fn((Display*)dpy, d, interval);
}

typedef int (*chippy_p_glXSwapIntervalMESA) (int interval);
int chippy_glXSwapIntervalMESA(void* p, int interval) {
	chippy_p_glXSwapIntervalMESA fn = (chippy_p_glXSwapIntervalMESA)p;
	return fn(interval);
}

typedef int (*chippy_p_glXSwapIntervalSGI) (int interval);
int chippy_glXSwapIntervalSGI(void* p, int interval) {
	chippy_p_glXSwapIntervalSGI fn = (chippy_p_glXSwapIntervalSGI)p;
	return fn(interval);
}

