package alg

import (
	"math"
	"testing"
)

func TestDegrees(t *testing.T) {
	t.Parallel()
	t.Run("DegreeToRadians", func(t *testing.T) {
		t.Parallel()
		var f Degree = 90
		if Radian(f*DegToRad) != f.Radians() {
			t.Fatalf("degree to radian identity failed")
		}
	})
	t.Run("RadianToDegrees", func(t *testing.T) {
		t.Parallel()
		var f2 Radian = math.Pi / 2
		if Degree(f2*RadToDeg) != f2.Degrees() {
			t.Fatalf("radian to degree identity failed")
		}
	})
}
