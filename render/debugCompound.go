package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"
)

// The DebugCompound type is intended for use to easily swap between multiple
// renderables that are drawn at the same position on the same layer.
// A common use case for this would be a character entitiy who switches
// their animation based on how they are moving or what they are doing.
//
// The DebugCompound type removes the need to repeatedly draw and undraw elements
// of a character, which has a tendency to leave nothing drawn for a draw frame.
type DebugCompound struct {
	LayeredPoint
	subRenderables map[string]Modifiable
	offsets        map[string]Point
	curRenderable  string
}

func NewDebugCompound(start string, m map[string]Modifiable) *DebugCompound {
	for k, v := range m {
		fmt.Println("Adding key", k, " and value", v, " to debug compound")
	}
	return &DebugCompound{
		subRenderables: m,
		curRenderable:  start,
		offsets:        make(map[string]Point),
	}
}

func (c *DebugCompound) Add(k string, v Modifiable) {
	c.subRenderables[k] = v
}

func (c *DebugCompound) Set(k string) {
	if _, ok := c.subRenderables[k]; !ok {
		panic("Unknown renderable for string " + k + " on compound")
	}
	c.curRenderable = k
}
func (c *DebugCompound) GetSub(s string) Modifiable {
	return c.subRenderables[s]
}
func (c *DebugCompound) Get() string {
	return c.curRenderable
}

func (c *DebugCompound) IsStatic() bool {
	switch c.subRenderables[c.curRenderable].(type) {
	case *Animation, *Sequence:
		return false
	case *Reverting:
		return c.subRenderables[c.curRenderable].(*Reverting).IsStatic()
	case *DebugCompound:
		return c.subRenderables[c.curRenderable].(*DebugCompound).IsStatic()
	}
	return true
}

func (c *DebugCompound) SetOffsets(k string, offsets Point) {
	c.offsets[k] = offsets
}

func (c *DebugCompound) Copy() Modifiable {
	newC := new(DebugCompound)
	*newC = *c
	newSubRenderables := make(map[string]Modifiable)
	for k, v := range c.subRenderables {
		newSubRenderables[k] = v.Copy()
	}
	newC.subRenderables = newSubRenderables
	return newC
}

func (c *DebugCompound) GetRGBA() *image.RGBA {
	return c.subRenderables[c.curRenderable].GetRGBA()
}

func (c *DebugCompound) ApplyColor(co color.Color) Modifiable {
	for _, rend := range c.subRenderables {
		rend.ApplyColor(co)
	}
	return c
}

func (c *DebugCompound) FillMask(img image.RGBA) Modifiable {
	for _, rend := range c.subRenderables {
		rend.FillMask(img)
	}
	return c
}

func (c *DebugCompound) ApplyMask(img image.RGBA) Modifiable {
	for _, rend := range c.subRenderables {
		rend.ApplyMask(img)
	}
	return c
}

func (c *DebugCompound) Rotate(degrees int) Modifiable {
	for _, rend := range c.subRenderables {
		rend.Rotate(degrees)
	}
	return c
}

func (c *DebugCompound) Scale(xRatio float64, yRatio float64) Modifiable {
	for _, rend := range c.subRenderables {
		rend.Scale(xRatio, yRatio)
	}
	return c
}

func (c *DebugCompound) FlipX() Modifiable {
	for _, rend := range c.subRenderables {
		rend.FlipX()
	}
	return c
}

func (c *DebugCompound) FlipY() Modifiable {
	for _, rend := range c.subRenderables {
		rend.FlipY()
	}
	return c
}
func (c *DebugCompound) Fade(alpha int) Modifiable {
	for _, rend := range c.subRenderables {
		rend.Fade(alpha)
	}
	return c
}

func (c *DebugCompound) Draw(buff draw.Image) {
	fmt.Println("Drawing DebugCompund:", c.subRenderables, c.curRenderable)
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Composite:
		t.Draw(buff)
		return
	case *Reverting:
		t.updateAnimation()
	case *Animation:
		t.updateAnimation()
	case *Sequence:
		t.update()
	default:
		fmt.Println("Type info of unknown element:", reflect.TypeOf(c.subRenderables[c.curRenderable]))
	}
	img := c.GetRGBA()
	if img == nil {
		fmt.Println("Type", reflect.TypeOf(c.subRenderables[c.curRenderable]), "Slipped through type switch")
	}
	drawX := int(c.X)
	drawY := int(c.Y)
	if offsets, ok := c.offsets[c.curRenderable]; ok {
		drawX += int(offsets.X)
		drawY += int(offsets.Y)
	}
	ShinyDraw(buff, img, drawX, drawY)
}

func (c *DebugCompound) SetPos(x, y float64) {
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Composite:
			t.SetPos(x, y)
		default:
			fmt.Println("Type info of unknown element:", reflect.TypeOf(c.subRenderables[c.curRenderable]))
		}
	}
	c.LayeredPoint.SetPos(x, y)
}

func (c *DebugCompound) Pause() {
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		t.Pause()
	case *Sequence:
		t.Pause()
	case *Reverting:
		t.Pause()
	}
}

func (c *DebugCompound) Unpause() {
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	case *Reverting:
		t.Unpause()
	}
}

func (c *DebugCompound) Revert(mod int) {
	for _, v := range c.subRenderables {
		switch t := v.(type) {
		case *Reverting:
			t.Revert(mod)
		}
	}
}
