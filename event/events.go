package event

import (
	"sync/atomic"
	"time"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

type UnsafeEventID int64

type EventID[T any] struct {
	UnsafeEventID
}

var (
	nextEventID int64
)

const NoEvent = 0

func RegisterEvent[T any]() EventID[T] {
	id := atomic.AddInt64(&nextEventID, 1)
	return EventID[T]{
		UnsafeEventID: UnsafeEventID(id),
	}
}

type NoPayload struct{}

// EnterPayload is the payload sent down to Enter bindings
type EnterPayload struct {
	FramesElapsed  int
	SinceLastFrame time.Duration
	TickPercent    float64
}

var (
	// Enter: the beginning of every logical frame.
	Enter = RegisterEvent[EnterPayload]()
	// AnimationEnd: Triggered on animations CIDs when they loop from the last to the first frame
	AnimationEnd = RegisterEvent[NoPayload]()
	// ViewportUpdate: Triggered when the position of of the viewport changes
	ViewportUpdate = RegisterEvent[intgeom.Point2]()
	// OnStop: Triggered when the engine is stopped.
	OnStop = RegisterEvent[NoPayload]()
	// FocusGain: Triggered when the window gains focus
	FocusGain = RegisterEvent[NoPayload]()
	// FocusLoss: Triggered when the window loses focus
	FocusLoss = RegisterEvent[NoPayload]()
)
