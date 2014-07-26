// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"azul3d.org/chippy.v1/internal/win32"
	"errors"
	"fmt"
	"runtime"
)

type W32GLContext struct {
	valid     bool
	destroyed bool
	hglrc     win32.HGLRC
}

func (w *W32GLContext) panicUnlessValid() {
	if !w.valid {
		panic("Invalid GLContext; did you attempt to create it yourself?")
	}
}

func (w *W32GLContext) panicIfDestroyed() {
	if w.destroyed {
		panic("GL Context is already destroyed!")
	}
}

type backend_GLConfig struct {
	index win32.Int
}

func (w *NativeWindow) olderGLConfigs() (configs []*GLConfig) {
	dispatch(func() {
		// Just to get started
		max := win32.Int(2)

		for index := win32.Int(1); index-1 < max; index++ {
			var pf *win32.PIXELFORMATDESCRIPTOR
			max, pf = win32.DescribePixelFormat(w.dc, index)
			if max == 0 {
				logger().Println("Unable to get GLConfig's; DescribePixelFormat():", win32.GetLastErrorString())
				return
			}

			// We can only use pixel formats who have PFD_SUPPORT_OPENGL and
			// PFD_DRAW_TO_WINDOW (otherwise it may be an offscreen pixel
			// format)
			drawToWindow := (pf.DwFlags & win32.PFD_DRAW_TO_WINDOW) > 0
			if !drawToWindow {
				continue
			}
			supportOpenGL := (pf.DwFlags & win32.PFD_SUPPORT_OPENGL) > 0
			if !supportOpenGL {
				continue
			}

			// We only want ones whose pixel type is PFD_TYPE_RGBA
			if pf.IPixelType != win32.PFD_TYPE_RGBA {
				continue
			}

			config := new(GLConfig)
			config.valid = true
			config.index = index

			config.RedBits = uint8(pf.CRedBits)
			config.GreenBits = uint8(pf.CGreenBits)
			config.BlueBits = uint8(pf.CBlueBits)
			config.AlphaBits = uint8(pf.CAlphaBits)

			config.AccumRedBits = uint8(pf.CAccumRedBits)
			config.AccumGreenBits = uint8(pf.CAccumGreenBits)
			config.AccumBlueBits = uint8(pf.CAccumBlueBits)
			config.AccumAlphaBits = uint8(pf.CAccumAlphaBits)

			config.AuxBuffers = uint8(pf.CAuxBuffers)

			if (pf.DwFlags&win32.PFD_GENERIC_ACCELERATED) == 0 && (pf.DwFlags&win32.PFD_GENERIC_FORMAT) > 0 {
				config.Accelerated = false
			} else {
				config.Accelerated = true
			}

			config.DoubleBuffered = (pf.DwFlags & win32.PFD_DOUBLEBUFFER) > 0
			config.StereoScopic = (pf.DwFlags & win32.PFD_STEREO) > 0
			config.DepthBits = uint8(pf.CDepthBits)
			config.StencilBits = uint8(pf.CStencilBits)

			// Even some 24-bit RGB modes support composition, but only
			// ones with alpha frame buffer bits will be transparent.
			if config.AlphaBits > 0 && config.Accelerated {
				config.Transparent = (pf.DwFlags & win32.PFD_SUPPORT_COMPOSITION) > 0
			}

			configs = append(configs, config)
		}
	})
	return
}

func (w *NativeWindow) GLConfigs() (configs []*GLConfig) {
	// Basically, for more advanced pixel formats like ones including MSAA we
	// must use the WGL_ARB_pixel_format extension, before using that extension
	// we must test for it's existance using wglGetExtensionsStringARB, and to
	// use wglGetExtensionsStringARB we must use wglGetProcAddress, which in
	// fact only works in the presence of an OpenGL context, which requires we
	// use a 'dummy' window or something of the sort.
	//
	// We get the 'older' pixel formats using w.olderGLConfigs() method, this
	// tells us the pixel formats not supporting advanced features (like MSAA).
	//
	// Then we hit GLSetConfig(), GLCreateContext(), and GLMakeCurrent() which
	// will give the window a valid OpenGL context.
	//
	// After that we can see if we support better MSAA pixel formats, and since
	// our GLSetConfig() will automatically rebuild the window for changing GL
	// pixel formats we are OK.
	configs = w.olderGLConfigs()

	runtime.LockOSThread()

	// Create dummy context
	dummyConfig := GLChooseConfig(configs, GLWorstConfig, GLBestConfig)
	if dummyConfig == nil {
		// Maybe we don't have OpenGL at all
		logger().Println("GLConfigs(): GLChooseConfig() returned nil!")
		return
	}

	w.GLSetConfig(dummyConfig)
	ctx, err := w.GLCreateContext(1, 0, GLCoreProfile, nil)
	if err != nil {
		// Couldn't create the context for some reason
		logger().Println("GLConfigs(): GLCreateContext():", err)
		return
	}
	defer w.GLDestroyContext(ctx)
	w.GLMakeCurrent(ctx)
	defer w.GLMakeCurrent(nil)

	extensions, ok := win32.WglGetExtensionsStringARB(w.dc)
	if !ok {
		// Can't get extension list
		logger().Println("GLConfigs(): wglGetExtensionsStringARB failed!")
		return
	}

	if !extSupported(extensions, "WGL_ARB_pixel_format") {
		// Don't have WGL_ARB_pixel_format
		logger().Println("GLConfigs(): WGL_ARB_pixel_format not supported.")
		return
	}

	// Takes a pixel format index and single attribute name and returns it's
	// value.
	singleAttrib := func(pixelFormat, attribName win32.Int) (result win32.Int) {
		attrs := []win32.Int{attribName}
		win32.WglGetPixelFormatAttribivARB(w.dc, pixelFormat, 0, attrs, &result)
		return
	}

	nPixelFormats := singleAttrib(0, win32.WGL_NUMBER_PIXEL_FORMATS_ARB)
	if nPixelFormats == 0 {
		// No pixel formats? That seems wrong.
		logger().Println("GLConfigs(): WGL_NUMBER_PIXEL_FORMATS_ARB == 0")
		return
	}

	// Empty configs slice
	configs = configs[:0]

	// Query advanced pixel formats
	for i := win32.Int(1); i < nPixelFormats; i++ {
		// We can only use pixel formats who have WGL_SUPPORT_OPENGL_ARB and
		// WGL_DRAW_TO_WINDOW_ARB (otherwise it may be an offscreen pixel
		// format)
		drawToWindow := singleAttrib(i, win32.WGL_DRAW_TO_WINDOW_ARB)
		if drawToWindow == 0 {
			continue
		}
		supportOpenGL := singleAttrib(i, win32.WGL_SUPPORT_OPENGL_ARB)
		if supportOpenGL == 0 {
			continue
		}

		// We only want ones whose pixel type is WGL_TYPE_RGBA_ARB
		pixelType := singleAttrib(i, win32.WGL_PIXEL_TYPE_ARB)
		if pixelType != win32.WGL_TYPE_RGBA_ARB {
			continue
		}

		config := new(GLConfig)
		config.valid = true
		config.index = i

		config.RedBits = uint8(singleAttrib(i, win32.WGL_RED_BITS_ARB))
		config.GreenBits = uint8(singleAttrib(i, win32.WGL_GREEN_BITS_ARB))
		config.BlueBits = uint8(singleAttrib(i, win32.WGL_BLUE_BITS_ARB))
		config.AlphaBits = uint8(singleAttrib(i, win32.WGL_ALPHA_BITS_ARB))

		config.AccumRedBits = uint8(singleAttrib(i, win32.WGL_ACCUM_RED_BITS_ARB))
		config.AccumGreenBits = uint8(singleAttrib(i, win32.WGL_ACCUM_GREEN_BITS_ARB))
		config.AccumBlueBits = uint8(singleAttrib(i, win32.WGL_ACCUM_BLUE_BITS_ARB))
		config.AccumAlphaBits = uint8(singleAttrib(i, win32.WGL_ACCUM_ALPHA_BITS_ARB))

		config.AuxBuffers = uint8(singleAttrib(i, win32.WGL_AUX_BUFFERS_ARB))

		accel := singleAttrib(i, win32.WGL_ACCELERATION_ARB)
		if accel == win32.WGL_GENERIC_ACCELERATION_ARB || accel == win32.WGL_FULL_ACCELERATION_ARB {
			config.Accelerated = true
		}

		config.DoubleBuffered = singleAttrib(i, win32.WGL_DOUBLE_BUFFER_ARB) == 1
		config.StereoScopic = singleAttrib(i, win32.WGL_STEREO_ARB) == 1
		config.DepthBits = uint8(singleAttrib(i, win32.WGL_DEPTH_BITS_ARB))
		config.StencilBits = uint8(singleAttrib(i, win32.WGL_STENCIL_BITS_ARB))

		if config.AlphaBits > 0 {
			config.Transparent = true
		}

		if singleAttrib(i, win32.WGL_SAMPLE_BUFFERS_ARB) > 0 {
			config.Samples = uint8(singleAttrib(i, win32.WGL_SAMPLES_ARB))
		}

		configs = append(configs, config)
	}

	return
}

func (w *NativeWindow) GLSetConfig(config *GLConfig) {
	if config == nil {
		panic("Invalid (nil) GLConfig; it must be an valid configuration!")
	}
	config.panicUnlessValid()

	w.glConfig = config

	dispatch(func() {
		if w.glPixelFormatSet {
			err := w.doRebuildWindow()
			if err != nil {
				panic(err)
			}
		}

		if !win32.SetPixelFormat(w.dc, config.index, nil) {
			logger().Println("GLSetConfig failed; SetPixelFormat():", win32.GetLastErrorString())
		}
		w.glPixelFormatSet = true
	})
}

func (w *NativeWindow) GLConfig() *GLConfig {
	return w.glConfig
}

func (w *NativeWindow) GLCreateContext(glVersionMajor, glVersionMinor uint, flags GLContextFlags, share GLContext) (GLContext, error) {
	if w.glConfig == nil {
		panic("Must call GLSetConfig() before GLCreateContext()!")
	}
	c := new(W32GLContext)
	c.valid = true

	var (
		swc    *W32GLContext
		shglrc win32.HGLRC
	)
	if share != nil {
		swc = share.(*W32GLContext)
		swc.panicUnlessValid()
		swc.panicIfDestroyed()
		shglrc = swc.hglrc
	}

	var err error
	dispatch(func() {
		// First, make an fake context to use for context creation
		fakeContext := win32.WglCreateContext(w.dc)
		if fakeContext == nil {
			err = errors.New(fmt.Sprintf("Unable to create OpenGL context; wglCreateContext(): %s", win32.GetLastErrorString()))
			return
		}
		if !win32.WglMakeCurrent(w.dc, fakeContext) {
			err = errors.New(fmt.Sprintf("Unable to create OpenGL context; wglMakeCurrent(): %s", win32.GetLastErrorString()))
			return
		}

		extensions, ok := win32.WglGetExtensionsStringARB(w.dc)
		if !ok {
			logger().Println("wglGetExtensionsStringARB failed!")
		}

		if extSupported(extensions, "WGL_ARB_create_context") {
			attribs := []win32.Int{}

			// Major version
			attribs = append(attribs, win32.WGL_CONTEXT_MAJOR_VERSION_ARB)
			attribs = append(attribs, win32.Int(glVersionMajor))

			// Minor version
			attribs = append(attribs, win32.WGL_CONTEXT_MINOR_VERSION_ARB)
			attribs = append(attribs, win32.Int(glVersionMinor))

			// Forward compat & debug
			wglFlags := win32.Int(0)
			if (flags & GLForwardCompatible) > 0 {
				wglFlags |= win32.WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB
			}
			if (flags & GLDebug) > 0 {
				wglFlags |= win32.WGL_CONTEXT_DEBUG_BIT_ARB
			}
			if wglFlags != 0 {
				attribs = append(attribs, win32.WGL_CONTEXT_FLAGS_ARB)
				attribs = append(attribs, wglFlags)
			}

			// Note: context creation will fail on nvidia drivers if trying to
			// select a profile and the requested version is < 3.2
			//
			// See: https://www.opengl.org/discussion_boards/showthread.php/177832-Small-NVIDIA-wglCreateContextAttribsARB-Bug
			if glVersionMajor >= 3 && glVersionMinor >= 2 {
				// Profile selection
				//
				// "GLCompatabilityProfile will be used if neither GLCoreProfile or GLCompatibilityProfile
				// are present, or if both are present."
				profileMask := win32.Int(win32.WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB)

				wantCoreProfile := (flags & GLCoreProfile) > 0
				wantCompatProfile := (flags & GLCompatibilityProfile) > 0
				if wantCoreProfile && !wantCompatProfile {
					profileMask = win32.WGL_CONTEXT_CORE_PROFILE_BIT_ARB
				}

				attribs = append(attribs, win32.WGL_CONTEXT_PROFILE_MASK_ARB)
				attribs = append(attribs, profileMask)
			}

			// Attribs list is zero terminated
			attribs = append(attribs, 0)

			c.hglrc, ok = win32.WglCreateContextAttribsARB(w.dc, shglrc, attribs)
			if !ok {
				// The wglCreateContextAttribsARB entry point is missing
				//
				// Fall back to old context.
				logger().Println("WGL_ARB_create_context supported -- but wglCreateContextAttribsARB is missing!")
				c.hglrc = fakeContext

				if share != nil && !win32.WglShareLists(c.hglrc, shglrc) {
					logger().Println("wglShareLists() failed:", win32.GetLastErrorString())
				}

			} else if c.hglrc == nil {
				// Context couldn't be created for some reason (likely the version is not supported).
				//
				// Fall back to old context.

				logger().Println("wglCreateContextAttribsARB() failed:", win32.GetLastErrorString())
				c.hglrc = fakeContext

				if share != nil && !win32.WglShareLists(c.hglrc, shglrc) {
					logger().Println("wglShareLists() failed:", win32.GetLastErrorString())
				}

			} else {
				// It worked! We got our context!
				//
				// Clean up the fake context
				win32.WglMakeCurrent(nil, nil)
				win32.WglDeleteContext(fakeContext)

				// So we can get the version below
				win32.WglMakeCurrent(w.dc, c.hglrc)
			}

		} else {
			// They have no WGL_ARB_create_context support.
			//
			// Fall back to old context.
			logger().Println("WGL_ARB_create_context is unavailable.")
			c.hglrc = fakeContext

			if share != nil && !win32.WglShareLists(c.hglrc, shglrc) {
				logger().Println("wglShareLists() failed:", win32.GetLastErrorString())
			}
		}

		defer win32.WglMakeCurrent(nil, nil)

		ver := win32.GlGetString(win32.GL_VERSION)
		if !versionSupported(ver, int(glVersionMajor), int(glVersionMinor)) {
			logger().Printf("GL_VERSION=%q; no support for OpenGL %v.%v found\n", ver, glVersionMajor, glVersionMinor)
			err = ErrGLVersionNotSupported
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return c, err
}

func (w *NativeWindow) GLDestroyContext(c GLContext) {
	wc := c.(*W32GLContext)
	if !wc.destroyed {
		wc.destroyed = true
		dispatch(func() {
			if !win32.WglDeleteContext(wc.hglrc) {
				logger().Println("Unable to destroy GL context; wglDeleteContext():", win32.GetLastErrorString())
			}
		})
	}
}

func (w *NativeWindow) GLMakeCurrent(c GLContext) {
	var hglrc win32.HGLRC

	if c != nil {
		wc := c.(*W32GLContext)
		wc.panicUnlessValid()
		wc.panicIfDestroyed()
		hglrc = wc.hglrc
	}

	// Note: Avoid the temptation, never dispatch()!
	if !win32.WglMakeCurrent(w.dc, hglrc) {
		logger().Println("Unable to make GL context current; wglMakeCurrent():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) GLSwapBuffers() {
	if w.r.Destroyed() {
		return
	}
	if w.glConfig.DoubleBuffered == false {
		return
	}
	if !win32.SwapBuffers(w.dc) {
		logger().Println("Unable to swap GL buffers; SwapBuffers():", win32.GetLastErrorString())
	}
}

func (w *NativeWindow) GLVerticalSync() VSyncMode {
	return w.glVSyncMode
}

func (w *NativeWindow) GLSetVerticalSync(mode VSyncMode) {
	if !mode.Valid() {
		panic("Invalid vertical sync constant specified.")
	}
	if w.r.Destroyed() {
		return
	}

	w.glVSyncMode = mode

	var v int
	switch mode {
	case NoVerticalSync:
		v = 0

	case VerticalSync:
		v = 1

	case AdaptiveVerticalSync:
		v = -1
	}

	if !win32.WglSwapIntervalEXT(v) {
		logger().Println("Unable to set vertical sync; wglSwapIntervalEXT():", win32.GetLastErrorString())
	}
}
