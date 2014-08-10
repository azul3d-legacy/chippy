// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"runtime"
	"sync"

	"azul3d.org/chippy.v1/internal/win32"
)

type NativeScreenMode struct {
	w32Bpp  win32.DWORD
	w32Mode *win32.DEVMODE
}

func newNativeScreenMode() *NativeScreenMode {
	m := &NativeScreenMode{}
	return m
}

type NativeScreen struct {
	access sync.RWMutex

	isDefaultScreen                             bool
	w32MonitorDeviceName, w32GraphicsDeviceName string
	dc                                          win32.HDC
	w32Position                                 win32.RECT
}

func newNativeScreen() *NativeScreen {
	s := &NativeScreen{}

	runtime.SetFinalizer(s, func(f *NativeScreen) {
		go dispatch(func() {
			// Do screen related cleanup here..
			if !win32.DeleteDC(f.dc) {
				logger().Println("Cannot delete DC; DeleteDC():", win32.GetLastErrorString())
			}
		})
	})

	return s
}

func (s *NativeScreen) setMode(newMode *ScreenMode) (err error) {
	s.access.Lock()
	defer s.access.Unlock()

	dispatch(func() {
		mode := newMode.NativeScreenMode.w32Mode

		//mode.SetDmFields(win32.DM_PELSWIDTH & win32.DM_PELSHEIGHT & win32.DM_DISPLAYFREQUENCY)

		ret := win32.ChangeDisplaySettingsEx(s.w32GraphicsDeviceName, mode, win32.CDS_TEST, nil)
		if ret == win32.DISP_CHANGE_BADMODE {
			logger().Println("Unable to set screen mode; ChangeDisplaySettingsEx(,,CDS_TEST,) reports bad mode.")
			err = ErrBadScreenMode
			return
		}

		ret = win32.ChangeDisplaySettingsEx(s.w32GraphicsDeviceName, mode, 0, nil)

		if ret == win32.DISP_CHANGE_BADDUALVIEW {
			logger().Println("Unable to set screen mode; Because the system is DualView capable.")
			err = ErrDualViewCapable
			return
		}
		if ret == win32.DISP_CHANGE_BADMODE {
			logger().Println("Unable to set screen mode; The graphics mode is not supported.")
			err = ErrBadScreenMode
			return
		}

		// highly unlikely errors:

		if ret == win32.DISP_CHANGE_BADFLAGS {
			logger().Println("Unable to set screen mode; An invalid set of flags was passed in.")
			err = ErrBadScreenMode
			return
		}
		if ret == win32.DISP_CHANGE_BADPARAM {
			logger().Println("Unable to set screen mode; Invalid parameter or invalid flag (or combination of)")
			err = ErrBadScreenMode
			return
		}
		if ret == win32.DISP_CHANGE_NOTUPDATED {
			logger().Println("Unable to set screen mode; Unable to write settings to the registry.")
			err = ErrBadScreenMode
			return
		}
		if ret == win32.DISP_CHANGE_RESTART {
			logger().Println("Unable to set screen mode; Windows requires restart to achieve specific mode.")
			err = ErrBadScreenMode
			return
		}
	})
	return
}

func fetchScreenModes(w32GraphicsDeviceName string) (modes []*ScreenMode, currentMode *ScreenMode) {
	mode := new(win32.DEVMODE)
	hasCurrentMode := win32.EnumDisplaySettings(w32GraphicsDeviceName, win32.ENUM_CURRENT_SETTINGS, mode)

	if hasCurrentMode {
		width := int(mode.DmPelsWidth)
		height := int(mode.DmPelsHeight)
		refreshRate := float32(mode.DmDisplayFrequency)
		bpp := int(mode.DmBitsPerPel)

		currentMode = newScreenMode(width, height, bpp, refreshRate)
		currentMode.NativeScreenMode.w32Bpp = win32.DWORD(mode.DmBitsPerPel)
		//currentMode.NativeScreenMode.w32Mode = mode // See: 'Assign it here' below
		modes = append(modes, currentMode)
	}

	hasNext := true
	i := 0
	for hasNext {
		mode := new(win32.DEVMODE)
		hasNext = win32.EnumDisplaySettings(w32GraphicsDeviceName, win32.DWORD(i), mode)
		i++
		if hasNext {
			// This one is an good one

			if mode.DmDisplayFixedOutput() != win32.DMDFO_STRETCH {
				// Skip all modes that are not specified to stretch across the screen (there is
				// always at least one for each resolution specified as stretching, so no worries)
				continue
			}

			width := int(mode.DmPelsWidth)
			height := int(mode.DmPelsHeight)
			refreshRate := float32(mode.DmDisplayFrequency)
			bpp := int(mode.DmBitsPerPel)
			screenMode := newScreenMode(width, height, bpp, refreshRate)
			screenMode.NativeScreenMode.w32Bpp = win32.DWORD(mode.DmBitsPerPel)
			screenMode.NativeScreenMode.w32Mode = mode

			if hasCurrentMode {
				cmWidth, cmHeight := currentMode.Resolution()
				cmBpp := int(currentMode.NativeScreenMode.w32Bpp)
				cmRefreshRate := currentMode.RefreshRate()

				if width == cmWidth && height == cmHeight && refreshRate == cmRefreshRate && bpp == cmBpp {
					// Assign it here, avoid issues later on with comparison just in case we ever use it.
					currentMode.NativeScreenMode.w32Mode = mode
					continue // We already appended this before
				}
			}

			modes = append(modes, screenMode)
		}
	}

	if !hasCurrentMode {
		// Just guess it. This shouldn't happen though.
		logger().Println("Failed to find current screen mode; guessing mode 0.")
		currentMode = modes[0]
	}

	return
}

func backend_doScreens() (screens []*Screen) {
	win2kOrBelow := (w32VersionMajor <= 5) || (w32VersionMinor <= 0)

	monitorNum := 0
	hasNext := true
	i := 0
	for hasNext {
		dd := new(win32.DISPLAY_DEVICE)
		hasNext = win32.EnumDisplayDevices("", win32.DWORD(i), dd, 0)
		i++
		if hasNext {
			// We're only interested in active devices (graphics cards)
			graphicsCardName := win32.String(dd.DeviceName[:])
			graphicsCardString := win32.String(dd.DeviceString[:])

			gflags := dd.StateFlags
			if (gflags & win32.DISPLAY_DEVICE_ACTIVE) > 0 {
				hasMoreMonitors := true
				j := 0
				for hasMoreMonitors {
					hasMoreMonitors = win32.EnumDisplayDevices(win32.String(dd.DeviceName[:]), 0, dd, 0)
					j++

					// We're only interested in active monitors, but windows 2000 and below
					// never sets the DISPLAY_DEVICE_ACTIVE flag.
					//
					flags := dd.StateFlags

					active := (flags & win32.DISPLAY_DEVICE_ACTIVE) > 0
					attached := (flags & win32.DISPLAY_DEVICE_ATTACHED) > 0
					if active || attached || win2kOrBelow {

						// It's difficult to get monitor name or model, and it's only available on Windows 7+
						// eventually, we should try to fix this, but it seems mingw is missing the proper headers
						// with these definitions.
						//
						// See: http://msdn.microsoft.com/en-us/library/windows/hardware/ff553903(v=vs.85).aspx
						//
						monitorNum++
						name := fmt.Sprintf("Monitor %d - %s", monitorNum, graphicsCardString)

						w32GraphicsDeviceName := graphicsCardName
						dc := win32.CreateDC(w32GraphicsDeviceName, "", nil)
						if dc != nil {
							physicalWidth := float32(win32.GetDeviceCaps(dc, win32.HORZSIZE))
							physicalHeight := float32(win32.GetDeviceCaps(dc, win32.VERTSIZE))

							modes, currentMode := fetchScreenModes(w32GraphicsDeviceName)
							screen := newScreen(name, physicalWidth, physicalHeight, modes, currentMode)
							screen.NativeScreen.w32GraphicsDeviceName = w32GraphicsDeviceName
							screen.NativeScreen.w32MonitorDeviceName = win32.String(dd.DeviceName[:])
							if (gflags&win32.DISPLAY_DEVICE_PRIMARY_DEVICE) > 0 && j == 1 {
								screen.NativeScreen.isDefaultScreen = true
							}
							screen.NativeScreen.dc = dc

							screens = append(screens, screen)
						} else {
							// This hopefully never happens, but if it does that means
							// there is something wrong with this screen most likely or
							// an graphics driver bug or something else, who knows?
							logger().Println("CreateDC() on screen failed! Unable to create device context!")
							logger().Println("^ Screen will be ignored!")
						}
					}
				}
			}
		}
	}

	// Find the correct MONITORINFO struct for each screen and assign their w32Position properties
	proc := func(hMonitor win32.HMONITOR, hdcMonitor win32.HDC, lprcMonitor *win32.RECT, dwData win32.LPARAM) bool {
		mi := new(win32.MONITORINFOEX)
		if !win32.GetMonitorInfo(hMonitor, mi) {
			logger().Println("Unable to detect monitor position; GetMonitorInfo():", win32.GetLastErrorString())
		} else {
			for _, screen := range screens {
				if screen.NativeScreen.w32GraphicsDeviceName == win32.String(mi.SzDevice[:]) {
					screen.NativeScreen.w32Position = mi.RcMonitor
				}
			}
		}

		return true
	}

	if !win32.EnumDisplayMonitors(nil, nil, proc, 0) {
		logger().Println("Unable to detect monitor positions; EnumDisplayMonitors():", win32.GetLastErrorString())
	}
	return
}

func backend_Screens() (screens []*Screen) {
	dispatch(func() {
		screens = backend_doScreens()
	})
	return
}

func backend_DefaultScreen() *Screen {
	screens := backend_Screens()
	for _, screen := range screens {
		if screen.NativeScreen.isDefaultScreen {
			return screen
		}
	}

	// Should never happen
	if len(screens) > 0 {
		logger().Println("Unable to find default screen; falling back to first screen as default.")
		return screens[0]
	}
	logger().Println("No screens available!")
	return nil

}
