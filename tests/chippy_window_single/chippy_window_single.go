// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens an single window on the specified screen.
package main

// Note: On Windows build with:
//   go install -ldflags "-H windowsgui" path/to/pkg
// to hide the command prompt

import (
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"azul3d.org/chippy.v1"
	"azul3d.org/keyboard.v1"
)

func program() {
	defer chippy.Exit()

	window := chippy.NewWindow()

	// Actually open the window
	screens := chippy.Screens()
	log.Printf("There are %d screens.\n", len(screens))
	log.Println("Default screen:", chippy.DefaultScreen())

	for i, screen := range screens {
		log.Printf("\nScreen %d - %s", i, screen)
	}

	fmt.Printf("Open window on screen: #")
	var screen int
	_, err := fmt.Scanln(&screen)
	if err != nil {
		log.Fatal(err)
	}

	if screen < 0 || screen > len(screens)-1 {
		log.Fatal("Incorrect screen number.")
	}
	chosenScreen := screens[screen]

	// Some events are sent before the window is opened. (Like caps lock state,
	// for instance)
	events := window.Events()
	defer window.CloseEvents(events)

	// Open the window
	err = window.Open(chosenScreen)
	if err != nil {
		log.Fatal(err)
	}

	// Print out what it currently has property-wise
	log.Println(window)

	for {
		ev := <-events
		log.Println(ev)

		typedEvent, ok := ev.(keyboard.TypedEvent)
		if ok {
			if typedEvent.Rune == 'd' {
				window.SetDecorated(!window.Decorated())
			}
			if typedEvent.Rune == 'p' {
				window.SetPosition(100, 100)
			}
			title := window.Title()
			if typedEvent.Rune == '\b' {
				// Backspace - remove one character from the end of the string.
				if len(title) > 0 {
					_, size := utf8.DecodeLastRune([]byte(title))
					window.SetTitle(title[:len(title)-size])
				}
			} else {
				window.SetTitle(title + string(typedEvent.Rune))
			}
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
