// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chippy

import (
	"io"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
)

var destroyCallbacks []*func()

// addDestroyCallback adds the given function as a destroyer. Each destroy
// callback is called in the order it was added when chippy.Exit() is called,
// and before backend_Destroy is called.
func addDestroyCallback(c *func()) {
	removeDestroyCallback(c) // In case it's already been added
	destroyCallbacks = append(destroyCallbacks, c)
}

// removeDestroyCallback removes the previously added destroy callback, c. Only
// exactly identical pointers can be removed, so take care to not write code
// like:
//  removeDestroyCallback(&c)
func removeDestroyCallback(c *func()) {
	for i := 0; i < len(destroyCallbacks); i++ {
		if destroyCallbacks[i] == c {
			// Remove it
			destroyCallbacks = append(destroyCallbacks[:i], destroyCallbacks[i+1:]...)
			break
		}
	}
}

var (
	globalLock sync.RWMutex

	// Tells whether chippy has been previously Init()
	isInit bool

	// Tells whether a previous call to Init() failed
	initError error

	mainLoopFrames chan func() bool
)

// Dispatches the function on the dispatcher thread, and waits for the operation to complete before
// returning.
func dispatch(f func()) {
	done := make(chan bool, 1)
	mainLoopFrames <- func() bool {
		f()
		done <- true
		return true
	}
	<-done
}

// IsInit returns whether Chippy has been initialized via a previous call to
// Init().
//
// IsInit() returns false if Destroy() was previously called.
func IsInit() bool {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return isInit
}

// Helper to panic unless previously initialized
func panicUnlessInit() {
	globalLock.RLock()
	defer globalLock.RUnlock()

	if !IsInit() {
		panic("Chippy must be initialized before calling this; Use Init() properly!")
	}
}

var theLogger *log.Logger

func logger() *log.Logger {
	globalLock.RLock()
	defer globalLock.RUnlock()

	return theLogger
}

// SetDebugOutput specifies the io.Writer that debug output will be written to (ioutil.Discard by
// default).
func SetDebugOutput(w io.Writer) {
	globalLock.Lock()
	defer globalLock.Unlock()

	theLogger = log.New(w, "chippy: ", log.Ltime|log.Lshortfile)
}

func init() {
	SetDebugOutput(ioutil.Discard)
}

// Init initializes Chippy, returning an error if there is a problem initializing some
// lower level part of Chippy, if an error was returned, it is disallowed to call any
// other Chippy functions. (And any attempt to do so will cause the program to panic.)
func Init() error {
	globalLock.Lock()
	defer globalLock.Unlock()

	if isInit == false {
		mainLoopFrames = make(chan func() bool, 32)

		// Now we try and initialize the backend, which may fail due to user configurations
		// or something of the sort (dumb user tries to run application on Linux box without
		// any working X11 server or something silly)
		err := backend_Init()
		if err != nil {
			initError = err
			return initError
		}

		// If we made it this far, Chippy should be loaded and ready, and everything is up to
		// the backend to handle things properly now
		isInit = true

		return nil
	}
	return nil
}

// Exit will exit the main loop previously entered via MainLoop(). All windows
// will be destroyed, etc.
//
// You may call Init() again after calling this function should you want to
// re-gain access to the API.
func Exit() {
	globalLock.Lock()
	defer globalLock.Unlock()

	if isInit == true {
		// Firstly, we call each destroy callback, chippyAccess is explicitly unlocked here
		globalLock.Unlock()
		for _, callback := range destroyCallbacks {
			(*callback)()
		}
		globalLock.Lock()
		backend_Destroy()
		isInit = false
		initError = nil
		destroyCallbacks = []*func(){}
	}

	mainLoopFrames <- func() bool {
		return false
	}
}

// MainLoopFrames returns an channel of functions which return an boolean
// status as to whether you should continue running the 'main loop'.
//
// Typically you would not use this function and would instead use the
// MainLoop() function.
//
// This is for advanced users where the main loop is required to be shared with
// some other external library. I.e. this allows for communicative main loop
// handling.
//
// See the MainLoop() function source code in chippy.go for an example of using
// this function properly.
//
// This function should only really be called once (the same channel is always
// returned).
func MainLoopFrames() chan func() bool {
	panicUnlessInit()
	return mainLoopFrames
}

// MainLoop enters chippy's main loop.
//
// This function *must* be called on the main thread (due to the restrictions
// that some platforms place, I.e. Cocoa on OS-X)
//
// It's best to place this function inside either your init or main function.
//
// This function will not return until chippy.Exit() is called.
//
// If chippy is not initialized (via an previous call to the Init() function)
// then a panic will occur.
func MainLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	frames := MainLoopFrames()
	for {
		frame := <-frames
		if !frame() {
			return
		}
	}
}
