// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tests

// Test - Prints the extents of an window's region.
package main

import (
	"azul3d.org/chippy.v1"
	"log"
	"os"
)

func program() {
	defer chippy.Exit()

	window := chippy.NewWindow()
	window.SetVisible(false)
	err := window.Open(chippy.DefaultScreen())
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	// Print what the window extents are
	log.Println(window.Extents())
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
