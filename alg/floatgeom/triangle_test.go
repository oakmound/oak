package floatgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriangleNormal(t *testing.T) {
	a := Tri3{
		Point3{0, 0, 0},
		Point3{1, 0, 0},
		Point3{0, 1, 0},
	}
	e := Point3{0, 0, 1}
	assert.Equal(t, e, a.Normal())
}

func TestTriangleBary(t *testing.T) {
	a := Tri3{
		Point3{0, 0, 0},
		Point3{1, 0, 0},
		Point3{0, 1, 0},
	}
	e := Point3{1, 1, -1}
	assert.Equal(t, e, a.Barycentric(1, 1))
	e = Point3{0.5, 0.5, 0}
	assert.Equal(t, e, a.Barycentric(.5, .5))
}
