package timing

import (
	"context"
	"time"
)

var (
	// ClearDelayCh is used to stop all ongoing delays. It should not be closed.
	ClearDelayCh = make(chan struct{})
)

// DoAfter wraps time calls in a select that will stop events from happening
// when ClearDelayCh pulls
func DoAfter(d time.Duration, f func()) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		f()
	case <-ClearDelayCh:
	}
}

// DoAfterContext executes the function if the context is completed.
// Clears out if the Delay Channel is cleared.
func DoAfterContext(ctx context.Context, f func()) {
	select {
	case <-ctx.Done():
		f()
	case <-ClearDelayCh:
	}
}
