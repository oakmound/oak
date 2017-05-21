package particle

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
)

// And chains together particle options into a single option
// for pre-baking option sets
func And(as ...func(Generator)) func(Generator) {
	return func(g Generator) {
		for _, a := range as {
			a(g)
		}
	}
}

// NewPerFrame sets how many particles should be produced per frame
func NewPerFrame(npf floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().NewPerFrame = npf
	}
}

// Pos sets the position of a generator
func Pos(x, y float64) func(Generator) {
	return func(g Generator) {
		g.SetPos(x, y)
	}
}

// LifeSpan sets how long a particle should last before dying
func LifeSpan(ls floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LifeSpan = ls
	}
}

// Angle sets the initial angle of a particle
func Angle(a floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Angle = a
	}
}

// Speed sets the initial speed of a particle
func Speed(s floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Speed = s
	}
}

// Spread sets how far from a generator's position a particle can spawn
func Spread(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Spread = physics.NewVector(x, y)
	}
}

// Duration sets how long a generator should produce particles for
func Duration(i intrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Duration = i
	}
}

// Rotation rotates particles by a variable amount per frame
func Rotation(a floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Rotation = a
	}
}

// Gravity sets how a particle should be shifted over time in either dimension
func Gravity(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Gravity = physics.NewVector(x, y)
	}
}

// SpeedDecay sets how the speed of a particle should decay
func SpeedDecay(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().SpeedDecay = physics.NewVector(x, y)
	}
}

// End sets what function should happen when a particle dies
func End(ef func(Particle)) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().EndFunc = ef
	}
}

// Layer sets a function to determine what draw layer a particle should exist on
func Layer(l func(physics.Vector) int) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LayerFunc = l
	}
}
