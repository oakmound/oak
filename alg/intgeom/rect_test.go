package intgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRectEq(t *testing.T) {
	r1 := NewRect(1, 1, 2, 2)
	r2 := NewRectWH(1, 1, 1, 1)
	assert.Equal(t, r1, r2)
}

func TestRectShuffle(t *testing.T) {
	rects := []Rect{
		NewRect(0, 0, 1, 1),
		NewRect(0, 0, 2, 2),
		NewRect(0, 0, 3, 3),
	}
	rects = ShuffleRects(rects)
	assert.Equal(t, len(rects), 3)
}
