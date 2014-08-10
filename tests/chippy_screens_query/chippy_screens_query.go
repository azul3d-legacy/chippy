// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Queries and lists all available screens, and their modes.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"azul3d.org/chippy.v1"
)

func program() {
	defer chippy.Exit()

	screens := chippy.Screens()
	log.Printf("There are %d screens.\n", len(screens))
	log.Println("Default screen:", chippy.DefaultScreen())

	for i, screen := range screens {
		log.Printf("\nScreen %d - %s", i, screen)

		for i, mode := range screen.Modes() {
			log.Printf("    -> ScreenMode %d - %s", i, mode)
		}
	}

	fmt.Printf("Change Screen: #")
	var screen int
	_, err := fmt.Scanln(&screen)
	if err != nil {
		log.Fatal(err)
	}

	if screen < 0 || screen > len(screens)-1 {
		log.Fatal("Incorrect screen number.")
	}
	chosenScreen := screens[screen]

	fmt.Printf("Change Screen #%d to mode: #", screen)
	var mode int
	_, err = fmt.Scanln(&mode)
	if err != nil {
		log.Fatal(err)
	}

	if mode < 0 || mode > len(screens[screen].Modes())-1 {
		log.Fatal("Incorrect screen number.")
	}
	chosenMode := chosenScreen.Modes()[mode]

	// Ensure that we restore the screen to it's original state before exiting
	defer chosenScreen.Restore()

	// Change screen mode
	chosenScreen.SetMode(chosenMode)

	log.Println("Waiting 15 seconds...")
	<-time.After(15 * time.Second)

	// Disable this with chippy.SetAutoRestoreOriginalScreenMode(false)
	log.Println("Original resolution should restore now")
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
