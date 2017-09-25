package floatgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	assert.Equal(t, a.GreaterOf(b), Point2{1, 1})

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	assert.Equal(t, c.GreaterOf(d), Point3{1, 2, 2})
}

func TestLesserOf(t *testing.T) {
	a := Point2{0, 1}
	b := Point2{1, 0}
	assert.Equal(t, a.LesserOf(b), Point2{0, 0})

	c := Point3{0, 1, 2}
	d := Point3{1, 2, 0}
	assert.Equal(t, c.LesserOf(d), Point3{0, 1, 0})
}

func TestDistance(t *testing.T) {

}
