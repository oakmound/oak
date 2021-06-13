// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package x11driver provides the X11 driver for accessing a screen.
package x11driver

// TODO: figure out what to say about the responsibility for users of this
// package to check any implicit dependencies' LICENSEs. For example, the
// driver might use third party software outside of golang.org/x, like an X11
// or OpenGL library.

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/xgb/render"
	"github.com/BurntSushi/xgb/shm"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"

	"github.com/oakmound/oak/v3/shiny/driver/internal/errscreen"
	"github.com/oakmound/oak/v3/shiny/screen"
)

// Main is called by the program's main function to run the graphical
// application.
//
// It calls f on the Screen, possibly in a separate goroutine, as some OS-
// specific libraries require being on 'the main thread'. It returns when f
// returns.
func Main(f func(screen.Screen)) {
	if err := main(f); err != nil {
		f(errscreen.Stub(err))
	}
}

var mainLock sync.Mutex

func main(f func(screen.Screen)) (retErr error) {
	xutil, err := xgbutil.NewConn()
	if err != nil {
		return fmt.Errorf("x11driver: xgb.NewConn failed: %v", err)
	}
	defer func() {
		if retErr != nil {
			xevent.Quit(xutil)
		}
	}()

	mainLock.Lock()
	if err := render.Init(xutil.Conn()); err != nil {
		return fmt.Errorf("x11driver: render.Init failed: %v", err)
	}
	if err := shm.Init(xutil.Conn()); err != nil {
		return fmt.Errorf("x11driver: shm.Init failed: %v", err)
	}
	mainLock.Unlock()

	s, err := newScreenImpl(xutil)
	if err != nil {
		return err
	}
	f(s)
	// TODO: tear down the s.run goroutine? It's probably not worth the
	// complexity of doing it cleanly, if the app is about to exit anyway.
	return nil
}
