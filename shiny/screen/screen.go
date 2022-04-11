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
	"image/draw"

	"golang.org/x/image/math/f64"
)

// Screen creates Images, Textures and Windows.
type Screen interface {
	// NewImage returns a new Image for this screen.
	NewImage(size image.Point) (Image, error)

	// NewTexture returns a new Texture for this screen.
	NewTexture(size image.Point) (Texture, error)

	// NewWindow returns a new Window for this screen.
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

	// Scale scales the sub-Texture defined by src and sr to the destination
	// (the method receiver), such that sr in src-space is mapped to dr in
	// dst-space.
	Scale(dr image.Rectangle, src Texture, sr image.Rectangle, op draw.Op)

	// Upload uploads the sub-Buffer defined by src and sr to the destination
	// (the method receiver), such that sr.Min in src-space aligns with dp in
	// dst-space. The destination's contents are overwritten; the draw operator
	// is implicitly draw.Src.
	//
	// It is valid to upload a Buffer while another upload of the same Buffer
	// is in progress, but a Buffer's image.RGBA pixel contents should not be
	// accessed while it is uploading. A Buffer is re-usable, in that its pixel
	// contents can be further modified, once all outstanding calls to Upload
	// have returned.
	//
	// TODO: make it optional that a Buffer's contents is preserved after
	// Upload? Undoing a swizzle is a non-trivial amount of work, and can be
	// redundant if the next paint cycle starts by clearing the buffer.
	//
	// When uploading to a Window, there will not be any visible effect until
	// Publish is called.
	Upload(dp image.Point, src Image, sr image.Rectangle)

	// Publish flushes any pending Upload and Draw calls to the window, and
	// swaps the back buffer to the front.
	Publish()
}

type SimpleDrawer interface {
	// Draw draws the sub-Texture defined by src and sr to the destination (the
	// method receiver). src2dst defines how to transform src coordinates to
	// dst coordinates. For example, if src2dst is the matrix
	//
	// m00 m01 m02
	// m10 m11 m12
	//
	// then the src-space point (sx, sy) maps to the dst-space point
	// (m00*sx + m01*sy + m02, m10*sx + m11*sy + m12).
	Draw(src2dst f64.Aff3, src Texture, sr image.Rectangle, op draw.Op)
}
