// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Opens two windows, changes their icon properties.
package main

import (
	"azul3d.org/chippy.v1"
	"image"
	_ "image/png"
	"log"
	"os"
	"time"
)

func program() {
	defer chippy.Exit()

	// Load the image that we'll use for the window icon
	file, err := os.Open("src/azul3d.org/v1/chippy/tests/data/icon_128x128.png")
	if err != nil {
		log.Fatal(err)
	}

	icon, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	window1 := chippy.NewWindow()
	window2 := chippy.NewWindow()

	window1.SetIcon(icon)

	// Actually open the windows
	screen := chippy.DefaultScreen()
	err = window1.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	err = window2.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	<-time.After(5 * time.Second)
	window1.SetIcon(icon)
	window2.SetIcon(icon)

	// Print out what they currently has property-wise
	log.Println(window1)
	log.Println(window2)

	log.Println("Waiting 5 seconds...")
	<-time.After(5 * time.Second)

	window1.SetIcon(nil)
	window2.SetIcon(nil)

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
