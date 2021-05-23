// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin !386,!amd64 ios
// +build !linux android
// +build !windows
// +build !openbsd

package gldriver

import (
	"fmt"
	"runtime"

	"github.com/oakmound/oak/v3/shiny/screen"
)

func newWindow(opts screen.WindowGenerator) (uintptr, error) { return 0, nil }

func moveWindow(w *windowImpl, opts screen.WindowGenerator) error { return nil }

const useLifecycler = true
const handleSizeEventsAtChannelReceive = true

func initWindow(id *windowImpl) {}
func showWindow(id *windowImpl) {}
func closeWindow(id uintptr)    {}
func drawLoop(w *windowImpl)    {}

func main(f func(screen.Screen)) error {
	return fmt.Errorf("gldriver: unsupported GOOS/GOARCH %s/%s", runtime.GOOS, runtime.GOARCH)
}
