// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"image"
	"time"
)

// DestroyedEvent is an event where the window was destroyed.
type DestroyedEvent struct {
	T time.Time
}

// String returns an string representation of this event.
func (ev DestroyedEvent) String() string {
	return fmt.Sprintf("DestroyedEvent(Time=%v)", ev.T)
}

// Time implements the generic event interface.
func (ev DestroyedEvent) Time() time.Time {
	return ev.T
}

// PaintEvent is an event where the pixels inside the window's client region
// where damaged and need to be redrawn.
type PaintEvent struct {
	// Rectangle that needs to be redrawn / has been damaged.
	Rectangle image.Rectangle

	T time.Time
}

// String returns an string representation of this event.
func (ev PaintEvent) String() string {
	return fmt.Sprintf("PaintEvent(Rectangle=%v, Time=%v)", ev.Rectangle, ev.T)
}

// Time implements the generic event interface.
func (ev PaintEvent) Time() time.Time {
	return ev.T
}

// CloseEvent is an event where the user attempted to close the window, using
// either the exit button, an quick-key combination like alt+F4, etc.
type CloseEvent struct {
	T time.Time
}

// String returns an string representation of this event.
func (ev CloseEvent) String() string {
	return fmt.Sprintf("CloseEvent(Time=%v)", ev.T)
}

// Time implements the generic event interface.
func (ev CloseEvent) Time() time.Time {
	return ev.T
}

// CursorPositionEvent is an event where the user moved the mouse cursor inside
// the window's client region.
type CursorPositionEvent struct {
	// Position of cursor relative to the window's client region.
	X, Y float64

	T time.Time
}

// String returns an string representation of this event.
func (ev CursorPositionEvent) String() string {
	return fmt.Sprintf("CursorPositionEvent(X=%v, Y=%v, Time=%v)", ev.X, ev.Y, ev.T)
}

// Time implements the generic event interface.
func (ev CursorPositionEvent) Time() time.Time {
	return ev.T
}

// CursorWithinEvent is an event where the user moved the mouse cursor inside
// or outside the window's client region.
type CursorWithinEvent struct {
	// Weather the mouse cursor is within the window's client region or not.
	Within bool

	T time.Time
}

// String returns an string representation of this event.
func (ev CursorWithinEvent) String() string {
	return fmt.Sprintf("CursorWithinEvent(Within=%v, Time=%v)", ev.Within, ev.T)
}

// Time implements the generic event interface.
func (ev CursorWithinEvent) Time() time.Time {
	return ev.T
}

// MaximizedEvent is an event where the user maximized (or un-maximized) the
// window.
type MaximizedEvent struct {
	// Weather or not the window is currently maximized.
	Maximized bool

	T time.Time
}

// String returns an string representation of this event.
func (ev MaximizedEvent) String() string {
	return fmt.Sprintf("MaximizedEvent(Maximized=%v, Time=%v)", ev.Maximized, ev.T)
}

// Time implements the generic event interface.
func (ev MaximizedEvent) Time() time.Time {
	return ev.T
}

// MinimizedEvent is an event where the user minimized (or un-minimized) the
// window.
type MinimizedEvent struct {
	// Weather or not the window is currently minimized.
	Minimized bool

	T time.Time
}

// String returns an string representation of this event.
func (ev MinimizedEvent) String() string {
	return fmt.Sprintf("MinimizedEvent(Minimized=%v, Time=%v)", ev.Minimized, ev.T)
}

// Time implements the generic event interface.
func (ev MinimizedEvent) Time() time.Time {
	return ev.T
}

// FocusedEvent is an event where the user changed the focus of the window.
type FocusedEvent struct {
	// Weather the window has focus or not.
	Focused bool

	T time.Time
}

// String returns an string representation of this event.
func (ev FocusedEvent) String() string {
	return fmt.Sprintf("FocusedEvent(Focused=%v, Time=%v)", ev.Focused, ev.T)
}

// Time implements the generic event interface.
func (ev FocusedEvent) Time() time.Time {
	return ev.T
}

// PositionEvent is an event where the user changed the position of the window.
type PositionEvent struct {
	// Position of the window's client region.
	X, Y int

	T time.Time
}

// String returns an string representation of this event.
func (ev PositionEvent) String() string {
	return fmt.Sprintf("PositionEvent(X=%v, Y=%v, Time=%v)", ev.X, ev.Y, ev.T)
}

// Time implements the generic event interface.
func (ev PositionEvent) Time() time.Time {
	return ev.T
}

// ResizedEvent is an event where the user changed the size of the window's
// client region.
type ResizedEvent struct {
	// Size of the window's client region.
	Width, Height int

	T time.Time
}

// String returns an string representation of this event.
func (ev ResizedEvent) String() string {
	return fmt.Sprintf("ResizedEvent(Width=%v, Height=%v, Time=%v)", ev.Width, ev.Height, ev.T)
}

// Time implements the generic event interface.
func (ev ResizedEvent) Time() time.Time {
	return ev.T
}

// ScreenChangedEvent is an event where the user moved the window onto an
// different screen.
type ScreenChangedEvent struct {
	// Screen that the window is located on.
	Screen *Screen

	T time.Time
}

// String returns an string representation of this event.
func (ev ScreenChangedEvent) String() string {
	return fmt.Sprintf("ScreenChangedEvent(Screen=%v, Time=%v)", ev.Screen, ev.T)
}

// Time implements the generic event interface.
func (ev ScreenChangedEvent) Time() time.Time {
	return ev.T
}
