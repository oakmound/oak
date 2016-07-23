package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
)

type Compound struct {
	x, y           float64
	subRenderables map[string]Modifiable
	curRenderable  string
	layer          int
}

func NewCompound(start string, m map[string]Modifiable) *Compound {
	return &Compound{
		subRenderables: m,
		curRenderable:  start,
	}
}

func (c *Compound) Add(k string, v Modifiable) {
	c.subRenderables[k] = v
}

func (c *Compound) Set(k string) {
	c.curRenderable = k
}

func (c *Compound) Copy() *Compound {
	newC := *c
	return &newC
}

func (c *Compound) GetRGBA() *image.RGBA {
	return c.subRenderables[c.curRenderable].GetRGBA()
}

func (c *Compound) ApplyColor(co color.Color) {
	for _, rend := range c.subRenderables {
		rend.ApplyColor(co)
	}
}

func (c *Compound) FillMask(img image.RGBA) {
	for _, rend := range c.subRenderables {
		rend.FillMask(img)
	}
}

func (c *Compound) ApplyMask(img image.RGBA) {
	for _, rend := range c.subRenderables {
		rend.ApplyMask(img)
	}
}

func (c *Compound) Rotate(degrees int) {
	for _, rend := range c.subRenderables {
		rend.Rotate(degrees)
	}
}

func (c *Compound) Scale(xRatio float64, yRatio float64) {
	for _, rend := range c.subRenderables {
		rend.Scale(xRatio, yRatio)
	}
}

func (c *Compound) FlipX() {
	for _, rend := range c.subRenderables {
		rend.FlipX()
	}
}

func (c *Compound) FlipY() {
	for _, rend := range c.subRenderables {
		rend.FlipY()
	}
}

func (c *Compound) SetPos(x, y float64) {
	c.x = x
	c.y = y
}

func (c *Compound) ShiftX(x float64) {
	c.x += x
}
func (c *Compound) ShiftY(y float64) {
	c.y += y
}

func (c *Compound) Draw(buff screen.Buffer) {
	img := c.GetRGBA()
	switch t := c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		t.updateAnimation()
	}
	ShinyDraw(buff, img, int(c.x), int(c.y))
}

func (c *Compound) GetLayer() int {
	return c.layer
}

func (c *Compound) SetLayer(l int) {
	c.layer = l
}

func (c *Compound) UnDraw() {
	c.layer = -1
}

func (c *Compound) Pause() {
	switch c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		c.subRenderables[c.curRenderable].(*Animation).playing = false
	}
}

func (c *Compound) Unpause() {
	switch c.subRenderables[c.curRenderable].(type) {
	case *Animation:
		c.subRenderables[c.curRenderable].(*Animation).playing = true
	}
}
