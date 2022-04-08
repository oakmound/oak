// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package windriver

import (
	"fmt"
	"image"

	"github.com/oakmound/oak/v3/shiny/driver/common"
	"github.com/oakmound/oak/v3/shiny/driver/internal/win32"
	"github.com/oakmound/oak/v3/shiny/screen"
)

type screenImpl struct {
	screenHWND win32.HWND
}

func newScreen(hwnd win32.HWND) *screenImpl {
	return &screenImpl{
		screenHWND: hwnd,
	}
}

func (*screenImpl) NewImage(size image.Point) (screen.Image, error) {
	return common.NewImage(size), nil
}

func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	return common.NewImage(size), nil
}

func (s *screenImpl) NewWindow(opts screen.WindowGenerator) (screen.Window, error) {
	w := &Window{}

	var err error
	w.hwnd, err = win32.NewWindow(s.screenHWND, opts)
	w.style = win32.WS_VISIBLE | win32.WS_CLIPSIBLINGS | win32.WS_OVERLAPPEDWINDOW
	w.exStyle = win32.WS_EX_WINDOWEDGE
	if opts.TopMost {
		w.exStyle |= win32.WS_EX_TOPMOST
		w.topMost = true
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	windowLock.Lock()
	allWindows[w.hwnd] = w
	windowLock.Unlock()

	err = win32.ResizeClientRect(w.hwnd, opts)
	if err != nil {
		return nil, err
	}

	if opts.Fullscreen {
		err = w.SetFullScreen(true)
		if err != nil {
			return nil, err
		}
	}
	if opts.Borderless {
		err = w.SetBorderless(true)
		if err != nil {
			return nil, err
		}
	}

	win32.Show(w.hwnd)
	return w, nil
}
