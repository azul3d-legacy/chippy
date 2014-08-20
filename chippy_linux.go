// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"errors"
	"fmt"
	"image"
	"sync"
	"unsafe"

	"azul3d.org/chippy.v1/internal/x11"
	"azul3d.org/chippy.v1/internal/xkbcommon"
)

var (
	xDisplay, glxDisplay     *x11.Display
	xConnection              *x11.Connection
	xDisplayName             string
	xrandrMajor, xrandrMinor int
	xinputMajor, xinputMinor int
	xDefaultScreenNumber     int
	clearCursor              *Cursor

	// There is no mention of thread-safety so we can only assume that context
	// probably does not use TLS but doesn't provide synchronization, so we
	// just use a simple mutex for synchronization.
	xkbContext struct {
		sync.Mutex
		ctx *xkbcommon.XkbContext
	}
	xkbDevice    int32
	xkbKeymap    *xkbcommon.XkbKeymap
	xkbState     *xkbcommon.XkbState
	xkbBaseEvent uint8

	ErrInvalidGLXVersion = errors.New("GLX version 1.4 is required but not available.")
	ErrInvalidXKBVersion = errors.New(fmt.Sprintf("XKB version %d.%d is required but not available.", xkbMinMajor, xkbMinMinor))
)

const (
	// We at least need 1.2 for multiple monitor support, etc, etc..
	xrandrMinMajor = 1
	xrandrMinMinor = 2

	// We need Xinput2 for raw mouse input
	xinputMinMajor = 2
	xinputMinMinor = 0

	xkbMinMajor = xkbcommon.XKB_X11_MIN_MAJOR_XKB_VERSION
	xkbMinMinor = xkbcommon.XKB_X11_MIN_MINOR_XKB_VERSION

	// We need GLX 1.4 for multisampling
	glxMinMajor = 1
	glxMinMinor = 4
)

// Converts an *x11.Connection to an *xkbcommon.Connection type.
func xkbcommonConn(c *x11.Connection) *xkbcommon.Connection {
	return (*xkbcommon.Connection)(unsafe.Pointer(c))
}

var (
	// EWMH atoms
	aNetRequestFrameExtents, aNetFrameExtents, aNetWmName, aNetWmState,
	aNetWmStateFullscreen, aNetWmStateAbove, aNetWmStateMaximizedHorz,
	aNetWmStateMaximizedVert, aNetWmStateDemandsAttention, aNetWmIcon,
	aUtf8String x11.Atom

	// MOTIF atoms
	aMotifWmHints x11.Atom

	// WM atoms
	aWmProtocols, aWmDeleteWindow, aWmChangeState x11.Atom
)

func initAtoms() {
	aNetRequestFrameExtents = xConnection.InternAtom(false, "_NET_REQUEST_FRAME_EXTENTS")
	aNetFrameExtents = xConnection.InternAtom(false, "_NET_FRAME_EXTENTS")
	aNetWmName = xConnection.InternAtom(false, "_NET_WM_NAME")
	aNetWmState = xConnection.InternAtom(false, "_NET_WM_STATE")
	aNetWmStateFullscreen = xConnection.InternAtom(false, "_NET_WM_STATE_FULLSCREEN")
	aNetWmStateAbove = xConnection.InternAtom(false, "_NET_WM_STATE_ABOVE")
	aNetWmStateMaximizedHorz = xConnection.InternAtom(false, "_NET_WM_STATE_MAXIMIZED_HORZ")
	aNetWmStateMaximizedVert = xConnection.InternAtom(false, "_NET_WM_STATE_MAXIMIZED_VERT")
	aNetWmStateDemandsAttention = xConnection.InternAtom(false, "_NET_WM_STATE_DEMANDS_ATTENTION")
	aNetWmIcon = xConnection.InternAtom(false, "_NET_WM_ICON")
	aUtf8String = xConnection.InternAtom(false, "UTF8_STRING")

	aMotifWmHints = xConnection.InternAtom(false, "_MOTIF_WM_HINTS")

	aWmProtocols = xConnection.InternAtom(false, "WM_PROTOCOLS")
	aWmDeleteWindow = xConnection.InternAtom(false, "WM_DELETE_WINDOW")
	aWmChangeState = xConnection.InternAtom(false, "WM_CHANGE_STATE")
}

func refreshKeyboardMapping() error {
	// Create keymap from the device.
	xkbContext.Lock()
	xkbKeymap = xkbcommon.XkbX11KeymapNewFromDevice(xkbContext.ctx, xkbcommonConn(xConnection), xkbDevice, 0)
	xkbContext.Unlock()
	if xkbKeymap == nil {
		return errors.New("XKB-common could not create keymap from device.")
	}

	// Create keyboard state manager from keymap and device.
	xkbContext.Lock()
	xkbState = xkbcommon.XkbX11StateNewFromDevice(xkbcommonConn(xConnection), xkbKeymap, xkbDevice)
	xkbContext.Unlock()
	if xkbState == nil {
		return errors.New("XKB-common could not create keyboard state manager from device.")
	}
	return nil
}

// SetDisplayName sets the string that will be passed into XOpenDisplay; equivalent to the DISPLAY
// environment variable on posix complaint systems.
//
// If set, this is used in place of the default DISPLAY environment variable.
//
// This function is only available on Linux.
func SetDisplayName(displayName string) {
	globalLock.Lock()
	defer globalLock.Unlock()
	xDisplayName = displayName
}

// DisplayName returns the display_name string, as it was passed into SetDisplayName.
//
// This function is only available on Linux.
func DisplayName() string {
	globalLock.RLock()
	defer globalLock.RUnlock()
	return xDisplayName
}

/*
func xGenericErrorHandler(display *x11.Display, event *x11.XErrorEvent) {
	logger().Println("Xlib Error:", x11.XGetErrorText(display, event.Code()))
}
*/

var (
	xWindowLookupAccess sync.RWMutex
	xWindowLookup       = make(map[x11.Window]*NativeWindow, 1)
)

func findWindow(w x11.Window) (*NativeWindow, bool) {
	xWindowLookupAccess.RLock()
	defer xWindowLookupAccess.RUnlock()

	nw, ok := xWindowLookup[w]
	return nw, ok
}

func xkbHandleEvent(e *x11.GenericEvent) {
	ev := (*x11.XkbAnyEvent)(unsafe.Pointer(e.EGenericEvent))
	if int32(ev.DeviceID) != xkbDevice {
		return
	}
	switch ev.XkbType {
	case x11.XKB_NEW_KEYBOARD_NOTIFY:
		refreshKeyboardMapping()

	case x11.XKB_MAP_NOTIFY:
		refreshKeyboardMapping()

	case x11.XKB_STATE_NOTIFY:
		ev := (*x11.XkbStateNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
		xkbContext.Lock()
		xkbState.UpdateMask(
			xkbcommon.XkbModMask(ev.BaseMods()),
			xkbcommon.XkbModMask(ev.LatchedMods()),
			xkbcommon.XkbModMask(ev.LockedMods()),
			xkbcommon.XkbLayoutIndex(ev.BaseGroup()),
			xkbcommon.XkbLayoutIndex(ev.LatchedGroup()),
			xkbcommon.XkbLayoutIndex(ev.LockedGroup()),
		)
		xkbContext.Unlock()
	}
}

var shutdownEventLoop = make(chan bool, 1)
var eventLoopReady = make(chan bool, 1)

func eventLoop() {
	readySent := false
	for {
		select {
		case <-shutdownEventLoop:
			return
		default:
			break
		}

		if !readySent {
			readySent = true
			eventLoopReady <- true
		}
		e := xConnection.WaitForEvent()
		if e == nil {
			// connection is closed
			return

		} else {
			switch e.ResponseType &^ 0x80 {
			case x11.KEY_PRESS:
				ev := (*x11.KeyPressEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.KEY_RELEASE:
				ev := (*x11.KeyReleaseEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.BUTTON_PRESS:
				ev := (*x11.ButtonPressEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.BUTTON_RELEASE:
				ev := (*x11.ButtonReleaseEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.MOTION_NOTIFY:
				ev := (*x11.MotionNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.ENTER_NOTIFY:
				ev := (*x11.EnterNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.LEAVE_NOTIFY:
				ev := (*x11.LeaveNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.FOCUS_IN:
				ev := (*x11.FocusInEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.FOCUS_OUT:
				ev := (*x11.FocusOutEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Event)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.EXPOSE:
				ev := (*x11.ExposeEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.VISIBILITY_NOTIFY:
				ev := (*x11.VisibilityNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.CREATE_NOTIFY:
				ev := (*x11.CreateNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.DESTROY_NOTIFY:
				ev := (*x11.DestroyNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			//case x11.MAPPING_NOTIFY:
			//ev := (*x11.MappingNotifyEvent)(unsafe.Pointer(e.EGenericEvent))

			case x11.CLIENT_MESSAGE:
				ev := (*x11.ClientMessageEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.PROPERTY_NOTIFY:
				ev := (*x11.PropertyNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.CONFIGURE_NOTIFY:
				ev := (*x11.ConfigureNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.REPARENT_NOTIFY:
				ev := (*x11.ReparentNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.MAP_NOTIFY:
				ev := (*x11.MapNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case x11.UNMAP_NOTIFY:
				ev := (*x11.UnmapNotifyEvent)(unsafe.Pointer(e.EGenericEvent))
				w, ok := findWindow(ev.Window)
				if ok {
					w.handleEvent(e, ev)
				}

			case 0:
				ev := (*x11.RequestError)(unsafe.Pointer(e.EGenericEvent))
				str := fmt.Sprintf("ErrorCode(%d)", ev.ErrorCode)
				switch ev.ErrorCode {
				case 1:
					str = "BadRequest"
				case 2:
					str = "BadValue"
				case 3:
					str = "BadWindow"
				case 4:
					str = "BadPixmap"
				case 5:
					str = "BadAtom"
				case 6:
					str = "BadCursor"
				case 7:
					str = "BadFont"
				case 8:
					str = "BadMatch"
				case 9:
					str = "BadDrawable"
				case 10:
					str = "BadAccess"
				case 11:
					str = "BadAlloc"
				case 12:
					str = "BadColor"
				case 13:
					str = "BadGC"
				case 14:
					str = "BadIDChoice"
				case 15:
					str = "BadName"
				case 16:
					str = "BadLength"
				case 17:
					str = "BadImplementation"
				}
				logger().Println("X Request Error:", str)
				logger().Printf("%+v\n", ev)

			default:
				if uint8(e.ResponseType) == xkbBaseEvent {
					// This is a XKB event.
					xkbHandleEvent(e)
				}

				// We still need to free the event.
				e.Free()
			}
		}
	}
}

func atLeastVersion(realMajor, realMinor, wantedMajor, wantedMinor int) bool {
	if realMajor != wantedMajor {
		return realMajor > wantedMajor
	}
	if realMinor != wantedMinor {
		return realMinor > wantedMinor
	}
	return true
}

func init() {
	x11.XInitThreads()
}

func backend_Init() (err error) {
	// It's not really safe to clear these at backend_Destroy() time.
	xDisplayName = ""
	xrandrMajor = 0
	xrandrMinor = 0

	if clearCursor == nil {
		clearCursor = &Cursor{
			Image: image.NewRGBA(image.Rect(0, 0, 16, 16)),
			X:     0,
			Y:     0,
		}
	}

	// We use a Xlib connection for GLX purely.
	glxDisplay = x11.XOpenDisplay(xDisplayName)
	if glxDisplay == nil {
		theLogger.Println("Unable to open X11 GLX display; Is the X server running?")
		return errors.New("Unable to open X11 GLX display; Is the X server running?")
	}

	x11.XSetErrorHandler(func(err string) {
		// Erm, if something below caused an error, we might hit a deadlock
		// here because we already have the lock used by logger().. so.. this
		// is kind of a cheap and easy fix (spawning a goroutine).
		go func() {
			logger().Println("X11 Error:", err)
		}()
	})

	// And we use an xcb connection for everything else.
	xDisplay = x11.XOpenDisplay(xDisplayName)
	if xDisplay == nil {
		theLogger.Println("Unable to open X11 XCB display; Is the X server running?")
		return errors.New("Unable to open X11 XCB display; Is the X server running?")
	}

	// We want XCB to own the event queue, not Xlib which does a poor job.
	xDisplay.XSetEventQueueOwner(x11.XCBOwnsEventQueue)
	xConnection = x11.XGetXCBConnection(xDisplay)
	xDefaultScreenNumber = xDisplay.XDefaultScreen()

	// Initialize atoms used
	initAtoms()

	// Setup Xkb-common extension, we need the xkbBaseEvent to identify XKB
	// events.
	var (
		ret                int
		xkbMajor, xkbMinor uint16
	)
	ret, xkbMajor, xkbMinor, xkbBaseEvent, _ = xkbcommon.XkbX11SetupXkbExtension(xkbcommonConn(xConnection), xkbMinMajor, xkbMinMinor, 0)
	if ret == 0 {
		theLogger.Printf("XKB version %d.%d exists, we require at least %d.%d\n", xkbMajor, xkbMinor, xkbMinMajor, xkbMinMinor)
		return ErrInvalidXKBVersion
	}

	// Create XKB context.
	xkbContext.Lock()
	xkbContext.ctx = xkbcommon.XkbContextNew(0)
	xkbContext.Unlock()
	if xkbContext.ctx == nil {
		return errors.New("XKB-common context could not be initialized!")
	}

	// Query the core/default keyboard device id.
	xkbDevice = xkbcommon.XkbX11GetCoreKeyboardDeviceId(xkbcommonConn(xConnection))
	if xkbDevice == -1 {
		return errors.New("XKB-common could not query core keyboard device")
	}

	// Refresh keyboard mapping.
	err = refreshKeyboardMapping()
	if err != nil {
		return err
	}

	// Select device events.
	var requiredEvents, requiredNknDetails, requiredMapParts, requiredStateDetails uint16
	requiredEvents |= x11.XKB_EVENT_TYPE_NEW_KEYBOARD_NOTIFY
	requiredEvents |= x11.XKB_EVENT_TYPE_MAP_NOTIFY
	requiredEvents |= x11.XKB_EVENT_TYPE_STATE_NOTIFY

	requiredNknDetails |= x11.XKB_NKN_DETAIL_KEYCODES

	requiredMapParts |= x11.XKB_MAP_PART_KEY_TYPES
	requiredMapParts |= x11.XKB_MAP_PART_KEY_SYMS
	requiredMapParts |= x11.XKB_MAP_PART_MODIFIER_MAP
	requiredMapParts |= x11.XKB_MAP_PART_EXPLICIT_COMPONENTS
	requiredMapParts |= x11.XKB_MAP_PART_KEY_ACTIONS
	requiredMapParts |= x11.XKB_MAP_PART_VIRTUAL_MODS
	requiredMapParts |= x11.XKB_MAP_PART_VIRTUAL_MOD_MAP

	requiredStateDetails |= x11.XKB_STATE_PART_MODIFIER_BASE
	requiredStateDetails |= x11.XKB_STATE_PART_MODIFIER_LATCH
	requiredStateDetails |= x11.XKB_STATE_PART_MODIFIER_LOCK
	requiredStateDetails |= x11.XKB_STATE_PART_GROUP_BASE
	requiredStateDetails |= x11.XKB_STATE_PART_GROUP_LATCH
	requiredStateDetails |= x11.XKB_STATE_PART_GROUP_LOCK

	details := &x11.XkbSelectEventDetails{
		AffectNewKeyboard:  requiredNknDetails,
		NewKeyboardDetails: requiredNknDetails,
		AffectState:        requiredStateDetails,
		StateDetails:       requiredStateDetails,
	}

	err = xConnection.RequestCheck(xConnection.XkbSelectEventsAuxChecked(
		x11.XkbDeviceSpec(xkbDevice), // deviceSpec
		requiredEvents,               // affectWhich
		0,                            // clear
		0,                            // selectAll
		requiredMapParts,             // affectMap
		requiredMapParts,             // map
		details,                      // details
	))
	if err != nil {
		logger().Println("XkbSelectEventsAux(): Failed to select keyboard events:", err)
		return err
	}

	// Ask to use detectable auto-repeat key events. This is more in sync with
	// the way that other platforms handle keyboard events and is not strange
	// to developers.
	wasSet, supported := xDisplay.XkbSetDetectableAutoRepeat(true)
	if !wasSet || !supported {
		logger().Println("XkbSetDetectableAutoRepeat(): cannot set or not supported.")
	}

	go eventLoop()
	<-eventLoopReady

	// See if we have xrandr support
	xrandrMajor = -1
	xrandrMinor = -1
	if xConnection.QueryExtension("RANDR") {
		reply, err := xConnection.RandrQueryVersionReply(xConnection.RandrQueryVersion(xrandrMinMajor, xrandrMinMinor))
		if err == nil {
			xrandrMajor = int(reply.MajorVersion)
			xrandrMinor = int(reply.MinorVersion)
		} else {
			theLogger.Println(err)
		}
	}

	// Tell what we're going to use
	if !atLeastVersion(xrandrMajor, xrandrMinor, xrandrMinMajor, xrandrMinMinor) {
		if xrandrMajor > 0 || xrandrMinor > 0 {
			theLogger.Printf("xrandr version %d.%d exists, we require at least %d.%d\n", xrandrMajor, xrandrMinor, xrandrMinMajor, xrandrMinMinor)
		} else {
			theLogger.Printf("xrandr extension is missing on X display.\n")
		}
		theLogger.Println("Falling back to pure X11; screen mode switching is impossible.")
	}

	var glxMajor, glxMinor x11.Int
	if !glxDisplay.GLXQueryVersion(&glxMajor, &glxMinor) || !atLeastVersion(int(glxMajor), int(glxMinor), glxMinMajor, glxMinMinor) {
		return ErrInvalidGLXVersion
	}

	return
}

func backend_Destroy() {
	shutdownEventLoop <- true
	xDisplay.Close()
	glxDisplay.Close()
	xConnection.Disconnect()
}
