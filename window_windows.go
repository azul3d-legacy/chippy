// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"azul3d.org/chippy.v1/internal/resize"
	"azul3d.org/chippy.v1/internal/win32"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
	"errors"
	"fmt"
	"image"
	"math"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

var windowsByHwnd = make(map[win32.HWND]*NativeWindow)

type loadedCursor struct {
	hCursor                             win32.HICON
	cursorColorBitmap, cursorMaskBitmap win32.HBITMAP
	cursorColorBits, cursorMaskBits     []uint32
}

type NativeWindow struct {
	access sync.RWMutex

	r *Window

	insideCallback bool

	cursors                                                                  map[*Cursor]*loadedCursor
	loadedCursor                                                             *loadedCursor
	hIcon, hSmIcon                                                           win32.HICON
	iconColorBitmap, iconMaskBitmap, smIconColorBitmap, smIconMaskBitmap     win32.HBITMAP
	iconColorBits, iconMaskBits, smIconColorBits, smIconMaskBits             []uint32
	dc, dcRender                                                             win32.HDC
	hwnd                                                                     win32.HWND
	windowClass                                                              string
	styleFlags                                                               win32.DWORD
	lastWmSizingLeft, lastWmSizingRight, lastWmSizingBottom, lastWmSizingTop int32
	lastCursorClip                                                           *win32.RECT

	// Blit things here
	blitLock     sync.Mutex
	blitBitmap   win32.HBITMAP
	blitBitmapDc win32.HDC
	blitBits     []uint32

	// OpenGL things here
	glVSyncMode      VSyncMode
	glConfig         *GLConfig
	glPixelFormatSet bool
}

func (w *NativeWindow) open(screen *Screen) (err error) {
	unlock := w.newAttemptUnlocker()
	defer unlock()

	w.cursors = make(map[*Cursor]*loadedCursor)

	w.r.trySetFocused(true)

	unlock()
	dispatch(func() {
		err = w.doRebuildWindow()
		if err != nil {
			return
		}

		// To let the application know of current toggle key states, we need to send them right now
		// otherwise they might already be on or off -- something the application might want to
		// know.
		virtualKeyMap := map[keyboard.Key]win32.Int{
			keyboard.CapsLock:   win32.VK_CAPITAL,
			keyboard.NumLock:    win32.VK_NUMLOCK,
			keyboard.ScrollLock: win32.VK_SCROLL,
			keyboard.LeftAlt:    win32.VK_LMENU,
			keyboard.RightAlt:   win32.VK_RMENU,
			keyboard.LeftCtrl:   win32.VK_LCONTROL,
			keyboard.RightCtrl:  win32.VK_RCONTROL,
			keyboard.LeftShift:  win32.VK_LSHIFT,
			keyboard.RightShift: win32.VK_RSHIFT,
			keyboard.LeftSuper:  win32.VK_LWIN,
			keyboard.RightSuper: win32.VK_RWIN,
		}

		var state keyboard.State
		for key, virtualKey := range virtualKeyMap {
			state = keyboard.On
			if (win32.GetAsyncKeyState(virtualKey) & 0x0001) == 0 {
				state = keyboard.Off
			}
			w.r.tryAddKeyboardStateEvent(key, uint64(virtualKey), state)
		}

		if w.r.Cursor() != nil {
			w.doSetCursor()
		}
	})

	if err == nil && w.r.Cursor() == nil {
		w.setCursor(nil)
	}
	return
}

func (w *NativeWindow) doRebuildWindow() (err error) {
	if w.hwnd != nil {
		w.doDestroy()
	}
	w.glPixelFormatSet = false

	// Make our window class
	w.windowClass = fmt.Sprintf("ChippyWindow%d", nextCounter())
	windowClass := win32.NewWNDCLASSEX()
	windowClass.SetLpfnWndProc()
	windowClass.SetHbrBackground(win32.CreateSolidBrush(0x00000000))

	// CS_OWNDC is needed to avoid some bugs with multiple windows, older
	// versions of windows, and some different (archaic?) graphics drivers.
	windowClass.SetStyle(win32.CS_OWNDC)

	//windowClass.SetHIcon(win32.LoadIcon(hInstance, szAppName))
	//windowClass.SetHCursor(win32.LoadCursor(nil, "IDC_ARROW"))
	//windowClass.SetHbrBackground(win32.IntToHBRUSH(win32.COLOR_WINDOW+2)) // Black background
	//windowClass.SetLpszMenuName(szAppName)

	windowClass.SetHInstance(hInstance)
	windowClass.SetLpszClassName(w.windowClass)

	classAtom := win32.RegisterClassEx(windowClass)
	if classAtom == 0 {
		err = errors.New(fmt.Sprintf("Unable to open window; RegisterClassEx(): %s", win32.GetLastErrorString()))
		return
	}

	// w.styleFlags will be updated to reflect current settings, that are passed into
	// CreateWindowEx to avoid some flicker
	w.doUpdateStyle()

	// SetPixelFormat() may only be called once -- so if we want to change any pixel format
	// values then our only option is to destroy the window and create it again.
	//
	// Since that would provide an largely noticable flicker to the user, we instead have an
	// 'rendering' window parented to our 'user managed' window, and we create the 'rendering'
	// window whenever we want to, thus bypassing the SetPixelFormat() issue noted above.
	//
	w.hwnd = win32.CreateWindowEx(0, w.windowClass, w.r.Title(), w.styleFlags, 0, 0, 0, 0, nil, nil, hInstance, nil)
	if w.hwnd == nil {
		err = errors.New(fmt.Sprintf("Unable to open window; CreateWindowEx(): %s", win32.GetLastErrorString()))
		return
	}
	w.dc = win32.GetDC(w.hwnd)
	if w.dc == nil {
		err = errors.New(fmt.Sprintf("Unable to get window DC; GetDC(): %s", win32.GetLastErrorString()))
		return
	}

	w.doUpdateTransparency()
	w.doSetWindowPos()

	// Make sure to enable opened now so that doUpdateStyle sets the new style properly
	w.doUpdateStyle()

	if w.r.Visible() {
		win32.ShowWindow(w.hwnd, win32.SW_SHOWDEFAULT)
		if w.r.Minimized() {
			win32.ShowWindow(w.hwnd, win32.SW_MINIMIZE)
		} else if w.r.Maximized() {
			win32.ShowWindow(w.hwnd, win32.SW_MAXIMIZE)
		}
	}

	windowsByHwnd[w.hwnd] = w
	win32.RegisterWndProc(w.hwnd, mainWindowProc)

	w.doMakeIcon()

	supportRawInput := w32VersionMajor >= 5 && w32VersionMinor >= 1
	if supportRawInput {
		rid := win32.RAWINPUTDEVICE{}
		rid.UsUsagePage = win32.HID_USAGE_PAGE_GENERIC
		rid.UsUsage = win32.HID_USAGE_GENERIC_MOUSE
		rid.DwFlags = win32.RIDEV_INPUTSINK
		rid.HwndTarget = w.hwnd
		win32.RegisterRawInputDevices(&rid, 1, win32.UINT(unsafe.Sizeof(rid)))
	}

	return
}

func (w *NativeWindow) doDestroy() {
	win32.UnregisterWndProc(w.hwnd)
	delete(windowsByHwnd, w.hwnd)

	if !win32.DestroyWindow(w.hwnd) {
		logger().Println("Unable to destroy window; DestroyWindow():", win32.GetLastErrorString())
	}

	if !win32.UnregisterClass(w.windowClass, hInstance) {
		logger().Println("Failed to unregister window class; UnregisterClass():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) destroy() {
	dispatch(func() {
		w.doDestroy()
	})
}

func (w *NativeWindow) notify() {
	blinkDelay := 3 * time.Second
	maxBlinks := 3

	timesBlinked := 0
	for {
		if w.r.Destroyed() {
			return
		}
		if timesBlinked >= maxBlinks {
			return
		}
		timesBlinked += 1
		dispatch(func() {
			win32.FlashWindow(w.hwnd, true)
		})
		<-time.After(blinkDelay)
	}
}

func (w *NativeWindow) PixelClear(rect image.Rectangle) {
	if !w.r.Opened() {
		return
	}

	if rect.Empty() {
		return
	}

	dispatch(func() {
		r := &win32.RECT{
			Left:   int32(rect.Min.X),
			Top:    int32(rect.Min.Y),
			Bottom: int32(rect.Max.Y),
			Right:  int32(rect.Max.X),
		}
		win32.FillRect(w.dc, r, win32.HBRUSH(win32.GetStockObject(win32.BLACK_BRUSH)))
	})
}

func (w *NativeWindow) PixelBlit(x, y uint, image *image.RGBA) {
	if !w.r.Opened() {
		return
	}

	w.blitLock.Lock()
	defer w.blitLock.Unlock()

	if w.blitBitmap != nil {
		if !win32.DeleteDC(w.blitBitmapDc) {
			logger().Println("Unable to delete blit bitmap; DeleteDC():", win32.GetLastErrorString())
		}

		if !win32.DeleteObject(win32.HGDIOBJ(w.blitBitmap)) {
			logger().Println("Unable to delete blit bitmap; DeleteObject():", win32.GetLastErrorString())
		}
	}

	sz := image.Bounds().Size()
	if sz.X <= 0 && sz.Y <= 0 {
		return
	}
	width := uint(sz.X)
	height := uint(sz.Y)

	w.blitBitmapDc = win32.CreateCompatibleDC(w.dc)
	w.blitBitmap = win32.CreateCompatibleBitmap(w.dc, win32.Int(width), win32.Int(height))
	if win32.SelectObject(w.blitBitmapDc, win32.HGDIOBJ(w.blitBitmap)) == nil {
		logger().Println("Unable to blit; SelectObject():", win32.GetLastErrorString())
		return
	}

	w.blitBits = make([]uint32, width*height)

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			r, g, b, a := image.At(x, y).RGBA()
			c := (uint32(a>>8) << 24) | (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | (uint32(b>>8) << 0)

			index := (int(height) - 1 - y) * int(width)
			index += x
			w.blitBits[index] = c
		}
	}

	bitmapInfo := win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			Size:          win32.DWORD(unsafe.Sizeof(win32.BITMAPINFOHEADER{})),
			Width:         win32.LONG(width),
			Height:        win32.LONG(height),
			Planes:        1,
			BitCount:      32,
			Compression:   win32.BI_RGB,
			SizeImage:     0,
			XPelsPerMeter: 0,
			YPelsPerMeter: 0,
			ClrUsed:       0,
			ClrImportant:  0,
		},
	}
	if win32.SetDIBits(w.dc, w.blitBitmap, 0, win32.UINT(height), unsafe.Pointer(&w.blitBits[0]), &bitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to blit; SetDiBits():", win32.GetLastErrorString())
		return
	}

	r := &win32.RECT{
		Left:   int32(x),
		Top:    int32(y),
		Right:  int32(x + width),
		Bottom: int32(y + height),
	}
	brush := win32.GetStockObject(win32.BLACK_BRUSH)
	if !win32.FillRect(w.dc, r, win32.HBRUSH(brush)) {
		logger().Println("PixelBlit(): FillRect():", win32.GetLastErrorString())
		return
	}

	if !win32.AlphaBlend(
		w.dc,
		win32.Int(x),
		win32.Int(y),
		win32.Int(width),
		win32.Int(height),
		w.blitBitmapDc,
		0, 0,
		win32.Int(width),
		win32.Int(height),
		&win32.BLENDFUNCTION{
			BlendOp:             win32.AC_SRC_OVER,
			SourceConstantAlpha: 255,
			AlphaFormat:         win32.AC_SRC_ALPHA,
		},
	) {
		logger().Println("PixelBlit(): AlphaBlend():", win32.GetLastErrorString())
		return
	}
}

func (w *NativeWindow) setTransparent(transparent bool) {
	dispatch(func() {
		w.doUpdateTransparency()
	})
}

func (w *NativeWindow) setTitle(title string) {
	dispatch(func() {
		if !win32.SetWindowText(w.hwnd, title) {
			logger().Println("Unable to set window title; SetWindowText():", win32.GetLastErrorString())
		}
	})
}

func (w *NativeWindow) setVisible(visible bool) {
	if visible {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_SHOW)
			win32.EnableWindow(w.hwnd, true)
		})
	} else {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_HIDE)
		})
	}
}

func (w *NativeWindow) setDecorated(decorated bool) {
	dispatch(func() {
		if w.r.Visible() {
			w.doUpdateStyle()
		}
	})
}

func (w *NativeWindow) setPosition(x, y int) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setSize(width, height int) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setMinimumSize(width, height int) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setMaximumSize(width, height int) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setAspectRatio(ratio float32) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setMinimized(minimized bool) {
	if minimized {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_MINIMIZE)

			win32.EnableWindow(w.hwnd, true)
			w.doSetWindowPos()
		})
	} else {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_RESTORE)

			win32.EnableWindow(w.hwnd, true)
			w.doSetWindowPos()
		})
	}
}

func (w *NativeWindow) setMaximized(maximized bool) {
	if maximized {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_MAXIMIZE)

			win32.EnableWindow(w.hwnd, true)
		})
	} else {
		dispatch(func() {
			win32.ShowWindow(w.hwnd, win32.SW_RESTORE)

			win32.EnableWindow(w.hwnd, true)
		})
	}
}

func (w *NativeWindow) setFullscreen(fullscreen bool) {
	dispatch(func() {
		w.doUpdateStyle()
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setAlwaysOnTop(alwaysOnTop bool) {
	dispatch(func() {
		w.doSetWindowPos()
	})
}

func (w *NativeWindow) setIcon(icon image.Image) {
	dispatch(func() {
		w.doMakeIcon()
	})
}

func (w *NativeWindow) doPrepareCursor(cursor *Cursor) {
	lc := new(loadedCursor)

	cursorWidth := win32.GetSystemMetrics(win32.SM_CXCURSOR)
	cursorHeight := win32.GetSystemMetrics(win32.SM_CYCURSOR)

	cursorImage := resize.Resize(cursor.Image, cursor.Image.Bounds(), int(cursorWidth), int(cursorHeight))

	cursorBitmapInfo := win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			Size:          win32.DWORD(unsafe.Sizeof(win32.BITMAPINFOHEADER{})),
			Width:         win32.LONG(cursorWidth),
			Height:        win32.LONG(cursorHeight),
			Planes:        1,
			BitCount:      32,
			Compression:   win32.BI_RGB,
			SizeImage:     0,
			XPelsPerMeter: 0,
			YPelsPerMeter: 0,
			ClrUsed:       0,
			ClrImportant:  0,
		},
	}

	lc.cursorColorBitmap = win32.CreateCompatibleBitmap(w.dc, cursorWidth, cursorHeight)
	lc.cursorColorBits = make([]uint32, cursorWidth*cursorHeight)
	for y := 0; y < int(cursorHeight); y++ {
		for x := 0; x < int(cursorWidth); x++ {
			r, g, b, _ := cursorImage.At(x, y).RGBA()
			c := (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | (uint32(b>>8) << 0)

			index := (int(cursorHeight) - 1 - y) * int(cursorWidth)
			index += x
			lc.cursorColorBits[index] = c //0xFF0000
		}
	}
	if win32.SetDIBits(w.dc, lc.cursorColorBitmap, 0, win32.UINT(cursorHeight), unsafe.Pointer(&lc.cursorColorBits[0]), &cursorBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set cursor; SetDiBits():", win32.GetLastErrorString())
		return
	}

	lc.cursorMaskBitmap = win32.CreateCompatibleBitmap(w.dc, cursorWidth, cursorHeight)
	lc.cursorMaskBits = make([]uint32, cursorWidth*cursorHeight)
	for y := 0; y < int(cursorHeight); y++ {
		for x := 0; x < int(cursorWidth); x++ {
			_, _, _, a := cursorImage.At(x, y).RGBA()
			c := uint32(0xFFFFFF)
			if a > 0 {
				c = 0
			}

			index := (int(cursorHeight) - 1 - y) * int(cursorWidth)
			index += x
			lc.cursorMaskBits[index] = c
		}
	}
	if win32.SetDIBits(w.dc, lc.cursorMaskBitmap, 0, win32.UINT(cursorHeight), unsafe.Pointer(&lc.cursorMaskBits[0]), &cursorBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set cursor; SetDiBits():", win32.GetLastErrorString())
		return
	}

	cursorInfo := win32.ICONINFO{
		FIcon:    0,
		XHotspot: win32.DWORD(cursor.X),
		YHotspot: win32.DWORD(cursor.Y),
		HbmMask:  lc.cursorMaskBitmap,
		HbmColor: lc.cursorColorBitmap,
	}

	lc.hCursor = win32.CreateIconIndirect(&cursorInfo)
	if lc.hCursor == nil {
		logger().Println("Unable to set cursor; CreateIconIndirect():", win32.GetLastErrorString())
		return
	}

	w.cursors[cursor] = lc
}

func (w *NativeWindow) prepareCursor(cursor *Cursor) {
	unlock := w.newAttemptUnlocker()
	defer unlock()

	_, ok := w.cursors[cursor]
	if ok {
		// It's already loaded!
		return
	}

	unlock()
	dispatch(func() {
		w.doPrepareCursor(cursor)
	})
}

func (w *NativeWindow) freeCursor(cursor *Cursor) {
	unlock := w.newAttemptUnlocker()
	defer unlock()

	lc, ok := w.cursors[cursor]
	if ok {
		delete(w.cursors, cursor)

		unlock()
		dispatch(func() {
			if !win32.DestroyCursor(win32.HCURSOR(lc.hCursor)) {
				logger().Println("Failed to free cursor; DestroyCursor():", win32.GetLastErrorString())
			}

			if !win32.DeleteObject(win32.HGDIOBJ(lc.cursorColorBitmap)) {
				logger().Println("Failed to free cursor; DeleteObject(cursorColorBitmap) failed!")
			}

			if !win32.DeleteObject(win32.HGDIOBJ(lc.cursorMaskBitmap)) {
				logger().Println("Failed to free cursor; DeleteObject(cursorMaskBitmap) failed!")
			}
		})
	}
}

func (w *NativeWindow) setCursor(cursor *Cursor) {
	//if cursor != nil {
	//	w.PrepareCursor(cursor)
	//}

	unlock := w.newAttemptUnlocker()
	w.loadedCursor = nil
	unlock()

	dispatch(func() {
		w.doSetCursor()
	})
}

func (w *NativeWindow) setCursorPosition(x, y int) {
	dispatch(func() {
		w.doSetCursorPos()
	})
}

func (w *NativeWindow) setCursorGrabbed(grabbed bool) {
	dispatch(func() {
		if grabbed {
			w.saveCursorClip()
			w.updateCursorClip()
		} else {
			w.restoreCursorClip()
		}

		w.doSetCursor()
		w.doSetCursorPos()
	})
}

// HWND returns the win32 handle to this Window, and it's child render window HWND.
//
// This is only useful when doing an few very select hack-ish things.
func (w *NativeWindow) HWND() win32.HWND {
	w.access.RLock()
	defer w.access.RUnlock()
	return w.hwnd
}

// Class returns the window class string of this Window (lpClassName), this is of course, Windows
// specific, and is only useful doing an small, select amount of things.
func (w *NativeWindow) Class() string {
	w.access.RLock()
	defer w.access.RUnlock()
	return w.windowClass
}

func (w *NativeWindow) newAttemptUnlocker() (unlock func()) {
	w.access.Lock()
	unlocked := false
	return func() {
		if !unlocked {
			unlocked = true
			w.access.Unlock()
		}
	}
}

func (w *NativeWindow) doMakeIcon() {
	///////////////////
	// Standard icon //
	///////////////////
	if w.hIcon != nil {
		if !win32.DestroyIcon(w.hIcon) {
			logger().Println("Failed to destroy icon; DestroyIcon():", win32.GetLastErrorString())
		}

		if !win32.DeleteObject(win32.HGDIOBJ(w.iconColorBitmap)) {
			logger().Println("Failed to destroy icon; DeleteObject(iconColorBitmap) failed!")
		}

		if !win32.DeleteObject(win32.HGDIOBJ(w.iconMaskBitmap)) {
			logger().Println("Failed to destroy icon; DeleteObject(iconMaskBitmap) failed!")
		}
	}

	iconWidth := win32.GetSystemMetrics(win32.SM_CXICON)
	iconHeight := win32.GetSystemMetrics(win32.SM_CYICON)

	icon := w.r.Icon()
	if icon == nil {
		return
	}
	iconImage := resize.Resize(icon, icon.Bounds(), int(iconWidth), int(iconHeight))

	iconBitmapInfo := win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			Size:          win32.DWORD(unsafe.Sizeof(win32.BITMAPINFOHEADER{})),
			Width:         win32.LONG(iconWidth),
			Height:        win32.LONG(iconHeight),
			Planes:        1,
			BitCount:      32,
			Compression:   win32.BI_RGB,
			SizeImage:     0,
			XPelsPerMeter: 0,
			YPelsPerMeter: 0,
			ClrUsed:       0,
			ClrImportant:  0,
		},
	}

	w.iconColorBitmap = win32.CreateCompatibleBitmap(w.dc, iconWidth, iconHeight)
	w.iconColorBits = make([]uint32, iconWidth*iconHeight)
	for y := 0; y < int(iconHeight); y++ {
		for x := 0; x < int(iconWidth); x++ {
			r, g, b, _ := iconImage.At(x, y).RGBA()
			c := (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | (uint32(b>>8) << 0)

			index := (int(iconHeight) - 1 - y) * int(iconWidth)
			index += x
			w.iconColorBits[index] = c //0xFF0000
		}
	}
	if win32.SetDIBits(w.dc, w.iconColorBitmap, 0, win32.UINT(iconHeight), unsafe.Pointer(&w.iconColorBits[0]), &iconBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set icon; SetDiBits():", win32.GetLastErrorString())
		return
	}

	w.iconMaskBitmap = win32.CreateCompatibleBitmap(w.dc, iconWidth, iconHeight)
	w.iconMaskBits = make([]uint32, iconWidth*iconHeight)
	for y := 0; y < int(iconHeight); y++ {
		for x := 0; x < int(iconWidth); x++ {
			_, _, _, a := iconImage.At(x, y).RGBA()
			c := uint32(0xFFFFFF)
			if a > 0 {
				c = 0
			}

			index := (int(iconHeight) - 1 - y) * int(iconWidth)
			index += x
			w.iconMaskBits[index] = c
		}
	}
	if win32.SetDIBits(w.dc, w.iconMaskBitmap, 0, win32.UINT(iconHeight), unsafe.Pointer(&w.iconMaskBits[0]), &iconBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set icon; SetDiBits():", win32.GetLastErrorString())
		return
	}

	iconInfo := win32.ICONINFO{
		FIcon:    1,
		XHotspot: 0,
		YHotspot: 0,
		HbmMask:  w.iconMaskBitmap,
		HbmColor: w.iconColorBitmap,
	}

	w.hIcon = win32.CreateIconIndirect(&iconInfo)
	if w.hIcon == nil {
		logger().Println("Unable to set icon; CreateIconIndirect():", win32.GetLastErrorString())
		return
	}

	////////////////
	// Small icon //
	////////////////
	if w.hSmIcon != nil {
		if !win32.DestroyIcon(w.hSmIcon) {
			logger().Println("Failed to destroy icon; DestroyIcon():", win32.GetLastErrorString())
		}

		if !win32.DeleteObject(win32.HGDIOBJ(w.smIconColorBitmap)) {
			logger().Println("Failed to destroy icon; DeleteObject(smIconColorBitmap) failed!")
		}

		if !win32.DeleteObject(win32.HGDIOBJ(w.smIconMaskBitmap)) {
			logger().Println("Failed to destroy icon; DeleteObject(smIconMaskBitmap) failed!")
		}
	}

	iconWidth = win32.GetSystemMetrics(win32.SM_CXSMICON)
	iconHeight = win32.GetSystemMetrics(win32.SM_CYSMICON)

	iconImage = resize.Resize(icon, icon.Bounds(), int(iconWidth), int(iconHeight))

	iconBitmapInfo = win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			Size:          win32.DWORD(unsafe.Sizeof(win32.BITMAPINFOHEADER{})),
			Width:         win32.LONG(iconWidth),
			Height:        win32.LONG(iconHeight),
			Planes:        1,
			BitCount:      32,
			Compression:   win32.BI_RGB,
			SizeImage:     0,
			XPelsPerMeter: 0,
			YPelsPerMeter: 0,
			ClrUsed:       0,
			ClrImportant:  0,
		},
	}

	w.smIconColorBitmap = win32.CreateCompatibleBitmap(w.dc, iconWidth, iconHeight)
	w.smIconColorBits = make([]uint32, iconWidth*iconHeight)
	for y := 0; y < int(iconHeight); y++ {
		for x := 0; x < int(iconWidth); x++ {
			r, g, b, _ := iconImage.At(x, y).RGBA()
			c := (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | (uint32(b>>8) << 0)

			index := (int(iconHeight) - 1 - y) * int(iconWidth)
			index += x
			w.smIconColorBits[index] = c //0xFF0000
		}
	}
	if win32.SetDIBits(w.dc, w.smIconColorBitmap, 0, win32.UINT(iconHeight), unsafe.Pointer(&w.smIconColorBits[0]), &iconBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set icon; SetDiBits():", win32.GetLastErrorString())
		return
	}

	w.smIconMaskBitmap = win32.CreateCompatibleBitmap(w.dc, iconWidth, iconHeight)
	w.smIconMaskBits = make([]uint32, iconWidth*iconHeight)
	for y := 0; y < int(iconHeight); y++ {
		for x := 0; x < int(iconWidth); x++ {
			_, _, _, a := iconImage.At(x, y).RGBA()
			c := uint32(0xFFFFFF)
			if a > 0 {
				c = 0
			}

			index := (int(iconHeight) - 1 - y) * int(iconWidth)
			index += x
			w.smIconMaskBits[index] = c
		}
	}
	if win32.SetDIBits(w.dc, w.smIconMaskBitmap, 0, win32.UINT(iconHeight), unsafe.Pointer(&w.smIconMaskBits[0]), &iconBitmapInfo, win32.DIB_RGB_COLORS) == 0 {
		logger().Println("Unable to set icon; SetDiBits():", win32.GetLastErrorString())
		return
	}

	iconInfo = win32.ICONINFO{
		FIcon:    1,
		XHotspot: 0,
		YHotspot: 0,
		HbmMask:  w.smIconMaskBitmap,
		HbmColor: w.smIconColorBitmap,
	}

	w.hSmIcon = win32.CreateIconIndirect(&iconInfo)
	if w.hSmIcon == nil {
		logger().Println("Unable to set icon; CreateIconIndirect():", win32.GetLastErrorString())
		return
	}
	w.doSetIcon()
}

func (w *NativeWindow) doSetIcon() {
	win32.SendMessage(w.hwnd, win32.WM_SETICON, win32.ICON_BIG, win32.LPARAM(uintptr(unsafe.Pointer(w.hIcon))))
	win32.SendMessage(w.hwnd, win32.WM_SETICON, win32.ICON_SMALL, win32.LPARAM(uintptr(unsafe.Pointer(w.hSmIcon))))
}

func (w *NativeWindow) doSetCursor() {
	if w.r.CursorWithin() {
		if w.r.CursorGrabbed() {
			win32.SetCursor(nil)
			return
		}

		if w.loadedCursor == nil {
			cursor := w.r.Cursor()
			if cursor != nil {
				// Cursor is just not prepared yet.
				w.doPrepareCursor(cursor)

				lc, ok := w.cursors[cursor]
				if ok {
					w.loadedCursor = lc
				}

			} else {
				// There is no cursor.
				lc := new(loadedCursor)
				lc.hCursor = win32.HICON(win32.LoadCursor(nil, "IDC_ARROW"))
				if lc.hCursor == nil {
					logger().Println("Unable to load default (IDC_ARROW) cursor! LoadCursor():", win32.GetLastErrorString())
				} else {
					w.loadedCursor = lc
				}
			}
		}

		if w.loadedCursor != nil {
			win32.SetCursor(win32.HCURSOR(w.loadedCursor.hCursor))
		}
	}
}

func (w *NativeWindow) doSetCursorPos() {
	if w.r.CursorGrabbed() {
		if !w.r.CursorWithin() {
			return
		}
		width, height := w.r.clampedSize()
		w.r.trySetCursorPosition(width/2, height/2)
	}

	x, y := w.r.Position()
	cursorX, cursorY := w.r.CursorPosition()
	if !win32.SetCursorPos(int32(x+cursorX), int32(y+cursorY)) {
		logger().Println("Unable to set cursor position: SetCursorPos():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) doUpdateTransparency() {
	bb := win32.DWM_BLURBEHIND{}
	bb.DwFlags = win32.DWM_BB_ENABLE | win32.DWM_BB_BLURREGION
	if w.r.Transparent() {
		bb.FEnable = 1
	} else {
		bb.FEnable = 0
	}
	rgn := win32.CreateRectRgn(0, 0, -1, -1)
	bb.HRgbBlur = rgn
	err := win32.DwmEnableBlurBehindWindow(w.hwnd, &bb)
	if err != nil {
		logger().Println(err)
	}
}

func (w *NativeWindow) doUpdateStyle() {
	originalStyle := win32.GetWindowLongPtr(w.hwnd, win32.GWL_STYLE)

	if w.r.Decorated() && !w.r.Fullscreen() {
		w.styleFlags = win32.WS_OVERLAPPEDWINDOW
	} else {
		w.styleFlags = win32.WS_SYSMENU | win32.WS_POPUP | win32.WS_CLIPCHILDREN | win32.WS_CLIPSIBLINGS
	}

	if w.r.Visible() {
		w.styleFlags |= win32.WS_VISIBLE
	}

	if w.r.Opened() {
		if win32.DWORD(originalStyle) != w.styleFlags {
			win32.SetWindowLongPtr(w.hwnd, win32.GWL_STYLE, win32.LONG_PTR(w.styleFlags))

			//win32.EnableWindow(w.hwnd, true)
			//if w.visible {
			//    win32.ShowWindow(w.hwnd, win32.SW_SHOWNA)
			//}
			w.doSetWindowPos()
		}
	}
}

func (w *NativeWindow) doSetWindowPos() {
	// win32.SWP_ASYNCWINDOWPOS|win32.SWP_FRAMECHANGED|win32.SWP_NOMOVE|win32.SWP_NOSIZE|win32.SWP_NOZORDER|win32.SWP_NOOWNERZORDER

	extentLeft := w.r.extentLeft
	extentRight := w.r.extentRight
	extentBottom := w.r.extentBottom
	extentTop := w.r.extentTop
	wX, wY := w.r.Position()
	x := win32.Int(wX - int(extentLeft))
	y := win32.Int(wY - int(extentTop))

	// Need to make the position relative to the original screen
	r := w.r.OriginalScreen().NativeScreen.w32Position
	x = win32.Int(r.Left) + x
	y = win32.Int(r.Top) + y

	if !w.r.Decorated() {
		x += win32.Int(extentLeft)
		y += win32.Int(extentTop)
	}

	wWidth, wHeight := w.r.clampedSize()
	width := float64(wWidth)
	height := float64(wHeight)

	if w.r.Decorated() {
		width += float64(extentLeft)
		width += float64(extentRight)
		height += float64(extentBottom)
		height += float64(extentTop)
	}

	ratio := w.r.AspectRatio()
	if ratio != 0.0 {
		if ratio > 1.0 {
			// Wider instead of taller
			width = float64(ratio * float32(height))
			width = float64(ratio * float32(height))
		} else {
			// Taller instead of wider
			height = float64((1.0 / ratio) * float32(width))
			height = float64((1.0 / ratio) * float32(width))
		}
	}

	insertAfter := win32.HWND_NOTOPMOST
	if w.r.AlwaysOnTop() {
		insertAfter = win32.HWND_TOPMOST
	}
	//		win32.SetWindowPos(w.hwnd, flag, 0, 0, 0, 0, win32.SWP_NOMOVE|win32.SWP_NOSIZE)

	if w.r.Fullscreen() {
		screen := w.r.Screen()
		r := screen.NativeScreen.w32Position
		x = win32.Int(r.Left)
		y = win32.Int(r.Top)
		sm := screen.Mode()
		screenWidth, screenHeight := sm.Resolution()
		width, height = float64(screenWidth), float64(screenHeight)

		wWidth, wHeight := w.r.Size()
		if wWidth != int(screenWidth) || wHeight != int(screenHeight) {
			w.r.trySetSize(int(screenWidth), int(screenHeight))
		}
	}

	// |win32.SWP_NOZORDER|win32.SWP_NOOWNERZORDER
	if !win32.SetWindowPos(w.hwnd, insertAfter, x, y, win32.Int(width), win32.Int(height), win32.SWP_FRAMECHANGED) {
		//logger().Println("Unable to set window position; SetWindowPos():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) translateKey(wParam win32.WPARAM) (key keyboard.Key) {
	switch wParam {
	// ? "Control-break processing" doesn't seem useful
	//case win32.VK_CANCEL:

	case win32.VK_BACK:
		key = keyboard.Backspace
	case win32.VK_TAB:
		key = keyboard.Tab
	case win32.VK_CLEAR:
		key = keyboard.Clear
	case win32.VK_RETURN:
		key = keyboard.Enter
	case win32.VK_PAUSE:
		key = keyboard.Pause
	case win32.VK_KANA:
		key = keyboard.Kana
	case win32.VK_JUNJA:
		key = keyboard.Junja
	case win32.VK_KANJI:
		key = keyboard.Kanji
	case win32.VK_ESCAPE:
		key = keyboard.Escape
	case win32.VK_SPACE:
		key = keyboard.Space
	case win32.VK_PRIOR:
		key = keyboard.PageUp
	case win32.VK_NEXT:
		key = keyboard.PageDown
	case win32.VK_END:
		key = keyboard.End
	case win32.VK_HOME:
		key = keyboard.Home
	case win32.VK_LEFT:
		key = keyboard.ArrowLeft
	case win32.VK_UP:
		key = keyboard.ArrowUp
	case win32.VK_RIGHT:
		key = keyboard.ArrowRight
	case win32.VK_DOWN:
		key = keyboard.ArrowDown
	case win32.VK_SELECT:
		key = keyboard.Select
	case win32.VK_PRINT:
		key = keyboard.Print
	case win32.VK_EXECUTE:
		key = keyboard.Execute
	case win32.VK_SNAPSHOT:
		key = keyboard.PrintScreen
	case win32.VK_INSERT:
		key = keyboard.Insert
	case win32.VK_DELETE:
		key = keyboard.Delete
	case win32.VK_HELP:
		key = keyboard.Help
	case win32.VK_UNDEF_0:
		key = keyboard.Zero
	case win32.VK_UNDEF_1:
		key = keyboard.One
	case win32.VK_UNDEF_2:
		key = keyboard.Two
	case win32.VK_UNDEF_3:
		key = keyboard.Three
	case win32.VK_UNDEF_4:
		key = keyboard.Four
	case win32.VK_UNDEF_5:
		key = keyboard.Five
	case win32.VK_UNDEF_6:
		key = keyboard.Six
	case win32.VK_UNDEF_7:
		key = keyboard.Seven
	case win32.VK_UNDEF_8:
		key = keyboard.Eight
	case win32.VK_UNDEF_9:
		key = keyboard.Nine
	case win32.VK_UNDEF_A:
		key = keyboard.A
	case win32.VK_UNDEF_B:
		key = keyboard.B
	case win32.VK_UNDEF_C:
		key = keyboard.C
	case win32.VK_UNDEF_D:
		key = keyboard.D
	case win32.VK_UNDEF_E:
		key = keyboard.E
	case win32.VK_UNDEF_F:
		key = keyboard.F
	case win32.VK_UNDEF_G:
		key = keyboard.G
	case win32.VK_UNDEF_H:
		key = keyboard.H
	case win32.VK_UNDEF_I:
		key = keyboard.I
	case win32.VK_UNDEF_J:
		key = keyboard.J
	case win32.VK_UNDEF_K:
		key = keyboard.K
	case win32.VK_UNDEF_L:
		key = keyboard.L
	case win32.VK_UNDEF_M:
		key = keyboard.M
	case win32.VK_UNDEF_N:
		key = keyboard.N
	case win32.VK_UNDEF_O:
		key = keyboard.O
	case win32.VK_UNDEF_P:
		key = keyboard.P

	case win32.VK_UNDEF_Q:
		key = keyboard.Q
	case win32.VK_UNDEF_R:
		key = keyboard.R
	case win32.VK_UNDEF_S:
		key = keyboard.S
	case win32.VK_UNDEF_T:
		key = keyboard.T
	case win32.VK_UNDEF_U:
		key = keyboard.U
	case win32.VK_UNDEF_V:
		key = keyboard.V
	case win32.VK_UNDEF_W:
		key = keyboard.W
	case win32.VK_UNDEF_X:
		key = keyboard.X
	case win32.VK_UNDEF_Y:
		key = keyboard.Y
	case win32.VK_UNDEF_Z:
		key = keyboard.Z

	//case win32.VK_CONVERT: key = keyboard.IMEConvert
	//case win32.VK_NONCONVERT: key = keyboard.IMENonConvert
	//case win32.VK_ACCEPT: key = keyboard.IMEAccept
	//case win32.VK_MODECHANGE: key = keyboard.IMEModeChange
	//case win32.VK_PROCESSKEY: key = keyboard.IMEProcess

	case win32.VK_LWIN:
		key = keyboard.LeftSuper
	case win32.VK_RWIN:
		key = keyboard.RightSuper
	case win32.VK_APPS:
		key = keyboard.Applications
	case win32.VK_SLEEP:
		key = keyboard.Sleep

	case win32.VK_NUMPAD0:
		key = keyboard.NumZero
	case win32.VK_NUMPAD1:
		key = keyboard.NumOne
	case win32.VK_NUMPAD2:
		key = keyboard.NumTwo
	case win32.VK_NUMPAD3:
		key = keyboard.NumThree
	case win32.VK_NUMPAD4:
		key = keyboard.NumFour
	case win32.VK_NUMPAD5:
		key = keyboard.NumFive
	case win32.VK_NUMPAD6:
		key = keyboard.NumSix
	case win32.VK_NUMPAD7:
		key = keyboard.NumSeven
	case win32.VK_NUMPAD8:
		key = keyboard.NumEight
	case win32.VK_NUMPAD9:
		key = keyboard.NumNine
	case win32.VK_MULTIPLY:
		key = keyboard.NumMultiply
	case win32.VK_ADD:
		key = keyboard.NumAdd
	case win32.VK_SEPARATOR:
		key = keyboard.NumComma
	case win32.VK_SUBTRACT:
		key = keyboard.NumSubtract
	case win32.VK_DECIMAL:
		key = keyboard.NumDecimal
	case win32.VK_DIVIDE:
		key = keyboard.NumDivide

	case win32.VK_F1:
		key = keyboard.F1
	case win32.VK_F2:
		key = keyboard.F2
	case win32.VK_F3:
		key = keyboard.F3
	case win32.VK_F4:
		key = keyboard.F4
	case win32.VK_F5:
		key = keyboard.F5
	case win32.VK_F6:
		key = keyboard.F6
	case win32.VK_F7:
		key = keyboard.F7
	case win32.VK_F8:
		key = keyboard.F8
	case win32.VK_F9:
		key = keyboard.F9
	case win32.VK_F10:
		key = keyboard.F10
	case win32.VK_F11:
		key = keyboard.F11
	case win32.VK_F12:
		key = keyboard.F12
	case win32.VK_F13:
		key = keyboard.F13
	case win32.VK_F14:
		key = keyboard.F14
	case win32.VK_F15:
		key = keyboard.F15
	case win32.VK_F16:
		key = keyboard.F16
	case win32.VK_F17:
		key = keyboard.F17
	case win32.VK_F18:
		key = keyboard.F18
	case win32.VK_F19:
		key = keyboard.F19
	case win32.VK_F20:
		key = keyboard.F20
	case win32.VK_F21:
		key = keyboard.F21
	case win32.VK_F22:
		key = keyboard.F22
	case win32.VK_F23:
		key = keyboard.F23
	case win32.VK_F24:
		key = keyboard.F24

	case win32.VK_BROWSER_BACK:
		key = keyboard.BrowserBack
	case win32.VK_BROWSER_FORWARD:
		key = keyboard.BrowserForward
	case win32.VK_BROWSER_REFRESH:
		key = keyboard.BrowserRefresh
	case win32.VK_BROWSER_STOP:
		key = keyboard.BrowserStop
	case win32.VK_BROWSER_SEARCH:
		key = keyboard.BrowserSearch
	case win32.VK_BROWSER_FAVORITES:
		key = keyboard.BrowserFavorites
	case win32.VK_BROWSER_HOME:
		key = keyboard.BrowserHome

	// User expects these to control windows volume -- I don't thing we should allow
	// intercepting these..
	//VK_VOLUME_MUTE
	//VK_VOLUME_DOWN
	//VK_VOLUME_UP

	case win32.VK_MEDIA_NEXT_TRACK:
		key = keyboard.MediaNext
	case win32.VK_MEDIA_PREV_TRACK:
		key = keyboard.MediaPrevious
	case win32.VK_MEDIA_STOP:
		key = keyboard.MediaStop
	case win32.VK_MEDIA_PLAY_PAUSE:
		key = keyboard.MediaPlayPause

	case win32.VK_LAUNCH_MAIL:
		key = keyboard.LaunchMail
	case win32.VK_LAUNCH_MEDIA_SELECT:
		key = keyboard.LaunchMedia
	case win32.VK_LAUNCH_APP1:
		key = keyboard.LaunchAppOne
	case win32.VK_LAUNCH_APP2:
		key = keyboard.LaunchAppTwo

	case win32.VK_OEM_PLUS:
		key = keyboard.Equals
	case win32.VK_OEM_COMMA:
		key = keyboard.Comma
	case win32.VK_OEM_MINUS:
		key = keyboard.Dash
	case win32.VK_OEM_PERIOD:
		key = keyboard.Period
	case win32.VK_OEM_1:
		key = keyboard.Semicolon
	case win32.VK_OEM_2:
		key = keyboard.ForwardSlash
	case win32.VK_OEM_3:
		key = keyboard.Tilde
	case win32.VK_OEM_4:
		key = keyboard.LeftBracket
	case win32.VK_OEM_5:
		key = keyboard.BackSlash
	case win32.VK_OEM_6:
		key = keyboard.RightBracket
	case win32.VK_OEM_7:
		key = keyboard.Apostrophe
	//case win32.VK_OEM_8:
	case win32.VK_OEM_102:
		key = keyboard.RightBracket

	case win32.VK_ATTN:
		key = keyboard.Attn
	case win32.VK_CRSEL:
		key = keyboard.CrSel
	case win32.VK_EXSEL:
		key = keyboard.ExSel
	case win32.VK_EREOF:
		key = keyboard.EraseEOF
	case win32.VK_PLAY:
		key = keyboard.Play
	case win32.VK_ZOOM:
		key = keyboard.Zoom
	//case win32.VK_PA1: key = keyboard.PA1
	case win32.VK_OEM_CLEAR:
		key = keyboard.Clear

	case win32.VK_SHIFT:
		key = keyboard.LeftShift
	case win32.VK_MENU:
		key = keyboard.LeftAlt
	case win32.VK_CONTROL:
		key = keyboard.LeftCtrl

	case win32.VK_CAPITAL:
		key = keyboard.CapsLock
	case win32.VK_NUMLOCK:
		key = keyboard.NumLock
	case win32.VK_SCROLL:
		key = keyboard.ScrollLock
	}
	return key
}

func (w *NativeWindow) saveCursorClip() {
	if w.lastCursorClip != nil {
		return
	}

	var ok bool
	w.lastCursorClip, ok = win32.GetClipCursor()
	if !ok {
		logger().Println("Unable to set clip cursor; GetClipCursor():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) restoreCursorClip() {
	if w.lastCursorClip == nil {
		return
	}

	win32.ClipCursor(w.lastCursorClip)

	// Clear it so we don't accidently restore it again later
	w.lastCursorClip = nil
}

func (w *NativeWindow) updateCursorClip() {
	// We don't use cursor clip if we have none to restore to (due to an failure).
	if w.lastCursorClip == nil {
		return
	}

	tl := new(win32.POINT)
	tl.X = 0
	tl.Y = 0
	if !win32.ClientToScreen(w.hwnd, tl) {
		logger().Println("Unable to set clip cursor; ClientToScreen():", win32.GetLastErrorString())
	}

	wWidth, wHeight := w.r.clampedSize()
	br := new(win32.POINT)
	br.X = int32(wWidth)
	br.Y = int32(wHeight)
	if !win32.ClientToScreen(w.hwnd, br) {
		logger().Println("Unable to set clip cursor; ClientToScreen():", win32.GetLastErrorString())
	}

	clip := &win32.RECT{
		Left:   int32(tl.X),
		Top:    int32(tl.Y),
		Right:  int32(br.X),
		Bottom: int32(br.Y),
	}

	if !win32.ClipCursor(clip) {
		logger().Println("Unable to set clip cursor; ClipCursor():", win32.GetLastErrorString())
	}
}

// Our MS windows event handler
//
// This is never executed under the pretence of an window's respective lock.
func mainWindowProc(hwnd win32.HWND, msg win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) (ret win32.LRESULT) {
	// Q: Why is this line here?
	// A: Windows API calls invoke this mainWindowProc, for instance an call to FooEx() might
	// invoke mainWindowProc(), and when it does the function, FooEx(), will block until each and
	// every message has been handled by mainWindowProc(). This line stops the message pump from
	// pausing other goroutines.
	runtime.Gosched()

	w, ok := windowsByHwnd[hwnd]
	if ok {
		var unlock func()
		if !w.insideCallback {
			unlock = w.newAttemptUnlocker()
			defer unlock()
			w.insideCallback = true
			defer func() {
				w.insideCallback = false
			}()
		} else {
			unlock = func() {
			}
		}

		switch {
		case msg == win32.WM_PAINT:
			rect := new(win32.RECT)
			if win32.GetUpdateRect(w.hwnd, rect, false) {
				win32.ValidateRect(w.hwnd, nil)

				x0 := int(rect.Left)
				y0 := int(rect.Top)
				x1 := int(rect.Right)
				y1 := int(rect.Bottom)

				if x0 <= x1 && y0 <= y1 {
					w.r.send(PaintEvent{
						T:         time.Now(),
						Rectangle: image.Rect(x0, y0, x1, y1),
					})
				} else {
					logger().Println("WARNING: Got non well-formed WM_PAINT rectangle, ignored!")
				}
			}
			return 0

		case msg == win32.WM_ERASEBKGND:
			//w.addPaintEvent()
			return 1

		case msg == win32.WM_GETMINMAXINFO:
			ratio := w.r.AspectRatio()

			// Add extents, so we operate on client region space only
			extentLeft, extentRight, extentBottom, extentTop := w.r.Extents()
			minWidth, minHeight := w.r.MinimumSize()
			maxWidth, maxHeight := w.r.MaximumSize()

			newMinWidth := minWidth + extentLeft + extentRight
			newMaxWidth := maxWidth + extentLeft + extentRight
			newMinHeight := minHeight + extentBottom + extentTop
			newMaxHeight := maxHeight + extentBottom + extentTop

			if ratio != 0.0 {
				if ratio > 1.0 {
					// Wider instead of taller
					newMinWidth = int(ratio * float32(newMinHeight))
					newMaxWidth = int(ratio * float32(newMaxHeight))
				} else {
					// Taller instead of wider
					newMinHeight = int((1.0 / ratio) * float32(newMinWidth))
					newMaxHeight = int((1.0 / ratio) * float32(newMaxWidth))
				}
			}

			// Set maximum and minimum window sizes, 0 means unlimited
			minMaxInfo := lParam.MINMAXINFO()

			if minWidth > 0 {
				minMaxInfo.PtMinTrackSize.X = int32(newMinWidth)
			}
			if minHeight > 0 {
				minMaxInfo.PtMinTrackSize.Y = int32(newMinHeight)
			}

			if maxWidth > 0 {
				minMaxInfo.PtMaxTrackSize.X = int32(newMaxWidth)
			}
			if maxHeight > 0 {
				minMaxInfo.PtMaxTrackSize.Y = int32(newMaxHeight)
			}
			return 0

		case msg == win32.WM_SIZING:
			ratio := w.r.AspectRatio()
			r := lParam.RECT()

			if ratio != 0 {
				width := r.Right - r.Left
				height := r.Bottom - r.Top

				newHeight := (1.0 / ratio) * float32(width)
				newWidth := ratio * float32(height)

				newRight := r.Left + int32(newWidth)
				//newLeft := r.Right - int32(newWidth)
				newBottom := r.Top + int32(newHeight)
				newTop := r.Bottom - int32(newHeight)

				if wParam == win32.WMSZ_RIGHT || wParam == win32.WMSZ_LEFT {
					r.Bottom = newBottom
				} else if wParam == win32.WMSZ_BOTTOM || wParam == win32.WMSZ_TOP {
					r.Right = newRight

				} else if wParam == win32.WMSZ_TOPLEFT || wParam == win32.WMSZ_TOPRIGHT {
					r.Top = newTop
				} else if wParam == win32.WMSZ_BOTTOMLEFT || wParam == win32.WMSZ_BOTTOMRIGHT {
					r.Bottom = newBottom
				}

				w.lastWmSizingLeft = r.Left
				w.lastWmSizingRight = r.Right
				w.lastWmSizingBottom = r.Bottom
				w.lastWmSizingTop = r.Top
			}

			newWidth := int(w.lastWmSizingRight - w.lastWmSizingLeft)
			newHeight := int(w.lastWmSizingBottom - w.lastWmSizingTop)

			if newWidth != 0 && newHeight != 0 && w.r.trySetSize(newWidth, newHeight) {
				if w.r.CursorGrabbed() {
					// Update our clip
					w.updateCursorClip()
				}
			}
			return 0

		case msg == win32.WM_SIZE:
			if wParam == win32.SIZE_MAXIMIZED {
				w.r.trySetMinimized(false)
				w.r.trySetMaximized(true)
			} else if wParam == win32.SIZE_MINIMIZED {
				w.r.trySetMinimized(true)
				w.r.trySetMaximized(false)
			} else {
				w.r.trySetMinimized(false)
				w.r.trySetMaximized(false)
			}

			if wParam != win32.SIZE_MINIMIZED {
				newWidth := int(lParam.LOWORD())
				newHeight := int(lParam.HIWORD())
				w.r.trySetSize(newWidth, newHeight)
			}
			return 0

		case msg == win32.WM_MOVE:
			xPos := int(int16(lParam))
			yPos := int(int16((uint32(lParam) >> 16) & 0xFFFF))

			if !win32.IsIconic(w.hwnd) {
				// Clamp when it goes onto an monitor to the left.. very unsure how to handle
				// window/screen interaction -- it's never very cross platform..

				w.r.trySetPosition(xPos, yPos)
			}

			return 0

		case msg == win32.WM_EXITSIZEMOVE:
			hMonitor := win32.MonitorFromWindow(w.hwnd, win32.MONITOR_DEFAULTTONEAREST)

			mi := new(win32.MONITORINFOEX)
			if !win32.GetMonitorInfo(hMonitor, mi) {
				logger().Println("Unable to detect monitor position; GetMonitorInfo():", win32.GetLastErrorString())
			} else {
				screens := backend_doScreens()
				for _, screen := range screens {
					if screen.NativeScreen.w32GraphicsDeviceName == win32.String(mi.SzDevice[:]) {
						if w.r.trySetScreen(screen) {
							screen.NativeScreen.w32Position = mi.RcMonitor
						}
						return 0
					}
				}
			}

			return 0

		case msg == win32.WM_ACTIVATE:
			if wParam.LOWORD() == win32.WA_INACTIVE || wParam.HIWORD() != 0 {
				if w.r.trySetFocused(false) {
					// If the window loses focus due to another window coming into the foreground,
					// like ctrl+alt+delete on 7+, or alt+tab to another application, then we never
					// receive mouse exit event, also we don't know to release mouse grab then.

					if w.r.trySetCursorWithin(false) {
						// Restore previous cursor clip
						w.restoreCursorClip()

						win32.ReleaseCapture()
					}

					// Also if keys, mouse buttons, etc were being held down
					// while the window did have focus, as it no longer does,
					// we need to release those keys or else the user can only
					// release them by putting them into down-state and then
					// releasing again!
					w.r.releaseDownedButtons()

					if w.r.Fullscreen() {
						// If the window is fullscreen and loses focus users will
						// expect it to minimize on it's own instead of hanging on
						// in the background.
						win32.ShowWindow(w.hwnd, win32.SW_MINIMIZE)
						w.doSetWindowPos()
						w.r.trySetMaximized(false)
						w.r.trySetMinimized(true)
					}
				}
			} else {
				w.r.trySetFocused(true)
			}
			return 0

		case msg == win32.WM_GETICON:
			switch wParam {
			case win32.ICON_BIG:
				if w.hIcon != nil {
					return win32.LRESULT(uintptr(unsafe.Pointer(w.hIcon)))
				}

			case win32.ICON_SMALL:
				if w.hSmIcon != nil {
					return win32.LRESULT(uintptr(unsafe.Pointer(w.hSmIcon)))
				}

			case win32.ICON_SMALL2:
				if w.hSmIcon != nil {
					return win32.LRESULT(uintptr(unsafe.Pointer(w.hSmIcon)))
				}
			}

		case msg == win32.WM_CHAR:
			w.r.send(keyboard.TypedEvent{
				T:    time.Now(),
				Rune: rune(wParam),
			})

		case msg == win32.WM_KEYDOWN || msg == win32.WM_SYSKEYDOWN || msg == win32.WM_KEYUP || msg == win32.WM_SYSKEYUP:
			if msg == win32.WM_KEYDOWN || msg == win32.WM_SYSKEYDOWN {
				keyRepeat := (lParam & 0x40000000) > 0
				if keyRepeat {
					return 0
				}
			}

			if (msg == win32.WM_SYSKEYDOWN || msg == win32.WM_SYSKEYUP) && wParam == win32.VK_F4 {
				altDown := (uint16(win32.GetAsyncKeyState(win32.VK_MENU)) & 0x8000) != 0
				if altDown {
					if msg == win32.WM_SYSKEYDOWN {
						// Trick: Consider this to be WM_CLOSE
						w.r.send(CloseEvent{
							T: time.Now(),
						})
					}
					return 0
				}
			}

			k := w.translateKey(wParam)

			var state keyboard.State
			if msg == win32.WM_KEYDOWN || msg == win32.WM_SYSKEYDOWN {
				state = keyboard.Down
			} else {
				state = keyboard.Up
			}

			switch k {
			case keyboard.LeftShift:
				leftShiftDown := (uint16(win32.GetAsyncKeyState(win32.VK_LSHIFT)) & 0x8000) != 0
				state = keyboard.Down
				if !leftShiftDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.LeftShift, uint64(win32.VK_LSHIFT), state)

				rightShiftDown := (uint16(win32.GetAsyncKeyState(win32.VK_RSHIFT)) & 0x8000) != 0
				state = keyboard.Down
				if !rightShiftDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.RightShift, uint64(win32.VK_RSHIFT), state)
				return 0

			case keyboard.LeftAlt:
				leftAltDown := (uint16(win32.GetAsyncKeyState(win32.VK_LMENU)) & 0x8000) != 0
				state = keyboard.Down
				if !leftAltDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.LeftAlt, uint64(win32.VK_LMENU), state)

				rightAltDown := (uint16(win32.GetAsyncKeyState(win32.VK_RMENU)) & 0x8000) != 0
				state = keyboard.Down
				if !rightAltDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.RightAlt, uint64(win32.VK_RMENU), state)
				return 0

			case keyboard.LeftCtrl:
				leftCtrlDown := (uint16(win32.GetAsyncKeyState(win32.VK_LCONTROL)) & 0x8000) != 0
				state = keyboard.Down
				if !leftCtrlDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.LeftCtrl, uint64(win32.VK_LCONTROL), state)

				rightCtrlDown := (uint16(win32.GetAsyncKeyState(win32.VK_RCONTROL)) & 0x8000) != 0
				state = keyboard.Down
				if !rightCtrlDown {
					state = keyboard.Up
				}
				w.r.tryAddKeyboardStateEvent(keyboard.RightCtrl, uint64(win32.VK_RCONTROL), state)
				return 0

			case keyboard.CapsLock:
				if (win32.GetKeyState(win32.VK_CAPITAL) & 0x0001) != 0 {
					state = keyboard.Down
				} else {
					state = keyboard.Up
				}

			case keyboard.NumLock:
				if (win32.GetKeyState(win32.VK_NUMLOCK) & 0x0001) != 0 {
					state = keyboard.Down
				} else {
					state = keyboard.Up
				}

			case keyboard.ScrollLock:
				if (win32.GetKeyState(win32.VK_SCROLL) & 0x0001) != 0 {
					state = keyboard.Down
				} else {
					state = keyboard.Up
				}
			}

			w.r.tryAddKeyboardStateEvent(k, uint64(wParam), state)
			return 0

		case msg == win32.WM_MOUSEMOVE:
			cursorX := int(int16(lParam))
			cursorY := int(int16((uint32(lParam) >> 16) & 0xFFFF))

			w.r.trySetCursorPosition(cursorX, cursorY)

			wWidth, wHeight := w.r.clampedSize()
			if cursorX >= wWidth || cursorY >= wHeight || cursorX <= 0 || cursorY <= 0 || !w.r.Focused() {
				// Better than WM_MOUSELEAVE
				if !w.r.CursorGrabbed() {
					if w.r.trySetCursorWithin(false) {
						// Restore previous cursor clip
						w.restoreCursorClip()

						win32.ReleaseCapture()
					}
				}
			} else {
				// Closest we'll get to WM_MOUSEENTER
				if w.r.trySetCursorWithin(true) {
					if w.r.CursorGrabbed() {
						// Store previous clipping
						w.saveCursorClip()

						// Update our clip
						w.updateCursorClip()
					}

					win32.SetCapture(w.hwnd)
				}
			}

			if w.r.CursorGrabbed() {
				supportRawInput := w32VersionMajor >= 5 && w32VersionMinor >= 1
				wWidth, wHeight := w.r.clampedSize()
				halfWidth := wWidth / 2
				halfHeight := wHeight / 2

				if w.r.preGrabCursorX == 0 && w.r.preGrabCursorY == 0 {
					w.r.preGrabCursorX = cursorX
					w.r.preGrabCursorY = cursorY
				}

				if cursorX != halfWidth || cursorY != halfHeight {
					if !supportRawInput {
						// If we have no support for raw mouse input, then we need to fall back to finding
						// the mouse movement on our own.
						diffX := cursorX - halfWidth
						diffY := cursorY - halfHeight
						if diffX != 0 || diffY != 0 {
							w.r.send(CursorPositionEvent{
								T: time.Now(),
								X: float64(diffX),
								Y: float64(diffY),
							})
						}
					}
					w.doSetCursorPos()
				}
			}
			w.doSetCursor()
			return 0

		case msg == win32.WM_INPUT:
			if w.r.CursorWithin() && w.r.CursorGrabbed() {
				var raw win32.RAWINPUT
				cbSize := win32.UINT(unsafe.Sizeof(raw))

				win32.GetRawInputData((win32.HRAWINPUT)(unsafe.Pointer(uintptr(lParam))), win32.RID_INPUT, unsafe.Pointer(&raw), &cbSize, win32.UINT(unsafe.Sizeof(win32.RAWINPUTHEADER{})))

				if raw.Header.DwType == win32.RIM_TYPEMOUSE {
					diffX := raw.Mouse().LLastX
					diffY := raw.Mouse().LLastY
					if diffX != 0 || diffY != 0 {
						w.r.send(CursorPositionEvent{
							T: time.Now(),
							X: float64(diffX),
							Y: float64(diffY),
						})
					}
				}
			}
			return 0

		// Mouse Buttons
		case msg == win32.WM_LBUTTONDOWN:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Left,
				State:  mouse.Down,
			})
			return 0

		case msg == win32.WM_LBUTTONUP:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Left,
				State:  mouse.Up,
			})
			return 0

		case msg == win32.WM_RBUTTONDOWN:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Right,
				State:  mouse.Down,
			})
			return 0

		case msg == win32.WM_RBUTTONUP:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Right,
				State:  mouse.Up,
			})
			return 0

		case msg == win32.WM_MBUTTONDOWN:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Wheel,
				State:  mouse.Down,
			})
			return 0

		case msg == win32.WM_MBUTTONUP:
			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: mouse.Wheel,
				State:  mouse.Up,
			})
			return 0

		case msg == win32.WM_XBUTTONDOWN:
			var button mouse.Button

			switch int16(wParam) {
			case win32.MK_XBUTTON1:
				button = mouse.Four

			case win32.MK_XBUTTON2:
				button = mouse.Five
			}

			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: button,
				State:  mouse.Down,
			})
			return 0

		case msg == win32.WM_XBUTTONUP:
			var button mouse.Button

			switch int16(wParam) {
			case win32.MK_XBUTTON1:
				button = mouse.Four

			case win32.MK_XBUTTON2:
				button = mouse.Five
			}

			w.r.send(mouse.Event{
				T:      time.Now(),
				Button: button,
				State:  mouse.Up,
			})
			return 0

		case msg == win32.WM_MOUSEWHEEL:
			delta := float64(int16((uint32(wParam) >> 16) & 0xFFFF))
			ticks := int(math.Abs(delta / 120))

			if delta > 0 {
				for i := 0; i < ticks; i++ {
					w.r.send(mouse.Event{
						T:      time.Now(),
						Button: mouse.Wheel,
						State:  mouse.ScrollForward,
					})
				}
			} else {
				for i := 0; i < ticks; i++ {
					w.r.send(mouse.Event{
						T:      time.Now(),
						Button: mouse.Wheel,
						State:  mouse.ScrollBack,
					})
				}
			}
			return 0

		case msg == win32.WM_CLOSE:
			w.r.send(CloseEvent{
				T: time.Now(),
			})
			return 0

		default:
			// We continue onto DefWindowProc(), which might call us again, so make sure to unlock.
			unlock()
		}
	}

	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func newNativeWindow(real *Window) *NativeWindow {
	w := new(NativeWindow)
	w.r = real

	// Get window extents
	titleHeight := win32.GetSystemMetrics(win32.SM_CYCAPTION)
	borderHeight := win32.GetSystemMetrics(win32.SM_CYSIZEFRAME)
	borderWidth := win32.GetSystemMetrics(win32.SM_CXSIZEFRAME)

	w.r.trySetExtents(
		int(borderWidth),
		int(borderWidth),
		int(borderHeight),
		int(borderHeight+titleHeight),
	)
	return w
}
