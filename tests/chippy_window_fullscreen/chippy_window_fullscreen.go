// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens an single window, changes it's fullscreen property.
package main

import (
	"log"
	"os"
	"time"

	"azul3d.org/chippy.v1"
)

func program() {
	defer chippy.Exit()

	window := chippy.NewWindow()

	window.SetFullscreen(true)

	// Actually open the windows
	screen := chippy.DefaultScreen()
	err := window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	// Print out what they currently has property-wise
	log.Println(window)

	log.Println("Waiting 10 seconds...")
	<-time.After(10 * time.Second)
	window.SetFullscreen(false)

	log.Println("Waiting 10 seconds...")
	<-time.After(10 * time.Second)
	window.SetFullscreen(true)

	log.Println("Waiting 5 seconds...")
	<-time.After(5 * time.Second)
	window.SetFullscreen(false)

	// Just wait an while so they can enjoy the window
	log.Println("Waiting 15 seconds...")
	<-time.After(15 * time.Second)
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
