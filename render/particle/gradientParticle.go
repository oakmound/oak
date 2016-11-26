package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type GradientGenerator struct {
	BaseGenerator
	StartColor, StartColorRand   color.Color
	StartColor2, StartColor2Rand color.Color
	EndColor, EndColorRand       color.Color
	EndColor2, EndColor2Rand     color.Color
	// The size, in pixel radius, of spawned particles
	Size, SizeRand int
	//
	// Some sort of particle type, for rendering triangles or squares or circles...
	Shape            ShapeFunction
	ProgressFunction func(x, y, w, h int) float64
}

func (gg *GradientGenerator) GenerateParticle(bp BaseParticle) Particle {
	return &GradientParticle{
		BaseParticle: bp,
		startColor:   randColor(gg.StartColor, gg.StartColorRand),
		endColor:     randColor(gg.EndColor, gg.EndColorRand),
		startColor2:  randColor(gg.StartColor2, gg.StartColor2Rand),
		endColor2:    randColor(gg.EndColor2, gg.EndColor2Rand),
		size:         gg.Size + intFromSpread(gg.SizeRand),
	}
}

func (gg *GradientGenerator) GetBaseGenerator() *BaseGenerator {
	return &gg.BaseGenerator
}

// Generate takes a generator and converts it into a source,
// drawing particles and binding functions for particle generation
// and rotation.
func (gg *GradientGenerator) Generate(layer int) *Source {

	// Convert rotation from degrees to radians
	gg.Rotation = gg.Rotation / 180 * math.Pi
	gg.RotationRand = gg.RotationRand / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: gg,
		particles: make([]Particle, 0),
	}

	// Bind things to that source:
	ps.Init()
	render.Draw(&ps, layer)

	return &ps
}
func (gp *GradientGenerator) GetParticleSize() (float64, float64, bool) {
	return 0, 0, true
}

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type GradientParticle struct {
	BaseParticle
	startColor  color.Color
	endColor    color.Color
	startColor2 color.Color
	endColor2   color.Color
	size        int
}

func (gp *GradientParticle) Draw(generator Generator, buff draw.Image) {
	gp.DrawOffset(generator, buff, 0, 0)
}

func (gp *GradientParticle) DrawOffset(generator Generator, buff draw.Image, xOff, yOff float64) {

	gen := generator.(*GradientGenerator)

	r, g, b, a := gp.startColor.RGBA()
	r2, g2, b2, a2 := gp.endColor.RGBA()
	progress := gp.life / gp.totalLife
	c1 := color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}
	r, g, b, a = gp.startColor2.RGBA()
	r2, g2, b2, a2 = gp.endColor2.RGBA()
	c2 := color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}
	r, g, b, a = c1.RGBA()
	r2, g2, b2, a2 = c2.RGBA()

	img := image.NewRGBA64(image.Rect(0, 0, gp.size, gp.size))

	for i := 0; i < gp.size; i++ {
		for j := 0; j < gp.size; j++ {
			if gen.Shape(i, j, gp.size) {
				progress := gen.ProgressFunction(i, j, gp.size, gp.size)
				c := color.RGBA64{
					uint16OnScale(r, r2, progress),
					uint16OnScale(g, g2, progress),
					uint16OnScale(b, b2, progress),
					uint16OnScale(a, a2, progress),
				}
				img.SetRGBA64(i, j, c)
			}
		}
	}

	halfSize := float64(gp.size / 2)

	render.ShinyDraw(buff, img, int((xOff+gp.x)-halfSize), int((yOff+gp.y)-halfSize))
}

func (gp *GradientParticle) GetBaseParticle() *BaseParticle {
	return &gp.BaseParticle
}

func (gp *GradientParticle) GetPos() (float64, float64) {
	fSize := float64(gp.size)
	return gp.x - fSize/2, gp.y - fSize/2
}
func (gp *GradientParticle) GetSize() (float64, float64) {
	fSize := float64(gp.size)
	return fSize, fSize
}
