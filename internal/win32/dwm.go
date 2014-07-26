// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package win32

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	dwmapi                        = syscall.NewLazyDLL("dwmapi.dll")
	pDwmEnableBlurBehindWindow    = dwmapi.NewProc("DwmEnableBlurBehindWindow")
	pDwmExtendFrameIntoClientArea = dwmapi.NewProc("DwmExtendFrameIntoClientArea")
)

type DWM_BLURBEHIND struct {
	DwFlags                uint32
	FEnable                int32
	HRgbBlur               HRGN
	FTransitionOnMaximized int32
}

const (
	DWM_BB_ENABLE                = 0x00000001 //     A value for the fEnable member has been specified.
	DWM_BB_BLURREGION            = 0x00000002 //     A value for the hRgnBlur member has been specified.
	DWM_BB_TRANSITIONONMAXIMIZED = 0x00000004 //     A value for the fTransitionOnMaximized member has been specified.
)

func DwmEnableBlurBehindWindow(hwnd HWND, pBlurBehind *DWM_BLURBEHIND) error {
	err := pDwmEnableBlurBehindWindow.Find()
	if err != nil {
		return err
	}

	cRet, _, _ := pDwmEnableBlurBehindWindow.Call(
		uintptr(unsafe.Pointer(hwnd)),
		uintptr(unsafe.Pointer(pBlurBehind)),
	)
	ret := int64(cRet)

	if ret >= 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("DwmEnableBlurBehindWindow(): HRESULT = %d", ret))
	}
}

type MARGINS struct {
	CxLeftWidth, CxRightWidth, CyTopHeight, CyBottomHeight int
}

func DwmExtendFrameIntoClientArea(hwnd HWND, pMarInset *MARGINS) error {
	err := pDwmExtendFrameIntoClientArea.Find()
	if err != nil {
		return err
	}

	cRet, _, _ := pDwmExtendFrameIntoClientArea.Call(
		uintptr(unsafe.Pointer(hwnd)),
		uintptr(unsafe.Pointer(pMarInset)),
	)
	ret := int64(cRet)

	if ret >= 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("DwmExtendFrameIntoClientArea(): HRESULT = %d", ret))
	}
}
