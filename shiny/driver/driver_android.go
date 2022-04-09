// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android

//go:build !nooswindow
// +build !nooswindow

package driver

import (
	"github.com/oakmound/oak/v3/shiny/driver/androiddriver"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	androiddriver.Main(f)
}

func monitorSize() (int, int) {
	// GetSystemMetrics syscall
	return 0, 0
}

type Window = androiddriver.Screen
