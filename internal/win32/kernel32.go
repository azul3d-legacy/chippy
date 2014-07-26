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
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	pFormatMessage   = kernel32.NewProc("FormatMessageW")
	pGetModuleHandle = kernel32.NewProc("GetModuleHandleW")
	pGetVersionEx    = kernel32.NewProc("GetVersionExW")
)

func FormatMessage(dwFlags DWORD, lpSource unsafe.Pointer, dwMessageId, dwLanguageId DWORD, lpBuffer unsafe.Pointer, nSize DWORD, args unsafe.Pointer) DWORD {
	if err := pFormatMessage.Find(); err != nil {
		panic("FormatMessageW missing!")
	}

	cRet, _, _ := pFormatMessage.Call(
		uintptr(dwFlags),
		uintptr(lpSource),
		uintptr(dwMessageId),
		uintptr(dwLanguageId),
		uintptr(lpBuffer),
		uintptr(nSize),
		uintptr(args),
	)
	return DWORD(cRet)
}

func GetModuleHandle(moduleName string) HMODULE {
	if err := pGetModuleHandle.Find(); err != nil {
		panic("GetModuleHandleW missing!")
	}

	var cModuleName *uint16
	if len(moduleName) > 0 {
		cModuleName, _ = syscall.UTF16PtrFromString(moduleName)
	}
	cRet, _, _ := pGetModuleHandle.Call(uintptr(unsafe.Pointer(cModuleName)))
	return HMODULE(unsafe.Pointer(cRet))
}

func GetVersionEx() (ret bool, vi *OSVERSIONINFOEX) {
	if err := pGetVersionEx.Find(); err != nil {
		panic("GetVersionExW missing!")
	}

	vi = new(OSVERSIONINFOEX)
	vi.DwOSVersionInfoSize = uint32(unsafe.Sizeof(OSVERSIONINFOEX{}))
	cRet, _, _ := pGetVersionEx.Call(
		uintptr(unsafe.Pointer(vi)),
	)
	return cRet != 0, vi
}
