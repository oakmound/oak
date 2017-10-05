package ray

import (
	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
)

var (
	// DefaultConeCaster is a global caster that all NewConeCaster()
	// calls are built on before options are applied.
	DefaultConeCaster = &ConeCaster{
		Caster:     DefaultCaster,
		CenterCone: true,
		Rays:       1,
	}
)

// SetDefaultConeCaster is analagous to SetDefaultCaster, however
// is equivalent to setting the global variable.
func SetDefaultConeCaster(coneCaster *ConeCaster) {
	DefaultConeCaster = coneCaster
}

// A ConeCaster will repeatedly Cast
// its underlying Caster in a cone shape.
type ConeCaster struct {
	*Caster
	CenterCone bool
	// ConeSpread is represented in radians
	ConeSpread float64
	Rays       float64
}

// A ConeCastOption represents a transformation on a ConeCaster.
type ConeCastOption func(*ConeCaster)

// NewConeCaster copies the DefaultConeCaster and modifies it with the input
// options, returning the modified Caster. Zero arguments is valid input.
func NewConeCaster(opts ...ConeCastOption) *ConeCaster {
	cc := DefaultConeCaster.Copy()
	for _, opt := range opts {
		opt(cc)
	}
	return cc
}

// CastTo casts a ray from origin to target, and otherwise acts as Cast.
func (cc *ConeCaster) CastTo(origin, target floatgeom.Point2) []collision.Point {
	return cc.Cast(origin, floatgeom.AnglePoint(origin.AngleTo(target)))
}

// Cast creates a ray from origin pointing at the given angle and returns
// some spaces collided with at the point of collision, given the settings of
// this ConeCaster. By default, all spaces hit will be returned. ConeCasters in
// addition will recast at progressive angles until they have cast up to their
// Rays value. Angles progress in counter-clockwise order.
func (cc *ConeCaster) Cast(origin, angle floatgeom.Point2) []collision.Point {
	points := make([]collision.Point, 0)
	if cc.Rays < 1 {
		return points
	}
	angleDelta := cc.ConeSpread / cc.Rays

	a := angle.ToRadians()
	if cc.CenterCone {
		a -= cc.ConeSpread / 2
	}

	for degrees := a; degrees <= a+cc.ConeSpread; degrees += angleDelta {
		points = append(points, cc.Caster.Cast(origin, floatgeom.RadianPoint(degrees))...)
	}
	return points
}

// Copy copies a ConeCaster.
func (cc *ConeCaster) Copy() *ConeCaster {
	cc2 := new(ConeCaster)
	*cc2 = *cc
	return cc2
}

// CenterCone sets whether the caster should center its cones around the
// input angles or progress out from those input angles. True by default.
func CenterCone(on bool) ConeCastOption {
	return func(cc *ConeCaster) {
		cc.CenterCone = on
	}
}

// ConeSpread sets how far a ConeCaster should progress its angles in degrees.
func ConeSpread(degrees float64) ConeCastOption {
	return func(cc *ConeCaster) {
		cc.ConeSpread = degrees * alg.DegToRad
	}
}

// ConeSpreadRadians sets how far a ConeCaster should progress its angles in radians.
func ConeSpreadRadians(radians float64) ConeCastOption {
	return func(cc *ConeCaster) {
		cc.ConeSpread = radians
	}
}

// ConeRays sets how many rays a ConeCaster should divide its spread into.
func ConeRays(rays int) ConeCastOption {
	return func(cc *ConeCaster) {
		cc.Rays = float64(rays)
	}
}
