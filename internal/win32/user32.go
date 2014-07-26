// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package win32

/*
#define UNICODE
#include <windows.h>
LPTSTR macro_MAKEINTRESOURCE(WORD wInteger);
*/
import "C"

import (
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	pEnumDisplayDevices      = user32.NewProc("EnumDisplayDevicesW")
	pEnumDisplaySettings     = user32.NewProc("EnumDisplaySettingsW")
	pChangeDisplaySettingsEx = user32.NewProc("ChangeDisplaySettingsExW")
	pRegisterRawInputDevices = user32.NewProc("RegisterRawInputDevices")
	pSetWindowLongPtr        = user32.NewProc("SetWindowLongPtrW")
	pGetWindowLongPtr        = user32.NewProc("GetWindowLongPtrW")
	pSetWindowLong           = user32.NewProc("SetWindowLongW")
	pGetWindowLong           = user32.NewProc("GetWindowLongW")
	pEnableWindow            = user32.NewProc("EnableWindow")
	pFlashWindow             = user32.NewProc("FlashWindow")
	pRegisterClassEx         = user32.NewProc("RegisterClassExW")
	pUnregisterClass         = user32.NewProc("UnregisterClassW")
	pLoadCursor              = user32.NewProc("LoadCursorW")
	pSetWindowText           = user32.NewProc("SetWindowTextW")
	pDestroyCursor           = user32.NewProc("DestroyCursor")
	pCreateIconIndirect      = user32.NewProc("CreateIconIndirect")
	pDestroyIcon             = user32.NewProc("DestroyIcon")
	pGetRawInputData         = user32.NewProc("GetRawInputData")
	pSetWindowPos            = user32.NewProc("SetWindowPos")
	pGetClipCursor           = user32.NewProc("GetClipCursor")
	pClipCursor              = user32.NewProc("ClipCursor")
	pGetUpdateRect           = user32.NewProc("GetUpdateRect")
	pValidateRect            = user32.NewProc("ValidateRect")
	pPeekMessage             = user32.NewProc("PeekMessageW")
	pTranslateMessage        = user32.NewProc("TranslateMessage")
)

const (
	// We define these manually because some MinGW versions don't define them.
	DISPLAY_DEVICE_ATTACHED_TO_DESKTOP = 0x00000001
	DISPLAY_DEVICE_MULTI_DRIVER        = 0x00000002
	DISPLAY_DEVICE_PRIMARY_DEVICE      = 0x00000004
	DISPLAY_DEVICE_MIRRORING_DRIVER    = 0x00000008
	DISPLAY_DEVICE_VGA_COMPATIBLE      = 0x00000010
	DISPLAY_DEVICE_REMOVABLE           = 0x00000020
	DISPLAY_DEVICE_ACC_DRIVER          = 0x00000040
	DISPLAY_DEVICE_TS_COMPATIBLE       = 0x00200000
	DISPLAY_DEVICE_UNSAFE_MODES_ON     = 0x00080000
	DISPLAY_DEVICE_MODESPRUNED         = 0x08000000
	DISPLAY_DEVICE_REMOTE              = 0x04000000
	DISPLAY_DEVICE_DISCONNECT          = 0x02000000
	DISPLAY_DEVICE_ACTIVE              = 0x00000001
	DISPLAY_DEVICE_ATTACHED            = 0x00000002
)

func EnumDisplayDevices(device string, iDevNum DWORD, pDisplayDevice *DISPLAY_DEVICE, dwFlags DWORD) bool {
	if err := pEnumDisplayDevices.Find(); err != nil {
		panic("EnumDisplayDevicesW missing!")
	}

	var cdevice *uint16
	if len(device) > 0 {
		cdevice, _ = syscall.UTF16PtrFromString(device)
	}

	pDisplayDevice.Cb = uint32(unsafe.Sizeof(DISPLAY_DEVICE{}))
	cRet, _, _ := pEnumDisplayDevices.Call(
		uintptr(unsafe.Pointer(cdevice)),
		uintptr(iDevNum),
		uintptr(unsafe.Pointer(pDisplayDevice)),
		uintptr(dwFlags),
	)
	if cRet == 0 {
		return false
	}
	return true
}

// See:
// http://msdn.microsoft.com/en-us/library/windows/desktop/dd183565(v=vs.85).aspx
// DmDisplayFixedOutput is inside a union.
type union_dmDisplayFixedOutput struct {
	DmPosition           POINTL
	DmDisplayOrientation DWORD
	DmDisplayFixedOutput DWORD
}

func (d *DEVMODE) DmDisplayFixedOutput() DWORD {
	u := (*union_dmDisplayFixedOutput)(unsafe.Pointer(&d.Anon0))
	return u.DmDisplayFixedOutput
}

const (
	// We define these manually because some MinGW versions don't define them.
	DM_BITSPERPEL       = 0x00040000
	DM_PELSWIDTH        = 0x00080000
	DM_PELSHEIGHT       = 0x00100000
	DM_DISPLAYFLAGS     = 0x00200000
	DM_DISPLAYFREQUENCY = 0x00400000
	DM_POSITION         = 0x00000020
)

func EnumDisplaySettings(deviceName string, iModeNum DWORD, pDevMode *DEVMODE) bool {
	if err := pEnumDisplaySettings.Find(); err != nil {
		panic("EnumDisplaySettingsW missing!")
	}

	var cDeviceName *uint16
	if len(deviceName) > 0 {
		cDeviceName, _ = syscall.UTF16PtrFromString(deviceName)
	}

	pDevMode.DmSize = uint16(unsafe.Sizeof(DEVMODE{}))
	cRet, _, _ := pEnumDisplaySettings.Call(
		uintptr(unsafe.Pointer(cDeviceName)),
		uintptr(iModeNum),
		uintptr(unsafe.Pointer(pDevMode)),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func ChangeDisplaySettingsEx(deviceName string, pDevMode *DEVMODE, dwFlags DWORD, lParam *VIDEOPARAMETERS) (ret LONG) {
	if err := pChangeDisplaySettingsEx.Find(); err != nil {
		panic("ChangeDisplaySettingsExW missing!")
	}

	var cDeviceName *uint16
	if len(deviceName) > 0 {
		cDeviceName, _ = syscall.UTF16PtrFromString(deviceName)
	}

	if pDevMode != nil {
		pDevMode.DmSize = uint16(unsafe.Sizeof(DEVMODE{}))
	}
	cRet, _, _ := pChangeDisplaySettingsEx.Call(
		uintptr(unsafe.Pointer(cDeviceName)),
		uintptr(unsafe.Pointer(pDevMode)),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(lParam)),
	)

	return LONG(cRet)
}

func RegisterRawInputDevices(pRawInputDevices *RAWINPUTDEVICE, uiNumDevices UINT, cbSize UINT) bool {
	if err := pRegisterRawInputDevices.Find(); err != nil {
		panic("RegisterRawInputDevices missing!")
	}
	cRet, _, _ := pRegisterRawInputDevices.Call(
		uintptr(unsafe.Pointer(pRawInputDevices)),
		uintptr(uiNumDevices),
		uintptr(cbSize),
	)
	if cRet == 0 {
		return false
	}
	return true
}

//func GetRegisteredRawInputDevices(pRawInputDevices *RAWINPUTDEVICE, puiNumDevices *UINT, cbSize UINT) UINT {
//	return UINT(C.GetRegisteredRawInputDevices((C.PRAWINPUTDEVICE)(unsafe.Pointer(pRawInputDevices)), (C.PUINT)(unsafe.Pointer(puiNumDevices)), C.UINT(cbSize)))
//}

func GetRawInputData(hRawInput HRAWINPUT, uiCommand UINT, pData unsafe.Pointer, pcbSize *UINT, cbSizeHeader UINT) UINT {
	if err := pGetRawInputData.Find(); err != nil {
		panic("GetRawInputData missing!")
	}
	cRet, _, _ := pGetRawInputData.Call(
		uintptr(hRawInput),
		uintptr(uiCommand),
		uintptr(pData),
		uintptr(unsafe.Pointer(pcbSize)),
		uintptr(cbSizeHeader),
	)
	return UINT(cRet)
}

func SetWindowLong(hwnd HWND, nIndex Int, dwNewLong LONG) (ret LONG) {
	if err := pSetWindowLong.Find(); err != nil {
		panic("SetWindowLongPtrW AND SetWindowLongW missing!")
	}
	cRet, _, _ := pSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(nIndex),
		uintptr(dwNewLong),
	)
	return LONG(cRet)
}

func SetWindowLongPtr(hwnd HWND, nIndex Int, dwNewLong LONG_PTR) (ret LONG_PTR) {
	if err := pSetWindowLongPtr.Find(); err != nil {
		return LONG_PTR(SetWindowLong(hwnd, nIndex, LONG(dwNewLong)))
	}
	cRet, _, _ := pSetWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(nIndex),
		uintptr(dwNewLong),
	)
	return LONG_PTR(cRet)
}

func GetWindowLong(hwnd HWND, nIndex Int) (ret LONG) {
	if err := pGetWindowLong.Find(); err != nil {
		panic("GetWindowLongPtrW AND GetWindowLongW missing!")
	}
	cRet, _, _ := pGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(nIndex),
	)
	return LONG(cRet)
}

func GetWindowLongPtr(hwnd HWND, nIndex Int) (ret LONG_PTR) {
	if err := pGetWindowLongPtr.Find(); err != nil {
		return LONG_PTR(GetWindowLong(hwnd, nIndex))
	}
	cRet, _, _ := pGetWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(nIndex),
	)
	return LONG_PTR(cRet)
}

func EnableWindow(hwnd HWND, enable bool) bool {
	if err := pEnableWindow.Find(); err != nil {
		panic("EnableWindow missing!")
	}
	var bEnable uintptr
	if enable {
		bEnable = 1
	}
	cRet, _, _ := pEnableWindow.Call(
		uintptr(hwnd),
		bEnable,
	)
	if cRet == 0 {
		return false
	}
	return true
}

func FlashWindow(hwnd HWND, invert bool) bool {
	if err := pFlashWindow.Find(); err != nil {
		panic("FlashWindow missing!")
	}
	var bInvert uintptr
	if invert {
		bInvert = 1
	}
	cRet, _, _ := pFlashWindow.Call(
		uintptr(hwnd),
		bInvert,
	)
	if cRet == 0 {
		return false
	}
	return true
}

func RegisterClassEx(wc *WNDCLASSEX) ATOM {
	if err := pRegisterClassEx.Find(); err != nil {
		panic("RegisterClassExW missing!")
	}
	cRet, _, _ := pRegisterClassEx.Call(
		uintptr(unsafe.Pointer(wc)),
	)
	return ATOM(cRet)
}

func UnregisterClass(class string, instance HINSTANCE) bool {
	if err := pUnregisterClass.Find(); err != nil {
		panic("UnregisterClassW missing!")
	}
	var cClass *uint16
	if len(class) > 0 {
		cClass, _ = syscall.UTF16PtrFromString(class)
	}
	cRet, _, _ := pUnregisterClass.Call(
		uintptr(unsafe.Pointer(cClass)),
		uintptr(instance),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func SetWindowText(hwnd HWND, str string) bool {
	if err := pSetWindowText.Find(); err != nil {
		panic("SetWindowTextW missing!")
	}
	var cString *uint16
	if len(str) > 0 {
		cString, _ = syscall.UTF16PtrFromString(str)
	}
	cRet, _, _ := pSetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(cString)),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func LoadCursor(instance HINSTANCE, cursorName string) HCURSOR {
	if err := pLoadCursor.Find(); err != nil {
		panic("LoadCursorW missing!")
	}
	var cCursorName unsafe.Pointer
	if cursorName == "IDC_ARROW" {
		cCursorName = unsafe.Pointer(C.macro_MAKEINTRESOURCE(32512))
	} else {
		if len(cursorName) > 0 {
			cptr, _ := syscall.UTF16PtrFromString(cursorName)
			cCursorName = unsafe.Pointer(cptr)
		}
	}
	cRet, _, _ := pLoadCursor.Call(
		uintptr(instance),
		uintptr(cCursorName),
	)
	return HCURSOR(cRet)
}

func DestroyCursor(cursor HCURSOR) bool {
	if err := pDestroyCursor.Find(); err != nil {
		panic("DestroyCursor missing!")
	}
	cRet, _, _ := pDestroyCursor.Call(
		uintptr(cursor),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func DestroyIcon(icon HICON) bool {
	if err := pDestroyIcon.Find(); err != nil {
		panic("DestroyIcon missing!")
	}
	cRet, _, _ := pDestroyIcon.Call(
		uintptr(icon),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func CreateIconIndirect(piconinfo *ICONINFO) HICON {
	if err := pCreateIconIndirect.Find(); err != nil {
		panic("CreateIconIndirect missing!")
	}
	cRet, _, _ := pCreateIconIndirect.Call(
		uintptr(unsafe.Pointer(piconinfo)),
	)
	return HICON(cRet)
}

func ClipCursor(rect *RECT) bool {
	if err := pClipCursor.Find(); err != nil {
		panic("ClipCursor missing!")
	}
	cRet, _, _ := pClipCursor.Call(
		uintptr(unsafe.Pointer(rect)),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func GetClipCursor() (clip *RECT, ok bool) {
	if err := pGetClipCursor.Find(); err != nil {
		panic("GetClipCursor missing!")
	}
	cRet, _, _ := pGetClipCursor.Call(
		uintptr(unsafe.Pointer(clip)),
	)
	if cRet != 0 {
		ok = true
		return
	}
	ok = false
	return
}

func GetUpdateRect(hwnd HWND, lpRect *RECT, bErase bool) bool {
	if err := pGetUpdateRect.Find(); err != nil {
		panic("GetUpdateRect missing!")
	}
	var cbool uintptr
	if bErase {
		cbool = 1
	}
	cRet, _, _ := pGetUpdateRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpRect)),
		cbool,
	)
	if cRet == 0 {
		return false
	}
	return true
}

func ValidateRect(hwnd HWND, rect *RECT) bool {
	if err := pValidateRect.Find(); err != nil {
		panic("ValidateRect missing!")
	}
	cRet, _, _ := pValidateRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(rect)),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func SetWindowPos(hwnd, hwndInsertAfter HWND, X, Y, cx, cy Int, uFlags UINT) (ret bool) {
	if err := pSetWindowPos.Find(); err != nil {
		panic("SetWindowPos missing!")
	}
	cRet, _, _ := pSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hwndInsertAfter),
		uintptr(X),
		uintptr(Y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func PeekMessage(msg *MSG, hwnd HWND, wMsgFilterMin, wMsgFilterMax, wRemoveMsg UINT) bool {
	if err := pPeekMessage.Find(); err != nil {
		panic("PeekMessageW missing!")
	}
	cRet, _, _ := pPeekMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		uintptr(wRemoveMsg),
	)
	if cRet == 0 {
		return false
	}
	return true
}

func TranslateMessage(msg *MSG) (ret bool) {
	if err := pTranslateMessage.Find(); err != nil {
		panic("TranslateMessage missing!")
	}
	cRet, _, _ := pTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)
	if cRet == 0 {
		return false
	}
	return true
}

/*
func SetTimer(hwnd HWND, nIDEvent UINT_PTR, uElapse UINT, lpTimerFunc TIMERPROC) (timer UINT_PTR) {
	timer = UINT_PTR(C.SetTimer(C.HWND(hwnd), C.UINT_PTR(nIDEvent), C.UINT(uElapse), C.TIMERPROC(lpTimerFunc)))
	return
}

func KillTimer(hwnd HWND, uIDEvent UINT_PTR) (ret bool) {
	ret = C.KillTimer(C.HWND(hwnd), C.UINT_PTR(uIDEvent)) != 0
	return
}

func GetWindowRect(hwnd HWND) (status bool, r *RECT) {
	var cr C.RECT
	status = C.GetWindowRect(C.HWND(hwnd), (C.LPRECT)(unsafe.Pointer(&cr))) != 0
	r = (*RECT)(&cr)
	return
}

func GetClientRect(hwnd HWND) (status bool, r *RECT) {
	var cr C.RECT
	status = C.GetClientRect(C.HWND(hwnd), (C.LPRECT)(unsafe.Pointer(&cr))) != 0
	r = (*RECT)(&cr)
	return
}

func MoveWindow(hwnd HWND, x, y, width, height Int, repaint bool) bool {
	cbool := C.WINBOOL(0)
	if repaint {
		cbool = 1
	}
	return C.MoveWindow(C.HWND(hwnd), C.int(x), C.int(y), C.int(width), C.int(height), cbool) != 0
}
*/
