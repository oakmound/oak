package btn

import "errors"

type reverting interface {
	Revert(n int)
}

// Revert will check that the given button's renderable
// can have modifications reverted, then revert the last
// n modifications.
func Revert(b Btn, n int) error {
	r, ok := b.GetRenderable().(reverting)
	if !ok {
		return errors.New("Button's renderable does not implement revert functionality")
	}
	r.Revert(n)
	return nil
}
