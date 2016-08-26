package entities

import (
	"math"
)

type Point struct {
	X, Y float64
}

func (p *Point) GetX() float64 {
	return p.X
}
func (p *Point) GetY() float64 {
	return p.Y
}
func (p *Point) SetPos(x, y float64) {
	p.SetLogicPos(x, y)
}
func (p *Point) GetLogicPos() (float64, float64) {
	return p.X, p.Y
}
func (p *Point) SetLogicPos(x, y float64) {
	p.X = x
	p.Y = y
}
func (p *Point) DistanceTo(x, y float64) float64 {
	return distance(p.X, p.Y, x, y)
}
func (p *Point) DistanceToPoint(p2 Point) float64 {
	return distance(p.X, p.Y, p2.X, p2.Y)
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(
		math.Pow(x1-x2, 2) +
			math.Pow(y1-y2, 2))
}
