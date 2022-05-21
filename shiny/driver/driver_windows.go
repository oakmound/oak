// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !nooswindow && !android
// +build !nooswindow,!android

package driver

import (
	"github.com/oakmound/oak/v4/shiny/driver/windriver"
	"github.com/oakmound/oak/v4/shiny/screen"
)

func main(f func(screen.Screen)) {
	windriver.Main(f)
}

type Window = windriver.Window
