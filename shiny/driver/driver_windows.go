// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !nooswindow

package driver

import (
	"github.com/oakmound/oak/v3/shiny/driver/windriver"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	windriver.Main(f)
}

func monitorSize() (int, int) {
	// GetSystemMetrics syscall
	return 0, 0
}
