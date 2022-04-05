//go:build js
// +build js

package jsdriver

import (
	"fmt"
	"image"
	"syscall/js"

	"github.com/oakmound/oak/v3/shiny/driver/common"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

func Main(f func(screen.Screen)) {
	f(&screenImpl{})
}

type screenImpl struct {
	windows []*Window
}

func (s *screenImpl) NewImage(size image.Point) (screen.Image, error) {
	return imageImpl{
		Image:  common.NewImage(size),
		screen: s,
	}, nil
}

func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	return &textureImpl{
		screen: s,
		size:   size,
	}, nil
}

func (s *screenImpl) NewWindow(opts screen.WindowGenerator) (screen.Window, error) {
	if opts.Width == 0 || opts.Height == 0 {
		return nil, fmt.Errorf("invalid width/height: %d/%d", opts.Width, opts.Height)
	}
	cvs := NewCanvas2d(opts.Width, opts.Height)

	w := &Window{
		cvs:    cvs,
		screen: s,
	}

	s.windows = append(s.windows, w)

	cvs.canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sendMouseEvent(args[0], mouse.DirNone)
		return nil
	}))
	cvs.canvas.Call("addEventListener", "mousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sendMouseEvent(args[0], mouse.DirPress)
		return nil
	}))
	cvs.canvas.Call("addEventListener", "mouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sendMouseEvent(args[0], mouse.DirRelease)
		return nil
	}))
	cvs.doc.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sendKeyEvent(args[0], key.DirPress)
		return nil
	}))
	cvs.doc.Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sendKeyEvent(args[0], key.DirRelease)
		return nil
	}))

	return w, nil
}
