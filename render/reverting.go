package render

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

// IsInterruptable returns if the underlying Modifiable for this reverting is interruptable.
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

// IsStatic returns if the underlying Modifiable for this reverting is static.
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

// Revert goes back n steps in this Reverting's history and displays that
// Modifiable
func (rv *Reverting) Revert(n int) {
	x := rv.GetX()
	y := rv.GetY()

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
func (rv *Reverting) RevertAndModify(n int, ms ...Modification) Modifiable {
	x := rv.GetX()
	y := rv.GetY()
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
func (rv *Reverting) Modify(ms ...Modification) Modifiable {
	next := rv.Modifiable.Copy().Modify(ms...)
	rv.rs = append(rv.rs, next)
	rv.Modifiable = rv.rs[len(rv.rs)-1]
	return rv
}

// Copy returns a copy of the Reverting
func (rv *Reverting) Copy() Modifiable {
	newRv := new(Reverting)
	newRv.rs = make([]Modifiable, len(rv.rs))
	for i, r := range rv.rs {
		newRv.rs[i] = r.Copy()
	}
	newRv.Modifiable = newRv.rs[len(rv.rs)-1]
	return newRv
}

func (rv *Reverting) updateAnimation() {
	switch t := rv.Modifiable.(type) {
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

// Set calls Set on underlying types below this Reverting that cat be Set
func (rv *Reverting) Set(k string) {
	switch t := rv.Modifiable.(type) {
	case *Compound:
		t.Set(k)
	}
	switch t := rv.rs[0].(type) {
	case *Compound:
		t.Set(k)
	}
}

// Pause ceases animating any renderable types that animate underneath this
func (rv *Reverting) Pause() {
	switch t := rv.Modifiable.(type) {
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

// Unpause resumes animating any renderable types that animate underneath this
func (rv *Reverting) Unpause() {
	switch t := rv.Modifiable.(type) {
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
