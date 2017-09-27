package render

import (
	"image"
	"image/draw"
	"sync"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/oakerr"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"
)

// The Switch type is intended for use to easily swap between multiple
// renderables that are drawn at the same position on the same layer.
// A common use case for this would be a character entitiy who switches
// their animation based on how they are moving or what they are doing.
//
// The Switch type removes the need to repeatedly draw and undraw elements
// of a character, which has a tendency to leave nothing drawn for a draw frame
// as the switch happens.
type Switch struct {
	LayeredPoint
	subRenderables map[string]Modifiable
	curRenderable  string
	lock           sync.RWMutex
}

// NewSwitch creates a new Switch from a map of names to modifiables
func NewSwitch(start string, m map[string]Modifiable) *Switch {
	return &Switch{
		LayeredPoint:   NewLayeredPoint(0, 0, 0),
		subRenderables: m,
		curRenderable:  start,
		lock:           sync.RWMutex{},
	}
}

// Add makes a new entry in the Switch's map
func (c *Switch) Add(k string, v Modifiable) (err error) {
	if _, ok := c.subRenderables[k]; ok {
		err = oakerr.ExistingElement{
			InputName:   "k",
			InputType:   "string",
			Overwritten: true,
		}
	}
	c.lock.Lock()
	c.subRenderables[k] = v
	c.lock.Unlock()
	return err
}

// Set sets the current renderable to the one specified
func (c *Switch) Set(k string) error {
	c.lock.RLock()
	if _, ok := c.subRenderables[k]; !ok {
		return oakerr.InvalidInput{InputName: "k"}
	}
	c.lock.RUnlock()
	c.curRenderable = k
	return nil
}

// GetSub returns a keyed Modifiable from this Switch's map
func (c *Switch) GetSub(s string) Modifiable {
	c.lock.RLock()
	m := c.subRenderables[s]
	c.lock.RUnlock()
	return m
}

// Get returns the Switch's current key
func (c *Switch) Get() string {
	return c.curRenderable
}

// SetOffsets sets the logical offset for the specified key
func (c *Switch) SetOffsets(k string, offsets physics.Vector) {
	c.lock.RLock()
	if r, ok := c.subRenderables[k]; ok {
		r.SetPos(offsets.X(), offsets.Y())
	}
	c.lock.RUnlock()
}

// Copy creates a copy of the Switch
func (c *Switch) Copy() Modifiable {
	newC := new(Switch)
	newC.LayeredPoint = c.LayeredPoint.Copy()
	newSubRenderables := make(map[string]Modifiable)
	c.lock.RLock()
	for k, v := range c.subRenderables {
		newSubRenderables[k] = v.Copy()
	}
	c.lock.RUnlock()
	newC.subRenderables = newSubRenderables
	newC.curRenderable = c.curRenderable
	newC.lock = sync.RWMutex{}
	return newC
}

//GetRGBA returns the current renderables rgba
func (c *Switch) GetRGBA() *image.RGBA {
	c.lock.RLock()
	rgba := c.subRenderables[c.curRenderable].GetRGBA()
	c.lock.RUnlock()
	return rgba
}

// Modify performs the input modifications on all elements of the Switch
func (c *Switch) Modify(ms ...mod.Mod) Modifiable {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Modify(ms...)
	}
	c.lock.RUnlock()
	return c
}

// Filter filters all elements of the Switch with fs
func (c *Switch) Filter(fs ...mod.Filter) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		r.Filter(fs...)
	}
	c.lock.RUnlock()
}

//DrawOffset draws the Switch at an offset from its logical location
func (c *Switch) DrawOffset(buff draw.Image, xOff float64, yOff float64) {
	c.lock.RLock()
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X()+xOff, c.Y()+yOff)
	c.lock.RUnlock()
}

//Draw draws the Switch at its logical location
func (c *Switch) Draw(buff draw.Image) {
	c.lock.RLock()
	c.subRenderables[c.curRenderable].DrawOffset(buff, c.X(), c.Y())
	c.lock.RUnlock()
}

// ShiftPos shifts the Switch's logical position
func (c *Switch) ShiftPos(x, y float64) {
	c.SetPos(c.X()+x, c.Y()+y)
}

// GetDims gets the current Renderables dimensions
func (c *Switch) GetDims() (int, int) {
	c.lock.RLock()
	w, h := c.subRenderables[c.curRenderable].GetDims()
	c.lock.RUnlock()
	return w, h
}

// Pause stops the current Renderable if possible
func (c *Switch) Pause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(CanPause); ok {
		cp.Pause()
	}
	c.lock.RUnlock()
}

// Unpause tries to unpause the current Renderable if possible
func (c *Switch) Unpause() {
	c.lock.RLock()
	if cp, ok := c.subRenderables[c.curRenderable].(CanPause); ok {
		cp.Unpause()
	}
	c.lock.RUnlock()
}

// IsInterruptable returns whether the current renderable is interruptable
func (c *Switch) IsInterruptable() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if i, ok := c.subRenderables[c.curRenderable].(NonInterruptable); ok {
		return i.IsInterruptable()
	}
	return true
}

// IsStatic returns whether the current renderable is static
func (c *Switch) IsStatic() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if s, ok := c.subRenderables[c.curRenderable].(NonStatic); ok {
		return s.IsStatic()
	}
	return true
}

// SetTriggerID sets the ID AnimationEnd will trigger on for animating subtypes.
// Todo: standardize this with the other interface Set functions so that it
// also only acts on the current subRenderable, or the other way around, or
// somehow offer both options
func (c *Switch) SetTriggerID(cid event.CID) {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		if t, ok := r.(Triggerable); ok {
			t.SetTriggerID(cid)
		}
	}
	c.lock.RUnlock()
}

// Revert will revert all parts of this Switch that can be reverted
func (c *Switch) Revert(mod int) {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.Revert(mod)
		}
	}
	c.lock.RUnlock()
}

// RevertAll will revert all parts of this Switch that can be reverted, back
// to their original state.
func (c *Switch) RevertAll() {
	c.lock.RLock()
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.RevertAll()
		}
	}
	c.lock.RUnlock()
}

func (c *Switch) update() {
	c.lock.RLock()
	for _, r := range c.subRenderables {
		if u, ok := r.(updates); ok {
			u.update()
		}
	}
	c.lock.RUnlock()
}
