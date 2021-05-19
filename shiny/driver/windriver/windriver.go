// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package windriver

import (
	"github.com/oakmound/oak/v3/shiny/driver/internal/errscreen"
	"github.com/oakmound/oak/v3/shiny/driver/internal/win32"
	"github.com/oakmound/oak/v3/shiny/screen"
)

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the Screen, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(screen.Screen)) {
	screenHWND, err := win32.NewScreen()
	if err != nil {
		f(errscreen.Stub(err))
		return
	}
	screen := newScreen(screenHWND)
	if err := win32.Main(screenHWND, func() { f(screen) }); err != nil {
		f(errscreen.Stub(err))
	}
}
