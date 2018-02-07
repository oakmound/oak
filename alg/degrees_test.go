package alg

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDegrees(t *testing.T) {
	var f Degree = 90
	assert.Equal(t, Radian(f*DegToRad), f.Radians())
	var f2 Radian = math.Pi / 2
	assert.Equal(t, Degree(f2*RadToDeg), f2.Degrees())
}
