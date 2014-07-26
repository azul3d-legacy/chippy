// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"azul3d.org/chippy.v1/internal/x11"
	"errors"
	"unsafe"
)

type X11GLContext struct {
	valid      bool
	destroyed  bool
	glxContext x11.GLXContext
}

func (w *X11GLContext) panicUnlessValid() {
	if !w.valid {
		panic("Invalid GLContext; did you attempt to create it yourself?")
	}
}

func (w *X11GLContext) panicIfDestroyed() {
	if w.destroyed {
		panic("GL Context is already destroyed!")
	}
}

type backend_GLConfig struct {
	glxConfig    x11.GLXFBConfig
	xVisual      x11.VisualId
	xVisualDepth uint8
}

func (w *NativeWindow) GLConfigs() (configs []*GLConfig) {
	screen := w.r.Screen()
	for _, glxConfig := range glxDisplay.GLXGetFBConfigs(screen.NativeScreen.xScreen) {
		config := new(GLConfig)
		config.valid = true
		config.glxConfig = glxConfig

		vi := glxDisplay.GLXGetVisualFromFBConfig(config.glxConfig)
		if vi == nil {
			// Doesn't have a matching X visual, we can't use it.
			continue
		}
		config.xVisual = x11.VisualId(vi.Visualid())
		config.xVisualDepth = uint8(vi.Depth())

		doubleBuffer, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_DOUBLEBUFFER)
		if err == 0 {
			config.DoubleBuffered = doubleBuffer == 1
		}

		stereo, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_STEREO)
		if err == 0 {
			config.StereoScopic = stereo == 1
		}

		auxBuffers, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_AUX_BUFFERS)
		if err == 0 {
			config.AuxBuffers = uint8(auxBuffers)
		}

		redBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_RED_SIZE)
		if err == 0 {
			config.RedBits = uint8(redBits)
		}

		greenBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_GREEN_SIZE)
		if err == 0 {
			config.GreenBits = uint8(greenBits)
		}

		blueBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_BLUE_SIZE)
		if err == 0 {
			config.BlueBits = uint8(blueBits)
		}

		alphaBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_ALPHA_SIZE)
		if err == 0 {
			config.AlphaBits = uint8(alphaBits)
		}

		depthBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_DEPTH_SIZE)
		if err == 0 {
			config.DepthBits = uint8(depthBits)
		}

		accumRedBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_ACCUM_RED_SIZE)
		if err == 0 {
			config.AccumRedBits = uint8(accumRedBits)
		}

		accumGreenBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_ACCUM_GREEN_SIZE)
		if err == 0 {
			config.AccumGreenBits = uint8(accumGreenBits)
		}

		accumBlueBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_ACCUM_BLUE_SIZE)
		if err == 0 {
			config.AccumBlueBits = uint8(accumBlueBits)
		}

		accumAlphaBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_ACCUM_ALPHA_SIZE)
		if err == 0 {
			config.AccumAlphaBits = uint8(accumAlphaBits)
		}

		sampleBuffers, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_SAMPLE_BUFFERS)
		if err == 0 && sampleBuffers > 0 {
			samples, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_SAMPLES)
			if err == 0 {
				config.Samples = uint8(samples)
			}
		}

		if vi.Depth() == 32 && vi.RedMask() == 0xFF0000 && vi.GreenMask() == 0x00FF00 && vi.BlueMask() == 0x0000FF {
			config.Transparent = true
		}

		stencilBits, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_STENCIL_SIZE)
		if err == 0 {
			config.StencilBits = uint8(stencilBits)
		}

		caveat, err := glxDisplay.GLXGetFBConfigAttrib(glxConfig, x11.GLX_CONFIG_CAVEAT)
		if err == 0 {
			config.Accelerated = caveat != x11.GLX_SLOW_CONFIG
		}

		x11.XFree(unsafe.Pointer(vi))
		configs = append(configs, config)
	}
	return
}

func (w *NativeWindow) GLSetConfig(config *GLConfig) {
	w.access.Lock()
	defer w.access.Unlock()

	if config == nil {
		panic("Invalid (nil) GLConfig; it must be an valid configuration!")
	}
	config.panicUnlessValid()

	// No need to set the GL config to one already active
	if w.glConfig != nil && w.glConfig.xVisual == config.xVisual {
		return
	}

	w.activeRenderMode = renderModeOpenGL
	w.glConfig = config

	// No need to set rebuild window if we already found a GLX compatible X
	// visual when we created the window.
	if config.xVisual == w.xVisual && config.xVisualDepth == w.xDepth {
		return
	}

	w.xDepth = config.xVisualDepth
	w.xVisual = config.xVisual
	w.doRebuildWindow()
	glxDisplay.XSync(false)
	xDisplay.XSync(false)
}

func (w *NativeWindow) GLConfig() *GLConfig {
	w.access.RLock()
	defer w.access.RUnlock()
	return w.glConfig
}

func (w *NativeWindow) GLCreateContext(glVersionMajor, glVersionMinor uint, flags GLContextFlags, share GLContext) (GLContext, error) {
	w.access.Lock()
	defer w.access.Unlock()
	xWindow := w.getXWindow()
	if xWindow == 0 {
		panic("Window is destroyed!")
	}
	if w.glConfig == nil {
		panic("Must call GLSetConfig() before GLCreateContext()!")
	}
	c := new(X11GLContext)
	c.valid = true

	var shareContext x11.GLXContext
	if share != nil {
		shareContext = share.(*X11GLContext).glxContext
	}

	// Pretty simple: if we want a GL 3.0 or higher context then we need the
	// GLX_ARB_create_context extension, otherwise we can't make the context.
	glxArbCreateContext := w.glxExtensionSupported("GLX_ARB_create_context")
	wantNewStyleContext := glVersionMajor >= 3

	if !glxArbCreateContext && wantNewStyleContext {
		return nil, ErrGLVersionNotSupported
	}

	if wantNewStyleContext {
		// Create a new-style GL 3+ context
		var attribs []x11.Int

		// Major version
		attribs = append(attribs, x11.GLX_CONTEXT_MAJOR_VERSION_ARB)
		attribs = append(attribs, x11.Int(glVersionMajor))

		// Minor version
		attribs = append(attribs, x11.GLX_CONTEXT_MINOR_VERSION_ARB)
		attribs = append(attribs, x11.Int(glVersionMinor))

		// Forward compat & debug
		glxFlags := x11.Int(0)
		if (flags & GLForwardCompatible) > 0 {
			glxFlags |= x11.GLX_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB
		}
		if (flags & GLDebug) > 0 {
			glxFlags |= x11.GLX_CONTEXT_DEBUG_BIT_ARB
		}
		if glxFlags != 0 {
			attribs = append(attribs, x11.GLX_CONTEXT_FLAGS_ARB)
			attribs = append(attribs, glxFlags)
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
			profileMask := x11.Int(x11.GLX_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB)

			wantCoreProfile := (flags & GLCoreProfile) > 0
			wantCompatProfile := (flags & GLCompatibilityProfile) > 0
			if wantCoreProfile && !wantCompatProfile {
				profileMask = x11.GLX_CONTEXT_CORE_PROFILE_BIT_ARB
			}

			attribs = append(attribs, x11.GLX_CONTEXT_PROFILE_MASK_ARB)
			attribs = append(attribs, profileMask)
		}

		// Attribs list is zero terminated
		attribs = append(attribs, 0)

		// Create the context
		c.glxContext = glxDisplay.GLXCreateContextAttribsARB(
			w.glConfig.glxConfig,
			shareContext,
			true, // direct
			&attribs[0],
		)
		if c.glxContext == nil {
			return nil, errors.New("glXCreateContextAttribsARB() failed!")
		}

	} else {
		c.glxContext = glxDisplay.GLXCreateNewContext(
			w.glConfig.glxConfig,
			x11.GLX_RGBA_TYPE,
			shareContext,
			true, // direct
		)
		if c.glxContext == nil {
			return nil, errors.New("glXCreateNewContext() failed!")
		}
	}

	glxDisplay.GLXMakeContextCurrent(x11.GLXDrawable(xWindow), x11.GLXDrawable(xWindow), c.glxContext)
	defer glxDisplay.GLXMakeContextCurrent(0, 0, nil)

	ver := x11.GlGetString(x11.GL_VERSION)
	if !versionSupported(ver, int(glVersionMajor), int(glVersionMinor)) {
		logger().Printf("GL_VERSION=%q; no support for OpenGL %v.%v found\n", ver, glVersionMajor, glVersionMinor)
		return nil, ErrGLVersionNotSupported
	}

	cb := func() {
		w.GLDestroyContext(c)
	}
	addDestroyCallback(&cb)

	return c, nil
}

func (w *NativeWindow) GLDestroyContext(c GLContext) {
	w.access.Lock()
	defer w.access.Unlock()
	wc := c.(*X11GLContext)
	wc.panicUnlessValid()
	if !wc.destroyed {
		wc.destroyed = true
		glxDisplay.GLXDestroyContext(wc.glxContext)
	}
}

func (w *NativeWindow) GLMakeCurrent(c GLContext) {
	w.access.Lock()
	defer w.access.Unlock()
	xWindow := w.getXWindow()
	if xWindow == 0 {
		return
	}

	var glxContext x11.GLXContext
	var glxWindow x11.GLXDrawable
	if c != nil {
		wc := c.(*X11GLContext)
		wc.panicUnlessValid()
		wc.panicIfDestroyed()
		glxContext = wc.glxContext
		glxWindow = x11.GLXDrawable(xWindow)
	}
	if glxDisplay.GLXMakeContextCurrent(glxWindow, glxWindow, glxContext) == 0 {
		logger().Println("glXMakeContextCurrent() failed!")
	}
}

func (w *NativeWindow) GLSwapBuffers() {
	w.access.Lock()
	defer w.access.Unlock()
	xWindow := w.getXWindow()
	if xWindow == 0 || w.glConfig.DoubleBuffered == false {
		return
	}

	glxDisplay.GLXSwapBuffers(x11.GLXDrawable(xWindow))
}

func (w *NativeWindow) GLVerticalSync() VSyncMode {
	return w.glVSyncMode
}

func (w *NativeWindow) glxExtensionSupported(wantedExt string) bool {
	if len(w.glxExtensionsString) == 0 {
		w.glxExtensionsString = glxDisplay.GLXQueryExtensionsString(w.r.Screen().xScreen)
	}

	return extSupported(w.glxExtensionsString, wantedExt)
}

func (w *NativeWindow) GLSetVerticalSync(mode VSyncMode) {
	if !mode.Valid() {
		panic("Invalid vertical sync constant specified.")
	}
	w.access.Lock()
	defer w.access.Unlock()
	xWindow := w.getXWindow()
	if xWindow == 0 {
		return
	}
	w.glVSyncMode = mode

	glxSwapControlEXT := w.glxExtensionSupported("GLX_EXT_swap_control")
	glxSwapControlMESA := w.glxExtensionSupported("GLX_MESA_swap_control")
	glxSwapControlSGI := w.glxExtensionSupported("GLX_SGI_swap_control")

	if !glxSwapControlEXT && !glxSwapControlMESA && !glxSwapControlSGI {
		if !glxSwapControlEXT {
			logger().Println("GLSetVerticalSync(): GLX_EXT_swap_control is not supported.")
		}
		if !glxSwapControlMESA {
			logger().Println("GLSetVerticalSync(): GLX_MESA_swap_control is not supported.")
		}
		if !glxSwapControlSGI {
			logger().Println("GLSetVerticalSync(): GLX_SGI_swap_control is not supported.")
		}
		logger().Println("GLSetVerticalSync(): Unable to set vertical sync due to above errors.")
		return
	}

	var v int
	switch mode {
	case NoVerticalSync:
		v = 0

	case VerticalSync:
		v = 1

	case AdaptiveVerticalSync:
		v = -1
	}
	if glxSwapControlEXT {
		glxDisplay.GLXSwapIntervalEXT(x11.GLXDrawable(xWindow), v)
	} else if glxSwapControlMESA {
		if mode == AdaptiveVerticalSync {
			logger().Println("GLSetVerticalSync(): Driver does not support adaptive vsync; using regular vsync.")
			v = 1
		}
		glxDisplay.GLXSwapIntervalMESA(v)
	}

	// Actually, glxSwapControlSGI does not let us use adaptive vsync (-1), nor turn off vsync (0).
	//
	// We can e.g. set vsync to 2, which would sync at 30hz for a 60hz display.
	//
	// It's not much use to us right now.
	/* else if glxSwapControlSGI {
		if mode == AdaptiveVerticalSync {
			logger().Println("GLSetVerticalSync(): Driver does not support adaptive vsync; using regular vsync.")
			v = 1
		} else if mode == NoVerticalSync {
			// Well, we can't set it to zero for off using the SGI extension.
			logger().Println("GLSetVerticalSync(): Driver does not support turning off vsync; sorry.")
			v = 1
		}
		glxDisplay.GLXSwapIntervalSGI(v)
	}*/
}
