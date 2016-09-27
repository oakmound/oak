package particle

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/render"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type ColorGenerator struct {
	BaseGenerator
	StartColor, StartColorRand color.Color
	EndColor, EndColorRand     color.Color
	// The size, in pixel radius, of spawned particles
	Size, SizeRand int
	//
	// Some sort of particle type, for rendering triangles or squares or circles...
	Shape ShapeFunction
}

func (cg *ColorGenerator) GenerateParticle(bp BaseParticle) Particle {
	return &ColorParticle{
		BaseParticle: bp,
		startColor:   randColor(cg.StartColor, cg.StartColorRand),
		endColor:     randColor(cg.EndColor, cg.EndColorRand),
		size:         cg.Size + intFromSpread(cg.SizeRand),
	}
}

func (cg *ColorGenerator) GetBaseGenerator() *BaseGenerator {
	return &cg.BaseGenerator
}

// Generate takes a generator and converts it into a source,
// drawing particles and binding functions for particle generation
// and rotation.
func (cg ColorGenerator) Generate(layer int) *Source {

	// Convert rotation from degrees to radians
	cg.Rotation = cg.Rotation / 180 * math.Pi
	cg.RotationRand = cg.Rotation / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: &cg,
		particles: make([]Particle, 0),
	}

	// Bind things to that source:
	ps.Init()
	render.Draw(&ps, layer)

	return &ps
}

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type ColorParticle struct {
	BaseParticle
	startColor color.Color
	endColor   color.Color
	size       int
}

func (cp *ColorParticle) Draw(generator Generator, buff draw.Image) {

	gen := generator.(*ColorGenerator)

	r, g, b, a := cp.startColor.RGBA()
	r2, g2, b2, a2 := cp.endColor.RGBA()
	progress := cp.life / cp.totalLife
	c := color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}

	img := image.NewRGBA64(image.Rect(0, 0, cp.size, cp.size))

	for i := 0; i < cp.size; i++ {
		for j := 0; j < cp.size; j++ {
			if gen.Shape(i, j, cp.size) {
				img.SetRGBA64(i, j, c)
			}
		}
	}

	halfSize := float64(cp.size / 2)

	render.ShinyDraw(buff, img, int(cp.x-halfSize), int(cp.y-halfSize))
}

func (cp *ColorParticle) GetBaseParticle() *BaseParticle {
	return &cp.BaseParticle
}
