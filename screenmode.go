// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
)

// We use this type for sorting the screen modes in order of highest
// resolution, bytes per pixel, and refresh rate.
type sortedScreenModes []*ScreenMode

func (s sortedScreenModes) Len() int {
	return len(s)
}

func (s sortedScreenModes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortedScreenModes) Less(i, j int) bool {
	iWidth, iHeight := s[i].Resolution()
	iResolution := iWidth + iHeight
	iRefreshRate := s[i].RefreshRate()
	iBytesPerPixel := s[i].BytesPerPixel()

	jWidth, jHeight := s[j].Resolution()
	jResolution := jWidth + jHeight
	jRefreshRate := s[j].RefreshRate()
	jBytesPerPixel := s[j].BytesPerPixel()

	// if resolution and bpp are the same, sort by refresh rate
	if iResolution == jResolution && iBytesPerPixel == jBytesPerPixel {
		return iRefreshRate > jRefreshRate

		// if resolution is the same, sort by bpp
	} else if iResolution == jResolution {
		return iBytesPerPixel > jBytesPerPixel
	}

	// otherwise just sort by resolution
	return iResolution > jResolution
}

// ScreenMode represents an single, unique, screen mode, with an resolution,
// refresh rate, and bpp.
//
// It is possible for multiple different ScreenMode's to exist with the same
// resolution, each with different refresh rates or bytes per pixel,
// respectively.
type ScreenMode struct {
	*NativeScreenMode

	width, height, bytesPerPixel int
	refreshRate                  float32
}

func newScreenMode(width, height, bytesPerPixel int, refreshRate float32) *ScreenMode {
	s := new(ScreenMode)
	s.NativeScreenMode = newNativeScreenMode()
	s.width = width
	s.height = height
	s.bytesPerPixel = bytesPerPixel
	s.refreshRate = refreshRate
	return s
}

// String returns an nice string representing this ScreenMode
func (m *ScreenMode) String() string {
	w, h := m.Resolution()
	return fmt.Sprintf("ScreenMode(%d by %dpx, %.1fhz, %dbpp)", w, h, m.RefreshRate(), m.BytesPerPixel())
}

// Equals compares two ScreenMode(s) for equality. It does this by comparing resolutions,
// refresh rates, and bytes per pixels.
func (m *ScreenMode) Equals(other *ScreenMode) bool {
	width, height := m.Resolution()
	otherWidth, otherHeight := other.Resolution()

	return (width == otherWidth) && (height == otherHeight) && (m.RefreshRate() == other.RefreshRate()) && (m.BytesPerPixel() == other.BytesPerPixel())
}

// Resolution returns the width and height of this ScreenMode, in pixels.
func (m *ScreenMode) Resolution() (width, height int) {
	return m.width, m.height
}

// RefreshRate returns the refresh rate of this ScreenMode, in hertz, or 0 if the refresh rate
// is unable to be determined.
func (m *ScreenMode) RefreshRate() float32 {
	return m.refreshRate
}

// BytesPerPixel returns the number of bytes that represent an single pixel of this ScreenMode,
// or 0 if the bytes per pixel is unable to be determined.
func (m *ScreenMode) BytesPerPixel() int {
	return m.bytesPerPixel
}
