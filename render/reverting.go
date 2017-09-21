package render

import (
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render/mod"
)

// The Reverting structure lets modifications be made to a Modifiable and then
// reverted, up to arbitrary history limits.
type Reverting struct {
	Modifiable
	rs []Modifiable
}

// NewReverting returns a Reverting type wrapped around the given modifiable
func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.rs = make([]Modifiable, 1)
	rv.rs[0] = m
	rv.Modifiable = m
	return rv
}

// Revert goes back n steps in this Reverting's history and displays that Modifiable
func (rv *Reverting) Revert(n int) {
	x := rv.X()
	y := rv.Y()

	if n >= len(rv.rs) {
		n = len(rv.rs) - 1
	}
	if n < 0 {
		return
	}

	rv.rs = rv.rs[:len(rv.rs)-n]
	rv.Modifiable = rv.rs[len(rv.rs)-1]
	rv.SetPos(x, y)
}

// RevertAll resets this reverting to its original Modifiable
func (rv *Reverting) RevertAll() {
	rv.Revert(len(rv.rs) - 1)
}

// RevertAndModify reverts n steps and then modifies this reverting. This
// is a separate function from Revert followed by Modify to prevent skipped
// draw frames.
func (rv *Reverting) RevertAndModify(n int, ms ...mod.Mod) Modifiable {
	x := rv.X()
	y := rv.Y()
	if n >= len(rv.rs) {
		n = len(rv.rs) - 1
	}
	if n > 0 {
		rv.rs = rv.rs[:len(rv.rs)-n]
	}
	rv.rs = append(rv.rs, rv.rs[len(rv.rs)-1].Copy().Modify(ms...))
	rv.Modifiable = rv.rs[len(rv.rs)-1]
	rv.SetPos(x, y)
	return rv
}

// Modify alters this reverting by the given modifications, appending the new
// modified renderable to it's list of modified versions and displaying it.
func (rv *Reverting) Modify(ms ...mod.Mod) Modifiable {
	next := rv.Modifiable.Copy().Modify(ms...)
	rv.rs = append(rv.rs, next)
	rv.Modifiable = rv.rs[len(rv.rs)-1]
	return rv
}

// Copy returns a copy of this Reverting
func (rv *Reverting) Copy() Modifiable {
	newRv := new(Reverting)
	newRv.rs = make([]Modifiable, len(rv.rs))
	for i, r := range rv.rs {
		newRv.rs[i] = r.Copy()
	}
	newRv.Modifiable = newRv.rs[len(rv.rs)-1]
	return newRv
}

// This might not ever be called?
func (rv *Reverting) update() {
	if u, ok := rv.Modifiable.(updates); ok {
		u.update()
	}
	if u, ok := rv.rs[0].(updates); ok {
		u.update()
	}
}

// SetTriggerID sets the ID AnimationEnd will trigger on for animating subtypes.
func (rv *Reverting) SetTriggerID(cid event.CID) {
	if t, ok := rv.Modifiable.(Triggerable); ok {
		t.SetTriggerID(cid)
	}
	if t, ok := rv.rs[0].(Triggerable); ok {
		t.SetTriggerID(cid)
	}
}

// Pause ceases animating any renderable types that animate underneath this
func (rv *Reverting) Pause() {
	if cp, ok := rv.Modifiable.(CanPause); ok {
		cp.Pause()
	}
	if cp, ok := rv.rs[0].(CanPause); ok {
		cp.Pause()
	}
}

// Unpause resumes animating any renderable types that animate underneath this
func (rv *Reverting) Unpause() {
	if cp, ok := rv.Modifiable.(CanPause); ok {
		cp.Unpause()
	}
	if cp, ok := rv.rs[0].(CanPause); ok {
		cp.Unpause()
	}
}

// IsInterruptable returns if whatever this reverting is currently dispalying is interruptable.
func (rv *Reverting) IsInterruptable() bool {
	if i, ok := rv.rs[0].(NonInterruptable); ok {
		return i.IsInterruptable()
	}
	return true
}

// IsStatic returns if whatever this reverting is currently displaying is static.
func (rv *Reverting) IsStatic() bool {
	if s, ok := rv.rs[0].(NonStatic); ok {
		return s.IsStatic()
	}
	return true
}

// Set calls Set on underlying types below this Reverting that cat be Set
// Todo: if Set becomes used by more types, this should use an interface like
// CanPause
func (rv *Reverting) Set(k string) error {
	var err error
	switch t := rv.Modifiable.(type) {
	case *Switch:
		err = t.Set(k)
		if err != nil {
			return err
		}
	}
	switch t := rv.rs[0].(type) {
	case *Switch:
		err = t.Set(k)
	}
	return err
}
