// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build windows

package win32

import (
	"syscall"
	"unsafe"
)

var (
	gdi32                   = syscall.NewLazyDLL("gdi32.dll")
	pCreateDC               = gdi32.NewProc("CreateDCW")
	pCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	pSetDIBits              = gdi32.NewProc("SetDIBits")
)

func CreateDC(driver string, device string, initData *DEVMODE) HDC {
	if err := pCreateDC.Find(); err != nil {
		panic("CreateDCW missing!")
	}

	var cdriver *uint16
	if len(driver) > 0 {
		cdriver, _ = syscall.UTF16PtrFromString(driver)
	}

	var cdevice *uint16
	if len(device) > 0 {
		cdevice, _ = syscall.UTF16PtrFromString(device)
	}

	cRet, _, _ := pCreateDC.Call(
		uintptr(unsafe.Pointer(cdriver)),
		uintptr(unsafe.Pointer(cdevice)),
		uintptr(0),
		uintptr(unsafe.Pointer(initData)),
	)
	return HDC(cRet)
}

func CreateCompatibleBitmap(hdc HDC, nWidth, nHeight Int) HBITMAP {
	if err := pCreateCompatibleBitmap.Find(); err != nil {
		panic("CreateCompatibleBitmap missing!")
	}
	cRet, _, _ := pCreateCompatibleBitmap.Call(
		uintptr(hdc),
		uintptr(nWidth),
		uintptr(nHeight),
	)
	return HBITMAP(cRet)
}

func SetDIBits(hdc HDC, hbmp HBITMAP, uStartScan, cScanLines UINT, lpvBits unsafe.Pointer, bmi *BITMAPINFO, fuColorUse UINT) Int {
	if err := pSetDIBits.Find(); err != nil {
		panic("SetDIBits missing!")
	}
	cRet, _, _ := pSetDIBits.Call(
		uintptr(hdc),
		uintptr(hbmp),
		uintptr(uStartScan),
		uintptr(cScanLines),
		uintptr(lpvBits),
		uintptr(unsafe.Pointer(bmi)),
		uintptr(fuColorUse),
	)
	return Int(cRet)
}
