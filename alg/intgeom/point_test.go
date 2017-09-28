package intgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointCompareOf(t *testing.T) {
	a := Point{0, 1}
	b := Point{1, 0}
	assert.Equal(t, Point{0, 0}, a.LesserOf(b))
	assert.Equal(t, Point{0, 0}, b.LesserOf(a))
	assert.Equal(t, Point{1, 1}, a.GreaterOf(b))
	assert.Equal(t, Point{1, 1}, b.GreaterOf(a))
}

func TestPointsBetween(t *testing.T) {
	a := Point{0, 0}
	b := Point{20, 20}
	ptsBtween := a.PointsBetween(b)
	assert.Equal(t, 21, len(ptsBtween))

	ptsBtween = b.PointsBetween(a)
	assert.Equal(t, 21, len(ptsBtween))
}
