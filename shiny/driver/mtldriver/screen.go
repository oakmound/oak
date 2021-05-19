// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package mtldriver

import (
	"image"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/oakmound/oak/v3/shiny/screen"
)

// screenImpl implements screen.Screen.
type screenImpl struct {
	newWindowCh chan newWindowReq
}

func (*screenImpl) NewImage(size image.Point) (screen.Image, error) {
	return &bufferImpl{
		rgba: image.NewRGBA(image.Rectangle{Max: size}),
	}, nil
}

func (*screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	return &textureImpl{
		rgba: image.NewRGBA(image.Rectangle{Max: size}),
	}, nil
}

func (s *screenImpl) NewWindow(opts screen.WindowGenerator) (screen.Window, error) {
	respCh := make(chan newWindowResp)
	s.newWindowCh <- newWindowReq{
		opts:   opts,
		respCh: respCh,
	}
	glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on newWindowCh.
	resp := <-respCh
	return resp.w, resp.err
}
