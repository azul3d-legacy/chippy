// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgo_export.h"

#define UNICODE
#include <windows.h>

typedef HGLRC (*chippy_p_wglCreateContextAttribsARB) (HDC, HGLRC, const int*);
HGLRC chippy_wglCreateContextAttribsARB(void* p, HDC hDC, HGLRC hshareContext, const int* attribList) {
	chippy_p_wglCreateContextAttribsARB fn = (chippy_p_wglCreateContextAttribsARB)p;
	return fn(hDC, hshareContext, attribList);
}

typedef char* (*chippy_p_wglGetExtensionsStringARB) (HDC);
char* chippy_wglGetExtensionsStringARB(void* p, HDC hdc) {
	chippy_p_wglGetExtensionsStringARB fn = (chippy_p_wglGetExtensionsStringARB)p;
	return fn(hdc);
}

typedef BOOL (*chippy_p_wglSwapIntervalEXT) (int);
BOOL chippy_wglSwapIntervalEXT(void* p, int interval) {
	chippy_p_wglSwapIntervalEXT fn = (chippy_p_wglSwapIntervalEXT)p;
	return fn(interval);
}

typedef BOOL (*chippy_p_wglGetPixelFormatAttribivARB) (HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, int* piValues);
BOOL chippy_wglGetPixelFormatAttribivARB(void* p, HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, int* piValues) {
	chippy_p_wglGetPixelFormatAttribivARB fn = (chippy_p_wglGetPixelFormatAttribivARB)p;
	return fn(hdc, iPixelFormat, iLayerPlane, nAttributes, piAttributes, piValues);
}


typedef BOOL (*chippy_p_wglGetPixelFormatAttribfvARB) (HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, FLOAT* pfValues);
BOOL chippy_wglGetPixelFormatAttribfvARB(void* p, HDC hdc, int iPixelFormat, int iLayerPlane, UINT nAttributes, const int* piAttributes, FLOAT* pfValues) {
	chippy_p_wglGetPixelFormatAttribfvARB fn = (chippy_p_wglGetPixelFormatAttribfvARB)p;
	return fn(hdc, iPixelFormat, iLayerPlane, nAttributes, piAttributes, pfValues);
}
