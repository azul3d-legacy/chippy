// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens an windows, animates the cursor property.
package main

import (
	"azul3d.org/chippy.v1"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"time"
)

func program() {
	defer chippy.Exit()

	window := chippy.NewWindow()
	window.SetTransparent(true)

	// Actually open the windows
	screen := chippy.DefaultScreen()
	err := window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	var cursors []*chippy.Cursor
	for i := 1; i < 25; i++ {
		// Load the image frame that we'll use for the animated cursor
		file, err := os.Open(fmt.Sprintf("src/azul3d.org/chippy.v1/tests/data/loading/%d.png", i))
		if err != nil {
			log.Fatal(err)
		}

		cursorImage, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		cursor := &chippy.Cursor{
			Image: cursorImage,
			X:     16,
			Y:     16,
		}
		window.PrepareCursor(cursor)
		cursors = append(cursors, cursor)
	}

	// Those cursors are cached for speed (and use up memory), if we wanted to stop using the
	// animated cursor, we could do the following:
	//
	// for _, cursor := range cursors {
	//     window.FreeCursor(cursor)
	// }

	// Print out what they currently has property-wise
	log.Println(window)

	events := window.Events()
	defer window.CloseEvents(events)

	frame := 0
	go func() {
		for {
			// Play back at 24 FPS
			time.Sleep((1000 / 24) * time.Millisecond)

			cursor := cursors[frame]
			window.SetCursor(cursor)
			frame += 1
			if frame >= 24 {
				frame = 0
			}
		}
	}()

	for {
		e := <-events
		switch e.(type) {
		case chippy.CloseEvent:
			return

		default:
			// We don't care about whatever event this is.
			break
		}
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
