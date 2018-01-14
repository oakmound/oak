package render

import (
	"image"
	"image/draw"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/timing"
)

// A Sequence is a series of modifiables drawn as an animation. It is more
// primitive than animation, but less efficient.
type Sequence struct {
	LayeredPoint
	pauseBool
	InterruptBool
	rs         []Modifiable
	lastChange time.Time
	sheetPos   int
	frameTime  int64
	cID        event.CID
}

// NewSequence returns a new sequence from the input modifiables, playing at
// fps rate
func NewSequence(fps float64, mods ...Modifiable) *Sequence {
	return &Sequence{
		LayeredPoint: NewLayeredPoint(0, 0, 0),
		pauseBool: pauseBool{
			playing: true,
		},
		InterruptBool: InterruptBool{
			Interruptable: true,
		},
		sheetPos:   0,
		frameTime:  timing.FPSToNano(fps),
		rs:         mods,
		lastChange: time.Now(),
	}
}

// Copy copies each modifiable inside this sequence in order to produce a new
// copied sequence
func (sq *Sequence) Copy() Modifiable {
	newSq := new(Sequence)
	*newSq = *sq

	newRs := make([]Modifiable, len(sq.rs))
	for i, r := range sq.rs {
		newRs[i] = r.Copy()
	}

	newSq.rs = newRs
	newSq.LayeredPoint = sq.LayeredPoint.Copy()
	return newSq
}

// SetTriggerID sets the ID that AnimationEnd will be triggered on when this
// sequence loops over from its last frame to its first
func (sq *Sequence) SetTriggerID(id event.CID) {
	sq.cID = id
}

func (sq *Sequence) update() {
	if sq.playing && time.Since(sq.lastChange).Nanoseconds() > sq.frameTime {
		sq.lastChange = time.Now()
		sq.sheetPos = (sq.sheetPos + 1) % len(sq.rs)
		if sq.sheetPos == (len(sq.rs)-1) && sq.cID != 0 {
			sq.cID.Trigger(event.AnimationEnd, nil)
		}
	}
}

// Get returns the Modifiable stored at this sequence's ith index. If the sequence
// does not have an ith index this returns nil
func (sq *Sequence) Get(i int) Modifiable {
	if i < 0 || i >= len(sq.rs) {
		return nil
	}
	return sq.rs[i]
}

// DrawOffset draws this sequence at +xOff, +yOff
func (sq *Sequence) DrawOffset(buff draw.Image, xOff, yOff float64) {
	sq.update()
	sq.rs[sq.sheetPos].DrawOffset(buff, sq.X()+xOff, sq.Y()+yOff)
}

// Draw draws this sequence to the input buffer
func (sq *Sequence) Draw(buff draw.Image) {
	sq.update()
	sq.rs[sq.sheetPos].DrawOffset(buff, sq.X(), sq.Y())
}

// GetRGBA returns the RGBA of the currently showing frame of this sequence
func (sq *Sequence) GetRGBA() *image.RGBA {
	return sq.rs[sq.sheetPos].GetRGBA()
}

// Modify alters each renderable in this sequence by the given
// modifications
func (sq *Sequence) Modify(ms ...mod.Mod) Modifiable {
	for _, r := range sq.rs {
		r.Modify(ms...)
	}
	return sq
}

// Filter filters each element in the sequence by the inputs
func (sq *Sequence) Filter(fs ...mod.Filter) {
	for _, r := range sq.rs {
		r.Filter(fs...)
	}
}

// IsStatic returns false for sequences
func (sq *Sequence) IsStatic() bool {
	return false
}

// TweenSequence returns a sequence that is the tweening between the input images
// at the given frame rate over the given frame count.
func TweenSequence(a, b image.Image, frames int, fps float64) *Sequence {
	images := Tween(a, b, frames)
	ms := make([]Modifiable, len(images))
	for i, v := range images {
		ms[i] = NewSprite(0, 0, v)
	}
	return NewSequence(fps, ms...)
}
