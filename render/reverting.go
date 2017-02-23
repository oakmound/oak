package render

import (
	"image"
	"image/color"
	"image/draw"
)

var (
	emptyMods = [8]interface{}{
		false,
		false,
		color.RGBA{0, 0, 0, 0},
		image.RGBA{Stride: 0},
		image.RGBA{Stride: 0},
		0,
		[2]float64{0, 0},
		0,
	}
)

type Reverting struct {
	root, current Modifiable
	mods          [8]interface{}
}

func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.root = m
	rv.current = m.Copy()
	rv.mods = emptyMods
	return rv
}

func (rv *Reverting) IsInterruptable() bool {
	switch t := rv.root.(type) {
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

func (rv *Reverting) IsStatic() bool {
	switch rv.root.(type) {
	case *Animation, *Sequence:
		return false
	case *Reverting:
		return rv.root.(*Reverting).IsStatic()
	case *Compound:
		return rv.root.(*Compound).IsStatic()
	}
	return true
}

func (rv *Reverting) Revert(mod int) {
	x := rv.current.GetX()
	y := rv.current.GetY()
	rv.mods[mod] = emptyMods[mod]
	rv.current = rv.root.Copy()
	rv.SetPos(x, y)
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
		case F_Fade:
			v := (in).(int)
			if v != 0 {
				rv.current.Fade(v)
			}
		}
	}
}

func (rv *Reverting) DrawOffset(buff draw.Image, xOff, yOff float64) {
	rv.current.DrawOffset(buff, xOff, yOff)
}
func (rv *Reverting) Draw(buff draw.Image) {
	rv.current.Draw(buff)
}
func (rv *Reverting) GetRGBA() *image.RGBA {
	return rv.current.GetRGBA()
}
func (rv *Reverting) ShiftX(x float64) {
	rv.current.ShiftX(x)
}
func (rv *Reverting) GetX() float64 {
	return rv.current.GetX()
}
func (rv *Reverting) GetY() float64 {
	return rv.current.GetY()
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

func (rv *Reverting) FlipX() Modifiable {
	rv.current.FlipX()
	rv.mods[F_FlipX] = !(rv.mods[F_FlipX]).(bool)
	return rv
}
func (rv *Reverting) FlipY() Modifiable {
	rv.current.FlipY()
	rv.mods[F_FlipY] = !(rv.mods[F_FlipY]).(bool)
	return rv
}
func (rv *Reverting) ApplyColor(c color.Color) Modifiable {
	rv.Revert(F_ApplyColor)
	rv.current.ApplyColor(c)
	rv.mods[F_ApplyColor] = c
	return rv
}

func (rv *Reverting) FillMask(img image.RGBA) Modifiable {
	rv.Revert(F_FillMask)
	rv.current.FillMask(img)
	rv.mods[F_FillMask] = img
	return rv
}
func (rv *Reverting) ApplyMask(img image.RGBA) Modifiable {
	rv.Revert(F_ApplyMask)
	rv.current.ApplyMask(img)
	rv.mods[F_ApplyMask] = img
	return rv
}
func (rv *Reverting) Rotate(degrees int) Modifiable {
	rv.Revert(F_Rotate)
	rv.current.Rotate(degrees)
	rv.mods[F_Rotate] = degrees
	return rv
}
func (rv *Reverting) Scale(xRatio float64, yRatio float64) Modifiable {
	rv.Revert(F_Scale)
	rv.current.Scale(xRatio, yRatio)
	rv.mods[F_Scale] = [2]float64{xRatio, yRatio}
	return rv
}
func (rv *Reverting) Fade(alpha int) Modifiable {
	rv.Revert(F_Fade)
	rv.current.Fade(alpha)
	rv.mods[F_Fade] = alpha
	return rv
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
	case *Sequence:
		t.update()
	}
	switch t := rv.root.(type) {
	case *Animation:
		t.updateAnimation()
	case *Sequence:
		t.update()
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
	case *Sequence:
		t.Pause()
	}
	switch t := rv.root.(type) {
	case *Animation:
		t.playing = false
	case *Compound:
		t.Pause()
	case *Sequence:
		t.Pause()
	}

}

func (rv *Reverting) Unpause() {
	switch t := rv.current.(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	}
	switch t := rv.root.(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	}
}

func (rv *Reverting) String() string {
	return rv.current.String()
}

func (rv *Reverting) AlwaysDirty() bool {
	return false
}
