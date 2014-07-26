// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"fmt"
	"math"
	"sort"
)

type sortedConfigs struct {
	c   []*GLConfig
	max *GLConfig
}

func (s sortedConfigs) Len() int      { return len(s.c) }
func (s sortedConfigs) Swap(i, j int) { s.c[i], s.c[j] = s.c[j], s.c[i] }
func (s sortedConfigs) Less(i, j int) bool {
	// Sort in order
	//
	// Accelerated
	// RedBits, GreenBits, BlueBits, AlphaBits
	// DoubleBuffered
	// Samples
	// DepthBits
	// StencilBits
	// StereoScopic
	// Transparent
	// AccumRedBits, AccumGreenBits, AccumBlueBits, AccumAlphaBits
	// AuxBuffers
	distCmp := func(a, b, target uint8) bool {
		return math.Abs(float64(a-target)) > math.Abs(float64(b-target))
	}

	boolCmp := func(a, b, target bool) bool {
		return a == target
	}

	a := s.c[i]
	b := s.c[j]
	if a.Accelerated != b.Accelerated {
		return boolCmp(a.Accelerated, b.Accelerated, s.max.Accelerated)
	}
	if a.RedBits != b.RedBits {
		return distCmp(a.RedBits, b.RedBits, s.max.RedBits)
	}
	if a.GreenBits != b.GreenBits {
		return distCmp(a.GreenBits, b.GreenBits, s.max.GreenBits)
	}
	if a.BlueBits != b.BlueBits {
		return distCmp(a.BlueBits, b.BlueBits, s.max.BlueBits)
	}
	if a.AlphaBits != b.AlphaBits {
		return distCmp(a.AlphaBits, b.AlphaBits, s.max.AlphaBits)
	}
	if a.DoubleBuffered != b.DoubleBuffered {
		return boolCmp(a.DoubleBuffered, b.DoubleBuffered, s.max.DoubleBuffered)
	}
	if a.Samples != b.Samples {
		return distCmp(a.Samples, b.Samples, s.max.Samples)
	}
	if a.DepthBits != b.DepthBits {
		return distCmp(a.DepthBits, b.DepthBits, s.max.DepthBits)
	}
	if a.StencilBits != b.StencilBits {
		return distCmp(a.StencilBits, b.StencilBits, s.max.StencilBits)
	}
	if a.StereoScopic != b.StereoScopic {
		return boolCmp(a.StereoScopic, b.StereoScopic, s.max.StereoScopic)
	}
	if a.Transparent != b.Transparent {
		return boolCmp(a.Transparent, b.Transparent, s.max.Transparent)
	}
	if a.AccumRedBits != b.AccumRedBits {
		return distCmp(a.AccumRedBits, b.AccumRedBits, s.max.AccumRedBits)
	}
	if a.AccumGreenBits != b.AccumGreenBits {
		return distCmp(a.AccumGreenBits, b.AccumGreenBits, s.max.AccumGreenBits)
	}
	if a.AccumBlueBits != b.AccumBlueBits {
		return distCmp(a.AccumBlueBits, b.AccumBlueBits, s.max.AccumBlueBits)
	}
	if a.AccumAlphaBits != b.AccumAlphaBits {
		return distCmp(a.AccumAlphaBits, b.AccumAlphaBits, s.max.AccumAlphaBits)
	}
	return distCmp(a.AuxBuffers, b.AuxBuffers, s.max.AuxBuffers)
}

// GLConfig represents an single opengl (frame buffer / pixel) configuration.
type GLConfig struct {
	backend_GLConfig

	// An GLConfig is 'valid' when the user did not create it themself.
	valid bool

	// Tells whether this configuration is hardware accelerated or uses some
	// software implementation version of OpenGL.
	//
	// Note: Most software implementations are very low OpenGL versions. (I.e.
	// GL 1.1)
	Accelerated bool

	// Tells whether or not this config will support the window being
	// transparent while using OpenGL rendering
	Transparent bool

	// The number of bits that represent an color per pixel in the frame buffer.
	RedBits, GreenBits, BlueBits, AlphaBits uint8

	// The number of bits that represent an color per pixel in the accumulation
	// buffer.
	//
	// Note: GLSL shaders can perform an much better job of anything you would
	// be trying to do with the accumulation buffer.
	AccumRedBits, AccumGreenBits, AccumBlueBits, AccumAlphaBits uint8

	// Number of Multi Sample Anti Aliasing samples this configuration supports
	// (e.g. 2 for 2x MSAA, 16 for 16x MSAA, etc)
	Samples uint8

	// The number of auxiliary buffers available.
	//
	// Note: Auxiliary buffers are very rarely supported on most OpenGL
	// implementations (or choosing a configuration with them causes
	// non-accelerated rendering).
	//
	// For more information about this see the following forum URL:
	//     http://www.opengl.org/discussion_boards/showthread.php/171060-auxiliary-buffers
	AuxBuffers uint8

	// The number of bits that represent an pixel in the depth buffer.
	DepthBits uint8

	// The number of bits that represent an pixel in the stencil buffer.
	StencilBits uint8

	// Weather this frame buffer configuration is double buffered.
	DoubleBuffered bool

	// Weather this frame buffer configuration is stereoscopic capable.
	StereoScopic bool
}

func (c *GLConfig) panicUnlessValid() {
	if !c.valid {
		panic("Invalid GLConfig; did you attempt to create it yourself?")
	}
}

func (c *GLConfig) String() string {
	return fmt.Sprintf("GLConfig(Accelerated=%t, %dbpp[%d,%d,%d,%d], AccumBits=[%d,%d,%d,%d], Samples=%d, AuxBuffers=%d, DepthBits=%d, StencilBits=%d, DoubleBuffered=%t, Transparent=%t, StereoScopic=%t)", c.Accelerated, c.RedBits+c.GreenBits+c.BlueBits+c.AlphaBits, c.RedBits, c.GreenBits, c.BlueBits, c.AlphaBits, c.AccumRedBits, c.AccumGreenBits, c.AccumBlueBits, c.AccumAlphaBits, c.Samples, c.AuxBuffers, c.DepthBits, c.StencilBits, c.DoubleBuffered, c.Transparent, c.StereoScopic)
}

// Equals tells whether this GLConfig equals the other GLFrameBufferConfig, by
// comparing each attribute.
func (c *GLConfig) Equals(other *GLConfig) bool {
	o := other

	if c.Accelerated != o.Accelerated || c.Transparent != o.Transparent || c.RedBits != o.RedBits || c.GreenBits != o.GreenBits || c.BlueBits != o.BlueBits || c.AlphaBits != o.AlphaBits || c.AccumRedBits != o.AccumRedBits || c.AccumGreenBits != o.AccumGreenBits || c.AccumBlueBits != o.AccumBlueBits || c.AccumAlphaBits != o.AccumAlphaBits || c.Samples != o.Samples || c.AuxBuffers != o.AuxBuffers || c.DepthBits != o.DepthBits || c.StencilBits != o.StencilBits || c.DoubleBuffered != o.DoubleBuffered || c.StereoScopic != o.StereoScopic {
		return false
	}
	return true
}

var (
	// Describes the worst possible OpenGL frame buffer configuration, this is
	// typically used as a parameter to GLChooseConfig.
	GLWorstConfig = &GLConfig{
		Accelerated:    false,
		Transparent:    false,
		RedBits:        0,
		GreenBits:      0,
		BlueBits:       0,
		AlphaBits:      0,
		AccumRedBits:   0,
		AccumGreenBits: 0,
		AccumBlueBits:  0,
		AccumAlphaBits: 0,
		Samples:        0,
		AuxBuffers:     0,
		DepthBits:      0,
		StencilBits:    0,
		DoubleBuffered: false,
		StereoScopic:   false,
	}

	// Describes the worst possible OpenGL frame buffer configuration while
	// still being hardware-accelerated by the graphics card, this is typically
	// used as a parameter to GLChooseConfig.
	GLWorstHWConfig = &GLConfig{
		Accelerated:    true,
		Transparent:    false,
		RedBits:        0,
		GreenBits:      0,
		BlueBits:       0,
		AlphaBits:      0,
		AccumRedBits:   0,
		AccumGreenBits: 0,
		AccumBlueBits:  0,
		AccumAlphaBits: 0,
		Samples:        0,
		AuxBuffers:     0,
		DepthBits:      0,
		StencilBits:    0,
		DoubleBuffered: false,
		StereoScopic:   false,
	}

	// Describes the best possible OpenGL frame buffer configuration that does
	// not have any transparency, this is typically used as a parameter to
	// GLChooseConfig.
	GLBestConfig = &GLConfig{
		Accelerated:    true,
		Transparent:    true,
		RedBits:        255,
		GreenBits:      255,
		BlueBits:       255,
		AlphaBits:      255,
		AccumRedBits:   255,
		AccumGreenBits: 255,
		AccumBlueBits:  255,
		AccumAlphaBits: 255,
		Samples:        255,
		AuxBuffers:     255,
		DepthBits:      255,
		StencilBits:    255,
		DoubleBuffered: true,
		StereoScopic:   false,
	}
)

// GLChooseConfig chooses an appropriate configuration from the slice of
// possible configurations.
//
// The returned configuration will have at least minConfig's attributes, or nil
// will be returned if there is no configuration that has at least minConfig's
// attributes.
//
// The returned configuration will have no greater than maxConfig's attributes,
// or nil will be returned if there is no configuration that is below
// maxConfig's attributes.
//
// After the selection process excludes configurations below minConfig and
// above maxConfig, the compatible configurations not excluded are sorted in
// order of closest-to maxConfig in the order of the following configuration
// attributes:
//
//  Accelerated
//  RedBits, GreenBits, BlueBits, AlphaBits
//  DoubleBuffered
//  Samples
//  DepthBits
//  StencilBits
//  StereoScopic
//  Transparent
//  AccumRedBits, AccumGreenBits, AccumBlueBits, AccumAlphaBits
//  AuxBuffers
//
// And the first configuration in the sorted list is returned. (I.e. you are
// most likely to recieve a Accelerated configuration, then one with high bytes
// per pixel, then one who is DoubleBuffered, Transparent, and so forth).
//
// You may use the predefined GLWorstConfig and GLBestConfig variables if they
// suite your case.
func GLChooseConfig(possible []*GLConfig, minConfig, maxConfig *GLConfig) *GLConfig {
	min := minConfig
	max := maxConfig

	// Remove any which are below minConfig
	var removed []*GLConfig
	for _, c := range possible {
		if c.RedBits < min.RedBits {
			continue
		}
		if c.GreenBits < min.GreenBits {
			continue
		}
		if c.BlueBits < min.BlueBits {
			continue
		}
		if c.AlphaBits < min.AlphaBits {
			continue
		}

		if c.AccumRedBits < min.AccumRedBits {
			continue
		}
		if c.AccumGreenBits < min.AccumGreenBits {
			continue
		}
		if c.AccumBlueBits < min.AccumBlueBits {
			continue
		}
		if c.AccumAlphaBits < min.AccumAlphaBits {
			continue
		}

		if c.Samples < min.Samples {
			continue
		}
		if c.AuxBuffers < min.AuxBuffers {
			continue
		}
		if c.DepthBits < min.DepthBits {
			continue
		}
		if c.StencilBits < min.StencilBits {
			continue
		}

		if min.Accelerated && !c.Accelerated {
			continue
		}
		if min.Transparent && !c.Transparent {
			continue
		}
		if min.DoubleBuffered && !c.DoubleBuffered {
			continue
		}
		if min.StereoScopic && !c.StereoScopic {
			continue
		}
		removed = append(removed, c)
	}
	possible = removed

	// Remove any which are above maxConfig
	removed = make([]*GLConfig, 0)
	for _, c := range possible {
		if c.RedBits > max.RedBits {
			continue
		}
		if c.GreenBits > max.GreenBits {
			continue
		}
		if c.BlueBits > max.BlueBits {
			continue
		}
		if c.AlphaBits > max.AlphaBits {
			continue
		}

		if c.AccumRedBits > max.AccumRedBits {
			continue
		}
		if c.AccumGreenBits > max.AccumGreenBits {
			continue
		}
		if c.AccumBlueBits > max.AccumBlueBits {
			continue
		}
		if c.AccumAlphaBits > max.AccumAlphaBits {
			continue
		}

		if c.Samples > max.Samples {
			continue
		}
		if c.AuxBuffers > max.AuxBuffers {
			continue
		}
		if c.DepthBits > max.DepthBits {
			continue
		}
		if c.StencilBits > max.StencilBits {
			continue
		}

		if c.Accelerated && !max.Accelerated {
			continue
		}
		if c.Transparent && !max.Transparent {
			continue
		}
		if c.DoubleBuffered && !max.DoubleBuffered {
			continue
		}
		if c.StereoScopic && !max.StereoScopic {
			continue
		}
		removed = append(removed, c)
	}
	possible = removed

	if len(possible) == 0 {
		return nil
	}

	sorted := sortedConfigs{possible, maxConfig}
	sort.Sort(sorted)
	return sorted.c[0]
}
