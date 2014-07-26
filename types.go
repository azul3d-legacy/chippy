// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"image"
	"time"
)

// Cursor represents the image and hotspot of an graphical mouse cursor.
type Cursor struct {
	// Cursor image, internally it will be converted to an RGBA image.
	Image image.Image

	// Cursor Hotspot
	X, Y uint
}

// Event represents an event of some sort. The only requirement is that the
// event specify the point in time at which it happened.
//
// Normally you will use an type assertion or type switch to retrieve more
// useful information from the underlying type.
type Event interface {
	Time() time.Time
}
