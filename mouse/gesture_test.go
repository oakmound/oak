package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"

	"golang.org/x/exp/shiny/gesture"
)

func TestGestureIdentity(t *testing.T) {
	tm := time.Now()
	ge := gesture.Event{
		Type:        gesture.TypeStart,
		Drag:        true,
		LongPress:   false,
		DoublePress: false,
		InitialPos:  gesture.Point{2.0, 3.0},
		CurrentPos:  gesture.Point{4.0, 5.0},
		Time:        tm,
	}
	mge := FromShinyGesture(ge)
	assert.Equal(t, mge.DoublePress, false)
	assert.Equal(t, mge.LongPress, false)
	assert.Equal(t, mge.Drag, true)
	assert.Equal(t, mge.InitialPos, physics.NewVector(2.0, 3.0))
	assert.Equal(t, mge.CurrentPos, physics.NewVector(4.0, 5.0))
	assert.Equal(t, mge.Time, tm)
}
