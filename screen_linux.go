// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"azul3d.org/chippy.v1/internal/x11"
	"errors"
	"fmt"
	"sync"
)

type NativeScreenMode struct {
	isCurrentMode    bool
	xrandrModeInfoId uint32
}

func newNativeScreenMode() *NativeScreenMode {
	n := new(NativeScreenMode)
	return n
}

func xrandrRefreshRateFromModeInfo(modeInfo *x11.RandrModeInfo) (refreshRate float32) {
	// Calculate refresh rate
	dotclock := float32(modeInfo.DotClock)
	vtotal := float32(modeInfo.VTotal)
	htotal := float32(modeInfo.HTotal)
	flags := modeInfo.ModeFlags
	if (flags & x11.RANDR_MODE_FLAG_DOUBLE_SCAN) > 0 {
		// doublescan doubles the number of lines
		vtotal *= 2
	}
	if (flags & x11.RANDR_MODE_FLAG_INTERLACE) > 0 {
		// interlace splits the frame into two fields
		// the field rate is what is typically reported by monitors
		vtotal /= 2
	}

	if htotal != 0 && vtotal != 0 {
		refreshRate = dotclock / (htotal * vtotal)
	}

	return dotclock / (float32(modeInfo.HTotal) * vtotal)
}

type NativeScreen struct {
	access sync.RWMutex

	isDefaultScreen      bool
	xPosition, yPosition int16

	xScreen    int
	xrandrCrtc x11.RandrCrtc
}

func (s *NativeScreen) setMode(m *ScreenMode) error {
	if atLeastVersion(xrandrMajor, xrandrMinor, xrandrMinMajor, xrandrMinMinor) {
		// Find the root window
		var root x11.Window
		screen := xConnection.ScreenOfDisplay(s.xScreen)
		if screen != nil {
			root = screen.Root
		} else {
			return errors.New("SetMode(): Failed to set screen mode; ScreenOfDisplay() failed!")
		}

		var modes *x11.RandrModeInfos
		if atLeastVersion(xrandrMajor, xrandrMinor, 1, 3) {
			// 1.3 API
			cookie := xConnection.RandrGetScreenResourcesCurrent(root)
			resourcesCurrent, err := xConnection.RandrGetScreenResourcesCurrentReply(cookie)
			if err != nil {
				return err
			}
			modes = xConnection.RandrGetScreenResourcesCurrentModes(resourcesCurrent)
		} else {
			// 1.2 API
			cookie := xConnection.RandrGetScreenResources(root)
			resources, err := xConnection.RandrGetScreenResourcesReply(cookie)
			if err != nil {
				return err
			}
			modes = xConnection.RandrGetScreenResourcesModes(resources)
		}

		// 1.2 API
		cookie := xConnection.RandrGetCrtcInfo(s.xrandrCrtc, x11.TIME_CURRENT_TIME)
		info, err := xConnection.RandrGetCrtcInfoReply(cookie)
		if err != nil {
			return err
		}
		infoOutputs := xConnection.RandrGetCrtcInfoReplyOutputs(info)

		for _, mode := range modes.Slice {
			if uint32(mode.Id) == m.NativeScreenMode.xrandrModeInfoId {
				// It's this one!
				cookie := xConnection.RandrSetCrtcConfig(
					s.xrandrCrtc,
					x11.TIME_CURRENT_TIME, x11.TIME_CURRENT_TIME,
					int16(info.X), int16(info.Y),
					x11.RandrMode(mode.Id),
					uint16(info.Rotation),
					infoOutputs.Slice,
				)
				_, err := xConnection.RandrSetCrtcConfigReply(cookie)
				if err != nil {
					return err
				}
				break
			}
		}

	} else {
		// No Xrandr extension available; we have no way to change the screen mode.
		return fmt.Errorf("SetMode(): X11 Randr extension not available.")
	}
	return nil
}

func newNativeScreen() *NativeScreen {
	n := new(NativeScreen)
	return n
}

func fetchScreenModes(xScreenNumber int, xrandrCrtc x11.RandrCrtc) (modes []*ScreenMode, current *ScreenMode) {
	if atLeastVersion(xrandrMajor, xrandrMinor, xrandrMinMajor, xrandrMinMinor) {
		// Find the root window
		var root x11.Window
		screen := xConnection.ScreenOfDisplay(xScreenNumber)
		if screen != nil {
			root = screen.Root
		} else {
			logger().Println("Screens(): ScreenOfDisplay() failed!")
			goto fallback
		}

		var resourceModes *x11.RandrModeInfos
		if atLeastVersion(xrandrMajor, xrandrMinor, 1, 3) {
			// 1.3 API
			cookie := xConnection.RandrGetScreenResourcesCurrent(root)
			resourcesCurrent, err := xConnection.RandrGetScreenResourcesCurrentReply(cookie)
			if err != nil {
				logger().Println("Screens(): err")
				goto fallback
			}
			resourceModes = xConnection.RandrGetScreenResourcesCurrentModes(resourcesCurrent)
		} else {
			// 1.2 API
			cookie := xConnection.RandrGetScreenResources(root)
			resources, err := xConnection.RandrGetScreenResourcesReply(cookie)
			if err != nil {
				logger().Println("Screens(): err")
				goto fallback
			}
			resourceModes = xConnection.RandrGetScreenResourcesModes(resources)
		}

		// 1.2 API
		cookie := xConnection.RandrGetCrtcInfo(xrandrCrtc, x11.TIME_CURRENT_TIME)
		info, err := xConnection.RandrGetCrtcInfoReply(cookie)
		if err != nil {
			logger().Println("Screens(): err")
			goto fallback
		}
		infoOutputs := xConnection.RandrGetCrtcInfoReplyOutputs(info)
		crtcModeId := info.Mode

		// The output's are the actual physical monitor devices
		for _, output := range infoOutputs.Slice {
			// 1.2 API
			cookie := xConnection.RandrGetOutputInfo(output, x11.TIME_CURRENT_TIME)
			outputInfo, err := xConnection.RandrGetOutputInfoReply(cookie)
			if err != nil {
				logger().Println("Screens():", err)
				logger().Printf("^ Dropping output #%d\n", output)
				continue
			}

			outputInfoModes := xConnection.RandrGetOutputInfoModes(outputInfo)
			for _, actualMode := range outputInfoModes.Slice {
				for _, modeInfo := range resourceModes.Slice {
					if actualMode == x11.RandrMode(modeInfo.Id) {
						// This mode is related to output

						screenMode := newScreenMode(int(modeInfo.Width), int(modeInfo.Height), 0, xrandrRefreshRateFromModeInfo(&modeInfo))
						screenMode.NativeScreenMode.xrandrModeInfoId = uint32(modeInfo.Id)

						if actualMode == crtcModeId {
							screenMode.isCurrentMode = true
							current = screenMode
						}
						modes = append(modes, screenMode)
					}
				}
			}
		}
	}

fallback:
	if len(modes) == 0 {
	}

	/*
		if len(modes) == 0 {
			// We have no way to determine the screen mode through just Xlib alone.. we need
			// the extensions, the only thing we can do is say there is only an single screen mode
			// available.
			xScreen := x11.XScreenOfDisplay(xDisplay, xScreenNumber)
			width := x11.XWidthOfScreen(xScreen)
			height := x11.XHeightOfScreen(xScreen)
			mode := newScreenMode(width, height, 0, 0)
			mode.NativeScreenMode.isCurrentMode = true
			modes = append(modes, mode)
			current = mode
		}
	*/

	return
}

func backend_Screens() (screens []*Screen) {
	if atLeastVersion(xrandrMajor, xrandrMinor, xrandrMinMajor, xrandrMinMinor) {
		// Xrandr is an big virtual screen made up of physical monitors arranged on the virtual one
		screenCount := xConnection.SetupRootsLength(xConnection.GetSetup())

		for screenNumber := 0; screenNumber < screenCount; screenNumber++ {
			screenNumber = 0

			// Find the root window
			var root x11.Window
			screen := xConnection.ScreenOfDisplay(screenNumber)
			if screen != nil {
				root = screen.Root
			} else {
				logger().Println("Screens(): ScreenOfDisplay() failed!")
				goto fallback
			}

			var (
				crtcs            *x11.RandrCrtcs
				cfgTime          x11.Timestamp
				resourcesCurrent *x11.RandrGetScreenResourcesCurrentReply
				resources        *x11.RandrGetScreenResourcesReply
				err              error
			)
			if atLeastVersion(xrandrMajor, xrandrMinor, 1, 3) {
				// 1.3 API
				cookie := xConnection.RandrGetScreenResourcesCurrent(root)
				resourcesCurrent, err = xConnection.RandrGetScreenResourcesCurrentReply(cookie)
				if err != nil {
					logger().Println("Screens():", err)
					goto fallback
				}
				cfgTime = resourcesCurrent.ConfigTimestamp
				crtcs = xConnection.RandrGetScreenResourcesCurrentCrtcs(resourcesCurrent)
			} else {
				// 1.2 API
				cookie := xConnection.RandrGetScreenResources(root)
				resources, err = xConnection.RandrGetScreenResourcesReply(cookie)
				if err != nil {
					logger().Println("Screens():", err)
					goto fallback
				}
				cfgTime = resources.ConfigTimestamp
				crtcs = xConnection.RandrGetScreenResourcesCrtcs(resources)
			}

			for _, crtc := range crtcs.Slice {
				// 1.2 API
				cookie := xConnection.RandrGetCrtcInfo(crtc, cfgTime)
				info, err := xConnection.RandrGetCrtcInfoReply(cookie)
				if err != nil {
					logger().Println("Screens():", err)
					goto fallback
				}

				// Check for Disabled/Inactive
				if info.Mode != 0 {
					// The output's are the actual physical monitor devices
					infoOutputs := xConnection.RandrGetCrtcInfoReplyOutputs(info)
					for _, output := range infoOutputs.Slice {
						// 1.2 API
						cookie := xConnection.RandrGetOutputInfo(output, x11.TIME_CURRENT_TIME)
						outputInfo, err := xConnection.RandrGetOutputInfoReply(cookie)
						if err != nil {
							logger().Println("Screens():", err)
							logger().Printf("^ Dropping output #%d\n", output)
							continue
						}

						modes, currentMode := fetchScreenModes(screenNumber, crtc)

						name := xConnection.RandrGetOutputInfoName(outputInfo)
						name = fmt.Sprintf("Screen %v - %s", screenNumber+1, name)
						screen := newScreen(name, float32(outputInfo.MMWidth), float32(outputInfo.MMHeight), modes, currentMode)
						screen.NativeScreen.xScreen = screenNumber
						screen.NativeScreen.xrandrCrtc = crtc

						screens = append(screens, screen)
					}
				}
			}
		}

		return
	}

fallback:
	// We don't have randr extension, so we can use pure X11 API's then.
	i := 0
	xConnection.RootsIterator(func(s *x11.Screen) bool {
		mode := newScreenMode(
			int(s.WidthInPixels),
			int(s.HeightInPixels),
			int(s.RootDepth),
			0,
		)
		mode.NativeScreenMode.isCurrentMode = true

		modes := []*ScreenMode{mode}
		currentMode := mode

		name := fmt.Sprintf("Screen %v", i+1)
		screen := newScreen(
			name,
			float32(s.WidthInMillimeters),
			float32(s.HeightInMillimeters),
			modes, currentMode,
		)
		screen.NativeScreen.xScreen = i
		screens = append(screens, screen)

		i++
		return true
	})

	/*
		// Otherwise we have no Xrandr, so we can only use Xlib
		screenCount := x11.XScreenCount(xDisplay)
		screens = make([]*Screen, screenCount)

		for i := 0; i < screenCount; i++ {
			var rr x11.RRCrtc
			modes, currentMode := fetchScreenModes(i, rr)

			xScreen := x11.XScreenOfDisplay(xDisplay, i)
			mmWidth := float32(x11.XWidthMMOfScreen(xScreen))
			mmHeight := float32(x11.XHeightMMOfScreen(xScreen))
			name := fmt.Sprintf("Screen %v", i)
			screen := newScreen(name, mmWidth, mmHeight, modes, currentMode)
			screen.NativeScreen.xScreen = i

			screens[i] = screen
		}
	*/
	return screens
}

func backend_DefaultScreen() *Screen {
	screens := backend_Screens()

	for _, screen := range screens {
		if screen.NativeScreen.xScreen == xDefaultScreenNumber {
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
