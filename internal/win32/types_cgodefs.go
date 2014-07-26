// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build ignore

package win32

/*
#define _UNICODE
#define UNICODE
#include <windows.h>

// Because some MinGW versions don't contain this.

typedef struct win32__VIDEOPARAMETERS {
  GUID  guid;
  ULONG dwOffset;
  ULONG dwCommand;
  ULONG dwFlags;
  ULONG dwMode;
  ULONG dwTVStandard;
  ULONG dwAvailableModes;
  ULONG dwAvailableTVStandard;
  ULONG dwFlickerFilter;
  ULONG dwOverScanX;
  ULONG dwOverScanY;
  ULONG dwMaxUnscaledX;
  ULONG dwMaxUnscaledY;
  ULONG dwPositionX;
  ULONG dwPositionY;
  ULONG dwBrightness;
  ULONG dwContrast;
  ULONG dwCPType;
  ULONG dwCPCommand;
  ULONG dwCPStandard;
  ULONG dwCPKey;
  ULONG bCP_APSTriggerBits;
  UCHAR bOEMCopyProtection[256];
} win32_VIDEOPARAMETERS;

typedef struct win32_tagMONITORINFOEX {
  DWORD cbSize;
  RECT  rcMonitor;
  RECT  rcWork;
  DWORD dwFlags;
  TCHAR szDevice[CCHDEVICENAME];
} win32_MONITORINFOEX;

typedef struct win32_tagRAWHID {
  DWORD dwSizeHid;
  DWORD dwCount;
  BYTE  bRawData[1];
} win32_RAWHID;

typedef struct win32_tagRAWMOUSE {
  USHORT usFlags;
  union {
    ULONG  ulButtons;
    struct {
      USHORT usButtonFlags;
      USHORT usButtonData;
    };
  };
  ULONG  ulRawButtons;
  LONG   lLastX;
  LONG   lLastY;
  ULONG  ulExtraInformation;
} win32_RAWMOUSE;

typedef struct win32_tagRAWINPUTHEADER {
  DWORD  dwType;
  DWORD  dwSize;
  HANDLE hDevice;
  WPARAM wParam;
} win32_RAWINPUTHEADER;

typedef struct win32_tagRAWKEYBOARD {
  USHORT MakeCode;
  USHORT Flags;
  USHORT Reserved;
  USHORT VKey;
  UINT   Message;
  ULONG  ExtraInformation;
} win32_RAWKEYBOARD;

typedef struct win32_tagRAWINPUT {
  win32_RAWINPUTHEADER header;
  union {
    win32_RAWMOUSE    mouse;
    win32_RAWKEYBOARD keyboard;
    win32_RAWHID      hid;
  } data;
} win32_RAWINPUT;

typedef struct win32_tagRAWINPUTDEVICE {
  USHORT usUsagePage;
  USHORT usUsage;
  DWORD  dwFlags;
  HWND   hwndTarget;
} win32_RAWINPUTDEVICE;

*/
import "C"

import "unsafe"

type (
	LONG_PTR  C.LONG_PTR
	UINT_PTR  C.UINT_PTR
	ULONG_PTR C.ULONG_PTR
	BYTE      C.BYTE
	Int       C.int
	LONG      C.LONG
	DWORD     C.DWORD
	UINT      C.UINT
	WORD      C.WORD
	USHORT    C.USHORT
	LRESULT   C.LRESULT
	TCHAR     C.TCHAR
	COLORREF  C.COLORREF

	HWND      unsafe.Pointer
	HDC       unsafe.Pointer
	HMENU     unsafe.Pointer
	HINSTANCE unsafe.Pointer
	HICON     unsafe.Pointer
	HCURSOR   unsafe.Pointer
	HBRUSH    unsafe.Pointer
	HRGN      unsafe.Pointer
	HMONITOR  unsafe.Pointer
	HHOOK     unsafe.Pointer
	HRAWINPUT unsafe.Pointer
)

const (
	ENUM_CURRENT_SETTINGS = C.ENUM_CURRENT_SETTINGS

	DISP_CHANGE_SUCCESSFUL  = C.DISP_CHANGE_SUCCESSFUL
	DISP_CHANGE_BADDUALVIEW = C.DISP_CHANGE_BADDUALVIEW
	DISP_CHANGE_BADFLAGS    = C.DISP_CHANGE_BADFLAGS
	DISP_CHANGE_BADMODE     = C.DISP_CHANGE_BADMODE
	DISP_CHANGE_BADPARAM    = C.DISP_CHANGE_BADPARAM
	DISP_CHANGE_FAILED      = C.DISP_CHANGE_FAILED
	DISP_CHANGE_NOTUPDATED  = C.DISP_CHANGE_NOTUPDATED
	DISP_CHANGE_RESTART     = C.DISP_CHANGE_RESTART

	CDS_TEST           = C.CDS_TEST
	CDS_UPDATEREGISTRY = C.CDS_UPDATEREGISTRY

	HORZSIZE      = C.HORZSIZE      // mm width
	VERTSIZE      = C.VERTSIZE      // mm height
	HORZRES       = C.HORZRES       // px width
	VERTRES       = C.VERTRES       // px height
	VREFRESH      = C.VREFRESH      // current refresh rate
	CM_GAMMA_RAMP = C.CM_GAMMA_RAMP // supports gamma ramps
)

type DISPLAY_DEVICE C.DISPLAY_DEVICEW
type DEVMODE C.DEVMODEW
type POINTL C.POINTL
type GUID C.GUID
type POINT C.POINT
type MINMAXINFO C.MINMAXINFO
type KBDLLHOOKSTRUCT C.KBDLLHOOKSTRUCT
type PIXELFORMATDESCRIPTOR C.PIXELFORMATDESCRIPTOR
type OSVERSIONINFOEX C.OSVERSIONINFOEX
type RECT C.RECT

// Below here is constants and types that MinGW does not define or cgo -godefs
// cannot handle.

const (
	RIDEV_INPUTSINK         = 0x00000100
	RID_INPUT               = 0x10000003
	RIM_TYPEMOUSE           = 0
	HID_USAGE_PAGE_GENERIC  = 0x01
	HID_USAGE_GENERIC_MOUSE = 0x02
	ICON_SMALL2             = 2
	WM_INPUT                = 0x00FF
	WM_XBUTTONDOWN          = 0x020B
	WM_XBUTTONUP            = 0x020C
	WM_XBUTTONDBLCLK        = 0x020D
	MK_CONTROL              = 0x0008
	MK_LBUTTON              = 0x0001
	MK_MBUTTON              = 0x0010
	MK_RBUTTON              = 0x0002
	MK_SHIFT                = 0x0004
	MK_XBUTTON1             = 0x0020
	MK_XBUTTON2             = 0x0040
	WM_MOUSEHWHEEL          = 0x020E

	DMDFO_DEFAULT = 0
	DMDFO_STRETCH = 1
	DMDFO_CENTER  = 2
)

type MONITORINFOEX C.win32_MONITORINFOEX
type RAWINPUTHEADER C.win32_RAWINPUTHEADER
type RAWMOUSE C.win32_RAWMOUSE
type RAWINPUT C.win32_RAWINPUT
type VIDEOPARAMETERS C.win32_VIDEOPARAMETERS

type RAWINPUTDEVICE struct {
	UsUsagePage uint16
	UsUsage     uint16
	DwFlags     uint32
	HwndTarget  HWND
}
