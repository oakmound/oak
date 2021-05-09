// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errscreen provides a stub Screen implementation.
package errscreen

import (
	"image"

	"github.com/oakmound/oak/v2/shiny/screen"
)

// Stub returns a Screen whose methods all return the given error.
func Stub(err error) screen.Screen {
	return stub{err}
}

type stub struct {
	err error
}

func (s stub) NewImage(size image.Point) (screen.Image, error)              { return nil, s.err }
func (s stub) NewTexture(size image.Point) (screen.Texture, error)          { return nil, s.err }
func (s stub) NewWindow(opts screen.WindowGenerator) (screen.Window, error) { return nil, s.err }
