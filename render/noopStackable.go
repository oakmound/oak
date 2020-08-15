package render

import (
	"image"
	"image/draw"
)

// NoopStackable is a Stackable element where all methods are no-ops.
// Use for tests to disable rendering.
type NoopStackable struct{}

// PreDraw on a NoopStackable does nothing.
func (ns NoopStackable) PreDraw() {}

// Add on a NoopStackable does nothing. The input Renderable is still returned.
func (ns NoopStackable) Add(r Renderable, _ ...int) Renderable {
	return r
}

// Replace on a NoopStackable does nothing.
func (ns NoopStackable) Replace(Renderable, Renderable, int) {}

// Copy on a NoopStackable returns itself.
func (ns NoopStackable) Copy() Stackable {
	return ns
}

// Todo (3.0): export draw method in interface
func (ns NoopStackable) draw(draw.Image, image.Point, int, int) {}
