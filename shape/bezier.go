package shape

import "github.com/oakmound/oak/oakerr"

// BezierCurve will form a Bezier on the given coordinates, expected in (x,y)
// pairs. If the inputs have an odd length, an error noting so is returned, and
// the Bezier returned is nil.
func BezierCurve(coords ...float64) (Bezier, error) {
	if len(coords) == 0 {
		return nil, oakerr.InsufficientInputs{AtLeast: 2, InputName: "coords"}
	}
	if len(coords)%2 != 0 {
		return nil, oakerr.IndivisibleInput{
			InputName:    "coords",
			IsList:       true,
			MustDivideBy: 2,
		}
	}
	pts := make([]Bezier, len(coords)/2)
	for i := 0; i < len(coords); i += 2 {
		pts[i/2] = BezierPoint{coords[i], coords[i+1]}
	}
	for len(pts) > 1 {
		for i := 0; i < len(pts)-1; i++ {
			pts[i] = BezierNode{pts[i], pts[i+1]}
		}
		pts = pts[:len(pts)-1]
	}
	return pts[0], nil
}

// A Bezier has a function indicating how far along a curve something is given
// some float64 progress between 0 and 1. This allows points, lines, and limitlessly complex
// bezier curves to be represented under this interface.
//
// Beziers will not necessarily break if given an input outside of 0 to 1, but the results
// shouldn't be relied upon.
type Bezier interface {
	Pos(progress float64) (x, y float64)
}

// A BezierNode ties together and find points between two other Beziers
type BezierNode struct {
	Left, Right Bezier
}

// Pos returns the a point progress percent between this node's left and
// right progress percent points.
func (bn BezierNode) Pos(progress float64) (x, y float64) {
	x1, y1 := bn.Left.Pos(progress)
	x2, y2 := bn.Right.Pos(progress)
	return x1 + ((x2 - x1) * progress), y1 + ((y2 - y1) * progress)
}

// A BezierPoint covers cases where only 1 point is supplied, and serve as roots.
// Consider: merging with floatgeom.Point2
type BezierPoint struct {
	X, Y float64
}

// Pos returns this point.
func (bp BezierPoint) Pos(float64) (x, y float64) {
	return bp.X, bp.Y
}
