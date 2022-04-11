// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !noop
// +build darwin,!noop

package driver

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/oakmound/oak/v3/shiny/driver/mtldriver"
	"github.com/oakmound/oak/v3/shiny/screen"
)

func main(f func(screen.Screen)) {
	mtldriver.Main(f)
}

type Window = mtldriver.Window
