// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ((linux && !android) || dragonfly || openbsd) && !nooswindow
// +build linux,!android dragonfly openbsd
// +build !nooswindow

package driver

import (
	"github.com/oakmound/oak/v4/shiny/driver/x11driver"
	"github.com/oakmound/oak/v4/shiny/screen"
)

func main(f func(screen.Screen)) {
	x11driver.Main(f)
}

type Window = x11driver.Window
