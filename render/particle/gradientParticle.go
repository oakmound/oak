package particle

import (
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/render"
)

// A GradientParticle has a gradient from one color to another
type GradientParticle struct {
	ColorParticle
	startColor2 color.Color
	endColor2   color.Color
}

// Draw redirects to DrawOffset
func (gp *GradientParticle) Draw(buff draw.Image) {
	gp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirects to DrawOffsetGen
func (gp *GradientParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	gp.DrawOffsetGen(gp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (gp *GradientParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {

	gen := generator.(*GradientGenerator)
	progress := gp.Life / gp.totalLife
	c1 := render.GradientColorAt(gp.startColor, gp.endColor, progress)
	c2 := render.GradientColorAt(gp.startColor2, gp.endColor2, progress)

	size := int(((1 - progress) * gp.size) + (progress * gp.endSize))

	halfSize := float64(size) / 2

	xOffi := int(xOff - halfSize)
	yOffi := int(yOff - halfSize)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if gen.Shape.In(i, j, size) {
				progress := gen.ProgressFunction(i, j, size, size)
				c := render.GradientColorAt(c1, c2, progress)
				render.ShinySet(buff, c, xOffi+i, yOffi+j)
			}
		}
	}
}
