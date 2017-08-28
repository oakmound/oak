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
