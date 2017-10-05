package particle

import (
	"image/color"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/render"
)

// A GradientGenerator is a ColorGenerator with a patterned gradient
// on its particles
type GradientGenerator struct {
	ColorGenerator
	StartColor2, StartColor2Rand color.Color
	EndColor2, EndColor2Rand     color.Color
	ProgressFunction             func(x, y, w, h int) float64
}

// NewGradientGenerator returns a new GradientGenerator
func NewGradientGenerator(options ...func(Generator)) Generator {
	g := new(GradientGenerator)
	g.setDefaults()

	for _, opt := range options {
		opt(g)
	}

	return g
}

func (gg *GradientGenerator) setDefaults() {
	gg.ColorGenerator.setDefaults()
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
		gg.Rotation = gg.Rotation.Mult(alg.DegToRad)
	}
	return NewSource(gg, layer)
}

// GenerateParticle creates a particle from a generator
func (gg *GradientGenerator) GenerateParticle(bp *baseParticle) Particle {
	return &GradientParticle{
		ColorParticle: ColorParticle{
			baseParticle: bp,
			startColor:   randColor(gg.StartColor, gg.StartColorRand),
			endColor:     randColor(gg.EndColor, gg.EndColorRand),
			size:         float64(gg.Size.Poll()),
		},
		startColor2: randColor(gg.StartColor2, gg.StartColor2Rand),
		endColor2:   randColor(gg.EndColor2, gg.EndColor2Rand),
	}
}

// Gradient Coloration
//

// SetStartColor2 satisfies Colorable2
func (gg *GradientGenerator) SetStartColor2(sc, scr color.Color) {
	gg.StartColor2 = sc
	gg.StartColor2Rand = scr
}

// SetEndColor2 satisfies Colorable2
func (gg *GradientGenerator) SetEndColor2(ec, ecr color.Color) {
	gg.EndColor2 = ec
	gg.EndColor2Rand = ecr
}

// A Progresses has a SetProgress function where a progress function
// returns how far between two colors a given coordinate in a space is
type Progresses interface {
	SetProgress(func(x, y, w, h int) float64)
}

// Progress sets a Progresses' Progress Function
func Progress(pf func(x, y, w, h int) float64) func(Generator) {
	return func(g Generator) {
		g.(Progresses).SetProgress(pf)
	}
}

// SetProgress satisfies Progresses
func (gg *GradientGenerator) SetProgress(pf func(x, y, w, h int) float64) {
	gg.ProgressFunction = pf
}
