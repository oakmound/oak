package render

type Reverting struct {
	Modifiable
	rs []Modifiable
}

func NewReverting(m Modifiable) *Reverting {
	rv := new(Reverting)
	rv.rs = make([]Modifiable, 1)
	rv.rs[0] = m
	rv.Modifiable = m
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

func (rv *Reverting) RevertAll() {
	rv.Revert(len(rv.rs) - 1)
}

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

func (rv *Reverting) Modify(ms ...Modification) Modifiable {
	next := rv.Modifiable.Copy().Modify(ms...)
	rv.rs = append(rv.rs, next)
	rv.Modifiable = rv.rs[len(rv.rs)-1]
	return rv
}

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
