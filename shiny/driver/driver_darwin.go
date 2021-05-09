// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package driver

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/oakmound/oak/v2/shiny/driver/gldriver"
	"github.com/oakmound/oak/v2/shiny/screen"
)

func main(f func(screen.Screen)) {
	gldriver.Main(f)
}

var (
	sysProfRegex = regexp.MustCompile(`Resolution: (\d)* x (\d)*`)
)

func monitorSize() (int, int) {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").CombinedOutput()
	if err != nil {
		return 0, 0
	}
	found := sysProfRegex.FindAll(out, -1)
	if len(found) == 0 {
		return 0, 0
	}
	if len(found) != 1 {
		fmt.Println("Found multiple screens", len(found))
	}
	first := found[0]
	first = bytes.TrimPrefix(first, []byte("Resolution: "))
	dims := bytes.Split(first, []byte(" x "))
	if len(dims) != 2 {
		return 0, 0
	}
	w, err := strconv.Atoi(string(dims[0]))
	if err != nil {
		return 0, 0
	}
	h, err := strconv.Atoi(string(dims[1]))
	if err != nil {
		return 0, 0
	}
	return w, h
}
