// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"errors"
	"image"
	"image/draw"
	"math"
	"reflect"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"azul3d.org/chippy.v1/internal/resize"
	"azul3d.org/chippy.v1/internal/x11"
	"azul3d.org/keyboard.v1"
	"azul3d.org/mouse.v1"
)

func toKeyboard(ks x11.Keysym) keyboard.Key {
	switch ks {
	case x11.XK_BackSpace:
		return keyboard.Backspace
	case x11.XK_Tab:
		return keyboard.Tab
	case x11.XK_ISO_Left_Tab:
		return keyboard.Tab
	case x11.XK_Return:
		return keyboard.Enter
	case x11.XK_Linefeed:
		return keyboard.Enter
	case x11.XK_Pause:
		return keyboard.Pause
	case x11.XK_Clear:
		return keyboard.Clear
	case x11.XK_Scroll_Lock:
		return keyboard.ScrollLock
	case x11.XK_Escape:
		return keyboard.Escape
	case x11.XK_Delete:
		return keyboard.Delete
	case x11.XK_minus:
		return keyboard.Dash
	case x11.XK_underscore:
		return keyboard.Dash
	case x11.XK_asciitilde:
		return keyboard.Tilde
	case x11.XK_grave:
		return keyboard.Tilde
	case x11.XK_Print:
		return keyboard.Print
	case x11.XK_Insert:
		return keyboard.Insert
	case x11.XK_3270_PrintScreen:
		return keyboard.PrintScreen

	case x11.XK_Select:
		return keyboard.Select
	case x11.XK_Execute:
		return keyboard.Execute
	case x11.XK_Help:
		return keyboard.Help
	case x11.XK_3270_Play:
		return keyboard.Play
	case x11.XK_3270_ExSelect:
		return keyboard.ExSel
	case x11.XK_3270_CursorSelect:
		return keyboard.CrSel
	case x11.XK_3270_Attn:
		return keyboard.Attn
	case x11.XK_3270_EraseEOF:
		return keyboard.EraseEOF
	case x11.XK_Kanji:
		return keyboard.Kanji

	case x11.XK_plus:
		return keyboard.Equals
	case x11.XK_equal:
		return keyboard.Equals
	case x11.XK_colon:
		return keyboard.Semicolon
	case x11.XK_semicolon:
		return keyboard.Semicolon
	case x11.XK_quotedbl:
		return keyboard.Apostrophe
	case x11.XK_apostrophe:
		return keyboard.Apostrophe
	case x11.XK_comma:
		return keyboard.Comma
	case x11.XK_less:
		return keyboard.Comma
	case x11.XK_period:
		return keyboard.Period
	case x11.XK_greater:
		return keyboard.Period
	case x11.XK_slash:
		return keyboard.ForwardSlash
	case x11.XK_question:
		return keyboard.ForwardSlash
	case x11.XK_backslash:
		return keyboard.BackSlash
	case x11.XK_bar:
		return keyboard.BackSlash
	case x11.XK_space:
		return keyboard.Space

	case x11.XK_0:
		return keyboard.Zero
	case x11.XK_1:
		return keyboard.One
	case x11.XK_2:
		return keyboard.Two
	case x11.XK_3:
		return keyboard.Three
	case x11.XK_4:
		return keyboard.Four
	case x11.XK_5:
		return keyboard.Five
	case x11.XK_6:
		return keyboard.Six
	case x11.XK_7:
		return keyboard.Seven
	case x11.XK_8:
		return keyboard.Eight
	case x11.XK_9:
		return keyboard.Nine

	case x11.XK_parenright:
		return keyboard.Zero
	case x11.XK_exclam:
		return keyboard.One
	case x11.XK_at:
		return keyboard.Two
	case x11.XK_numbersign:
		return keyboard.Three
	case x11.XK_dollar:
		return keyboard.Four
	case x11.XK_percent:
		return keyboard.Five
	case x11.XK_asciicircum:
		return keyboard.Six
	case x11.XK_ampersand:
		return keyboard.Seven
	case x11.XK_asterisk:
		return keyboard.Eight
	case x11.XK_parenleft:
		return keyboard.Nine

	// Lower case alphabet maps to same thing in chippy's virtual keys
	case x11.XK_a:
		return keyboard.A
	case x11.XK_b:
		return keyboard.B
	case x11.XK_c:
		return keyboard.C
	case x11.XK_d:
		return keyboard.D
	case x11.XK_e:
		return keyboard.E
	case x11.XK_f:
		return keyboard.F
	case x11.XK_g:
		return keyboard.G
	case x11.XK_h:
		return keyboard.H
	case x11.XK_i:
		return keyboard.I
	case x11.XK_j:
		return keyboard.J
	case x11.XK_k:
		return keyboard.K
	case x11.XK_l:
		return keyboard.L
	case x11.XK_m:
		return keyboard.M
	case x11.XK_n:
		return keyboard.N
	case x11.XK_o:
		return keyboard.O
	case x11.XK_p:
		return keyboard.P
	case x11.XK_q:
		return keyboard.Q
	case x11.XK_r:
		return keyboard.R
	case x11.XK_s:
		return keyboard.S
	case x11.XK_t:
		return keyboard.T
	case x11.XK_u:
		return keyboard.U
	case x11.XK_v:
		return keyboard.V
	case x11.XK_w:
		return keyboard.W
	case x11.XK_x:
		return keyboard.X
	case x11.XK_y:
		return keyboard.Y
	case x11.XK_z:
		return keyboard.Z

	case x11.XK_A:
		return keyboard.A
	case x11.XK_B:
		return keyboard.B
	case x11.XK_C:
		return keyboard.C
	case x11.XK_D:
		return keyboard.D
	case x11.XK_E:
		return keyboard.E
	case x11.XK_F:
		return keyboard.F
	case x11.XK_G:
		return keyboard.G
	case x11.XK_H:
		return keyboard.H
	case x11.XK_I:
		return keyboard.I
	case x11.XK_J:
		return keyboard.J
	case x11.XK_K:
		return keyboard.K
	case x11.XK_L:
		return keyboard.L
	case x11.XK_M:
		return keyboard.M
	case x11.XK_N:
		return keyboard.N
	case x11.XK_O:
		return keyboard.O
	case x11.XK_P:
		return keyboard.P
	case x11.XK_Q:
		return keyboard.Q
	case x11.XK_R:
		return keyboard.R
	case x11.XK_S:
		return keyboard.S
	case x11.XK_T:
		return keyboard.T
	case x11.XK_U:
		return keyboard.U
	case x11.XK_V:
		return keyboard.V
	case x11.XK_W:
		return keyboard.W
	case x11.XK_X:
		return keyboard.X
	case x11.XK_Y:
		return keyboard.Y
	case x11.XK_Z:
		return keyboard.Z

	case x11.XK_F1:
		return keyboard.F1
	case x11.XK_F2:
		return keyboard.F2
	case x11.XK_F3:
		return keyboard.F3
	case x11.XK_F4:
		return keyboard.F4
	case x11.XK_F5:
		return keyboard.F5
	case x11.XK_F6:
		return keyboard.F6
	case x11.XK_F7:
		return keyboard.F7
	case x11.XK_F8:
		return keyboard.F8
	case x11.XK_F9:
		return keyboard.F9
	case x11.XK_F10:
		return keyboard.F10
	case x11.XK_F11:
		return keyboard.F11
	case x11.XK_F12:
		return keyboard.F12
	case x11.XK_F13:
		return keyboard.F13
	case x11.XK_F14:
		return keyboard.F14
	case x11.XK_F15:
		return keyboard.F15
	case x11.XK_F16:
		return keyboard.F16
	case x11.XK_F17:
		return keyboard.F17
	case x11.XK_F18:
		return keyboard.F18
	case x11.XK_F19:
		return keyboard.F19
	case x11.XK_F20:
		return keyboard.F20
	case x11.XK_F21:
		return keyboard.F21
	case x11.XK_F22:
		return keyboard.F22
	case x11.XK_F23:
		return keyboard.F23
	case x11.XK_F24:
		return keyboard.F24

	case x11.XK_Shift_L:
		return keyboard.LeftShift
	case x11.XK_Shift_R:
		return keyboard.RightShift
	case x11.XK_Control_L:
		return keyboard.LeftCtrl
	case x11.XK_Control_R:
		return keyboard.RightCtrl
	case x11.XK_Alt_L:
		return keyboard.LeftAlt
	case x11.XK_Alt_R:
		return keyboard.RightAlt
	case x11.XK_Super_L:
		return keyboard.LeftSuper
	case x11.XK_Super_R:
		return keyboard.RightSuper

	case x11.XK_braceleft:
		return keyboard.LeftBracket
	case x11.XK_bracketleft:
		return keyboard.LeftBracket

	case x11.XK_braceright:
		return keyboard.RightBracket
	case x11.XK_bracketright:
		return keyboard.RightBracket

	case x11.XK_Home:
		return keyboard.Home
	case x11.XK_Left:
		return keyboard.ArrowLeft
	case x11.XK_Right:
		return keyboard.ArrowRight
	case x11.XK_Down:
		return keyboard.ArrowDown
	case x11.XK_Up:
		return keyboard.ArrowUp
	case x11.XK_Page_Up:
		return keyboard.PageUp
	case x11.XK_Page_Down:
		return keyboard.PageDown
	case x11.XK_End:
		return keyboard.End

	case x11.XK_Caps_Lock:
		return keyboard.CapsLock
	case x11.XK_Num_Lock:
		return keyboard.NumLock

	case x11.XK_KP_Enter:
		return keyboard.NumEnter
	case x11.XK_KP_Multiply:
		return keyboard.NumMultiply
	case x11.XK_KP_Divide:
		return keyboard.NumDivide
	case x11.XK_KP_Subtract:
		return keyboard.NumSubtract
	case x11.XK_KP_Separator:
		return keyboard.NumComma
	case x11.XK_KP_Decimal:
		return keyboard.NumDecimal
	case x11.XK_KP_Add:
		return keyboard.NumAdd

	case x11.XK_KP_Delete:
		return keyboard.NumDecimal

	case x11.XK_KP_0:
		return keyboard.NumZero
	case x11.XK_KP_1:
		return keyboard.NumOne
	case x11.XK_KP_2:
		return keyboard.NumTwo
	case x11.XK_KP_3:
		return keyboard.NumThree
	case x11.XK_KP_4:
		return keyboard.NumFour
	case x11.XK_KP_5:
		return keyboard.NumFive
	case x11.XK_KP_6:
		return keyboard.NumSix
	case x11.XK_KP_7:
		return keyboard.NumSeven
	case x11.XK_KP_8:
		return keyboard.NumEight
	case x11.XK_KP_9:
		return keyboard.NumNine

	case x11.XK_KP_Insert:
		return keyboard.NumZero
	case x11.XK_KP_End:
		return keyboard.NumOne
	case x11.XK_KP_Down:
		return keyboard.NumTwo
	case x11.XK_KP_Page_Down:
		return keyboard.NumThree
	case x11.XK_KP_Left:
		return keyboard.NumFour
	case x11.XK_KP_Begin:
		return keyboard.NumFive
	case x11.XK_KP_Right:
		return keyboard.NumSix
	case x11.XK_KP_Home:
		return keyboard.NumSeven
	case x11.XK_KP_Up:
		return keyboard.NumEight
	case x11.XK_KP_Page_Up:
		return keyboard.NumNine
	}
	return keyboard.Invalid
}

type loadedCursor struct {
	cursor x11.Cursor
}

const (
	renderModeOpenGL = iota
	renderModePixelBlit
)

var (
	cachedExtents       []int32
	cachedExtentsAccess sync.RWMutex
)

type NativeWindow struct {
	access, extentsAccess sync.RWMutex

	r *Window

	xWindowAccess sync.RWMutex
	xWindow       x11.Window

	waitForMap, waitForUnmap, waitForFrameExtents, waitForMotifHints,
	waitForNetWmAllowedActions chan bool

	extents                  []int32
	xVisual                  x11.VisualId
	xDepth                   uint8
	xGC                      x11.GContext
	xicon                    []uint32
	pixmapFmt32, pixmapFmt24 *x11.Format
	cursors                  map[*Cursor]loadedCursor
	activeRenderMode         int

	last struct {
		sync.RWMutex
		cursorX, cursorY int
	}

	can struct {
		sync.RWMutex
		sendPositionEvents     bool
		sendSizeEvents         bool
		sendRelativeMoveEvents bool
	}

	// OpenGL things here
	glVSyncMode         VSyncMode
	glConfig            *GLConfig
	glxExtensionsString string
}

func (w *NativeWindow) open(screen *Screen) (err error) {
	w.access.Lock()
	defer w.access.Unlock()

	// When we first open the window, we can just use PixelBlit mode.
	w.activeRenderMode = renderModePixelBlit
	return w.doRebuildWindow()
}

func (w *NativeWindow) doRebuildWindow() (err error) {
	if w.r.Opened() {
		w.doDestroy()
	}
	w.clearLastCursorPosition()
	screen := w.r.Screen()
	width, height := w.r.clampedSize()
	x, y := w.r.Position()

	xScreen := xConnection.ScreenOfDisplay(screen.NativeScreen.xScreen)

	xWindow := xConnection.GenerateId()
	w.xWindowAccess.Lock()
	w.xWindow = xWindow
	w.xWindowAccess.Unlock()

	// For event management
	xWindowLookupAccess.Lock()
	xWindowLookup[xWindow] = w
	xWindowLookupAccess.Unlock()

	var eventMask uint32
	eventMask |= x11.EVENT_MASK_NO_EVENT
	eventMask |= x11.EVENT_MASK_KEY_PRESS
	eventMask |= x11.EVENT_MASK_KEY_RELEASE
	eventMask |= x11.EVENT_MASK_BUTTON_PRESS
	eventMask |= x11.EVENT_MASK_BUTTON_RELEASE
	eventMask |= x11.EVENT_MASK_POINTER_MOTION
	eventMask |= x11.EVENT_MASK_ENTER_WINDOW
	eventMask |= x11.EVENT_MASK_LEAVE_WINDOW
	eventMask |= x11.EVENT_MASK_FOCUS_CHANGE
	eventMask |= x11.EVENT_MASK_EXPOSURE
	eventMask |= x11.EVENT_MASK_PROPERTY_CHANGE
	eventMask |= x11.EVENT_MASK_STRUCTURE_NOTIFY

	// Find proper visual
	switch w.activeRenderMode {
	case renderModePixelBlit:
		w.doFindProperVisual()

	case renderModeOpenGL:
		// Don't touch the already-set w.xVisual or w.xDepth
		break

	default:
		panic("Invalid render mode.")
	}

	// Create colormap
	cmap := x11.Colormap(xConnection.GenerateId())
	createColormap := xConnection.CreateColormapChecked(false, cmap, xScreen.Root, w.xVisual)
	tmpErr := xConnection.RequestCheck(createColormap)
	if tmpErr != nil {
		return errors.New("x11_create_colormap(): " + tmpErr.Error())
	}

	// MUST keep order:
	// XCB_CW_BACK_PIXMAP       = 1L<<0,
	// XCB_CW_BACK_PIXEL        = 1L<<1,
	// XCB_CW_BORDER_PIXMAP     = 1L<<2,
	// XCB_CW_BORDER_PIXEL      = 1L<<3,
	// XCB_CW_BIT_GRAVITY       = 1L<<4,
	// XCB_CW_WIN_GRAVITY       = 1L<<5,
	// XCB_CW_BACKING_STORE     = 1L<<6,
	// XCB_CW_BACKING_PLANES    = 1L<<7,
	// XCB_CW_BACKING_PIXEL     = 1L<<8,
	// XCB_CW_OVERRIDE_REDIRECT = 1L<<9,
	// XCB_CW_SAVE_UNDER        = 1L<<10,
	// XCB_CW_EVENT_MASK        = 1L<<11,
	// XCB_CW_DONT_PROPAGATE    = 1L<<12,
	// XCB_CW_COLORMAP          = 1L<<13,
	// XCB_CW_CURSOR            = 1L<<14

	values := []uint32{
		0,            // CW_BORDER_PIXEL
		eventMask,    // CW_EVENT_MASK
		uint32(cmap), // CW_COLORMAP
	}

	// Note: We also set the window position after mapping, as some WM's like
	// Unity seem to ignore initial placement requests.
	tmpErr = xConnection.RequestCheck(xConnection.CreateWindowChecked(
		w.xDepth,
		xWindow,
		xScreen.Root,
		int16(x), int16(y),
		uint16(width), uint16(height),
		0, // border_width
		x11.WINDOW_CLASS_INPUT_OUTPUT,
		w.xVisual,
		x11.CW_BORDER_PIXEL|x11.CW_EVENT_MASK|x11.CW_COLORMAP,
		&values[0], // masks
	))
	if tmpErr != nil {
		return errors.New("x11_create_window(): " + tmpErr.Error())
	}

	var gmask uint32
	gmask |= x11.GC_FOREGROUND
	gmask |= x11.GC_BACKGROUND
	gvalues := []uint32{0, 0}
	w.xGC = x11.GContext(xConnection.GenerateId())
	tmpErr = xConnection.RequestCheck(xConnection.CreateGCChecked(
		w.xGC,
		x11.Drawable(xWindow),
		gmask,
		&gvalues[0],
	))
	if tmpErr != nil {
		logger().Println("x11_create_gc():", tmpErr)
		logger().Println("Some features will not work due to above errors.")
		w.xGC = 0
	}

	cursor := w.r.Cursor()
	if cursor != nil {
		// Note: We can't specify cursor before load, because we need a GC to
		// create the cursor, and to create the GC we need a window! No biggie.
		w.doSetCursor(cursor)
	}

	// Update window title
	w.doSetTitle(w.r.Title())

	// Configure size now, since ICCCM WM's need min/max and aspect ratios now.
	w.doConfigureSize(width, height)

	// Since our window is not mapped (and may never be) we can ask the WM to
	// guess the window extents.
	xConnection.SendClientMessage(xWindow, xScreen.Root, aNetRequestFrameExtents, 0, nil)
	xConnection.Flush()

	var cExtents []int32
	cachedExtentsAccess.RLock()
	copy(cExtents, cachedExtents)
	cachedExtentsAccess.RUnlock()

	haveCExtents := len(cExtents) == 4
	haveCExtentsButZero := haveCExtents && cExtents[0] == 0 && cExtents[1] == 0 && cExtents[2] == 0 && cExtents[3] == 0
	if haveCExtents && !haveCExtentsButZero {
		w.extentsAccess.Lock()
		w.extents = cExtents
		w.extentsAccess.Unlock()
	} else {
		// Wait a bit to see if we can get _NET_REQUEST_FRAME_EXTENTS to respond.
		select {
		case <-time.After(1 * time.Second):
			logger().Println("Timed out waiting for _NET_REQUEST_FRAME_EXTENTS request.")
			w.fetchExtents()
			break

		case <-w.waitForFrameExtents:
			break
		}
	}

	w.extentsAccess.RLock()
	haveExtents := len(w.extents) == 4
	haveExtentsButZero := haveExtents && w.extents[0] == 0 && w.extents[1] == 0 && w.extents[2] == 0 && w.extents[3] == 0
	w.extentsAccess.RUnlock()

	if !haveExtents || haveExtentsButZero {
		w.doSetVisible(true)

		// We wait to see if we get frame extents for a bit, since some WM's need
		// time to do this.
		select {
		case <-w.waitForFrameExtents:
			break
		case <-time.After(1 * time.Second):
			logger().Println("Timed out waiting for _NET_FRAME_EXTENTS PropertyNotify event.")
			break
		}
	}

	if haveExtents {
		// extents is [l, r, t, b], but trySetExtents takes [l, r, b, t]
		w.extentsAccess.RLock()
		l := int(w.extents[0])
		r := int(w.extents[1])
		t := int(w.extents[2])
		b := int(w.extents[3])
		w.extentsAccess.RUnlock()
		w.r.trySetExtents(l, r, b, t)
	}

	// Locate 24-bit pixmap format
	w.pixmapFmt24 = xConnection.FindPixmapFormat(24, 32)
	if w.pixmapFmt24 == nil {
		logger().Println("Could not locate 24-bit depth pixmap format!")
		logger().Println("Some features will not work due to the above error!")
	}

	// Locate 32-bit pixmap format
	w.pixmapFmt32 = xConnection.FindPixmapFormat(32, 32)
	if w.pixmapFmt32 == nil {
		logger().Println("Could not locate 32-bit depth pixmap format!")
		logger().Println("Some features will not work due to the above error!")
	}

	// Finally show/hide the window
	w.doSetVisible(w.r.Visible())

	// Send indicator events
	w.refreshIndicators()

	if aWmProtocols != 0 && aWmDeleteWindow != 0 {
		// We accept WM_DELETE_WINDOW messages.
		xConnection.ChangeProperty(
			x11.PROP_MODE_REPLACE,
			xWindow,
			aWmProtocols,
			x11.ATOM_ATOM,
			32,
			1,
			unsafe.Pointer(&aWmDeleteWindow),
		)
	}

	return
}

func (w *NativeWindow) getXWindow() x11.Window {
	w.xWindowAccess.RLock()
	xWindow := w.xWindow
	w.xWindowAccess.RUnlock()
	return xWindow
}

func (w *NativeWindow) doFindProperVisual() {
	xScreen := xConnection.ScreenOfDisplay(w.r.Screen().NativeScreen.xScreen)

	w.xVisual = xScreen.RootVisual
	w.xDepth = uint8(x11.COPY_FROM_PARENT)

	transparent := w.r.Transparent()

	// See if we can use the visual that is associated with GLBestConfig
	// because then we would not require rebuilding the window should they try
	// to turn the window into an OpenGL one.
	glxCompatible := GLChooseConfig(w.r.GLConfigs(), GLWorstConfig, GLBestConfig)

	var (
		foundCompatibleVisual bool
		bestDepth             uint8
		bestVisual            x11.VisualId
	)
	xConnection.ScreenAllowedDepthsIterator(xScreen, func(d *x11.Depth) bool {
		d.Iterate(func(vis *x11.VisualType) bool {
			depth := uint8(d.Depth)
			depthMatchesTransparency := (transparent && depth == 32) || (!transparent && depth == 24)
			if depthMatchesTransparency {
				bestDepth = depth
				bestVisual = vis.VisualId
			}

			if glxCompatible != nil {
				// We're searching for a GLX compatible X visual as well.
				if vis.VisualId == glxCompatible.xVisual && depthMatchesTransparency {
					// This one is compatible! Huzzah!
					foundCompatibleVisual = true
					w.xDepth = depth
					w.xVisual = vis.VisualId
				}
			}
			return true
		})
		return true
	})

	if !foundCompatibleVisual {
		// Could not find 24/32 bit-depth compatible GLBestConfig or
		// GLBestTransparentConfig visual. This is alright, but GLSetConfig
		// has no choice but to rebuild the window in-place now.
		w.xDepth = bestDepth
		w.xVisual = bestVisual
	}
}

func (w *NativeWindow) refreshIndicators() {
	state, status := xDisplay.XkbGetIndicatorState(x11.XkbUseCoreKbd)
	if status == x11.Success {
		// Caps Lock
		key := keyboard.CapsLock
		os := uint64(x11.XK_Caps_Lock)
		keyState := keyboard.Off
		if (state & 0x01) > 0 {
			keyState = keyboard.On
		}
		w.r.tryAddKeyboardStateEvent(key, os, keyState)

		// Num Lock
		key = keyboard.NumLock
		os = uint64(x11.XK_Num_Lock)
		keyState = keyboard.Off
		if (state & 0x02) > 0 {
			keyState = keyboard.On
		}
		w.r.tryAddKeyboardStateEvent(key, os, keyState)

		// Scroll Lock
		key = keyboard.ScrollLock
		os = uint64(x11.XK_Scroll_Lock)
		keyState = keyboard.Off
		if (state & 0x04) > 0 {
			keyState = keyboard.On
		}
		w.r.tryAddKeyboardStateEvent(key, os, keyState)
	}
}

func (w *NativeWindow) handleEvent(ref *x11.GenericEvent, e interface{}) {
	//logger().Println(reflect.TypeOf(e))
	//logger().Printf("%+v\n", e)
	switch ev := e.(type) {
	case *x11.KeyPressEvent:
		xkbContext.Lock()
		keycode := x11.XkbKeycode(ev.Detail)
		keysym := xkbState.KeyGetOneSym(keycode)
		r := keysym.Rune()
		xkbContext.Unlock()

		kb := toKeyboard(x11.Keysym(keysym))
		if kb != keyboard.Invalid {
			if kb == keyboard.CapsLock || kb == keyboard.NumLock || kb == keyboard.ScrollLock {
				w.refreshIndicators()
			} else {
				w.r.tryAddKeyboardStateEvent(kb, uint64(keysym), keyboard.Down)
			}
		} else {
			logger().Println("Unknown X keysym", keysym)
		}

		escape := '\u001b'
		deleteKey := '\u007f'
		if r != utf8.RuneError && r != 0 && r != escape && r != deleteKey {
			w.r.send(keyboard.TypedEvent{
				T:    time.Now(),
				Rune: r,
			})
		}

	case *x11.KeyReleaseEvent:
		xkbContext.Lock()
		keycode := x11.XkbKeycode(ev.Detail)
		keysym := xkbState.KeyGetOneSym(keycode)
		xkbContext.Unlock()

		kb := toKeyboard(x11.Keysym(keysym))
		if kb != keyboard.Invalid {
			if kb == keyboard.CapsLock || kb == keyboard.NumLock || kb == keyboard.ScrollLock {
				w.refreshIndicators()
			} else {
				w.r.tryAddKeyboardStateEvent(kb, uint64(keysym), keyboard.Up)
			}
		} else {
			logger().Println("Unknown X keysym", keysym)
		}

	case *x11.ButtonPressEvent:
		var (
			button mouse.Button
			state  = mouse.Down
		)

		switch ev.Detail {
		case 1:
			button = mouse.Left
		case 2:
			button = mouse.Wheel
		case 3:
			button = mouse.Right
		case 4:
			button = mouse.Wheel
			state = mouse.ScrollForward
		case 5:
			button = mouse.Wheel
			state = mouse.ScrollBack
		case 6:
			button = mouse.Wheel
			state = mouse.ScrollLeft
		case 7:
			button = mouse.Wheel
			state = mouse.ScrollRight
		default:
			logger().Printf("Unknown button press event; Detail=%v\n", ev.Detail)
			return
		}
		w.r.send(mouse.Event{
			T:      time.Now(),
			Button: button,
			State:  state,
		})

	case *x11.ButtonReleaseEvent:
		var (
			button mouse.Button
			state  = mouse.Up
		)

		switch ev.Detail {
		case 1:
			button = mouse.Left
		case 2:
			button = mouse.Wheel
		case 3:
			button = mouse.Right
		case 4:
			button = mouse.Wheel
			state = mouse.ScrollForward
		case 5:
			button = mouse.Wheel
			state = mouse.ScrollBack
		case 6:
			button = mouse.Wheel
			state = mouse.ScrollLeft
		case 7:
			button = mouse.Wheel
			state = mouse.ScrollRight
		default:
			logger().Printf("Unknown button release event; Detail=%v\n", ev.Detail)
			return
		}
		w.r.send(mouse.Event{
			T:      time.Now(),
			Button: button,
			State:  state,
		})

	case *x11.MotionNotifyEvent:
		x := int(ev.EventX)
		y := int(ev.EventY)
		if w.r.CursorGrabbed() && w.r.Focused() {
			// Find relative movement
			w.last.Lock()
			diffX := float64(x - w.last.cursorX)
			diffY := float64(y - w.last.cursorY)

			// Check if the last cursor position was not known, if so reject
			// this relative movement.
			if w.last.cursorX == -1 || w.last.cursorY == -1 {
				diffX = 0
				diffY = 0
				// Update cursor position to this one (warpPointer does it for
				// us below but we won't get to it in this case):
				w.last.cursorX = x
				w.last.cursorY = y
			}
			w.last.Unlock()
			if diffX != 0 || diffY != 0 {
				w.can.Lock()
				if w.can.sendRelativeMoveEvents {
					w.r.send(CursorPositionEvent{
						T: time.Now(),
						X: diffX,
						Y: diffY,
					})
				}
				w.can.sendRelativeMoveEvents = true
				w.can.Unlock()

				// Event though we have the cursor grabbed, it might still
				// hit the border of the window, in which case we can't
				// determine relative movement!
				//
				// Best thing we can do is set it back to the center of the
				// window.
				wWidth, wHeight := w.r.Size()
				w.warpPointer(wWidth/2, wHeight/2)
				return
			}
		} else {
			w.last.Lock()
			w.last.cursorX = x
			w.last.cursorY = y
			w.last.Unlock()
			w.r.trySetCursorPosition(x, y)
		}

	case *x11.EnterNotifyEvent:
		w.r.trySetCursorWithin(true)
		w.clearLastCursorPosition()

		// The cursor is inside the window, and the window has focus so we
		// should grab the cursor.
		if w.r.CursorGrabbed() && w.r.Focused() {
			go func() {
				w.access.Lock()
				defer w.access.Unlock()
				w.doSetCursorGrabbed(
					true,  // grabbed
					false, // restore position
				)
				xConnection.Flush()
			}()
		}

	case *x11.LeaveNotifyEvent:
		w.r.trySetCursorWithin(false)

	case *x11.FocusInEvent:
		w.r.trySetFocused(true)

		if w.r.CursorGrabbed() {
			if w.r.CursorWithin() {
				// Window in focus, cursor inside, grabbing is allowed.
				go func() {
					w.access.Lock()
					defer w.access.Unlock()
					w.doSetCursorGrabbed(
						true,  // grabbed
						false, // restore position
					)
					xConnection.Flush()
				}()
			} else {
				// Window in focus, cursor not inside, no grabbing allowed.
				go func() {
					w.access.Lock()
					defer w.access.Unlock()
					w.doSetCursorGrabbed(
						false, // grabbed
						false, // restore position
					)
					xConnection.Flush()
				}()
			}
		}

	case *x11.FocusOutEvent:
		w.r.trySetFocused(false)

		// Release downed mouse buttons or else they will become stuck and
		// annoy the user.
		w.r.releaseDownedButtons()

		// The window has lost focus, ungrab the cursor now otherwise the user
		// won't have a mouse cursor (i.e. if they alt-tab away from the app or
		// switch screens).
		if w.r.CursorGrabbed() {
			// We only restore the cursor position if it's actually inside the
			// window.
			restorePosition := w.r.CursorWithin()
			go func() {
				w.access.Lock()
				defer w.access.Unlock()
				w.doSetCursorGrabbed(
					false,           // grabbed
					restorePosition, // restore position
				)
				xConnection.Flush()
			}()
		}

	case *x11.ClientMessageEvent:
		atom := *(*x11.Atom)(unsafe.Pointer(&ev.Data[0]))
		if atom == aWmDeleteWindow {
			w.r.send(CloseEvent{
				T: time.Now(),
			})
		}

	case *x11.PropertyNotifyEvent:
		// Copy evente atom because ev will be free'd below.
		evAtom := ev.Atom
		go func() {
			if evAtom == aNetFrameExtents {
				w.fetchExtents()
				select {
				case w.waitForFrameExtents <- true:
					break
				case <-time.After(5 * time.Second):
					break
				}

			} else if evAtom == aMotifWmHints {
				select {
				case w.waitForMotifHints <- true:
					break
				case <-time.After(5 * time.Second):
					break
				}
			}
		}()

	case *x11.ConfigureNotifyEvent:
		w.can.RLock()
		if ev.Width != 0 && ev.Height != 0 {
			if w.can.sendSizeEvents {
				w.r.trySetSize(int(ev.Width), int(ev.Height))
			}
		}

		if ev.X != 0 || ev.Y != 0 {
			if w.can.sendPositionEvents {
				x := int(ev.X)
				y := int(ev.Y)
				w.r.trySetPosition(x, y)
			}
		}
		w.can.RUnlock()

	case *x11.ExposeEvent:
		go func() {
			select {
			case w.waitForMap <- true:
				break
			case <-time.After(1 * time.Second):
				break
			}
		}()

	case *x11.MapNotifyEvent:
		go func() {
			select {
			case w.waitForMap <- true:
				break
			case <-time.After(5 * time.Second):
				break
			}
		}()

	case *x11.UnmapNotifyEvent:
		go func() {
			select {
			case w.waitForUnmap <- true:
				break
			case <-time.After(5 * time.Second):
				break
			}
		}()

	case *x11.ReparentNotifyEvent:
		break

	default:
		logger().Println(reflect.TypeOf(ev))
	}

	// Free the event reference.
	ref.Free()
}

func (w *NativeWindow) fetchExtents() {
	_, ptr, ptrLen, err := xConnection.GetProperty(
		false, // delete
		w.getXWindow(),
		aNetFrameExtents,
		x11.ATOM_CARDINAL,
		0,              // offset
		math.MaxUint32, // length
	)
	if err == nil && ptrLen == 16 {
		var extents []int32
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&extents))
		sliceHeader.Cap = 4
		sliceHeader.Len = 4
		sliceHeader.Data = uintptr(ptr)

		w.extentsAccess.Lock()
		defer w.extentsAccess.Unlock()

		// Note: We don't perform trySetExtents here but only in open(),
		// because window extents should not change once the window is opened.
		cpy := make([]int32, len(extents))
		copy(cpy, extents)
		w.extents = cpy

		if len(w.extents) == 4 {
			cachedExtentsAccess.Lock()
			copy(cachedExtents, w.extents)
			cachedExtentsAccess.Unlock()
		}

	} else {
		if err != nil {
			logger().Println("GetProperty(_NET_FRAME_EXTENTS):", err)
		}
	}
}

func (w *NativeWindow) configurePosition(windowX, windowY int) {
	w.access.RLock()
	defer w.access.RUnlock()

	w.doConfigurePosition(windowX, windowY)
}

func (w *NativeWindow) doConfigurePosition(windowX, windowY int) {
	var left, top int

	w.extentsAccess.RLock()
	if len(w.extents) > 0 && w.r.Decorated() {
		left = int(w.extents[0])
		top = int(w.extents[2])
	}
	w.extentsAccess.RUnlock()

	windowX -= left
	windowY -= top

	values := []uint32{uint32(windowX), uint32(windowY)}
	xConnection.ConfigureWindow(w.getXWindow(), x11.CONFIG_WINDOW_X|x11.CONFIG_WINDOW_Y, values)
}

func (w *NativeWindow) configureSize(width, height int) {
	w.access.RLock()
	defer w.access.RUnlock()

	w.doConfigureSize(width, height)
}

func (w *NativeWindow) doConfigureSize(width, height int) {
	// We don't need to adjust width/height by window extents, because X
	// specifies width/height as proper *client area* size.

	xWindow := w.getXWindow()

	values := []uint32{uint32(width), uint32(height)}
	xConnection.ConfigureWindow(xWindow, x11.CONFIG_WINDOW_WIDTH|x11.CONFIG_WINDOW_HEIGHT, values)

	// For ICCCM complient WM's
	hints := new(x11.SizeHints)

	hints.Flags |= x11.SIZE_HINT_P_SIZE
	hints.Width = x11.Int32(width)
	hints.Height = x11.Int32(height)

	// Clamp window to maximum and minimum sizes if they are specified.
	maxWidth, maxHeight := w.r.MaximumSize()
	minWidth, minHeight := w.r.MinimumSize()

	if minWidth != 0 && minHeight != 0 {
		hints.Flags |= x11.SIZE_HINT_P_MIN_SIZE
		hints.MinWidth = x11.Int32(minWidth)
		hints.MinHeight = x11.Int32(minHeight)
	}

	if maxWidth != 0 && maxHeight != 0 {
		hints.Flags |= x11.SIZE_HINT_P_MAX_SIZE
		hints.MaxWidth = x11.Int32(maxWidth)
		hints.MaxHeight = x11.Int32(maxHeight)
	}

	// This is our best bet for getting a specific aspect ratio on a window.
	//
	// Sadly, at least Unity WM does not seem to respect this. We *could* force
	// users to not be able to resize past aspect ratio in ConfigureNotifyEvent
	// but we might risk hitting a recursive descent or an undefined behavior
	// on other WM's, and it's not very pretty either, as such we don't do it.
	aspectRatio := w.r.AspectRatio()
	if aspectRatio != 0 {
		hints.Flags |= x11.SIZE_HINT_P_ASPECT
		fWidth := aspectRatio * 1000.0
		fHeight := fWidth * (1.0 / aspectRatio)
		hints.MinAspectNum = x11.Int32(fWidth)
		hints.MinAspectDen = x11.Int32(fHeight)
		hints.MaxAspectNum = x11.Int32(fWidth)
		hints.MaxAspectDen = x11.Int32(fHeight)
	}

	xConnection.SetWmSizeHints(xWindow, hints)
}

func (w *NativeWindow) PixelClear(rect image.Rectangle) {
	img := image.NewRGBA(rect)
	w.PixelBlit(uint(rect.Min.X), uint(rect.Min.Y), img)
}

func (w *NativeWindow) PixelBlit(x, y uint, img *image.RGBA) {
	if !w.r.Opened() {
		return
	}

	w.access.RLock()
	if w.activeRenderMode != renderModePixelBlit {
		w.access.RUnlock()

		// Enter write lock
		w.access.Lock()

		w.activeRenderMode = renderModePixelBlit
		w.doRebuildWindow()

		w.access.Unlock()
	} else {
		w.access.RUnlock()
	}

	if w.xGC == 0 {
		return
	}

	chosenFmt := w.pixmapFmt24
	if w.r.Transparent() {
		chosenFmt = w.pixmapFmt32
	}
	if chosenFmt == nil {
		return
	}
	defer xConnection.Flush()

	bounds := img.Bounds()
	sz := bounds.Size()
	if sz.X <= 0 && sz.Y <= 0 {
		return
	}
	width := uint16(sz.X)
	height := uint16(sz.Y)

	// Convert to RGBA
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)

	// RGBA to BGRA
	for x := 0; x < sz.X; x++ {
		for y := 0; y < sz.Y; y++ {
			p := (y-rgba.Rect.Min.Y)*rgba.Stride + (x-rgba.Rect.Min.X)*4
			r := rgba.Pix[p]
			g := rgba.Pix[p+1]
			b := rgba.Pix[p+2]
			a := rgba.Pix[p+3]

			rgba.Pix[p] = b
			rgba.Pix[p+1] = g
			rgba.Pix[p+2] = r
			rgba.Pix[p+3] = a
		}
	}

	setup := xConnection.GetSetup()
	x11Image := x11.ImageCreate(
		width,
		height,
		x11.IMAGE_FORMAT_Z_PIXMAP,
		uint8(chosenFmt.ScanlinePad),
		uint8(chosenFmt.Depth),
		uint8(chosenFmt.BitsPerPixel),
		32,
		int(setup.ImageByteOrder),
		x11.IMAGE_ORDER_LSB_FIRST,
		unsafe.Pointer(&rgba.Pix[0]),
		uint32(len(rgba.Pix)),
		&rgba.Pix[0],
	)

	nativeImage := xConnection.ImageNative(x11Image, true)
	xConnection.ImagePut(x11.Drawable(w.getXWindow()), w.xGC, nativeImage, int16(x), int16(y), 0)
}

func (w *NativeWindow) setTitle(title string) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetTitle(title)
}

func (w *NativeWindow) doSetTitle(title string) {
	defer xConnection.Flush()

	xWindow := w.getXWindow()

	if len(title) == 0 {
		// Some WM's don't like empty title strings, so we can just use a space..
		title = " "
	}

	titleBytes := []byte(title)
	cTitle := unsafe.Pointer(&titleBytes[0])

	// _NET_WM_NAME is used by EWMH compliant WM's for UTF-8 encoded window
	// titles.
	xConnection.ChangeProperty(
		x11.PROP_MODE_REPLACE,
		xWindow,
		aNetWmName,
		aUtf8String,
		8, uint32(len(title)),
		cTitle,
	)

	// WM_NAME is used for non-EWMH compliant WM's but only support ? LATIN-1
	// or ASCII encoding, I don't know. Most WM's are EWMH compliant though.
	xConnection.ChangeProperty(
		x11.PROP_MODE_REPLACE,
		xWindow,
		x11.ATOM_WM_NAME,
		x11.ATOM_STRING,
		8, uint32(len(title)),
		cTitle,
	)
}

func (w *NativeWindow) setVisible(visible bool) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetVisible(visible)
}

func (w *NativeWindow) doSetVisible(visible bool) {
	// Some window managers will allow us to specify window position and size
	// *before* being mapped, and then the window won't flicker/jump on the
	// user's screen.
	//
	// But sadly other window managers will only respect window positions and
	// sizes *after* being mapped. We try to get the absolute best of both
	// worlds by sending the position requests twice.

	xWindow := w.getXWindow()

	if visible {
		// Set the window position and size, before the window is mapped.
		windowX, windowY := w.r.Position()
		w.doConfigurePosition(windowX, windowY)
		width, height := w.r.Size()
		w.doConfigureSize(width, height)

		// Required before mapping.
		w.doUpdateNetWmState()
		w.doUpdateMotifWMHints(w.r.Decorated())
		w.doUpdateNetWmIcon()

		// Map the window.
		xConnection.MapWindow(xWindow)

		// Now set the window position and size, after it has been mapped.
		w.doConfigurePosition(windowX, windowY)
		w.doConfigureSize(width, height)
		xConnection.Flush()

		// Wait for an map notify event (or timeout if the server does not send it)
		select {
		case <-w.waitForMap:
			break
		case <-time.After(1 * time.Second):
			logger().Println("Timed out waiting for MapNotifyEvent.")
			break
		}

		// If the window is minimized, minimized it now.
		//
		// NOTE: restoring a minimized window involves calling this setVisible
		// function, so be careful here!
		if w.r.Minimized() {
			w.r.doSetMinimized(true)
		}

		// Lastly, we should set our mouse cursor
		w.doSetCursor(w.r.Cursor())

		// Since we need to switch to an invisible cursor if the cursor is
		// currently grabbed.
		w.doSetCursorGrabbed(w.r.CursorGrabbed(), false)

		xConnection.Flush()

		// Start receiving position and size events now that configurePosition
		// is completed.
		w.can.Lock()
		w.can.sendPositionEvents = true
		w.can.sendSizeEvents = true
		w.can.Unlock()

	} else {
		// Stop receiving position and size events as some invalid ones come in
		// after our UnmapWindow request.
		w.can.Lock()
		w.can.sendPositionEvents = false
		w.can.sendSizeEvents = false
		w.can.Unlock()

		err := xConnection.RequestCheck(xConnection.UnmapWindow(xWindow))
		if err != nil {
			logger().Println("UnmapWindow", err)
		}
		xConnection.Flush()
	}
}

func (w *NativeWindow) setTransparent(transparent bool) {
	w.access.Lock()
	defer w.access.Unlock()

	if transparent && w.xDepth == 32 {
		// We already have a 32-bit depth visual
		return
	}

	if !transparent && w.xDepth == 24 {
		// We already have a 24-bit depth visual
		return
	}

	w.doFindProperVisual()
	w.doRebuildWindow()
}

func (w *NativeWindow) doUpdateMotifWMHints(decorated bool) {
	hints := &x11.MotifWMHints{
		Flags:       2, // Changing decorations
		Decorations: 0,
	}
	if decorated {
		hints.Decorations = 1
	}

	err := xConnection.RequestCheck(xConnection.ChangePropertyChecked(
		x11.PROP_MODE_REPLACE,
		w.getXWindow(),
		aMotifWmHints,
		aMotifWmHints,
		8, uint32(unsafe.Sizeof(*hints)),
		unsafe.Pointer(hints),
	))
	if err != nil {
		logger().Println("ChangeProperty(_MOTIF_WM_HINTS):", err)
	}
}

func (w *NativeWindow) setDecorated(decorated bool) {
	w.access.Lock()
	defer w.access.Unlock()

	defer xConnection.Flush()

	x, y := w.r.Position()

	w.can.Lock()
	w.can.sendPositionEvents = false
	w.can.sendSizeEvents = false
	w.can.Unlock()

	// Update motif window hints
	w.doUpdateMotifWMHints(decorated)

	// Hopefully this covers most WM's okay.
	timeout := time.After(750 * time.Millisecond)
l:
	for {
		select {
		case <-timeout:
			break l
		default:
			time.Sleep(60 * time.Millisecond)
			w.doConfigurePosition(x, y)
			xConnection.Flush()
		}
	}

	w.can.Lock()
	w.can.sendPositionEvents = true
	w.can.sendSizeEvents = true
	w.can.Unlock()
}

func (w *NativeWindow) setAspectRatio(aspectRatio float32) {
	defer xConnection.Flush()
	width, height := w.r.clampedSize()
	w.configureSize(width, height)
}

func (w *NativeWindow) setSize(width, height int) {
	defer xConnection.Flush()
	w.configureSize(width, height)
}

func (w *NativeWindow) setMinimumSize(x, y int) {
	defer xConnection.Flush()
	width, height := w.r.clampedSize()
	w.configureSize(width, height)
}

func (w *NativeWindow) setMaximumSize(x, y int) {
	defer xConnection.Flush()
	width, height := w.r.clampedSize()
	w.configureSize(width, height)
}

func (w *NativeWindow) setCursorGrabbed(grabbed bool) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetCursorGrabbed(grabbed, true)
	xConnection.Flush()
}

func (w *NativeWindow) doSetCursorGrabbed(grabbed, restorePosition bool) {
	if grabbed {
		xWindow := w.getXWindow()
		cookie := xConnection.GrabPointer(
			1,                   // boolean owner events
			xWindow,             // grab events
			0,                   // event mask (0 because owner events is on)
			x11.GRAB_MODE_ASYNC, // pointer mode
			x11.GRAB_MODE_ASYNC, // keyboard mode
			xWindow,             // confine to
			0,
			x11.TIME_CURRENT_TIME,
		)
		_, err := xConnection.GrabPointerReply(cookie)
		if err != nil {
			logger().Println(err)
		}

		w.doSetCursor(clearCursor)
		w.can.Lock()
		w.can.sendRelativeMoveEvents = false
		w.can.Unlock()

	} else {
		xConnection.UngrabPointer(x11.TIME_CURRENT_TIME)

		if restorePosition {
			// Restore cursor to the original position it was in before the grab.
			x, y := w.r.CursorPosition()
			if x != 0 && y != 0 {
				w.warpPointer(x, y)
			}
		}

		w.doSetCursor(w.r.Cursor())
	}
}

func (w *NativeWindow) setPosition(x, y int) {
	defer xConnection.Flush()
	w.configurePosition(x, y)
}

func (w *NativeWindow) setCursorPosition(x, y int) {
	if w.r.CursorGrabbed() && w.r.Focused() {
		return
	}
	w.warpPointer(x, y)
}

func (w *NativeWindow) clearLastCursorPosition() {
	w.last.Lock()
	w.last.cursorX = -1
	w.last.cursorY = -1
	w.last.Unlock()
}

func (w *NativeWindow) warpPointer(x, y int) {
	w.last.Lock()
	w.last.cursorX = x
	w.last.cursorY = y
	w.last.Unlock()
	xConnection.WarpPointer(0, w.getXWindow(), 0, 0, 0, 0, int16(x), int16(y))
	xConnection.Flush()
}

func (w *NativeWindow) setCursor(cursor *Cursor) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetCursor(cursor)
}

func (w *NativeWindow) doSetCursor(cursor *Cursor) {
	var x11Cursor x11.Cursor

	if cursor != nil {
		// Prepare the cursor
		lc := w.doPrepareCursor(cursor)
		if lc == nil {
			// Couldn't prepare it
			return
		}
		x11Cursor = lc.cursor
	}

	// Set the cursor
	xConnection.ChangeWindowAttributes(w.getXWindow(), x11.CW_CURSOR, (*uint32)(unsafe.Pointer(&x11Cursor)))
	xConnection.Flush()
}

func (w *NativeWindow) freeCursor(cursor *Cursor) {
	w.access.Lock()
	defer w.access.Unlock()

	lc, ok := w.cursors[cursor]
	if !ok {
		return
	}

	xConnection.FreeCursor(lc.cursor)
	xConnection.Flush()

	delete(w.cursors, cursor)
}

func (w *NativeWindow) doPrepareCursor(cursor *Cursor) *loadedCursor {
	if cursor.Image == nil {
		return nil
	}

	lc, ok := w.cursors[cursor]
	if ok {
		return &lc
	}

	if w.xGC == 0 {
		return nil
	}

	chosenFmt := w.pixmapFmt32
	if chosenFmt == nil {
		return nil
	}
	defer xConnection.Flush()

	bounds := cursor.Image.Bounds()
	sz := bounds.Size()
	if sz.X <= 0 && sz.Y <= 0 {
		return nil
	}

	width := uint16(sz.X)
	height := uint16(sz.Y)

	// Locate picture formats
	cookie := xConnection.RenderQueryPictFormats()
	reply, err := xConnection.RenderQueryPictFormatsReply(cookie)
	if err != nil {
		logger().Println("RenderQueryPictFormats():", err)
		return nil
	}

	// Find ARGB32 format
	pictFormat := x11.RenderUtilFindStandardFormat(reply, x11.PICT_STANDARD_ARGB_32)

	// Convert to RGBA
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), cursor.Image, bounds.Min, draw.Src)

	// RGBA to BGRA
	for x := 0; x < sz.X; x++ {
		for y := 0; y < sz.Y; y++ {
			p := (y-rgba.Rect.Min.Y)*rgba.Stride + (x-rgba.Rect.Min.X)*4
			r := rgba.Pix[p]
			g := rgba.Pix[p+1]
			b := rgba.Pix[p+2]
			a := rgba.Pix[p+3]

			rgba.Pix[p] = b
			rgba.Pix[p+1] = g
			rgba.Pix[p+2] = r
			rgba.Pix[p+3] = a
		}
	}

	// Create an image
	setup := xConnection.GetSetup()
	x11Image := x11.ImageCreate(
		width,
		height,
		x11.IMAGE_FORMAT_Z_PIXMAP,
		uint8(chosenFmt.ScanlinePad),
		uint8(chosenFmt.Depth),
		uint8(chosenFmt.BitsPerPixel),
		32,
		int(setup.ImageByteOrder),
		x11.IMAGE_ORDER_LSB_FIRST,
		unsafe.Pointer(&rgba.Pix[0]),
		uint32(len(rgba.Pix)),
		&rgba.Pix[0],
	)

	// Convert to native format
	nativeImage := xConnection.ImageNative(x11Image, true)

	// Locate the screen
	xScreen := xConnection.ScreenOfDisplay(w.r.Screen().NativeScreen.xScreen)

	// Create a pixmap
	pixmap := x11.Pixmap(xConnection.GenerateId())
	err = xConnection.RequestCheck(xConnection.CreatePixmapChecked(
		uint8(chosenFmt.Depth),
		pixmap,
		x11.Drawable(xScreen.Root),
		width,
		height,
	))
	if err != nil {
		logger().Println("x11_create_pixmap():", err)
		return nil
	}
	defer xConnection.FreePixmap(pixmap)

	// Create a GC for the pixmap on the screen
	cursorGC := x11.GContext(xConnection.GenerateId())
	err = xConnection.RequestCheck(xConnection.CreateGCChecked(cursorGC, x11.Drawable(pixmap), 0, nil))
	if err != nil {
		logger().Println("can't prepare cursor; x11_create_gc():", err)
		return nil
	}
	defer xConnection.FreeGC(cursorGC)

	// Put into pixmap
	xConnection.ImagePut(x11.Drawable(pixmap), cursorGC, nativeImage, 0, 0, 0)

	// Create a picture
	pict := x11.RenderPicture(xConnection.GenerateId())
	err = xConnection.RequestCheck(xConnection.RenderCreatePictureChecked(
		pict,
		x11.Drawable(pixmap),
		pictFormat.Id,
		0,
		nil,
	))
	if err != nil {
		logger().Println("x11_render_create_picture():", err)
		return nil
	}
	defer xConnection.RenderFreePicture(pict)

	// Create the cursor
	lc.cursor = x11.Cursor(xConnection.GenerateId())
	err = xConnection.RequestCheck(xConnection.RenderCreateCursorChecked(
		lc.cursor,
		pict,
		uint16(cursor.X),
		uint16(cursor.Y),
	))
	if err != nil {
		logger().Println("x11_render_create_cursor():", err)
		return nil
	}

	// Store the cursor
	w.cursors[cursor] = lc
	return &lc
}

func (w *NativeWindow) prepareCursor(cursor *Cursor) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doPrepareCursor(cursor)
}

func (w *NativeWindow) setMinimized(minimized bool) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetMinimized(minimized)
}

func (w *NativeWindow) doSetMinimized(minimized bool) {
	if minimized {
		atoms := make([]x11.Atom, 5)
		atoms[0] = x11.WM_STATE_ICONIC

		xScreen := xConnection.ScreenOfDisplay(w.r.Screen().NativeScreen.xScreen)

		cme := new(x11.ClientMessageEvent)
		cme.ResponseType = x11.CLIENT_MESSAGE
		cme.Format = 32
		cme.Window = w.getXWindow()
		cme.Type = aWmChangeState
		cme.Data = *(*[20]byte)(unsafe.Pointer(&atoms[0]))
		xConnection.SendEvent(
			false, //propagate
			xScreen.Root,
			x11.EVENT_MASK_STRUCTURE_NOTIFY|x11.EVENT_MASK_SUBSTRUCTURE_REDIRECT,
			unsafe.Pointer(cme),
		)

		xConnection.Flush()
	} else {
		// Map the window
		w.doSetVisible(w.r.Visible())
	}
}

func (w *NativeWindow) doUpdateNetWmIcon() {
	iconImage := w.r.Icon()
	if iconImage == nil {
		return
	}
	icon, ok := iconImage.(*image.RGBA)
	if !ok {
		bounds := iconImage.Bounds()
		icon = image.NewRGBA(bounds)
		draw.Draw(icon, icon.Bounds(), iconImage, bounds.Min, draw.Src)
	}

	if w.xicon == nil {
		put := func(img *image.RGBA) {
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			buf := make([]uint32, width*height)
			for x := 0; x < width; x++ {
				for y := 0; y < width; y++ {
					r, g, b, a := img.At(x, y).RGBA()
					c := (uint32(a>>8) << 24) | (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | (uint32(b>>8) << 0)
					buf[(y*width)+x] = c
				}
			}

			w.xicon = append(w.xicon, uint32(width))
			w.xicon = append(w.xicon, uint32(height))
			w.xicon = append(w.xicon, buf...)
		}

		size := func(img image.Image, w, h int) *image.RGBA {
			bounds := img.Bounds()
			if bounds.Dx() < w && bounds.Dy() < h {
				// Resample is cheaper
				return resize.Resample(img, bounds, w, h).(*image.RGBA)
			}
			return resize.Resize(img, bounds, w, h).(*image.RGBA)
		}

		// Silly, most standard window managers simply ignore high-resolution
		// icons, and instead opt for something like high-resolution icons being
		// specified in .desktop files. For instance see:
		//
		// http://askubuntu.com/questions/90845/pygtk-application-icon-blurred-in-unity
		//
		// We still provide 128x128, 64x64, 48x48, 32x32, and 16x16 icons just in
		// case some WM will use them. We also put them in that order (largest
		// first) in case some WM might assume the first to be the *best* naively.

		// Whatever size our specified icon is (E.g. 256x256px).
		put(icon)

		// Now other sizes.
		bounds := iconImage.Bounds()
		iw, ih := bounds.Dx(), bounds.Dy()
		if iw != 128 || ih != 128 {
			put(size(icon, 128, 128))
		}
		if iw != 64 || ih != 64 {
			put(size(icon, 64, 64))
		}
		if iw != 48 || ih != 48 {
			put(size(icon, 48, 48))
		}
		if iw != 32 || ih != 32 {
			put(size(icon, 32, 32))
		}
		if iw != 16 || ih != 16 {
			put(size(icon, 16, 16))
		}
	}

	xConnection.ChangeProperty(
		x11.PROP_MODE_REPLACE,
		w.getXWindow(),
		aNetWmIcon,
		x11.ATOM_CARDINAL,
		32,
		uint32(len(w.xicon)),
		unsafe.Pointer(&w.xicon[0]),
	)
}

func (w *NativeWindow) setIcon(img image.Image) {
	w.access.Lock()
	defer w.access.Unlock()

	w.doSetIcon(img)
}

func (w *NativeWindow) doSetIcon(img image.Image) {
	w.xicon = nil
	w.doUpdateNetWmIcon()
	xConnection.Flush()
}

func (w *NativeWindow) doUpdateNetWmState() {
	var atoms []x11.Atom

	if w.r.Fullscreen() {
		atoms = append(atoms, aNetWmStateFullscreen)
	}

	if w.r.Maximized() {
		atoms = append(atoms, aNetWmStateMaximizedVert)
		atoms = append(atoms, aNetWmStateMaximizedHorz)
	}

	if w.r.AlwaysOnTop() {
		atoms = append(atoms, aNetWmStateAbove)
	}

	var catoms unsafe.Pointer
	if len(atoms) > 0 {
		catoms = unsafe.Pointer(&atoms[0])
	}
	xConnection.ChangeProperty(
		x11.PROP_MODE_REPLACE,
		w.getXWindow(),
		aNetWmState,
		x11.ATOM_ATOM,
		32,
		uint32(len(atoms)),
		catoms,
	)
}

func (w *NativeWindow) netWMState(set bool, prop x11.Atom) {
	var action x11.Atom
	if set {
		action = 1 //_NET_WM_STATE_ADD
	} else {
		action = 0 //_NET_WM_STATE_REMOVE
	}

	atoms := []x11.Atom{
		action,
		0,
		prop,
		1, // source indication: normal application
		0, // must be here or cme.Data below is invalid pointer (as 20 / sizeof(Atom) == 5)
	}
	xScreen := xConnection.ScreenOfDisplay(w.r.Screen().NativeScreen.xScreen)

	cme := new(x11.ClientMessageEvent)
	cme.ResponseType = x11.CLIENT_MESSAGE
	cme.Format = 32
	cme.Window = w.getXWindow()
	cme.Type = aNetWmState
	cme.Data = *(*[20]byte)(unsafe.Pointer(&atoms[0]))
	xConnection.SendEvent(
		false, //propagate
		xScreen.Root,
		x11.EVENT_MASK_STRUCTURE_NOTIFY|x11.EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		unsafe.Pointer(cme),
	)
}

func (w *NativeWindow) setFullscreen(fullscreen bool) {
	w.netWMState(fullscreen, aNetWmStateFullscreen)
	xConnection.Flush()
}

func (w *NativeWindow) setMaximized(maximized bool) {
	w.netWMState(maximized, aNetWmStateMaximizedVert)
	w.netWMState(maximized, aNetWmStateMaximizedHorz)
	xConnection.Flush()
}

func (w *NativeWindow) setAlwaysOnTop(alwaysOnTop bool) {
	w.netWMState(alwaysOnTop, aNetWmStateAbove)
	xConnection.Flush()
}

func (w *NativeWindow) notify() {
	w.netWMState(true, aNetWmStateDemandsAttention)
	xConnection.Flush()
}

func (w *NativeWindow) doDestroy() {
	w.xWindowAccess.Lock()
	xWindow := w.xWindow
	w.xWindow = 0
	w.xWindowAccess.Unlock()

	xWindowLookupAccess.Lock()
	delete(xWindowLookup, xWindow)
	xWindowLookupAccess.Unlock()
	xConnection.DestroyWindow(xWindow)
	xConnection.FreeGC(w.xGC)
	xConnection.Flush()
}

func (w *NativeWindow) destroy() {
	w.access.Lock()
	w.doDestroy()
	w.access.Unlock()
}

func newNativeWindow(real *Window) *NativeWindow {
	w := new(NativeWindow)
	w.r = real
	w.waitForMap = make(chan bool)
	w.waitForUnmap = make(chan bool)
	w.waitForFrameExtents = make(chan bool)
	w.waitForMotifHints = make(chan bool)
	w.waitForNetWmAllowedActions = make(chan bool)
	w.cursors = make(map[*Cursor]loadedCursor)
	return w
}
