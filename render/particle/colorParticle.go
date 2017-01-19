package particle

import (
	"image"
	"image/color"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
	"bitbucket.org/oakmoundstudio/oak/render"
)

// Color Particles are particles with a defined color gradient and size
type ColorParticle struct {
	*BaseParticle
	startColor color.Color
	endColor   color.Color
	size       int
}

func (cp *ColorParticle) Draw(buff draw.Image) {
	cp.DrawOffset(buff, 0, 0)
}

func (cp *ColorParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

func (cp *ColorParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*ColorGenerator)

	r, g, b, a := cp.startColor.RGBA()
	r2, g2, b2, a2 := cp.endColor.RGBA()
	progress := cp.Life / cp.totalLife
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

	render.ShinyDraw(buff, img, int((xOff+cp.Pos.X)-halfSize), int((yOff+cp.Pos.Y)-halfSize))
}

func (cp *ColorParticle) GetBaseParticle() *BaseParticle {
	return cp.BaseParticle
}

func (cp *ColorParticle) GetPos() *physics.Vector {
	fSize := float64(cp.size)
	return physics.NewVector(cp.Pos.X-fSize/2, cp.Pos.Y-fSize/2)
}
func (cp *ColorParticle) GetSize() (float64, float64) {
	fSize := float64(cp.size)
	return fSize, fSize
}
