package collision

import (
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
)

type Caster2 struct {
	Filters   []CastFilter
	Limits    []CastLimit
	PointSize floatgeom.Point2
	PointSpan float64
	Rays      int
	// ConeSpread is represented in radians
	ConeSpread   float64
	CastDistance float64
	Tree         *Tree
	CenterCone   bool
	CenterPoints bool
}

var (
	DefaultCaster2 = &Caster2{
		PointSize:    floatgeom.Point2{.1, .1},
		PointSpan:    1.0,
		CastDistance: 300,
		Rays:         1,
	}
)

func NewCaster(opts ...CastOption) *Caster2 {
	c := DefaultCaster2.Copy()
	if c.Tree == nil {
		c.Tree = DefTree
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Caster2) Cast(origin, angle floatgeom.Point2) []Point {
	angleDelta := c.ConeSpread / float64(c.Rays)
	points := make([]Point, 0)
	resultHash := make(map[*Space]bool)

	x := origin.X()
	y := origin.Y()

	a := angle.ToRadians()

	if c.CenterCone {
		a -= c.ConeSpread / 2
	}

	for degrees := a; degrees <= a+c.ConeSpread; degrees += angleDelta {

		sin := math.Sin(degrees)
		cos := math.Cos(degrees)

		for i := 0.0; i < c.CastDistance; i += c.PointSpan {

			hits := c.Tree.SearchIntersect(NewRect(x, y, c.PointSize.X(), c.PointSize.Y()))

		hitLoop:
			for k := 0; k < len(hits); k++ {
				next := hits[k]
				if _, ok := resultHash[next]; !ok {
					resultHash[next] = true

					for _, f := range c.Filters {
						if !f(next) {
							continue hitLoop
						}
					}

					points = append(points, NewPoint(next, x, y))

					for _, l := range c.Limits {
						if !l(points) {
							return points
						}
					}
				}
			}
			x += cos
			y += sin
		}
	}
	return points
}

func (c *Caster2) Transform(opts ...CastOption) *Caster2 {
	c2 := c.Copy()
	for _, opt := range opts {
		opt(c2)
	}
	return c2
}

func (c *Caster2) Copy() *Caster2 {
	c2 := new(Caster2)
	*c2 = *c
	return c2
}

type CastOption func(*Caster2)

// func Cone(coneSpan float64, rays int) CastOption {
// Todo
// }
