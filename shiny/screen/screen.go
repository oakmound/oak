// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package screen provides interfaces for portable two-dimensional graphics and
// input events.
//
// Screens are not created directly. Instead, driver packages provide access to
// the screen through a Main function that is designed to be called by the
// program's main function. The golang.org/x/exp/shiny/driver package provides
// the default driver for the system, such as the X11 driver for desktop Linux,
// but other drivers, such as the OpenGL driver, can be explicitly invoked by
// calling that driver's Main function. To use the default driver:
//
//	package main
//
//	import (
//		"github.com/oakmound/oak/v3/shiny/driver"
//		"github.com/oakmound/oak/v3/shiny/screen"
//		"golang.org/x/mobile/event/lifecycle"
//	)
//
//	func main() {
//		driver.Main(func(s screen.Screen) {
//			w, err := s.NewWindow(nil)
//			if err != nil {
//				handleError(err)
//				return
//			}
//			defer w.Release()
//
//			for {
//				switch e := w.NextEvent().(type) {
//				case lifecycle.Event:
//					if e.To == lifecycle.StageDead {
//						return
//					}
//					etc
//				case etc:
//					etc
//				}
//			}
//		})
//	}
//
// Each driver package provides Screen, Image, Texture and Window
// implementations that work together. Such types are interface types because
// this package is driver-independent, but those interfaces aren't expected to
// be implemented outside of drivers. For example, a driver's Window
// implementation will generally work only with that driver's Image
// implementation, and will not work with an arbitrary type that happens to
// implement the Image methods.
package screen

import (
	"image"
)

// Screen creates Images, Textures and Windows.
type Screen interface {
	// NewImage returns a new Image for this screen.
	NewImage(size image.Point) (Image, error)

	// NewTexture returns a new Texture for this screen.
	NewTexture(size image.Point) (Texture, error)

	// NewWindow returns a new Window for this screen.
	//
	// A nil opts is valid and means to use the default option values.
	NewWindow(opts WindowGenerator) (Window, error)
}

// Window is a top-level, double-buffered GUI window.
type Window interface {
	// Release closes the window.
	//
	// The behavior of the Window after Release, whether calling its methods or
	// passing it as an argument, is undefined.
	Release()

	EventDeque

	Drawer

	// Publish flushes any pending Upload and Draw calls to the window, and
	// swaps the back buffer to the front.
	Publish() PublishResult
}

// PublishResult is the result of an Window.Publish call.
type PublishResult struct {
	// BackBufferPreserved is whether the contents of the back buffer was
	// preserved. If false, the contents are undefined.
	BackBufferPreserved bool
}
