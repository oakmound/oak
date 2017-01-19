package particle

import (
	"bitbucket.org/oakmoundstudio/oak/alg"
	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

var (
	// Inf represents Infinite duration
	Inf = alg.Infinite{}
)

type Generator interface {
	GetBaseGenerator() *BaseGenerator
	GenerateParticle(*BaseParticle) Particle
	Generate(int) *Source
	GetParticleSize() (float64, float64, bool)
	ShiftX(float64)
	ShiftY(float64)
	SetPos(float64, float64)
	GetPos() (float64, float64)
}

// A BaseGenerator fulfills all basic requirements to generate
// particles
// Modeled after Parcycle
type BaseGenerator struct {
	render.Point
	// This float is currently forced to an integer
	// at new particle rotation. This should be changed
	// to something along the lines of 'new per 30 frames',
	// or allow low fractional values to be meaningful,
	// so that more fine-tuned particle generation speeds are possible.
	NewPerFrame alg.FloatRange
	// The number of frames each particle should persist
	// before being removed.
	LifeSpan alg.FloatRange
	// 0 - between quadrant 1 and 4
	// 90 - between quadrant 2 and 1
	Angle  alg.FloatRange
	Speed  alg.FloatRange
	Spread *physics.Vector
	// Duration in milliseconds for the particle source.
	// After this many milliseconds have passed, it will
	// stop sending out new particles. Old particles will
	// not be removed until their individual lifespans run
	// out.
	// A duration of -1 represents never stopping.
	Duration alg.IntRange
	// Rotational acceleration, to change angle over time
	Rotation alg.FloatRange
	// Gravity X and Gravity Y represent particle acceleration per frame.
	Gravity    *physics.Vector
	SpeedDecay *physics.Vector
	EndFunc    func(Particle)
	LayerFunc  func(*physics.Vector) int
}

func (bg *BaseGenerator) SetDefaults() {
	*bg = BaseGenerator{
		Point:       render.Point{X: 0, Y: 0},
		NewPerFrame: alg.Constantf(1),
		LifeSpan:    alg.Constantf(60),
		Angle:       alg.Constantf(0),
		Speed:       alg.Constantf(1),
		Spread:      physics.NewVector(0, 0),
		Duration:    Inf,
		Rotation:    nil,
		Gravity:     nil,
		SpeedDecay:  nil,
		EndFunc:     nil,
		LayerFunc:   func(*physics.Vector) int { return 1 },
	}
}

func NewPerFrame(npf alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().NewPerFrame = npf
	}
}

func Pos(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().SetPos(x, y)
	}
}

func LifeSpan(ls alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LifeSpan = ls
	}
}

func Angle(a alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Angle = a
	}
}

func Speed(s alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Speed = s
	}
}

func Spread(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Spread = physics.NewVector(x, y)
	}
}

func Duration(i alg.IntRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Duration = i
	}
}

func Rotation(a alg.FloatRange) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Rotation = a
	}
}

func Gravity(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Gravity = physics.NewVector(x, y)
	}
}

func SpeedDecay(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().SpeedDecay = physics.NewVector(x, y)
	}
}

func End(ef func(Particle)) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().EndFunc = ef
	}
}

func Layer(l func(*physics.Vector) int) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LayerFunc = l
	}
}
