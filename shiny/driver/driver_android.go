// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android && !nooswindow
// +build android,!nooswindow

package driver

import (
	"github.com/oakmound/oak/v4/shiny/driver/androiddriver"
	"github.com/oakmound/oak/v4/shiny/screen"
)

func main(f func(screen.Screen)) {
	androiddriver.Main(f)
}

type Window = androiddriver.Screen
