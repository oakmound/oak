package shape

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRect(t *testing.T) {
	shapes := []Shape{Square, Rectangle, Diamond, Circle, Checkered, Heart, JustIn(NotIn(Diamond.In))}

	w, h := 10, 10
	for _, s := range shapes {
		r := InToRect(s.In)(w, h)
		assert.Equal(t, r, s.Rect(w, h))
		for i := 0; i < w; i++ {
			for j := 0; j < h; j++ {
				assert.Equal(t, r[i][j], s.In(i, j, w, h))
			}
		}
	}
}

func TestRectangleIn(t *testing.T) {
	assert.False(t, Rectangle.In(10, 10, 5, 5))
}

func TestStrictRect(t *testing.T) {
	sr := NewStrictRect(5, 5)
	for x := 0; x < 6; x++ {
		for y := 0; y < 6; y++ {
			assert.False(t, sr.In(x, y))
		}
	}
	r := sr.Rect()
	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			assert.False(t, r[x][y])
		}
	}

	sr[3][3] = true

	out, err := sr.Outline()
	assert.Nil(t, err)
	assert.Len(t, out, 1)
}
