package render

import (
	"errors"
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/dlog"
)

var (
	GlobalDrawStack = &DrawStack{
		as: []Addable{NewHeap(false)},
	}
	initialDrawStack = GlobalDrawStack
)

type DrawStack struct {
	as     []Addable
	toPush []Addable
	toPop  int
}

type Addable interface {
	PreDraw()
	Add(Renderable, int) Renderable
	Replace(Renderable, Renderable, int)
	Copy() Addable
	draw(draw.Image, image.Point, int, int)
}

func SetDrawStack(as ...Addable) {
	GlobalDrawStack = &DrawStack{as: as}
	initialDrawStack = GlobalDrawStack.Copy()
	dlog.Info("Global", GlobalDrawStack)
	dlog.Info("Initial", initialDrawStack)
}

func ResetDrawStack() {
	GlobalDrawStack = initialDrawStack.Copy()
}

func (ds *DrawStack) Draw(world draw.Image, view image.Point, w, h int) {

	for _, a := range ds.as {
		// If we had concurrent operations, we'd do it here
		// in that case each draw call would return to us something
		// to composite onto the window / world
		a.draw(world, view, w, h)
	}
}

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

func (ds *DrawStack) Push(a Addable) {
	ds.toPush = append(ds.toPush, a)

}

func (ds *DrawStack) Pop() {
	ds.toPop++
}

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

func PreDraw() {
	if resetDraw {
		ResetDrawStack()
		resetDraw = false
	} else {
		GlobalDrawStack.PreDraw()
	}
}
