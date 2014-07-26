// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chippy implements cross platform window management, and window
// rendering access.
//
// Thread Safety
//
// Chippy is thread safe, and can be fully used from within multiple
// goroutines without any worry about operating system threads, or locks, etc.
//
// It should be explicitly noted that while Chippy and it's API's are thread
// safe, anything to do with OpenGL needs special care regarding thread
// safety.
//
// OpenGL Support
//
// Creating both new and old style OpenGL contexts is supported (this allows
// creating an OpenGL context of any version). Many platform specific OpenGL
// functions are abstracted away for you (such as WGL, GLX, etc extensions).
// Shared OpenGL contexts, multisampling, vertical sync toggling, etc are all
// supported.
//
// Chippy works with all OpenGL wrappers, it does not provide any OpenGL
// wrappers itself (although azul3d.org/v1/native/gl has some good ones).
//
// Although we handle the platform-specific parts of OpenGL for you, no magic
// is performed: OpenGL still uses thread local storage so when working with
// OpenGL's API you'll need to utilize runtime.LockOSThread() properly.
//
// Microsoft Windows FAQ
//
// What versions of Windows are supported?
//  Chippy requires Windows XP or higher.
//
//  It might also work on Windows 2000 Professional/Server editions, but
//  support for these version is not tested actively.
//
// How do I add an application icon to my program?
//  You can place .syso files with the source of your main package, and the
//  6l/8l linker will link that file with your program.
//
//  Take an look at the "app.rc" file inside the chippy/tests/data folder for
//  more information. Also look at the single window test located in the
//  chippy/tests/chippy_window_single directory for an example of this.
//
// How do I stop the command prompt from appearing when my application starts?
//  You can stop the terminal from appearing by using the 8l/6l linker flag
//  "-H windowsgui" on your 'go install' command, like so:
//
//  go install -ldflags "-H windowsgui" path/to/pkg
//
// Linux-X11 FAQ
//
// What X extensions are needed?
//  GLX 1.4 (required, for OpenGL tasks, 1.4 is needed for multisampling)
//  XRandR 1.2 (optional, used for screen-mode switching, etc)
//  XInput 2.0 (required, for raw mouse input)
//  XKB 1.0 (required, for keyboard input)
//
// What about a pure Wayland client?
//  A pure Wayland implementation would be interesting and could be enabled
//  using a build-tag until Wayland becomes more main-stream, but for now
//  Wayland does have the ability to still run X applications so Chippy does
//  still work on Wayland.
//
//  A pure Wayland client is not a priority right now, but we are open to
//  working with a contributor who would like to add this feature.
//
package chippy
