package render

import (
	"image/draw"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/oakmound/oak/v2/oakerr"

	"github.com/oakmound/oak/v2/dlog"
)

var (
	// GlobalDrawStack is the stack that all draw calls are parsed through.
	GlobalDrawStack  = NewDrawStack(NewDynamicHeap())
	initialDrawStack = GlobalDrawStack
)

//The DrawStack is a stack with a safe adding mechanism that creates isolation between draw steps via predraw
type DrawStack struct {
	as     []Stackable
	toPush []Stackable
	toPop  int
}

// An Stackable manages Renderables
type Stackable interface {
	PreDraw()
	Add(Renderable, ...int) Renderable
	Replace(Renderable, Renderable, int)
	Copy() Stackable
	DrawToScreen(draw.Image, intgeom.Point2, int, int)
}

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
	initialDrawStack = GlobalDrawStack.Copy()
}

//ResetDrawStack resets the Global stack back to the initial stack
func ResetDrawStack() {
	GlobalDrawStack = initialDrawStack.Copy()
}

// DrawToScreen on a stack will render its contents to the input buffer, for a screen
// of w,h dimensions, from a view point of view.
func (ds *DrawStack) DrawToScreen(world draw.Image, view intgeom.Point2, w, h int) {
	for _, a := range ds.as {
		// If we had concurrent operations, we'd do it here
		// in that case each draw call would return to us something
		// to composite onto the window / world
		// TODO v3: add 'DrawConcurrency'? Bake in the background drawing? Benchmark if done
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

func (d *DrawStack) Draw(r Renderable, layers ...int) (Renderable, error) {
	if r == nil {
		return nil, oakerr.NilInput{InputName: "r"}
	}
	if len(d.as) == 1 {
		return d.as[0].Add(r, layers...), nil
	}
	if len(layers) > 0 {
		stackLayer := layers[0]
		if stackLayer < 0 || stackLayer >= len(d.as) {
			dlog.Error("Layer", stackLayer, "does not exist on global draw stack")
			return nil, oakerr.InvalidInput{InputName: "layers"}
		}
		return d.as[stackLayer].Add(r, layers[1:]...), nil
	}
	return d.as[0].Add(r), nil
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
