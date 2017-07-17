package render

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/physics"
)

// The Compound type is intended for use to easily swap between multiple
// renderables that are drawn at the same position on the same layer.
// A common use case for this would be a character entitiy who switches
// their animation based on how they are moving or what they are doing.
//
// The Compound type removes the need to repeatedly draw and undraw elements
// of a character, which has a tendency to leave nothing drawn for a draw frame.
type Compound struct {
	LayeredPoint
	subRenderables map[string]Modifiable
	curRenderable  string
}

//NewCompound creates a new compound from a map of names to modifiables
func NewCompound(start string, m map[string]Modifiable) *Compound {
	return &Compound{
		LayeredPoint:   NewLayeredPoint(0, 0, 0),
		subRenderables: m,
		curRenderable:  start,
	}
}

//Add makes a new entry in the Compounds map
func (c *Compound) Add(k string, v Modifiable) {
	c.subRenderables[k] = v
}

//Set sets the current renderable to the one specified
func (c *Compound) Set(k string) {
	if _, ok := c.subRenderables[k]; !ok {
		panic("Unknown renderable for string " + k + " on compound")
	}
	c.curRenderable = k
}

//GetSub returns a given subrenderable from the map
func (c *Compound) GetSub(s string) Modifiable {
	return c.subRenderables[s]
}

//Get returns the Compounds current Renderable
func (c *Compound) Get() string {
	return c.curRenderable
}

//IsInterruptable returns whether the current renderable is interruptable
func (c *Compound) IsInterruptable() bool {
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		return t.Interruptable
	case *Sequence:
		return t.Interruptable
	case *Reverting:
		return t.IsInterruptable()
	case *Compound:
		return t.IsInterruptable()
	}
	return true
}

//IsStatic returns whether the current renderable is static
func (c *Compound) IsStatic() bool {
	switch c.subRenderables[c.curRenderable].(type) {
	case *Animation, *Sequence:
		return false
	case *Reverting:
		return c.subRenderables[c.curRenderable].(*Reverting).IsStatic()
	case *Compound:
		return c.subRenderables[c.curRenderable].(*Compound).IsStatic()
	}
	return true
}

//SetOffsets sets the logical offset for the specified subrenderable
func (c *Compound) SetOffsets(k string, offsets physics.Vector) {
	if r, ok := c.subRenderables[k]; ok {
		r.SetPos(offsets.X(), offsets.Y())
	}
}

//Copy creates a copy of the Compound
func (c *Compound) Copy() Modifiable {
	newC := new(Compound)
	newC.LayeredPoint = c.LayeredPoint.Copy()
	newSubRenderables := make(map[string]Modifiable)
	for k, v := range c.subRenderables {
		newSubRenderables[k] = v.Copy()
	}
	newC.subRenderables = newSubRenderables
	newC.curRenderable = c.curRenderable
	return newC
}

//GetRGBA returns the current renderables rgba
func (c *Compound) GetRGBA() *image.RGBA {
	return c.subRenderables[c.curRenderable].GetRGBA()
}

//Modify performs a series of modifications on the Compound
func (c *Compound) Modify(ms ...Modification) Modifiable {
	for _, r := range c.subRenderables {
		r.Modify(ms...)
	}
	return c
}

//DrawOffset draws the Compound at an offset from its logical location
func (c *Compound) DrawOffset(buff draw.Image, xOff float64, yOff float64) {
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X()+xOff, c.Y()+yOff)
}

//Draw draws the Compound at its logical location
func (c *Compound) Draw(buff draw.Image) {
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X(), c.Y())
}

//ShiftPos shifts the Compounds logical position
func (c *Compound) ShiftPos(x, y float64) {
	c.SetPos(c.X()+x, c.Y()+y)
}

//ShiftY shifts the Compounds logical y position
func (c *Compound) ShiftY(y float64) {
	c.SetPos(c.X(), c.Y()+y)
}

//ShiftX shifts the Compounds logical x position
func (c *Compound) ShiftX(x float64) {
	c.SetPos(c.X()+x, c.Y())
}

//SetPos sets the Compound's logical position
func (c *Compound) SetPos(x, y float64) {
	c.LayeredPoint.SetPos(x, y)
}

//GetDims gets the current Renderables dimensions
func (c *Compound) GetDims() (int, int) {
	return c.subRenderables[c.curRenderable].GetDims()
}

//Pause stops the current Renderable if possible
func (c *Compound) Pause() {
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		t.Pause()
	case *Sequence:
		t.Pause()
	case *Reverting:
		t.Pause()
	}
}

// Unpause tries to unpause the current Renderable if possible
func (c *Compound) Unpause() {
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	case *Reverting:
		t.Unpause()
	}
}

//Revert tries to revert the current Renderable if possible
func (c *Compound) Revert(mod int) {
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.Revert(mod)
		}
	}
}

//RevertAll tries to revert the all sub-Renderables if possible
func (c *Compound) RevertAll() {
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.RevertAll()
		}
	}
}
