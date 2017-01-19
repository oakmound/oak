package particle

import (
	"image/color"
	"math"

	"bitbucket.org/oakmoundstudio/oak/render"
)

type GradientGenerator struct {
	ColorGenerator
	StartColor2, StartColor2Rand color.Color
	EndColor2, EndColor2Rand     color.Color
	ProgressFunction             func(x, y, w, h int) float64
}

func NewGradientGenerator(options ...func(Generator)) Generator {
	g := new(GradientGenerator)
	g.SetDefaults()

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (gg *GradientGenerator) SetDefaults() {
	gg.ColorGenerator.SetDefaults()
	gg.StartColor2 = color.RGBA{0, 0, 0, 0}
	gg.StartColor2Rand = color.RGBA{0, 0, 0, 0}
	gg.EndColor2 = color.RGBA{0, 0, 0, 0}
	gg.EndColor2Rand = color.RGBA{0, 0, 0, 0}
	gg.ProgressFunction = render.HorizontalProgress
}

// Generate takes a generator and converts it into a source,
// drawing particles and binding functions for particle generation
// and rotation.
func (gg *GradientGenerator) Generate(layer int) *Source {
	// Convert rotation from degrees to radians
	if gg.Rotation != nil {
		gg.Rotation = gg.Rotation.Mult(math.Pi / 180)
	}
	return NewSource(gg)
}

func (gg *GradientGenerator) GenerateParticle(bp *BaseParticle) Particle {
	return &GradientParticle{
		BaseParticle: bp,
		startColor:   randColor(gg.StartColor, gg.StartColorRand),
		endColor:     randColor(gg.EndColor, gg.EndColorRand),
		startColor2:  randColor(gg.StartColor2, gg.StartColor2Rand),
		endColor2:    randColor(gg.EndColor2, gg.EndColor2Rand),
		size:         gg.Size.Poll(),
	}
}

func (gg *GradientGenerator) GetBaseGenerator() *BaseGenerator {
	return &gg.BaseGenerator
}

func (gp *GradientGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

// Gradient Coloration
//

type Colorable2 interface {
	SetStartColor2(color.Color, color.Color)
	SetEndColor2(color.Color, color.Color)
}

func Color2(sc, scr, ec, ecr color.Color) func(Generator) {
	return func(g Generator) {
		c := g.(Colorable2)
		c.SetStartColor2(sc, scr)
		c.SetEndColor2(ec, ecr)
	}
}

func (gp *GradientGenerator) SetStartColor2(sc, scr color.Color) {
	gp.StartColor2 = sc
	gp.StartColor2Rand = scr
}

func (gp *GradientGenerator) SetEndColor2(ec, ecr color.Color) {
	gp.EndColor2 = ec
	gp.EndColor2Rand = ecr
}

type Progresses interface {
	SetProgress(func(x, y, w, h int) float64)
}

func Progress(pf func(x, y, w, h int) float64) func(Generator) {
	return func(g Generator) {
		g.(Progresses).SetProgress(pf)
	}
}

func (gp *GradientGenerator) SetProgress(pf func(x, y, w, h int) float64) {
	gp.ProgressFunction = pf
}
