package render

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/oakerr"

	"github.com/oakmound/oak/dlog"
)

var (
	// GlobalDrawStack is the stack that all draw calls are parsed through.
	GlobalDrawStack = &DrawStack{
		as: []Stackable{NewHeap(false)},
	}
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
	draw(draw.Image, image.Point, int, int)
}

// SetDrawStack takes in a set of Addables which act as the set of Drawstacks available
// and resets how calls to Draw will act. If this is called mid scene,
// all elements on the existing draw stack will be lost.
func SetDrawStack(as ...Stackable) {
	GlobalDrawStack = &DrawStack{as: as}
	initialDrawStack = GlobalDrawStack.Copy()
	dlog.Info("Global draw stack", GlobalDrawStack)
	dlog.Info("Initial draw stack", initialDrawStack)
}

//ResetDrawStack resets the Global stack back to the initial stack
func ResetDrawStack() {
	GlobalDrawStack = initialDrawStack.Copy()
}

// Draw on a stack will render its contents to the input buffer, for a screen
// of w,h dimensions, from a view point of view.
func (ds *DrawStack) Draw(world draw.Image, view image.Point, w, h int) {
	for _, a := range ds.as {
		// If we had concurrent operations, we'd do it here
		// in that case each draw call would return to us something
		// to composite onto the window / world
		a.draw(world, view, w, h)
	}
}

// Draw adds the given renderable to the global draw stack.
//
// If the draw stack has only one stackable, the item will be added to that
// stackable with the input layers as its argument. Otherwise, the item will be added
// to the l[0]th stackable, with remaining layers supplied to the stackable
// as arguments.
//
// If zero layers are provided, it will add to the zeroth stack layer and
// give nothing to the stackable's argument.
//
//
func Draw(r Renderable, layers ...int) (Renderable, error) {
	if r == nil {
		dlog.Error("Tried to draw nil")
		return nil, oakerr.NilInput{InputName: "r"}
	}
	if len(GlobalDrawStack.as) == 1 {
		return GlobalDrawStack.as[0].Add(r, layers...), nil
	}
	if len(layers) > 0 {
		stackLayer := layers[0]
		if stackLayer < 0 || stackLayer >= len(GlobalDrawStack.as) {
			dlog.Error("Layer", stackLayer, "does not exist on global draw stack")
			return nil, oakerr.InvalidInput{InputName: "layers"}
		}
		return GlobalDrawStack.as[stackLayer].Add(r, layers[1:]...), nil
	}
	return GlobalDrawStack.as[0].Add(r), nil
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

// PreDraw tries to reset the GlobalDrawStack or performs the GlobalDrawStack's predraw functions
func PreDraw() {
	if resetDraw {
		ResetDrawStack()
		resetDraw = false
	} else {
		GlobalDrawStack.PreDraw()
	}
}
