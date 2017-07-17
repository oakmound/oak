package entities

import (
	"strconv"

	"github.com/oakmound/oak/physics"
)

// A Point is a wrapper around a physics vector.
type Point struct {
	physics.Vector
}

// NewPoint returns a new point
func NewPoint(x, y float64) Point {
	return Point{physics.NewVector(x, y)}
}

// GetLogicPos returns the logical position of an entity. See SetLogicPos.
func (p *Point) GetLogicPos() (float64, float64) {
	return p.X(), p.Y()
}

// SetLogicPos is an explicit declaration for setting just the logical
// position of an entity. On a Point there is no distinction as there is nothing
// but the logical position but this is important for other entity types
func (p *Point) SetLogicPos(x, y float64) {
	p.Vector.SetPos(x, y)
}

// DistanceTo returns the euclidean distance to (x,y)
func (p *Point) DistanceTo(x, y float64) float64 {
	return p.Distance(physics.NewVector(x, y))
}

// DistanceToPoint returns the euclidean distance to p2.GetLogicPos()
func (p *Point) DistanceToPoint(p2 Point) float64 {
	return p.Distance(p2.Vector)
}

func (p *Point) String() string {
	x := strconv.FormatFloat(p.X(), 'f', 2, 32)
	y := strconv.FormatFloat(p.Y(), 'f', 2, 32)
	return "X(): " + x + ", Y(): " + y
}
