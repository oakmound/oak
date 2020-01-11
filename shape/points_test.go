package shape

import (
	"testing"

	"github.com/oakmound/oak/v2/alg/intgeom"
	"github.com/stretchr/testify/assert"
)

var (
	testPoints = NewPoints(
		intgeom.Point2{1, 1}, intgeom.Point2{2, 1}, intgeom.Point2{3, 1},
		intgeom.Point2{1, 2}, intgeom.Point2{3, 2},
		intgeom.Point2{1, 3}, intgeom.Point2{2, 3}, intgeom.Point2{3, 3},
	)
)

func TestPointsIn(t *testing.T) {
	assert.True(t, testPoints.In(1, 3, 1, 1))
	assert.False(t, testPoints.In(10, 10, 1, 1))
}

func TestPointsOutline(t *testing.T) {
	testOutline, _ := testPoints.Outline(4, 4)
	assert.Equal(t, intgeom.Point2{3, 2}, testOutline[3])
}

func TestPointsRect(t *testing.T) {
	testRect := testPoints.Rect(4, 4)
	assert.False(t, testRect[0][0])
	assert.False(t, testRect[2][2])
	assert.True(t, testRect[1][1])
}
