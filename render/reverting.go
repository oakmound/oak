package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
)

type Reverting struct {
	root, current Modifiable
	mods          map[int]interface{}
}

func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.root = m
	rv.current = m
	rv.mods = make(map[int]interface{})
	return rv
}

func (rv *Reverting) Revert(mod int) {
	delete(rv.mods, mod)
	rv.current = rv.root.Copy()
	for mod, in := range rv.mods {
		switch mod {
		case F_FlipX:
			v := (in).(bool)
			if v {
				rv.current.FlipX()
			}
		case F_FlipY:
			v := (in).(bool)
			if v {
				rv.current.FlipY()
			}
		case F_ApplyColor:
			v := (in).(color.Color)
			rv.current.ApplyColor(v)
		case F_FillMask:
			v := (in).(image.RGBA)
			rv.current.FillMask(v)
		case F_ApplyMask:
			v := (in).(image.RGBA)
			rv.current.ApplyMask(v)
		case F_Rotate:
			v := (in).(int)
			rv.current.Rotate(v)
		case F_Scale:
			v := (in).([2]float64)
			rv.current.Scale(v[0], v[1])
		}
	}
}

func (rv *Reverting) Draw(buff screen.Buffer) {
	rv.current.Draw(buff)
}
func (rv *Reverting) GetRGBA() *image.RGBA {
	return rv.current.GetRGBA()
}
func (rv *Reverting) ShiftX(x float64) {
	rv.current.ShiftX(x)
}
func (rv *Reverting) ShiftY(y float64) {
	rv.current.ShiftY(y)
}
func (rv *Reverting) SetPos(x, y float64) {
	rv.current.SetPos(x, y)
}
func (rv *Reverting) GetLayer() int {
	return rv.current.GetLayer()
}
func (rv *Reverting) SetLayer(l int) {
	rv.current.SetLayer(l)
}
func (rv *Reverting) UnDraw() {
	rv.current.UnDraw()
}

func (rv *Reverting) FlipX() {
	rv.current.FlipX()
	if v, ok := rv.mods[F_FlipX]; ok {
		rv.mods[F_FlipX] = !(v).(bool)
	} else {
		rv.mods[F_FlipX] = true
	}
}
func (rv *Reverting) FlipY() {
	rv.current.FlipY()
	if v, ok := rv.mods[F_FlipY]; ok {
		rv.mods[F_FlipY] = !(v).(bool)
	} else {
		rv.mods[F_FlipY] = true
	}
}
func (rv *Reverting) ApplyColor(c color.Color) {
	if _, ok := rv.mods[F_ApplyColor]; ok {
		rv.Revert(F_ApplyColor)
	}
	rv.current.ApplyColor(c)
	rv.mods[F_ApplyColor] = c
}

func (rv *Reverting) FillMask(img image.RGBA) {
	if _, ok := rv.mods[F_FillMask]; ok {
		rv.Revert(F_FillMask)
	}
	rv.current.FillMask(img)
	rv.mods[F_FillMask] = img
}
func (rv *Reverting) ApplyMask(img image.RGBA) {
	if _, ok := rv.mods[F_ApplyMask]; ok {
		rv.Revert(F_ApplyMask)
	}
	rv.current.ApplyMask(img)
	rv.mods[F_ApplyMask] = img
}
func (rv *Reverting) Rotate(degrees int) {
	if _, ok := rv.mods[F_Rotate]; ok {
		rv.Revert(F_Rotate)
	}
	rv.current.Rotate(degrees)
	rv.mods[F_Rotate] = degrees
}
func (rv *Reverting) Scale(xRatio float64, yRatio float64) {
	if _, ok := rv.mods[F_Scale]; ok {
		rv.Revert(F_Scale)
	}
	rv.current.Scale(xRatio, yRatio)
	rv.mods[F_Scale] = [2]float64{xRatio, yRatio}
}
func (rv *Reverting) Copy() Modifiable {
	newRv := new(Reverting)
	newRv.root = rv.root.Copy()
	newRv.current = rv.current.Copy()
	newRv.mods = rv.mods
	return newRv
}

func (rv *Reverting) updateAnimation() {
	switch t := rv.current.(type) {
	case *Animation:
		t.updateAnimation()
	}
}
