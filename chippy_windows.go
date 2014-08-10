// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"azul3d.org/chippy.v1/internal/win32"
	"azul3d.org/keyboard.v1"
)

func eventLoop() {
	for {
		hasMessage := true

		var msg win32.MSG
		dispatch(func() {
			for hasMessage {
				hasMessage = win32.PeekMessage(&msg, nil, 0, 0, win32.PM_REMOVE)
				if hasMessage {
					win32.TranslateMessage(&msg)
					win32.DispatchMessage(&msg)
				}
			}
		})

		if !hasMessage {
			// let thread idle
			time.Sleep(10 * time.Millisecond)
		}
	}
}

var classNameCounter = 0
var classNameCounterAccess sync.Mutex

func nextCounter() int {
	classNameCounterAccess.Lock()
	defer classNameCounterAccess.Unlock()
	classNameCounter++
	return classNameCounter
}

var windowsKeyDisabled bool

func SetWindowsKeyDisabled(disabled bool) {
	globalLock.Lock()
	defer globalLock.Unlock()

	windowsKeyDisabled = disabled
}

func WindowsKeyDisabled() bool {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return windowsKeyDisabled
}

var hKeyboardHook win32.HHOOK

func keyboardHook(nCode win32.Int, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	if nCode < 0 || nCode != win32.HC_ACTION {
		return win32.CallNextHookEx(hKeyboardHook, nCode, wParam, lParam)
	}

	eatKeystroke := false
	if wParam == win32.WM_KEYDOWN || wParam == win32.WM_KEYUP {
		if WindowsKeyDisabled() {
			p := (*win32.KBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))

			keysToEat := []win32.DWORD{
				win32.VK_LWIN,
				win32.VK_RWIN,
			}

			anyKeysToEat := false
			for _, k := range keysToEat {
				if k == win32.DWORD(p.VkCode) {
					anyKeysToEat = true
					break
				}
			}

			if anyKeysToEat {
				for _, window := range windowsByHwnd {
					if window.r.Focused() {
						eatKeystroke = true

						// Send the event to the window
						state := keyboard.Down
						if wParam == win32.WM_KEYUP {
							state = keyboard.Up
						}

						switch p.VkCode {
						case win32.VK_LWIN:
							window.r.tryAddKeyboardStateEvent(
								keyboard.LeftSuper,
								uint64(win32.VK_LWIN),
								state,
							)

						case win32.VK_RWIN:
							window.r.tryAddKeyboardStateEvent(
								keyboard.RightSuper,
								uint64(win32.VK_RWIN),
								state,
							)
						}
					}
				}
			}
		}
	}

	if eatKeystroke {
		return 1
	}
	return win32.CallNextHookEx(hKeyboardHook, nCode, wParam, lParam)
}

//var classAtom win32.ATOM
//var windowClass *win32.WNDCLASSEX
var hInstance win32.HINSTANCE

var w32VersionMajor, w32VersionMinor win32.DWORD

func backend_Init() error {
	windowsKeyDisabled = true

	go dispatch(func() {
		hInstance = win32.HINSTANCE(win32.GetModuleHandle(""))
		if hInstance == nil {
			logger().Printf(fmt.Sprintf("Unable to determine hInstance; GetModuleHandle():", win32.GetLastErrorString()))
		}

		// Get OS version, we use this to do some hack-ish fixes for different windows versions
		ret, vi := win32.GetVersionEx()
		if ret {
			w32VersionMajor = win32.DWORD(vi.DwMajorVersion)
			w32VersionMinor = win32.DWORD(vi.DwMinorVersion)
		} else {
			logger().Printf("Unable to determine windows version information; GetVersionEx():", win32.GetLastErrorString())
		}

		// It's not safe to set a low level keyboard hook if our process is
		// 32-bit and we are running in a 64-bit OS, additionally it seems that
		// there is a bug with 32-bit versions of windows and keyboard hooks
		// that causes a crash (we should look into it).
		if runtime.GOARCH != "386" {
			hKeyboardHook = win32.SetLowLevelKeyboardHook(keyboardHook, hInstance, 0)
			if hKeyboardHook == nil {
				logger().Println("Failed to disable keyboard shortcuts; SetWindowsHookEx():", win32.GetLastErrorString())
			}
		}
	})

	go eventLoop()

	return nil
}

func backend_Destroy() {
	dispatch(func() {
		if hKeyboardHook != nil {
			if !win32.UnhookWindowsHookEx(hKeyboardHook) {
				logger().Println("Failed to unhook keyboard hook; UnhookWindowsHookEx():", win32.GetLastErrorString())
			}
		}
	})

	classNameCounter = 0
}
