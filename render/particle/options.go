package particle

import (
	"bitbucket.org/oakmoundstudio/oak/alg"
	"bitbucket.org/oakmoundstudio/oak/physics"
)

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
