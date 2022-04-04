// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package mtldriver

import (
	"fmt"
	"image"
	"log"

	"dmitri.shuralyov.com/gpu/mtl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/driver/internal/lifecycler"
	"github.com/oakmound/oak/v3/shiny/driver/mtldriver/internal/coreanim"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/mobile/event/size"
)

// Window implements screen.Window.
type Window struct {
	device mtl.Device
	window *glfw.Window
	chans  windowRequestChannels
	ml     coreanim.MetalLayer
	cq     mtl.CommandQueue

	event.Deque
	lifecycler lifecycler.State

	bgra    *BGRA
	texture mtl.Texture // Used in Publish.

	title      string
	fullscreen bool
	borderless bool

	w, h int
	x, y int
}

func (w *Window) HideCursor() error {
	w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	return nil
}

func (w *Window) SetBorderless(borderless bool) error {
	if w.borderless == borderless {
		return nil
	}
	w.borderless = borderless
	respCh := make(chan struct{})
	w.chans.updateCh <- updateWindowReq{
		setBorderless: &borderless,
		window:        w.window,
		x:             w.x,
		y:             w.y,
		width:         w.w,
		height:        w.h,
		respCh:        respCh,
	}
	glfw.PostEmptyEvent()
	<-respCh
	return nil
}

func (w *Window) SetFullScreen(full bool) error {
	if w.fullscreen == full {
		return nil
	}
	w.fullscreen = full
	if full {
		w.x, w.y = w.window.GetPos()
	}
	respCh := make(chan struct{})
	w.chans.updateCh <- updateWindowReq{
		setFullscreen: &full,
		window:        w.window,
		x:             w.x,
		y:             w.y,
		width:         w.w,
		height:        w.h,
		respCh:        respCh,
	}
	glfw.PostEmptyEvent()
	<-respCh
	return nil
}

func (w *Window) MoveWindow(x, y, width, height int) error {
	respCh := make(chan struct{})
	w.x = x
	w.y = y
	w.w = width
	w.h = height
	w.chans.updateCh <- updateWindowReq{
		window: w.window,
		setPos: true,
		x:      w.x,
		y:      w.y,
		width:  w.w,
		height: w.h,
		respCh: respCh,
	}
	glfw.PostEmptyEvent()
	<-respCh
	return nil
}

func (w *Window) GetCursorPosition() (x, y float64) {
	return w.window.GetCursorPos()
}

func (w *Window) Release() {
	respCh := make(chan struct{})
	w.chans.releaseCh <- releaseWindowReq{
		window: w.window,
		respCh: respCh,
	}
	glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on releaseWindowCh.
	<-respCh
}

func (w *Window) SetTitle(title string) error {
	respCh := make(chan struct{})
	w.chans.updateCh <- updateWindowReq{
		window: w.window,
		title:  &title,
		respCh: respCh,
	}
	glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on releaseWindowCh.
	<-respCh
	return nil
}

type attribPair struct {
	key glfw.Hint
	val int
}

func (w *Window) SetTopMost(topMost bool) error {
	respCh := make(chan struct{})
	val := glfw.True
	if !topMost {
		val = glfw.False
	}
	w.chans.updateCh <- updateWindowReq{
		window: w.window,
		attribs: []attribPair{{
			key: glfw.Floating,
			val: val,
		}},
		respCh: respCh,
	}
	glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on releaseWindowCh.
	<-respCh
	return nil
}

func (w *Window) SetIcon(image.Image) error {
	// TODO: the problem here is that this ^ takes a path, because windows
	// wants a path, where glfw wants an image.Image (or set of them).
	// for v4, standardize this interface.
	// w.window.SetIcon()
	return fmt.Errorf("unimplemented")
}

func (w *Window) NextEvent() interface{} {
	e := w.Deque.NextEvent()
	if sz, ok := e.(size.Event); ok {
		// TODO(dmitshur): this is the best place/time/frequency to do this
		//                 I've found so far, but see if it can be even better

		// Set drawable size, create backing image and texture.
		w.ml.SetDrawableSize(sz.WidthPx, sz.HeightPx)
		w.bgra = NewBGRA(image.Rectangle{Max: image.Point{X: sz.WidthPx, Y: sz.HeightPx}})
		w.texture = w.device.MakeTexture(mtl.TextureDescriptor{
			PixelFormat: mtl.PixelFormatRGBA8UNorm,
			Width:       sz.WidthPx,
			Height:      sz.HeightPx,
			StorageMode: mtl.StorageModeManaged,
		})
	}
	return e
}

func (w *Window) Publish() screen.PublishResult {
	// Copy w.rgba pixels into a texture.
	region := mtl.RegionMake2D(0, 0, w.texture.Width, w.texture.Height)
	bytesPerRow := 4 * w.texture.Width
	w.texture.ReplaceRegion(region, 0, &w.bgra.Pix[0], uintptr(bytesPerRow))

	drawable, err := w.ml.NextDrawable()
	if err != nil {
		log.Println("Window.Publish: couldn't get the next drawable:", err)
		return screen.PublishResult{}
	}

	cb := w.cq.MakeCommandBuffer()

	// Copy the texture into the drawable.
	bce := cb.MakeBlitCommandEncoder()
	bce.CopyFromTexture(
		w.texture, 0, 0, mtl.Origin{}, mtl.Size{
			Width:  w.texture.Width,
			Height: w.texture.Height,
			Depth:  1,
		},
		drawable.Texture(), 0, 0, mtl.Origin{})
	bce.EndEncoding()

	cb.PresentDrawable(drawable)
	cb.Commit()

	return screen.PublishResult{}
}
