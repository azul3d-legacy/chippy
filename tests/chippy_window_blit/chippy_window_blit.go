// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens an windows, uses blitting to copy pixels on to it.
package main

import (
	"azul3d.org/chippy.v1"
	"flag"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func program() {
	defer chippy.Exit()

	var err error

	// Load the image that we'll use for the window icon
	file, err := os.Open("src/azul3d.org/chippy.v1/tests/data/chippy_720x320.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	rgba, ok := img.(*image.RGBA)
	if !ok {
		// Need to convert to RGBA image
		b := img.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	}

	window := chippy.NewWindow()
	window.SetSize(720, 320)
	window.SetTransparent(true)
	window.SetDecorated(false)
	window.SetAlwaysOnTop(true)

	// Actually open the windows
	screen := chippy.DefaultScreen()

	// Center the window on the screen
	window.SetPositionCenter(screen)

	err = window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	// Print out what they currently has property-wise
	log.Println(window)

	events := window.Events()
	defer window.CloseEvents(events)

	var (
		measureBlitSpeed   = time.After(10 * time.Second)
		measuringBlitSpeed bool
		numBlits           int
		totalBlitTime      time.Duration
	)

	for {
		// In order to clear a rectangle on the window, this is much faster
		// than PixelBlit() with an fully transparent image:
		//
		// window.PixelClear(image.Rect(0, 0, 30, 30))

		// Blit the image to the window, at x=0, y=0, blitting the entire image
		start := time.Now()
		window.PixelBlit(0, 0, rgba)
		blitTime := time.Since(start)

		numBlits++
		totalBlitTime += blitTime

		log.Printf("PixelBlit(): %v (%v average)", blitTime, totalBlitTime/time.Duration(numBlits))

		if measuringBlitSpeed {
			// Say that PixelBlit() above takes more time than it takes for
			// your high precision mouse to send two or more events -- we would
			// end up with the event buffer filling and missing events. The for
			// statement below solves this.
			//
			// Also, "why not just use a simple for loop or range, like so?":
			//
			//  for len(events) > 0
			//
			// If your high precision mouse was to constantly send position
			// events (because you are constantly wiggling the mouse) will the
			// for loop ever end and let you render? No. Storing the variable
			// n below lets us know for certain it will eventually end and let
			// us PixelBlit() to the window once again in the near future once
			// we've read at max cap(events) from the channel.
			//
			// In practice, you would not tie your rendering loop and your
			// event loop together for this exact reason. It would be much
			// better to handle events in a seperate goroutine and communicate
			// with the render loop over a channel.
			for n := len(events); n > 0; n-- {
				e := <-events
				switch e.(type) {
				case chippy.CloseEvent:
					chippy.Exit()
					goto stats
				}
			}
			continue
		}

		// Wait for an paint event
		gotPaintEvent := false
	loop:
		for !gotPaintEvent {
			select {
			case <-measureBlitSpeed:
				measuringBlitSpeed = true
				break loop

			case e := <-events:
				switch e.(type) {
				case chippy.PaintEvent:
					log.Println(e)
					gotPaintEvent = true

				case chippy.CloseEvent:
					chippy.Exit()
					goto stats

				default:
					// We don't care about whatever event this is.
					break
				}
			}
		}
	}

stats:
	log.Printf("%d PixelBlit() over %v\n", numBlits, totalBlitTime)
	log.Printf("Average blit time: %v\n", totalBlitTime/time.Duration(numBlits))
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
