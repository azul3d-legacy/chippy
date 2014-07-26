// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "_cgo_export.h"

WORD win32_MAKELANGID(USHORT usPrimaryLanguage, USHORT usSubLanguage) {
    return MAKELANGID(usPrimaryLanguage, usSubLanguage);
}

LPTSTR macro_MAKEINTRESOURCE(WORD wInteger) {
	return MAKEINTRESOURCE(wInteger);
}

//MONITORENUMPROC win32_MonitorEnumProcCallbackHandle = (MONITORENUMPROC)MonitorEnumProcCallback;
//HOOKPROC win32_LowLevelKeyboardHookCallbackHandle = (HOOKPROC)LowLevelKeyboardHookCallback;

