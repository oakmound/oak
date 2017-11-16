package render

import (
	"errors"
	"image"
	"image/draw"

	"github.com/oakmound/oak/dlog"
)

var (
	// GlobalDrawStack is the stack that all draw calls are parsed through.
	GlobalDrawStack = &DrawStack{
		as: []Addable{NewHeap(false)},
	}
	initialDrawStack = GlobalDrawStack
)

//The DrawStack is a stack with a safe adding mechanism that creates isolation between draw steps via predraw
type DrawStack struct {
	as     []Addable
	toPush []Addable
	toPop  int
}

// An Addable manages Renderables
type Addable interface {
	PreDraw()
	Add(Renderable, int) Renderable
	Replace(Renderable, Renderable, int)
	Copy() Addable
	draw(draw.Image, image.Point, int, int)
}

// SetDrawStack takes in a set of Addables which act as the set of Drawstacks available
// and resets how calls to Draw will act. If this is called mid scene,
// all elements on the existing draw stack will be lost.
func SetDrawStack(as ...Addable) {
	GlobalDrawStack = &DrawStack{as: as}
	initialDrawStack = GlobalDrawStack.Copy()
	dlog.Info("Global draw stack", GlobalDrawStack)
	dlog.Info("Initial draw stack", initialDrawStack)
}

//ResetDrawStack resets the Global stack back to the initial stack
func ResetDrawStack() {
	GlobalDrawStack = initialDrawStack.Copy()
}

//Draw actively draws the onto the actual screen
func (ds *DrawStack) Draw(world draw.Image, view image.Point, w, h int) {
	for _, a := range ds.as {
		// If we had concurrent operations, we'd do it here
		// in that case each draw call would return to us something
		// to composite onto the window / world
		a.draw(world, view, w, h)
	}
}

//Draw accesses the global draw stack
func Draw(r Renderable, l int) (Renderable, error) {
	if r == nil {
		dlog.Error("Tried to draw nil")
		return nil, errors.New("Tried to draw nil")
	}
	// If there's only one element, l refers to the layer
	// within that element.
	if len(GlobalDrawStack.as) == 1 {
		return GlobalDrawStack.as[0].Add(r, l), nil

		// Otherwise, l refers to the index within the DrawStack.
	}
	if l < 0 || l >= len(GlobalDrawStack.as) {
		dlog.Error("Layer", l, "does not exist on global draw stack")
		return nil, errors.New("Layer does not exist on stack")
	}
	return GlobalDrawStack.as[l].Add(r, r.GetLayer()), nil
}

// ReplaceDraw will undraw r1 and draw r2 after the next draw frame
// Useful for not working
func ReplaceDraw(r1, r2 Renderable, stackLayer, layer int) {
	if r1 == nil || r2 == nil {
		dlog.Error("Tried to draw nil")
		return
	}
	if stackLayer < 0 || stackLayer >= len(GlobalDrawStack.as) {
		dlog.Error("Layer", stackLayer, "does not exist on global draw stack")
		return
	}
	r2.SetLayer(layer)
	GlobalDrawStack.as[stackLayer].Replace(r1, r2, layer)
}

// Push appends an addable to the draw stack during the next PreDraw.
func (ds *DrawStack) Push(a Addable) {
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
		ds.toPush = []Addable{}
	}
	for _, a := range ds.as {
		a.PreDraw()
	}
}

// Copy creates a new deep copy of a Drawstack
func (ds *DrawStack) Copy() *DrawStack {
	ds2 := new(DrawStack)
	ds2.as = make([]Addable, len(ds.as))
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
