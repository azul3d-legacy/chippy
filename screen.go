// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"sort"
	"sync"
)

// Screen represents an single physical screen device. It is only possible to
// get an screen from either the Screens or DefaultScreen functions.
type Screen struct {
	*NativeScreen

	access                        sync.RWMutex
	name                          string
	physicalWidth, physicalHeight float32
	modes                         sortedScreenModes
	mode, originalMode            *ScreenMode
}

func newScreen(name string, physicalWidth, physicalHeight float32, modes []*ScreenMode, currentMode *ScreenMode) *Screen {
	s := new(Screen)
	s.NativeScreen = newNativeScreen()
	s.name = name
	s.physicalWidth = physicalWidth
	s.physicalHeight = physicalHeight
	s.modes = sortedScreenModes(modes)
	s.mode = currentMode
	s.originalMode = currentMode
	sort.Sort(s.modes)
	return s
}

// String returns an nice string representation of this Screen
func (s *Screen) String() string {
	w, h := s.PhysicalSize()
	return fmt.Sprintf("Screen(\"%s\", %.0f by %.0fmm)", s.Name(), w, h)
}

// Equals compares this screen with the other screen, determining whether or
// not they are the same physical screen.
func (s *Screen) Equals(other *Screen) bool {
	if s.String() == other.String() {
		otherModes := other.Modes()
		for i, mode := range s.Modes() {
			if !mode.Equals(otherModes[i]) {
				return false
			}
		}
		return true
	}

	return false
}

// Name returns an formatted string of the screens name, this is something that
// the user should be able to relate on their own to the actual physical screen
// device, this typically includes device brand name or model etc..
func (s *Screen) Name() string {
	return s.name
}

// PhysicalSize returns the physical width and height of this screen, in
// millimeters, or zero as both width and height in the event there is no way
// to determine the physical size of this screen.
func (s *Screen) PhysicalSize() (width float32, height float32) {
	return s.physicalWidth, s.physicalHeight
}

// OriginalMode returns the original screen mode of this screen, as it was when
// this screen was created.
func (s *Screen) OriginalMode() *ScreenMode {
	return s.originalMode
}

// Modes returns all available screen modes on this screen, sorted by highest
// resolution, then highest bytes per pixel, then highest refresh rate.
func (s *Screen) Modes() []*ScreenMode {
	cpy := make([]*ScreenMode, len(s.modes))
	copy(cpy, s.modes)
	return cpy
}

// SetMode switches this screen to the specified screen mode, or returns an
// error in the event that we where unable to switch the screen to the specified
// screen mode.
//
// The newMode parameter must be an screen mode that originally came from one
// of the methods Modes(), Mode(), or OriginalMode(), or else this function
// will panic.
//
// If an error is returned, it will be either ErrBadScreenMode, or
// ErrDualViewCapable.
func (s *Screen) SetMode(newMode *ScreenMode) error {
	s.access.Lock()
	defer s.access.Unlock()

	if s.mode.Equals(newMode) {
		// We're already using this mode -- avoid flicker etc.
		return nil
	}

	s.access.Unlock()
	err := s.NativeScreen.setMode(newMode)
	s.access.Lock()

	s.mode = newMode
	return err
}

// Mode returns the current screen mode in use by this screen, this will be
// either the last screen mode set via SetMode, or the original screen mode
// from OriginalMode in the event that no screen mode was previously set on
// this screen.
func (s *Screen) Mode() *ScreenMode {
	s.access.RLock()
	defer s.access.RUnlock()

	return s.mode
}

// Restore is short-hand for:
//
//  s.SetMode(s.OriginalMode())
//
func (s *Screen) Restore() {
	s.SetMode(s.OriginalMode())
}

var (
	screenCacheLock       sync.RWMutex
	cachedScreens         []*Screen
	cachedDefaultScreen   *Screen
	queriedScreensAlready bool
)

// Screens returns all available, attached, and activated screens possible.
// Once this function is called, the result is cached such that future calls to
// this function are faster and return the cached result.
//
// To update the internal screen cache, see the RefreshScreens function.
func Screens() []*Screen {
	screenCacheLock.RLock()
	defer screenCacheLock.RUnlock()

	if !queriedScreensAlready {
		screenCacheLock.RUnlock()
		RefreshScreens()
		screenCacheLock.RLock()
	}
	return cachedScreens
}

// RefreshScreens queries for all available screens, and updates the internal
// cache returned by the Screens() function, such that the Screens() function
// returns new set of attatched/detatched Screen devices.
func RefreshScreens() {
	screenCacheLock.Lock()
	defer screenCacheLock.Unlock()

	queriedScreensAlready = true
	cachedScreens = backend_Screens()
	cachedDefaultScreen = backend_DefaultScreen()
}

// DefaultScreen returns the 'default' screen, this is determined by either the
// window manager itself (as per an user's personal setup and configuration) or
// will be best-guessed by Chippy.
//
// It is possible for this function to return nil, in the highly unlikely event
// that Screens() returns no screens at all, due to an user having none plugged
// in or activated.
func DefaultScreen() *Screen {
	screenCacheLock.RLock()
	defer screenCacheLock.RUnlock()

	if !queriedScreensAlready {
		screenCacheLock.RUnlock()
		RefreshScreens()
		screenCacheLock.RLock()
	}

	return cachedDefaultScreen
}
