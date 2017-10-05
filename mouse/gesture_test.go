package mouse

import (
	"testing"
	"time"

	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"

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
	assert.Equal(t, mge.DoublePress, false)
	assert.Equal(t, mge.LongPress, false)
	assert.Equal(t, mge.Drag, true)
	assert.Equal(t, mge.InitialPos, physics.NewVector(2.0, 3.0))
	assert.Equal(t, mge.CurrentPos, physics.NewVector(4.0, 5.0))
	assert.Equal(t, mge.Time, tm)
}
