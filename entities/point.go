package entities

import (
	"strconv"

	"bitbucket.org/oakmoundstudio/oak/physics"
)

type Point struct {
	physics.Vector
}

func NewPoint(x, y float64) Point {
	return Point{physics.NewVector(x, y)}
}

func (p *Point) GetX() float64 {
	return p.X()
}
func (p *Point) GetY() float64 {
	return p.Y()
}
func (p *Point) SetPos(x, y float64) {
	p.SetLogicPos(x, y)
}
func (p *Point) GetLogicPos() (float64, float64) {
	return p.X(), p.Y()
}
func (p *Point) SetLogicPos(x, y float64) {
	p.Vector.SetPos(x, y)
}
func (p *Point) DistanceTo(x, y float64) float64 {
	return p.Distance(physics.NewVector(x, y))
}
func (p *Point) DistanceToPoint(p2 Point) float64 {
	return p.Distance(p2.Vector)
}

func (p *Point) String() string {
	x := strconv.FormatFloat(p.X(), 'f', 2, 32)
	y := strconv.FormatFloat(p.Y(), 'f', 2, 32)
	return "X(): " + x + ", Y(): " + y
}
