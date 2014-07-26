// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"image"
)

type Blitable interface {
	// PixelBlit blits the specified RGBA image onto the window, at the given X
	// and Y coordinates.
	PixelBlit(x, y uint, pixels *image.RGBA)

	// PixelClear clears the given rectangle on the window's client region.
	PixelClear(rect image.Rectangle)
}
