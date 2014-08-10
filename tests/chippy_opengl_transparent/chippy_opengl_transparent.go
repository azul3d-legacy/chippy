// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens an transparent window and uses OpenGL 1.5 rendering.
package main

import (
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"azul3d.org/chippy.v1"
	"azul3d.org/clock.v1"
	"azul3d.org/keyboard.v1"
	opengl "azul3d.org/native/gl.v1"
)

var (
	gl      *opengl.Context
	rot     float64
	window  *chippy.Window
	glClock *clock.Clock
)

// Alternative for gluPerspective.
func gluPerspective(gl *opengl.Context, fovY, aspect, zNear, zFar float64) {
	fH := math.Tan(fovY/360*math.Pi) * zNear
	fW := fH * aspect
	gl.Frustum(-fW, fW, -fH, fH, zNear, zFar)
}

func resizeScene(width, height int) {
	gl.Viewport(0, 0, uint32(width), uint32(height)) // Reset The Current Viewport And Perspective Transformation
	gl.MatrixMode(opengl.PROJECTION)
	gl.LoadIdentity()
	gluPerspective(gl, 45.0, float64(width)/float64(height), 0.1, 100.0)
	gl.MatrixMode(opengl.MODELVIEW)
}

func initScene() {
	gl.Enable(opengl.BLEND)
	gl.Enable(opengl.ALPHA_TEST)
	gl.Enable(opengl.DEPTH_TEST)

	gl.BlendFunc(opengl.SRC_ALPHA, opengl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 0.3)

	gl.ClearDepth(1.0)
	gl.ShadeModel(opengl.SMOOTH)

	width, height := window.Size()
	resizeScene(int(width), int(height))
}

func renderScene() {
	// Clear The Screen And The Depth Buffer
	gl.Clear(uint32(opengl.COLOR_BUFFER_BIT | opengl.DEPTH_BUFFER_BIT))
	gl.LoadIdentity() // Reset The View

	// Move into the screen 6.0 units.
	gl.Translatef(0, 0, -6.0)

	// We have smooth color mode on, this will blend across the vertices.
	// Draw a triangle rotated on the Y axis.
	gl.Rotatef(float32(rot), 0.0, 1.0, 0.0) // Rotate
	gl.Begin(opengl.POLYGON)                // Start drawing a polygon
	gl.Color3f(1.0, 0.0, 0.0)               // Red
	gl.Vertex3f(0.0, 1.0, 0.0)              // Top
	gl.Color3f(0.0, 1.0, 0.0)               // Green
	gl.Vertex3f(1.0, -1.0, 0.0)             // Bottom Right
	gl.Color3f(0.0, 0.0, 1.0)               // Blue
	gl.Vertex3f(-1.0, -1.0, 0.0)            // Bottom Left
	gl.End()                                // We are done with the polygon

	gl.Flush()

	// Determine time since frame began
	delta := glClock.Delta()

	// Increase the rotation by 90 degrees each second
	rot += 90.0 * delta.Seconds()

	// Clamp the result to 360 degrees
	if rot >= 360 {
		rot = 0
	}
	if rot < 0 {
		rot = 360
	}
}

func toggleVerticalSync() {
	vsync := window.GLVerticalSync()

	switch vsync {
	case chippy.NoVerticalSync:
		vsync = chippy.VerticalSync

	case chippy.VerticalSync:
		vsync = chippy.AdaptiveVerticalSync

	case chippy.AdaptiveVerticalSync:
		vsync = chippy.NoVerticalSync
	}

	log.Println(vsync)
	window.GLSetVerticalSync(vsync)
}

var MSAA = true

func toggleMSAA() {
	if MSAA {
		MSAA = false
		gl.Disable(opengl.MULTISAMPLE)
	} else {
		MSAA = true
		gl.Enable(opengl.MULTISAMPLE)
	}
	log.Println("MSAA enabled?", MSAA)
}

func program() {
	defer chippy.Exit()

	window = chippy.NewWindow()

	// Make it transparent
	window.SetTransparent(true)

	// Actually open the windows
	screen := chippy.DefaultScreen()
	err := window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	// Print some instructions for the user
	log.Println("Instructions:")
	log.Println("v key - Toggle Vertical Sync")
	log.Println("m key - Toggle Multi Sample Anti Aliasing")
	log.Println("b key - Toggle OpenGL call batching")

	// Choose an buffer format, these include things like double buffering, bytes per pixel, number of depth bits, etc.
	configs := window.GLConfigs()

	// See documentation for this function and vars to see how it determines the 'best' format
	transparentCfg := *chippy.GLWorstConfig
	transparentCfg.Transparent = true
	bestConfig := chippy.GLChooseConfig(configs, &transparentCfg, chippy.GLBestConfig)
	if bestConfig == nil {
		log.Fatal("Could not find a proper transparent OpenGL configuration.")
	}
	window.GLSetConfig(bestConfig)

	// Print out all the formats, and which one we determined to be the 'best'.
	log.Println("\nChosen configuration:")
	log.Println(bestConfig)

	// All OpenGL related calls must occur in the same OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Create an OpenGL context with the OpenGL version we wish
	context, err := window.GLCreateContext(1, 5, chippy.GLCoreProfile, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make the context current in this OS thread
	window.GLMakeCurrent(context)

	// Create an opengl.Context (which provides API access to an existing OpenGL context), for each
	// OpenGL context you wish to interace.
	//
	// We only make one here (as we are only using one context).
	gl = opengl.New()
	if gl == nil {
		log.Fatal("You have no support for OpenGL 1.5!")
	}
	log.Println(gl.GetError())

	// Initialize some things
	initScene()

	// We'll use this glClock for timing things
	glClock = clock.New()

	// Start an goroutine to display statistics
	go func() {
		for {
			<-time.After(1 * time.Second)

			// Print our FPS and average FPS
			log.Printf("FPS: %4.3f\tAverage: %4.3f\tDeviation: %f\n", glClock.FrameRate(), glClock.AverageFrameRate(), glClock.FrameRateDeviation())
		}
	}()

	events := window.Events()
	defer window.CloseEvents(events)

	// Begin our rendering loop
	for !window.Destroyed() {
		// Inform the clock that an new frame has begun
		glClock.Tick()

		for i := 0; i < len(events); i++ {
			e := <-events
			switch ev := e.(type) {
			case chippy.ResizedEvent:
				resizeScene(ev.Width, ev.Height)

			case keyboard.StateEvent:
				if ev.State == keyboard.Down {
					switch ev.Key {
					case keyboard.V:
						toggleVerticalSync()
					case keyboard.M:
						toggleMSAA()
					case keyboard.B:
						gl.SetBatching(!gl.Batching())
						log.Println("Batching?", gl.Batching())
					}
				}

			case chippy.CloseEvent:
				return
			}
		}

		// Render the scene
		renderScene()

		// Swap the display buffers
		window.GLSwapBuffers()
	}
}

func main() {
	log.SetFlags(0)

	// Enable debug output
	chippy.SetDebugOutput(os.Stdout)

	// Initialize Chippy
	err := chippy.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Start program
	go program()

	// Enter main loop
	chippy.MainLoop()
}
