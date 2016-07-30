package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
)

var (
	emptyMods = [7]interface{}{
		false,
		false,
		color.RGBA{0, 0, 0, 0},
		image.RGBA{Stride: 0},
		image.RGBA{Stride: 0},
		0,
		[2]float64{0, 0},
	}
)

type Reverting struct {
	root, current Modifiable
	mods          [7]interface{}
}

func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.root = m
	rv.current = m.Copy()
	rv.mods = emptyMods
	return rv
}

func (rv *Reverting) Revert(mod int) {
	rv.mods[mod] = emptyMods[mod]
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
			if v != (color.RGBA{0, 0, 0, 0}) {
				rv.current.ApplyColor(v)
			}
		case F_FillMask:
			v := (in).(image.RGBA)
			if v.Stride != 0 {
				rv.current.FillMask(v)
			}
		case F_ApplyMask:
			v := (in).(image.RGBA)
			if v.Stride != 0 {
				rv.current.ApplyMask(v)
			}
		case F_Rotate:
			v := (in).(int)
			if v != 0 {
				rv.current.Rotate(v)
			}
		case F_Scale:
			v := (in).([2]float64)
			if v[0] != 0 && v[1] != 0 {
				rv.current.Scale(v[0], v[1])
			}
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
	rv.mods[F_FlipX] = !(rv.mods[F_FlipX]).(bool)
}
func (rv *Reverting) FlipY() {
	rv.current.FlipY()
	rv.mods[F_FlipY] = !(rv.mods[F_FlipY]).(bool)
}
func (rv *Reverting) ApplyColor(c color.Color) {
	rv.Revert(F_ApplyColor)
	rv.current.ApplyColor(c)
	rv.mods[F_ApplyColor] = c
}

func (rv *Reverting) FillMask(img image.RGBA) {
	rv.Revert(F_FillMask)
	rv.current.FillMask(img)
	rv.mods[F_FillMask] = img
}
func (rv *Reverting) ApplyMask(img image.RGBA) {
	rv.Revert(F_ApplyMask)
	rv.current.ApplyMask(img)
	rv.mods[F_ApplyMask] = img
}
func (rv *Reverting) Rotate(degrees int) {
	rv.Revert(F_Rotate)
	rv.current.Rotate(degrees)
	rv.mods[F_Rotate] = degrees
}
func (rv *Reverting) Scale(xRatio float64, yRatio float64) {
	rv.Revert(F_Scale)
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
	switch t := rv.root.(type) {
	case *Animation:
		t.updateAnimation()
	}
}

func (rv *Reverting) Set(k string) {
	switch t := rv.current.(type) {
	case *Compound:
		t.Set(k)
	}
	switch t := rv.root.(type) {
	case *Compound:
		t.Set(k)
	}
}

func (rv *Reverting) Pause() {
	switch t := rv.current.(type) {
	case *Animation:
		t.playing = false
	case *Compound:
		t.Pause()
	}
	switch t := rv.root.(type) {
	case *Animation:
		t.playing = false
	case *Compound:
		t.Pause()
	}

}

func (rv *Reverting) Unpause() {
	switch t := rv.current.(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	}
	switch t := rv.root.(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	}
}
