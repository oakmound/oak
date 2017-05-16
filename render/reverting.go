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
	rs []Modifiable
}

func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.rs = make([]Modifiable, 1)
	rv.rs[0] = m
	return rv
}

func (rv *Reverting) IsInterruptable() bool {
	switch t := rv.rs[0].(type) {
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
	switch t := rv.rs[0].(type) {
	case *Animation, *Sequence:
		return false
	case *Reverting:
		return t.IsStatic()
	case *Compound:
		return t.IsStatic()
	}
	return true
}

func (rv *Reverting) Revert(n int) {
	x := rv.current().GetX()
	y := rv.current().GetY()

	if n >= len(rv.rs) {
		n = len(rv.rs) - 1
	}

	rv.rs = rv.rs[:len(rv.rs)-n]
	rv.SetPos(x, y)
}

func (rv *Reverting) RevertAll() {
	rv.Revert(len(rv.rs) - 1)
}

func (rv *Reverting) DrawOffset(buff draw.Image, xOff, yOff float64) {
	rv.current().DrawOffset(buff, xOff, yOff)
}
func (rv *Reverting) Draw(buff draw.Image) {
	rv.current().Draw(buff)
}
func (rv *Reverting) GetRGBA() *image.RGBA {
	return rv.current().GetRGBA()
}
func (rv *Reverting) ShiftX(x float64) {
	rv.current().ShiftX(x)
}
func (rv *Reverting) GetX() float64 {
	return rv.current().GetX()
}
func (rv *Reverting) GetY() float64 {
	return rv.current().GetY()
}
func (rv *Reverting) ShiftY(y float64) {
	rv.current().ShiftY(y)
}
func (rv *Reverting) SetPos(x, y float64) {
	rv.current().SetPos(x, y)
}
func (rv *Reverting) GetDims() (int, int) {
	return rv.current().GetDims()
}
func (rv *Reverting) GetLayer() int {
	return rv.current().GetLayer()
}
func (rv *Reverting) SetLayer(l int) {
	rv.current().SetLayer(l)
}
func (rv *Reverting) UnDraw() {
	rv.current().UnDraw()
}

func (rv *Reverting) current() Modifiable {
	return rv.rs[len(rv.rs)-1]
}

func (rv *Reverting) Modify(ms ...Modification) Modifiable {
	next := rv.current().Copy().Modify(ms...)
	rv.rs = append(rv.rs, next)
	return rv
}

func (rv *Reverting) Copy() Modifiable {
	newRv := new(Reverting)
	newRv.rs = make([]Modifiable, len(rv.rs))
	for i, r := range rv.rs {
		newRv.rs[i] = r.Copy()
	}
	return newRv
}

func (rv *Reverting) updateAnimation() {
	switch t := rv.current().(type) {
	case *Animation:
		t.updateAnimation()
	case *Sequence:
		t.update()
	}
	switch t := rv.rs[0].(type) {
	case *Animation:
		t.updateAnimation()
	case *Sequence:
		t.update()
	}
}

func (rv *Reverting) Set(k string) {
	switch t := rv.current().(type) {
	case *Compound:
		t.Set(k)
	}
	switch t := rv.rs[0].(type) {
	case *Compound:
		t.Set(k)
	}
}

func (rv *Reverting) Pause() {
	switch t := rv.current().(type) {
	case *Animation:
		t.playing = false
	case *Compound:
		t.Pause()
	case *Sequence:
		t.Pause()
	}
	switch t := rv.rs[0].(type) {
	case *Animation:
		t.playing = false
	case *Compound:
		t.Pause()
	case *Sequence:
		t.Pause()
	}

}

func (rv *Reverting) Unpause() {
	switch t := rv.current().(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	}
	switch t := rv.rs[0].(type) {
	case *Animation:
		t.playing = true
	case *Compound:
		t.Unpause()
	case *Sequence:
		t.Unpause()
	}
}

func (rv *Reverting) String() string {
	return rv.current().String()
}
