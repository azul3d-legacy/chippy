// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrGLVersionNotSupported = errors.New("opengl version is not supported")
)

// Do not use for multiple extensions as it splits the string and searches it slowly..
//
// I.e. do not expose to end users..
func extSupported(str, ext string) bool {
	for _, s := range strings.Split(str, " ") {
		if s == ext {
			return true
		}
	}
	return false
}

func versionSupported(ver string, wantedMajor, wantedMinor int) bool {
	if len(ver) > 0 {
		// Version string must not be empty
		var (
			major, minor  int
			versionString string
			err           error
		)

		// According to http://www.opengl.org/sdk/docs/man/xhtml/glGetString.xml
		//
		// the string returned may be 'major.minor' or 'major.minor.release'
		// and may be following by a space and any vendor specific information.

		// First locate a proper version string without vendor specific
		// information.
		if strings.Contains(ver, " ") {
			// It must have vendor information.
			split := strings.Split(ver, " ")
			if len(split) > 0 || len(split[0]) > 0 {
				// Everything looks good.
				versionString = split[0]
			} else {
				// Something must be wrong with their vendor string.
				return false
			}
		} else {
			// No vendor information.
			versionString = ver
		}

		// We have a proper version string now without vendor information.
		dots := strings.Count(versionString, ".")
		if dots == 1 {
			// It's a 'major.minor' style string
			versions := strings.Split(versionString, ".")
			if len(versions) == 2 {
				major, err = strconv.Atoi(versions[0])
				if err != nil {
					return false
				}

				minor, err = strconv.Atoi(versions[1])
				if err != nil {
					return false
				}

			} else {
				return false
			}

		} else if dots == 2 {
			// It's a 'major.minor.release' style string
			versions := strings.Split(versionString, ".")
			if len(versions) == 3 {
				major, err = strconv.Atoi(versions[0])
				if err != nil {
					return false
				}

				minor, err = strconv.Atoi(versions[1])
				if err != nil {
					return false
				}
			} else {
				return false
			}
		}

		if major > wantedMajor {
			return true
		} else if major == wantedMajor && minor >= wantedMinor {
			return true
		}
	}
	return false
}

type GLContextFlags uint8

const (
	GLDebug GLContextFlags = iota
	GLForwardCompatible
	GLCoreProfile
	GLCompatibilityProfile
)

// VSyncMode represents an single vertical reresh rate sync mode.
type VSyncMode uint8

const (
	VerticalSync VSyncMode = iota
	NoVerticalSync
	AdaptiveVerticalSync
)

// Valid tells if the vertical sync mode is one of the predefined constants
// defined in this package or not.
func (mode VSyncMode) Valid() bool {
	switch mode {
	case VerticalSync:
		return true

	case NoVerticalSync:
		return true

	case AdaptiveVerticalSync:
		return true
	}
	return false
}

// String returns a string representation of this vertical sync mode.
func (mode VSyncMode) String() string {
	switch mode {
	case VerticalSync:
		return "VerticalSync"

	case NoVerticalSync:
		return "NoVerticalSync"

	case AdaptiveVerticalSync:
		return "AdaptiveVerticalSync"
	}
	return fmt.Sprintf("VSyncMode(%d)", mode)
}

// GLContext represents an OpenGL contect; although it represents any value it
// represents an important idea of what it's data actually is.
type GLContext interface {
}

type GLRenderable interface {
	// GLConfigs returns all possible OpenGL configurations, these are valid
	// configurations that may be used in an call to GLSetConfig.
	GLConfigs() []*GLConfig

	// GLSetConfig sets the OpenGL framebuffer configuration, unlike other
	// window management libraries, this action may be performed multiple
	// times.
	//
	// The config parameter must be an *GLConfig that originally came from the
	// GLConfigs() function mainly do to the fact that it must be initialized
	// internally.
	GLSetConfig(config *GLConfig)

	// GLConfig returns the currently in use *GLConfig or nil, as it was
	// previously set via an call to GLSetConfig()
	GLConfig() *GLConfig

	// GLCreateContext creates an OpenGL context for the specified OpenGL
	// version, or returns an error in the event that we cannot create an
	// context for that version.
	//
	// The flags parameter may be any combination of the predifined flags, as
	// follows:
	//
	// GLDebug, you will receive an OpenGL debug context. *
	//
	// GLForwardCompatible, you will receive an OpenGL forward compatible
	// context. *
	//
	// GLCoreProfile, you will receive an OpenGL core context.
	//
	// GLCompatibilityProfile, you will receive an OpenGL compatibility
	// context.
	//
	// Only one of GLCoreProfile or GLCompatibilityProfile should be present.
	//
	// GLCompatabilityProfile will be used if neither GLCoreProfile or
	// GLCompatibilityProfile are present, or if both are present.
	//
	// * = It is not advised to use this flag in production.
	//
	// You must call GLSetConfig() before calling this function.
	//
	// If the error returned is not nil, it will either be an undefined error
	// (which may indicate an error with the user's graphics card, or drivers),
	// or will be ErrGLVersionNotSupported, which indicates that the requested
	// OpenGL version is not available.
	GLCreateContext(major, minor uint, flags GLContextFlags, share GLContext) (GLContext, error)

	// GLDestroyContext destroys the specified OpenGL context.
	//
	// The context to destroy must not be active in any thread, period.
	GLDestroyContext(c GLContext)

	// GLMakeCurrent makes the specified context the current, active OpenGL
	// context in the current operating system thread.
	//
	// To make the OpenGL context inactive, you may call this function using
	// nil as the context, which will release the context.
	//
	// This function may be called from any thread, but an OpenGL context may
	// only be active inside one thread at an time.
	GLMakeCurrent(c GLContext)

	// GLSwapBuffers swaps the front and back buffers of this Renderable.
	//
	// This function may only be called in the presence of an active OpenGL
	// context.
	//
	// If the GLConfig set previously via GLSetConfig() is not DoubleBuffered,
	// then this function is no-op.
	GLSwapBuffers()

	// GLSetVerticalSync sets the vertical refresh rate sync mode (vsync).
	//
	// This function should only be called in the presence of an active OpenGL
	// context or else the call may fail due to drivers or platforms that
	// require an active context (e.g. Mesa).
	GLSetVerticalSync(mode VSyncMode)
}
