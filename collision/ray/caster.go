package ray

import (
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
)

var (
	// DefaultCaster is a global caster that all
	// NewCaster() calls are built on before options are applied.
	DefaultCaster = &Caster{
		PointSize: floatgeom.Point2{.1, .1},
		PointSpan: 1.0,
		// CastDistance needs to be defined, but
		// there isn't a reasonable default.
		// Consider: Cast() could take in distance as well.
		CastDistance: 200,
	}
)

// SetDefaultCaster sets the global caster to be the input, and
// sets the caster behind the global cone caster as well.
func SetDefaultCaster(caster *Caster) {
	DefaultCaster = caster
	DefaultConeCaster.Caster = DefaultCaster
}

// A Caster can cast rays and return the colliding collision points
// of rays cast from points at angles. This behavior is customizable
// through CastOptions.
type Caster struct {
	Filters      []CastFilter
	Limits       []CastLimit
	PointSize    floatgeom.Point2
	PointSpan    float64
	CastDistance float64
	Tree         *collision.Tree
	CenterPoints bool
}

// A CastOption represents a transformation to a ray caster.
type CastOption func(*Caster)

// NewCaster will copy and modify the DefaultCaster by the input options
// and return the modified Caster. Giving no inputs is valid.
func NewCaster(opts ...CastOption) *Caster {
	c := DefaultCaster.Copy()
	if c.Tree == nil {
		c.Tree = collision.DefTree
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// CastTo casts a ray from origin to target, and otherwise acts as Cast.
func (c *Caster) CastTo(origin, target floatgeom.Point2) []collision.Point {
	return c.Cast(origin, floatgeom.AnglePoint(origin.AngleTo(target)))
}

// Cast creates a ray from origin pointing at the given angle and returns
// some spaces collided with at the point of collision, given the settings of
// this Caster. By default, all spaces hit will be returned.
func (c *Caster) Cast(origin, angle floatgeom.Point2) []collision.Point {
	points := make([]collision.Point, 0)
	resultHash := make(map[*collision.Space]bool)

	x := origin.X()
	y := origin.Y()

	degrees := angle.ToRadians()
	sin := math.Sin(degrees)
	cos := math.Cos(degrees)

	for i := 0.0; i < c.CastDistance; i += c.PointSpan {

		hits := c.Tree.SearchIntersect(
			collision.NewRect(x, y, c.PointSize.X(), c.PointSize.Y()),
		)

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

				points = append(points, collision.NewPoint(next, x, y))

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
	return points
}

// Copy copies a Caster.
func (c *Caster) Copy() *Caster {
	c2 := new(Caster)
	*c2 = *c
	return c2
}

// Tree sets the collision tree of a Caster.
func Tree(t *collision.Tree) CastOption {
	return func(c *Caster) {
		c.Tree = t
	}
}

// CenterPoints sets whether a Caster should center its collision points that
// form its ray. This is by default false, and is only significant if said
// points' dimensions are significantly large.
func CenterPoints(on bool) CastOption {
	return func(c *Caster) {
		c.CenterPoints = on
	}
}
