// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"image"
	"sync"
	"time"

	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
)

// Window represents an single window, it will be non-visible untill the Open
// function is called.
type Window struct {
	*NativeWindow

	Keyboard *keyboard.Watcher
	Mouse    *mouse.Watcher

	access sync.RWMutex

	events map[chan Event]bool

	icon   image.Image
	cursor *Cursor

	opened, destroyed, focused, visible, decorated, minimized, maximized,
	fullscreen, alwaysOnTop, cursorGrabbed, cursorWithin, transparent bool

	extentLeft, extentRight, extentBottom, extentTop, width, height,
	preFullscreenWidth, preFullscreenHeight, minWidth, minHeight, maxWidth,
	maxHeight, x, y, preFullscreenX, preFullscreenY, cursorX, cursorY,
	lastCursorX, lastCursorY, preGrabCursorX, preGrabCursorY int

	aspectRatio float32

	originalScreen, screen *Screen

	title string
}

// String returns an string representation of this window.
func (w *Window) String() string {
	w.access.RLock()
	defer w.access.RUnlock()
	return fmt.Sprintf("Window(Title=%q, Size=%dx%d, Position=%dx%d)", w.title, w.width, w.height, w.x, w.y)
}

// Open opens the window using the current settings, on the specified
// screen, or returns an error in the event that we are unable to open the
// window for some reason (the error will be descriptive).
//
// If the window is already open this function is no-op.
//
// If the window is destroyed, it's NativeWindow struct will be replaced with a
// new one and the window will become valid again.
func (w *Window) Open(screen *Screen) error {
	w.access.RLock()

	if w.opened && !w.destroyed {
		w.access.RUnlock()
		return nil
	}

	// Enter write lock
	w.access.RUnlock()
	w.access.Lock()

	var err error
	if w.destroyed {
		// Swap the *NativeWindow out with an fresh one.
		w.NativeWindow = newNativeWindow(w)
		w.destroyed = false
	}

	w.originalScreen = screen
	w.screen = screen

	w.access.Unlock()

	err = w.NativeWindow.open(screen)
	if err != nil {
		return err
	}

	w.access.Lock()
	w.focused = true
	w.opened = true
	w.access.Unlock()

	return nil
}

// Destroy destroys the window. It is closed, and is then considered to be in
// an destroyed state.
//
// After calling this function, the window is considered destroyed.
//
// If the window is not currently open or is already destroyed then this
// function is no-op.
func (w *Window) Destroy() {
	w.access.RLock()

	if !w.opened || w.destroyed {
		w.access.RUnlock()
		return
	}

	// Enter write lock
	w.access.RUnlock()

	w.access.Lock()
	w.opened = false
	w.destroyed = true
	w.access.Unlock()

	w.send(DestroyedEvent{
		T: time.Now(),
	})
	w.NativeWindow.destroy()
}

// Notify causes the window to notify the user that an event has happened
// with the application, and they should look at the application.
//
// Typically this is an small flashing animation, etc.
func (w *Window) Notify() {
	w.access.RLock()

	if !w.opened || w.destroyed {
		w.access.RUnlock()
		return
	}

	w.access.RUnlock()
	go w.NativeWindow.notify()
}

// SetIcon specifies the window icon which should be displayed anywhere
// that an window icon is needed, this typically includes in the title bar
// decoration, or in the icon tray.
//
// If the icon is nil; the default 'chippy' icon is restored.
func (w *Window) SetIcon(icon image.Image) {
	w.access.RLock()

	if w.icon != icon {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.icon = icon
		if w.opened {
			if icon == nil {
				icon = defaultIcon
			}
			go w.NativeWindow.setIcon(icon)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetTitle sets the title of the window, this is shown anywhere where
// there needs to be an string representation, typical places include the
// window's Title Bar decoration, and in the icon tray (which displays
// minimized windows, etc).
//
// If the window is destroyed, this function will panic.
func (w *Window) SetTitle(title string) {
	w.access.RLock()

	if w.title != title {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.title = title
		if w.opened {
			go w.NativeWindow.setTitle(title)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetVisible specifies whether this window should be visibly seen by the
// user, if false the window will appear simply gone (even though it
// actually exists, and you may render to it, and at an later time show the
// window again).
//
// If the window is destroyed, this function will panic.
func (w *Window) SetVisible(visible bool) {
	w.access.RLock()

	if w.visible != visible {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.visible = visible
		if w.opened {
			go w.NativeWindow.setVisible(visible)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetTransparent specifies whether this window should be transparent, used
// in things like splash screens, etc.
//
// Default: false
func (w *Window) SetTransparent(transparent bool) {
	w.access.RLock()

	if w.transparent != transparent {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.transparent = transparent
		if w.opened {
			go w.NativeWindow.setTransparent(transparent)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetDecorated specifies whether this window should have window
// decorations, this includes the title bar, exit buttons, borders, system
// menu buttons, icons, etc.
//
// If the window is destroyed, this function will panic.
func (w *Window) SetDecorated(decorated bool) {
	w.access.RLock()

	if w.decorated != decorated {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.decorated = decorated
		if w.opened {
			go w.NativeWindow.setDecorated(decorated)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetAlwaysOnTop specifies whether the window should be always on top of
// other windows.
//
// If the window is destroyed, this function will panic.
func (w *Window) SetAlwaysOnTop(alwaysOnTop bool) {
	w.access.RLock()

	if w.alwaysOnTop != alwaysOnTop {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.alwaysOnTop = alwaysOnTop
		if w.opened {
			go w.NativeWindow.setAlwaysOnTop(alwaysOnTop)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetAspectRatio specifies the aspect ratio that the window should try to
// keep when the user resizes the window.
//
// If the ratio is zero, then the window will be allowed to resize freely,
// without being restricted to an aspect ratio.
func (w *Window) SetAspectRatio(ratio float32) {
	w.access.RLock()

	if w.aspectRatio != ratio {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.aspectRatio = ratio
		if w.opened {
			go w.NativeWindow.setAspectRatio(ratio)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

func (w *Window) trySetSize(width, height int) bool {
	w.access.RLock()
	if w.width != width || w.height != height {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.width = width
		w.height = height
		w.access.Unlock()

		w.send(ResizedEvent{
			T:      time.Now(),
			Width:  width,
			Height: height,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// w.access lock must not currently be held
func (w *Window) clampedSize() (width, height int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.doClampedSize()
}

// w.access lock must currently be held
func (w *Window) doClampedSize() (width, height int) {
	width = w.width
	height = w.height

	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}

	if w.minWidth != 0 && w.minHeight != 0 {
		if width < w.minWidth {
			width = w.minWidth
		}
		if height < w.minHeight {

			height = w.minHeight
		}
	}

	if w.maxWidth != 0 && w.maxHeight != 0 {
		if width > w.maxWidth {
			width = w.maxWidth
		}
		if height > w.maxHeight {
			height = w.maxHeight
		}
	}
	return
}

// SetSize specifies the new width and height of this window's client
// region, in pixels.
//
// The window's size will be clamped such that it is always 1px wide/tall; and
// never exceeds the bounds of the minimum or maximum size of the window if one
// is specified.
//
// If w.Size() is later called, it will return the identical (non-clamped)
// values you provide here.
func (w *Window) SetSize(width, height int) {
	w.access.RLock()

	if w.width != width || w.height != height {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.width = width
		w.height = height
		if w.opened {
			clampedWidth, clampedHeight := w.doClampedSize()
			go w.NativeWindow.setSize(clampedWidth, clampedHeight)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetCursorGrabbed specifies whether the mouse cursor should be grabbed,
// this means the cursor will be invisible, and will be forced to stay
// within the client region of the window. This behavior is the same as you
// would typically see in first person shooter games.
//
// If the cursor is being released (false), then the original cursor
// position will be restored to where it was originally at the time of the
// last call to SetCursorGrabbed(true).
func (w *Window) SetCursorGrabbed(grabbed bool) {
	w.access.RLock()

	if w.cursorGrabbed != grabbed {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.cursorGrabbed = grabbed

		if w.cursorGrabbed {
			w.preGrabCursorX = w.cursorX
			w.preGrabCursorY = w.cursorY
		} else {
			w.cursorX = w.preGrabCursorX
			w.cursorY = w.preGrabCursorY
			w.preGrabCursorX = 0
			w.preGrabCursorY = 0
		}

		if w.opened && w.cursorWithin {
			go w.NativeWindow.setCursorGrabbed(grabbed)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetCursor specifies the cursor to become the default cursor while the
// mouse is inside this window's client region.
//
// If the cursor does not exist within the internal cache already -- it
// will be cached using PrepareCursor() automatically. Once you are done
// using the cursor, you should use the FreeCursor() function to remove the
// cursor from the internal cache (the cursor can still be displayed again
// after using FreeCursor(), it just will have to be loaded into the cache
// again).
func (w *Window) SetCursor(cursor *Cursor) {
	w.access.RLock()

	if w.cursor != cursor {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.cursor = cursor
		if w.opened {
			go w.NativeWindow.setCursor(cursor)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// FreeCursor removes the specified cursor from the internal cache. Cursors
// are cached because it allows for SetCursor() operations to perform more
// quickly without performing multiple copy operations underneath.
//
// This allows you to simply use SetCursor() as often as you wish, for
// instance creating cursor animations and such.
//
// If the specified cursor is nil; an panic will occur.
//
// If the specified cursor is the active cursor (previously set via
// the SetCursor function); then the default cursor (nil) will be set and then
// the specified cursor free'd.
func (w *Window) FreeCursor(cursor *Cursor) {
	if cursor == nil {
		panic("FreeCursor(): Cannot free nil cursor!")
	}

	// Special case: We're using this cursor right now
	if cursor == w.Cursor() {
		// Restore the default cursor first
		w.SetCursor(nil)
	}

	w.access.RLock()
	defer w.access.RUnlock()

	if w.opened && !w.destroyed {
		go w.NativeWindow.freeCursor(cursor)
	}
}

// PrepareCursor prepares the specified cursor to be displayed, but does
// not display it. This is useful when you wish to load each frame for an
// cursor animation, but not cause the cursor to flicker while loading them
// into the internal cache.
//
// If the specified cursor is nil; an panic will occur.
//
// If the window is not open or is destroyed; the cursor cannot be prepared and
// this function is no-op.
func (w *Window) PrepareCursor(cursor *Cursor) {
	if cursor == nil {
		panic("PrepareCursor(): Cannot prepare nil cursor!")
	}

	w.access.RLock()
	defer w.access.RUnlock()

	if w.opened && !w.destroyed {
		go w.NativeWindow.prepareCursor(cursor)
	}
}

// SetFullscreen specifies whether the window should be full screen,
// consuming the entire screen's size, and being the only thing displayed
// on the screen.
//
// If the window is destroyed, this function will panic.
func (w *Window) SetFullscreen(fullscreen bool) {
	w.access.RLock()

	if w.fullscreen != fullscreen {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.fullscreen = fullscreen

		if fullscreen {
			// Entering fullscreen: Save window position and size
			w.preFullscreenX = w.x
			w.preFullscreenY = w.y
			w.preFullscreenWidth = w.width
			w.preFullscreenHeight = w.height
		} else {
			// Leaving fullscreen: Restore original size
			w.x = w.preFullscreenX
			w.y = w.preFullscreenY
			w.width = w.preFullscreenWidth
			w.height = w.preFullscreenHeight
		}

		if w.opened {
			go w.NativeWindow.setFullscreen(fullscreen)
			go w.NativeWindow.setSize(w.width, w.height)
			go w.NativeWindow.setPosition(w.x, w.y)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

func (w *Window) trySetMinimized(minimized bool) bool {
	w.access.RLock()
	if w.minimized != minimized {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.minimized = minimized
		w.access.Unlock()

		w.send(MinimizedEvent{
			T:         time.Now(),
			Minimized: minimized,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// SetMinimized specifies whether the window should currently be minimized.
func (w *Window) SetMinimized(minimized bool) {
	w.access.RLock()

	if w.minimized != minimized {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.minimized = minimized
		if w.opened && w.visible {
			go w.NativeWindow.setMinimized(minimized)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

func (w *Window) trySetMaximized(maximized bool) bool {
	w.access.RLock()
	if w.maximized != maximized {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.maximized = maximized
		w.access.Unlock()

		w.send(MaximizedEvent{
			T:         time.Now(),
			Maximized: maximized,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// SetMaximized specifies whether the window should currently be maximized.
func (w *Window) SetMaximized(maximized bool) {
	w.access.RLock()

	if w.maximized != maximized {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.maximized = maximized
		if w.opened && w.visible {
			go w.NativeWindow.setMaximized(maximized)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

func (w *Window) trySetCursorPosition(x, y int) bool {
	w.access.RLock()
	if w.cursorX != x || w.cursorY != y {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.cursorX = x
		w.cursorY = y
		w.access.Unlock()

		if !w.CursorGrabbed() {
			w.send(CursorPositionEvent{
				T: time.Now(),
				X: float64(x),
				Y: float64(y),
			})
		}

		return true
	}
	w.access.RUnlock()
	return false
}

// SetCursorPosition sets the mouse cursor to the new position x and y,
// specified in pixels relative to the client region of this window.
//
// It is possible to move the cursor outside both the client region and
// window region, either by specifying an negative number, or an positive
// number larger than the window region.
func (w *Window) SetCursorPosition(x, y int) {
	w.access.RLock()

	if w.cursorX != x || w.cursorY != y {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.cursorX = x
		w.cursorY = y
		if w.opened {
			go w.NativeWindow.setCursorPosition(x, y)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetMinimumSize specifies the minimum width and height that this windows
// client region is allowed to have, the user will be disallowed to resize
// the window any smaller than this specified size.
//
// If either width or height are zero, then there will be no maximum size
// restriction placed.
//
// If the size passed into both SetMinimumSize and SetMaximumSize are the
// same, then the window will be non-resizable.
func (w *Window) SetMinimumSize(width, height int) {
	w.access.RLock()

	if w.minWidth != width || w.minHeight != height {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.minWidth = width
		w.minHeight = height
		if w.opened {
			go w.NativeWindow.setMinimumSize(width, height)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetMaximumSize specifies the maximum width and height that this windows
// client region is allowed to have, the user will be disallowed to resize
// the window any larger than this specified size.
//
// If the size passed into both SetMaximumSize and SetMinimumSize are the
// same, then the window will be non-resizable.
//
// If either width or height are zero, then there will be no maximum size
// restriction placed.
//
// If the window is destroyed, this function will panic.
func (w *Window) SetMaximumSize(width, height int) {
	w.access.RLock()

	if w.maxWidth != width || w.maxHeight != height {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.maxWidth = width
		w.maxHeight = height
		if w.opened {
			go w.NativeWindow.setMaximumSize(width, height)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

func (w *Window) trySetPosition(x, y int) bool {
	w.access.RLock()
	if w.x != x || w.y != y {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.x = x
		w.y = y
		w.access.Unlock()

		w.send(PositionEvent{
			T: time.Now(),
			X: x,
			Y: y,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// SetPosition specifies the new x and y position of this window's client
// region, relative to the top-left corner of the screen, in pixels.
//
// If the window is destroyed, this function will panic.
func (w *Window) SetPosition(x, y int) {
	w.access.RLock()

	if w.x != x || w.y != y {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.x = x
		w.y = y
		if w.opened {
			go w.NativeWindow.setPosition(x, y)
		}
		w.access.Unlock()
		return
	}

	w.access.RUnlock()
}

// SetPositionCenter sets the window position such that it is perfectly in
// the center of the specified screen.
func (w *Window) SetPositionCenter(screen *Screen) {
	screenWidth, screenHeight := screen.Mode().Resolution()
	windowWidth, windowHeight := w.Size()
	halfScreenWidth := int(screenWidth / 2)
	halfScreenHeight := int(screenHeight / 2)
	halfWindowWidth := int(windowWidth / 2)
	halfWindowHeight := int(windowHeight / 2)
	w.SetPosition(halfScreenWidth-halfWindowWidth, halfScreenHeight-halfWindowHeight)
}

// Opened tells whether the window is currently open.
//
// If the window is destroyed, the returned value will be false.
func (w *Window) Opened() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.opened
}

// OriginalScreen returns the screen that this window was created on at the
// time Open() was called.
func (w *Window) OriginalScreen() *Screen {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.originalScreen
}

func (w *Window) trySetScreen(screen *Screen) bool {
	w.access.RLock()
	if !w.screen.Equals(screen) {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.screen = screen
		w.access.Unlock()

		w.send(ScreenChangedEvent{
			T:      time.Now(),
			Screen: screen,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// Screen returns the current screen that this window is currently residing on.
//
// This function will return the original screen the window was created on in
// the event that we are unable to determine the current screen.
func (w *Window) Screen() *Screen {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.screen
}

// Destroyed tells whether there was an previous call to the Destroy function.
func (w *Window) Destroyed() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.destroyed
}

// Position tells what the current x and y position of this window's client
// region.
func (w *Window) Position() (x, y int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.x, w.y
}

// Size tells the current width and height of this window, as set
// previously by an call to the SetSize function, or due to the user
// resizing the window through the window manager itself.
//
// Both width and height will be at least 1.
func (w *Window) Size() (width, height int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.width, w.height
}

// MinimumSize tells the current minimum width and height of this windows
// client region, as set previously via the SetMinimumSize function, or the
// default values of width=150, height=150.
//
// Both width and height will be at least 1.
func (w *Window) MinimumSize() (width, height int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.minWidth, w.minHeight
}

// MaximumSize tells the current maximum width and height of this windows
// client region, as set previously via the SetMaximumSize function, or the
// default values of width=0, height=0
//
// Both width and height will be at least 1.
func (w *Window) MaximumSize() (width, height int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.maxWidth, w.maxHeight
}

// AspectRatio tells the aspect ratio that the window should try and keep
// when the user resizes the window, as previously set via SetAspectRatio,
// or the default of 0.
//
// Note: If you want to determine the aspect ratio of the window, you
// should instead calculate it from the Size() function, by dividing width
// by height.
//
// (Because if there was no previous call to SetAspectRatio, this function
// will return 0, which is not the actual window aspect ratio.)
func (w *Window) AspectRatio() float32 {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.aspectRatio
}

// Minimized tells whether the window is currently minimized, as previously
// set via an call to the SetMinimized function, or due to the user
// changing the minimized status of the window directly through the window
// manager, or the default value of false.
func (w *Window) Minimized() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.minimized
}

// Maximized tells whether the window is currently maximized, as previously
// set via an call to the SetMaximized function, or due to the user
// changing the maximized status of the window directly through the window
// manager, or the default value of false.
func (w *Window) Maximized() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.maximized
}

// Fullscreen tells whether the window is currently full screen, as
// previously set by an call to the SetFullscreen function.
func (w *Window) Fullscreen() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.fullscreen
}

// AlwaysOnTop tells whether the window is currently always on top of other
// windows, due to an previous call to the SetAlwaysOnTop function, or due
// to the user changing the always on top state directly through the window
// manager itself.
func (w *Window) AlwaysOnTop() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.alwaysOnTop
}

// Icon returns the currently in use icon image, as previously set via an
// call to SetIcon.
//
// Changes made to this Image *after* an initial call to SetIcon will not
// be reflected by the window unless you call SetIcon again.
func (w *Window) Icon() image.Image {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.icon
}

// Cursor returns the currently in use cursor image, as previously set via
// an call to SetCursor.
//
// Changes made to this Image *after* an initial call to SetCursor will not
// be reflected by the window unless you call SetCursor again.
func (w *Window) Cursor() *Cursor {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.cursor
}

// CursorPosition tells the current mouse cursor position, both x and y,
// relative to the client region of this window (specified in pixels)
func (w *Window) CursorPosition() (x, y int) {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.cursorX, w.cursorY
}

func (w *Window) trySetCursorWithin(within bool) bool {
	w.access.RLock()
	if w.cursorWithin != within {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.cursorWithin = within
		w.access.Unlock()

		w.send(CursorWithinEvent{
			T:      time.Now(),
			Within: within,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// CursorWithin tells whether the mouse cursor is inside the window's
// client region or not.
func (w *Window) CursorWithin() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.cursorWithin
}

// CursorGrabbed tells whether the mouse cursor is currently grabbed, as
// previously set via an call to the SetCursorGrabbed function.
func (w *Window) CursorGrabbed() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.cursorGrabbed
}

func (w *Window) trySetExtents(left, right, bottom, top int) {
	w.access.Lock()
	defer w.access.Unlock()

	w.extentLeft = left
	w.extentRight = right
	w.extentBottom = bottom
	w.extentTop = top
}

// Extents returns how far the window region extends outward from the client
// region of this window, in pixels.
//
// If the window's extents are unknown, [0, 0, 0, 0] is returned.
//
// If the window is not open yet, [0, 0, 0, 0] is returned.
//
// If the window is destroyed, [0, 0, 0, 0] is returned.
//
// None of the extents will ever be less than zero.
func (w *Window) Extents() (left, right, bottom, top int) {
	w.access.RLock()
	defer w.access.RUnlock()

	if w.opened && !w.destroyed {
		return w.extentLeft, w.extentRight, w.extentBottom, w.extentTop
	}
	return 0, 0, 0, 0
}

func (w *Window) trySetFocused(focused bool) bool {
	w.access.RLock()
	if w.focused != focused {
		// Enter write lock
		w.access.RUnlock()

		w.access.Lock()
		w.focused = focused
		w.access.Unlock()

		w.send(FocusedEvent{
			T:       time.Now(),
			Focused: focused,
		})

		return true
	}
	w.access.RUnlock()
	return false
}

// Focused tells whether this window currently has focus, and is therefor
// the current window that is being interacted with by the user.
func (w *Window) Focused() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.focused
}

// Transparent tells whether this window is transparent, via an previous
// call to SetTransparent()
func (w *Window) Transparent() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.transparent
}

// Title returns the title of the window, as it was set by SetTitle, or the
// default title: "Chippy Window".
func (w *Window) Title() string {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.title
}

// Visible tells whether this window is currently visible to the user, as
// previously set by the SetVisible function, or the default value of true
// (visible).
func (w *Window) Visible() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.visible
}

// Decorations tells whether this window has window decorations on, as
// previously set by the SetDecorations function, or the default value of
// true (on).
func (w *Window) Decorated() bool {
	w.access.RLock()
	defer w.access.RUnlock()

	return w.decorated
}

// send sends an event to each current event buffer returned from the
// EventBuffer function.
//
// Sending is guaranteed to not block, and is done in an single goroutine.
func (w *Window) send(ev Event) {
	w.access.Lock()
	defer w.access.Unlock()

	for ch, _ := range w.events {
		select {
		case ch <- ev:
			break
		default:
			logger().Println("Warning: Event buffer full; Missed event!")
		}
	}
}

func (w *Window) tryAddKeyboardStateEvent(k keyboard.Key, raw uint64, s keyboard.State) {
	// Prevent accidentially repeating an identical state event
	if w.Keyboard.State(k) == s && w.Keyboard.RawState(raw) == s {
		return
	}

	w.Keyboard.SetState(k, s)
	w.Keyboard.SetRawState(raw, s)
	w.send(keyboard.StateEvent{
		T:     time.Now(),
		Key:   k,
		Raw:   raw,
		State: s,
	})
}

func (w *Window) releaseDownedButtons() {
	// Release the keys that we know of (defined constants) first.
	for key, state := range w.Keyboard.States() {
		if state == keyboard.Down {
			// Change Key state
			state = keyboard.Up
			w.Keyboard.SetState(key, state)

			// Send event
			w.send(keyboard.StateEvent{
				T:     time.Now(),
				Key:   key,
				State: state,
			})
		}
	}

	// Release the keys that we are unknown to us.
	for raw, state := range w.Keyboard.RawStates() {
		if state == keyboard.Down {
			// Change Key state
			state = keyboard.Up
			w.Keyboard.SetRawState(raw, state)

			// Send event
			w.send(keyboard.StateEvent{
				T:     time.Now(),
				Raw:   raw,
				State: state,
			})
		}
	}

	// Release mouse buttons
	for button, state := range w.Mouse.States() {
		if state == mouse.Down {
			// Set new state
			state = mouse.Up

			// Assign button state to the window
			w.Mouse.SetState(button, state)

			// Add mouse event
			w.send(mouse.Event{
				T:      time.Now(),
				Button: button,
				State:  state,
			})
		}
	}
}

// EventsBuffer returns an new read-only channel over which window events will
// be sent.
//
// Events are sent in an non-blocking fashion. As such once the amount of
// buffered items in the channel reaches the maximum, then events will stop
// being sent.
//
// You should ensure that the buffer size is large enough for you to read the
// events in an short amount of time.
func (w *Window) EventsBuffer(bufferSize int) chan Event {
	w.access.Lock()
	defer w.access.Unlock()

	ch := make(chan Event, bufferSize)
	w.events[ch] = true
	return ch
}

// Events is short-hand for:
//
//  w.EventsBuffer(64)
//
func (w *Window) Events() chan Event {
	return w.EventsBuffer(64)
}

// CloseEvents stops the specified channel from receiving any more window
// events.
//
// Call this function whenever you are done receiving window events from
// the given channel.
//
// This is important for channels which have an very large buffer size to
// reduce memory consumption.
func (w *Window) CloseEvents(ch chan Event) {
	w.access.RLock()

	_, ok := w.events[ch]
	if ok {
		// Exit read lock
		w.access.RUnlock()

		// Enter write lock
		w.access.Lock()
		delete(w.events, ch)
		w.access.Unlock()
		return
	}
	w.access.RUnlock()
}

func NewWindow() *Window {
	w := new(Window)
	w.NativeWindow = newNativeWindow(w)
	w.Keyboard = keyboard.NewWatcher()
	w.Mouse = mouse.NewWatcher()

	w.events = make(map[chan Event]bool)

	w.SetIcon(defaultIcon)

	w.SetTitle("Chippy Window")
	w.SetVisible(true)
	w.SetDecorated(true)
	w.SetPosition(100, 100)
	w.SetSize(640, 480)
	//w.SetMaximumSize(0, 0)
	w.SetMinimumSize(150, 150)
	w.SetAspectRatio(0)
	//w.SetMinimized(false)
	//w.SetMaximized(false)
	//w.SetFullscreen(false)
	//w.SetAlwaysOnTop(false)
	return w
}
