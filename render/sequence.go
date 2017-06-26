package render

import (
	"image"
	"image/draw"
	"time"

	"bitbucket.org/oakmoundstudio/oak/event"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/timing"
)

type Sequence struct {
	LayeredPoint
	rs            []Modifiable
	lastChange    time.Time
	sheetPos      int
	frameTime     int64
	cID           event.CID
	playing       bool
	Interruptable bool
}

func NewSequence(mods []Modifiable, fps float64) *Sequence {
	return &Sequence{
		LayeredPoint: LayeredPoint{
			Vector: physics.NewVector(0, 0),
		},
		sheetPos:      0,
		frameTime:     timing.FPSToNano(fps),
		rs:            mods,
		lastChange:    time.Now(),
		playing:       true,
		Interruptable: true,
	}
}

func (sq *Sequence) Copy() Modifiable {
	newSq := new(Sequence)
	*newSq = *sq

	newRs := make([]Modifiable, len(sq.rs))
	for i, r := range sq.rs {
		newRs[i] = r.Copy()
	}

	newSq.rs = newRs
	return newSq
}

func (sq *Sequence) SetTriggerID(id event.CID) {
	sq.cID = id
}

func (sq *Sequence) update() {
	if sq.playing && time.Since(sq.lastChange).Nanoseconds() > sq.frameTime {
		sq.lastChange = time.Now()
		sq.sheetPos = (sq.sheetPos + 1) % len(sq.rs)
		if sq.sheetPos == (len(sq.rs)-1) && sq.cID != 0 {
			sq.cID.Trigger("AnimationEnd", nil)
		}
	}
}

func (sq *Sequence) Get(i int) Modifiable {
	return sq.rs[i]
}

func (sq *Sequence) DrawOffset(buff draw.Image, xOff, yOff float64) {
	sq.update()
	sq.rs[sq.sheetPos].DrawOffset(buff, sq.X()+xOff, sq.Y()+yOff)
}

func (sq *Sequence) Draw(buff draw.Image) {
	sq.update()
	sq.rs[sq.sheetPos].DrawOffset(buff, sq.X(), sq.Y())
}

func (sq *Sequence) GetRGBA() *image.RGBA {
	return sq.rs[sq.sheetPos].GetRGBA()
}

func (sq *Sequence) Modify(ms ...Modification) Modifiable {
	for _, r := range sq.rs {
		r.Modify(ms...)
	}
	return sq
}

func (sq *Sequence) Pause() {
	sq.playing = false
}

func (sq *Sequence) Unpause() {
	sq.playing = true
}

func TweenSequence(a, b image.Image, frames int, fps float64) *Sequence {
	images := Tween(a, b, frames)
	ms := make([]Modifiable, len(images))
	for i, v := range images {
		ms[i] = NewSprite(0, 0, v)
	}
	return NewSequence(ms, fps)
}
