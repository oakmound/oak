package jsdriver

import (
	"image"
	"syscall/js"
)

// Adapted from Mark Farnan's go-canvas library (github.com/markfarnan/go-canvas)
type Canvas2D struct {
	// DOM properties
	window js.Value
	doc    js.Value
	body   js.Value

	// Canvas properties
	canvas  js.Value
	ctx     js.Value
	imgData js.Value

	copybuff js.Value
}

func NewCanvas2d(width int, height int) *Canvas2D {
	var c Canvas2D
	c.window = js.Global()
	c.doc = c.window.Get("document")
	c.body = c.doc.Get("body")

	canvas := c.doc.Call("createElement", "canvas")

	canvas.Set("height", height)
	canvas.Set("width", width)
	// TODO: screen position
	c.body.Call("appendChild", canvas)

	c.canvas = canvas

	// Setup the 2D Drawing context
	c.ctx = c.canvas.Call("getContext", "2d", map[string]interface{}{"alpha": false})
	c.imgData = c.ctx.Call("createImageData", width, height) // Note Width, then Height
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c.copybuff = js.Global().Get("Uint8Array").New(len(img.Pix)) // Static JS buffer for copying data out to JS. Defined once and re-used to save on un-needed allocations

	return &c
}
