package particle

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
)

func And(as ...func(Generator)) func(Generator) {
	return func(g Generator) {
		for _, a := range as {
			a(g)
		}
	}
}

func NewPerFrame(npf floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().NewPerFrame = npf
	}
}

func Pos(x, y float64) func(Generator) {
	return func(g Generator) {
		g.SetPos(x, y)
	}
}

func LifeSpan(ls floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LifeSpan = ls
	}
}

func Angle(a floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Angle = a
	}
}

func Speed(s floatrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Speed = s
	}
}

func Spread(x, y float64) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Spread = physics.NewVector(x, y)
	}
}

func Duration(i intrange.Range) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().Duration = i
	}
}

func Rotation(a floatrange.Range) func(Generator) {
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

func Layer(l func(physics.Vector) int) func(Generator) {
	return func(g Generator) {
		g.GetBaseGenerator().LayerFunc = l
	}
}
