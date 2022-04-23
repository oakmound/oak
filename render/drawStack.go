package render

import (
	"image/draw"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/oakerr"
)

var (
	// GlobalDrawStack is the stack that all draw calls are sent through.
	GlobalDrawStack = NewDrawStack(NewDynamicHeap())
)

//The DrawStack is a stack with a safe adding mechanism that creates isolation between draw steps via predraw
type DrawStack struct {
	as     []Stackable
	toPush []Stackable
	toPop  int
}

// A Stackable can be put onto a draw stack. It usually manages how a subset of renderables
// are drawn.
type Stackable interface {
	PreDraw()
	Add(Renderable, ...int) Renderable
	Replace(Renderable, Renderable, int)
	Copy() Stackable
	DrawToScreen(draw.Image, *intgeom.Point2, int, int)
	Clear()
}

// NewDrawStack creates a DrawStack with the given stackable items, drawn in descending index order.
func NewDrawStack(stack ...Stackable) *DrawStack {
	return &DrawStack{
		as: stack,
	}
}

// SetDrawStack takes in a set of Stackables which act as the Drawstack available
// and resets how calls to Draw will act. If this is called mid scene,
// all elements on the existing draw stack will be lost.
func SetDrawStack(stackLayers ...Stackable) {
	GlobalDrawStack = NewDrawStack(stackLayers...)
}

// Clear clears all stackables in a draw stack. This should revert the stack to contain
// no renderable components.
func (ds *DrawStack) Clear() {
	for _, stackable := range ds.as {
		stackable.Clear()
	}
}

// DrawToScreen on a stack will render its contents to the input buffer, for a screen
// of w,h dimensions, from a view point of view.
func (ds *DrawStack) DrawToScreen(world draw.Image, view *intgeom.Point2, w, h int) {
	for _, a := range ds.as {
		// If we had concurrent operations, we'd do it here
		// in that case each draw call would return to us something
		// to composite onto the window / world
		a.DrawToScreen(world, view, w, h)
	}
}

// Draw adds the given renderable to the global draw stack.
//
// If the draw stack has only one stackable, the item will be added to that
// stackable with the input layers as its argument. Otherwise, the item will be added
// to the layers[0]th stackable, with remaining layers supplied to the stackable
// as arguments.
//
// If zero layers are provided, it will add to the zeroth stack layer and
// give nothing to the stackable's argument.
func Draw(r Renderable, layers ...int) (Renderable, error) {
	return GlobalDrawStack.Draw(r, layers...)
}

// Draw adds the given renderable to the draw stack at the appropriate position based
// on the input layers. See render.Draw.
func (ds *DrawStack) Draw(r Renderable, layers ...int) (Renderable, error) {
	if r == nil {
		return nil, oakerr.NilInput{InputName: "r"}
	}
	if len(ds.as) == 1 {
		return ds.as[0].Add(r, layers...), nil
	}
	if len(layers) > 0 {
		stackLayer := layers[0]
		if stackLayer < 0 || stackLayer >= len(ds.as) {
			return nil, oakerr.InvalidInput{InputName: "layers"}
		}
		return ds.as[stackLayer].Add(r, layers[1:]...), nil
	}
	return ds.as[0].Add(r), nil
}

// Push appends a Stackable to the draw stack during the next PreDraw.
func (ds *DrawStack) Push(a Stackable) {
	ds.toPush = append(ds.toPush, a)

}

// Pop pops an element from the stack at the next PreDraw call.
func (ds *DrawStack) Pop() {
	ds.toPop++
}

// PreDraw performs whatever processes need to occur before this can be
// drawn. In the case of the stack, it enacts previous Push and Pop calls,
// and signals to elements on the stack to also prepare to be drawn.
func (ds *DrawStack) PreDraw() {
	if ds.toPop > 0 {
		ds.as = ds.as[0 : len(ds.as)-ds.toPop]
		ds.toPop = 0
	}
	if len(ds.toPush) > 0 {
		ds.as = append(ds.as, ds.toPush...)
		// Should use two toPush lists, for this and
		// draw heaps, so this call won't ever drop anything
		ds.toPush = []Stackable{}
	}
	for _, a := range ds.as {
		a.PreDraw()
	}
}

// Copy creates a new deep copy of a Drawstack
func (ds *DrawStack) Copy() *DrawStack {
	ds2 := new(DrawStack)
	ds2.as = make([]Stackable, len(ds.as))
	for i, a := range ds.as {
		ds2.as[i] = a.Copy()
	}
	ds2.toPop = ds.toPop
	ds2.toPush = ds.toPush
	return ds2
}
