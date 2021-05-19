// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!android dragonfly openbsd
// +build !noop

package driver

import (
	"github.com/oakmound/oak/v3/shiny/driver/x11driver"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	x11driver.Main(f)
}

func monitorSize() (int, int) {
	return 0, 0
}
