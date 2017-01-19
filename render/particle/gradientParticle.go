package particle

import (
	"image"
	"image/color"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// A particle is a colored pixel at a given position, moving in a certain direction.
// After a while, it will dissipate.
type GradientParticle struct {
	*BaseParticle
	startColor  color.Color
	endColor    color.Color
	startColor2 color.Color
	endColor2   color.Color
	size        int
}

func (gp *GradientParticle) Draw(buff draw.Image) {
	gp.DrawOffset(buff, 0, 0)
}

func (gp *GradientParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	gp.DrawOffsetGen(gp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

func (gp *GradientParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {

	gen := generator.(*GradientGenerator)

	r, g, b, a := gp.startColor.RGBA()
	r2, g2, b2, a2 := gp.endColor.RGBA()
	progress := gp.Life / gp.totalLife
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

	render.ShinyDraw(buff, img, int((xOff+gp.Pos.X)-halfSize), int((yOff+gp.Pos.Y)-halfSize))
}

func (gp *GradientParticle) GetBaseParticle() *BaseParticle {
	return gp.BaseParticle
}

func (gp *GradientParticle) GetPos() *physics.Vector {
	fSize := float64(gp.size)
	return physics.NewVector(gp.Pos.X-fSize/2, gp.Pos.Y-fSize/2)
}
func (gp *GradientParticle) GetSize() (float64, float64) {
	fSize := float64(gp.size)
	return fSize, fSize
}
