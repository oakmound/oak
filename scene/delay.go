package scene

import (
	"context"
	"time"

	"github.com/oakmound/oak/v3/render"
)

// DoAfter will execute the given function after some duration. When the scene
// ends, DoAfter will exit without calling f. This call blocks until one of those
// conditions is reached.
func (c *Context) DoAfter(d time.Duration, f func()) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		f()
	case <-c.Done():
	}
}

// DoAfterContext will execute the given function once the passed in context is closed.
// When the scene ends, DoAfterContext will exit without calling f. This call blocks until
// one of those conditions is reached.
func (c *Context) DoAfterContext(ctx context.Context, f func()) {
	select {
	case <-ctx.Done():
		f()
	case <-c.Done():
	}
}

// DrawForTime draws, and after d, undraws an element
func (c *Context) DrawForTime(r render.Renderable, d time.Duration, layers ...int) error {
	_, err := c.DrawStack.Draw(r, layers...)
	if err != nil {
		return err
	}
	go func(r render.Renderable, d time.Duration) {
		c.DoAfter(d, func() {
			r.Undraw()
		})
	}(r, d)
	return nil
}
