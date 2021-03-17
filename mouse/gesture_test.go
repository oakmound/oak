package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/v2/physics"

	"github.com/oakmound/shiny/gesture"
)

func TestGestureIdentity(t *testing.T) {
	tm := time.Now()
	ge := gesture.Event{
		Type:        gesture.TypeStart,
		Drag:        true,
		LongPress:   false,
		DoublePress: false,
		InitialPos:  gesture.Point{X: 2.0, Y: 3.0},
		CurrentPos:  gesture.Point{X: 4.0, Y: 5.0},
		Time:        tm,
	}
	mge := FromShinyGesture(ge)
	if mge.DoublePress != false {
		t.Fatalf("got %v expected %v", mge.DoublePress, false)
	}
	if mge.LongPress != false {
		t.Fatalf("got %v expected %v", mge.LongPress, false)
	}
	if mge.Drag != true {
		t.Fatalf("got %v expected %v", mge.Drag, true)
	}
	if mge.InitialPos.X() != 2 || mge.InitialPos.Y() != 3.0 {
		t.Fatalf("got %v expected %v", mge.InitialPos, physics.NewVector(2.0, 3.0))
	}
	if mge.CurrentPos.X() != 4.0 || mge.CurrentPos.Y() != 5.0 {
		t.Fatalf("got %v expected %v", mge.CurrentPos, physics.NewVector(4.0, 5.0))
	}
	if mge.Time != tm {
		t.Fatalf("got %v expected %v", mge.Time, tm)
	}
}
