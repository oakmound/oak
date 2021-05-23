// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package mtldriver

import (
	"image"
	"image/color"
	"log"

	"dmitri.shuralyov.com/gpu/mtl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/oakmound/oak/v3/shiny/driver/internal/drawer"
	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/driver/internal/lifecycler"
	"github.com/oakmound/oak/v3/shiny/driver/mtldriver/internal/coreanim"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/size"
)

// windowImpl implements screen.Window.
type windowImpl struct {
	device          mtl.Device
	window          *glfw.Window
	releaseWindowCh chan releaseWindowReq
	moveWindowCh    chan moveWindowReq
	ml              coreanim.MetalLayer
	cq              mtl.CommandQueue

	event.Deque
	lifecycler lifecycler.State

	rgba    *image.RGBA
	texture mtl.Texture // Used in Publish.
}

func (w *windowImpl) MoveWindow(x, y, width, height int32) error {
	w.moveWindowCh <- moveWindowReq{
		window: w.window,
		x:      int(x),
		y:      int(y),
		width:  int(width),
		height: int(height),
	}
	return nil
}

func (w *windowImpl) Release() {
	respCh := make(chan struct{})
	w.releaseWindowCh <- releaseWindowReq{
		window: w.window,
		respCh: respCh,
	}
	glfw.PostEmptyEvent() // Break main loop out of glfw.WaitEvents so it can receive on releaseWindowCh.
	<-respCh
}

func (w *windowImpl) NextEvent() interface{} {
	e := w.Deque.NextEvent()
	if sz, ok := e.(size.Event); ok {
		// TODO(dmitshur): this is the best place/time/frequency to do this
		//                 I've found so far, but see if it can be even better

		// Set drawable size, create backing image and texture.
		w.ml.SetDrawableSize(sz.WidthPx, sz.HeightPx)
		w.rgba = image.NewRGBA(image.Rectangle{Max: image.Point{X: sz.WidthPx, Y: sz.HeightPx}})
		w.texture = w.device.MakeTexture(mtl.TextureDescriptor{
			PixelFormat: mtl.PixelFormatRGBA8UNorm,
			Width:       sz.WidthPx,
			Height:      sz.HeightPx,
			StorageMode: mtl.StorageModeManaged,
		})
	}
	return e
}

func (w *windowImpl) Publish() screen.PublishResult {
	// Copy w.rgba pixels into a texture.
	region := mtl.RegionMake2D(0, 0, w.texture.Width, w.texture.Height)
	bytesPerRow := 4 * w.texture.Width
	w.texture.ReplaceRegion(region, 0, &w.rgba.Pix[0], uintptr(bytesPerRow))

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

func (w *windowImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle) {
	// This panics
	//draw.Draw(w.rgba, sr.Sub(sr.Min).Add(dp), src.RGBA(), sr.Min, draw.Src)
}

func (w *windowImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	draw.Draw(w.rgba, dr, &image.Uniform{src}, image.Point{}, op)
}

func (w *windowImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op) {
	draw.NearestNeighbor.Transform(w.rgba, src2dst, src.(*textureImpl).rgba, sr, op, nil)
}

func (w *windowImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {
	draw.NearestNeighbor.Transform(w.rgba, src2dst, &image.Uniform{src}, sr, op, nil)
}

func (w *windowImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Copy(w, dp, src, sr, op)
}

func (w *windowImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	drawer.Scale(w, dr, src, sr, op)
}
