// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <X11/Xlib-xcb.h>

#include "_cgo_export.h"

int chippy_xlib_error(Display* d, XErrorEvent* e) {
	chippy_xlib_error_callback(d, e);
	return 0;
}

