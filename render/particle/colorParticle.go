package particle

import (
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

// A ColorParticle is a particle with a defined color and size
type ColorParticle struct {
	*baseParticle
	startColor color.Color
	endColor   color.Color
	size       float64
	endSize    float64
}

// Draw redirects to DrawOffset
func (cp *ColorParticle) Draw(buff draw.Image) {
	cp.DrawOffset(buff, 0, 0)
}

// DrawOffset redirects to DrawOffsetGen
func (cp *ColorParticle) DrawOffset(buff draw.Image, xOff, yOff float64) {
	cp.DrawOffsetGen(cp.GetBaseParticle().Src.Generator, buff, xOff, yOff)
}

// DrawOffsetGen draws a particle with it's generator's variables
func (cp *ColorParticle) DrawOffsetGen(generator Generator, buff draw.Image, xOff, yOff float64) {
	gen := generator.(*ColorGenerator)

	// Hmm. this is expensive.
	// This work should be done by the Source because if the draw rate is faster
	// than the enter frame rate than this is doing duplicate work
	// does that mean every particle is the same struct (
	//	baseParticle + image
	//)
	// and different particle types are just different update functions?
	// -No- because we still need to keep track of variable things on these particles
	// but it -does- mean that particles should track an image that they all have a function
	// to create instead of these Draw functions which should just be provided by
	// baseParticle

	r, g, b, a := cp.startColor.RGBA()
	r2, g2, b2, a2 := cp.endColor.RGBA()
	progress := cp.Life / cp.totalLife
	if progress < 0 {
		progress = 0
	}

	size := int(((progress) * cp.size) + ((1 - progress) * cp.endSize))

	c := color.RGBA64{
		uint16OnScale(r, r2, progress),
		uint16OnScale(g, g2, progress),
		uint16OnScale(b, b2, progress),
		uint16OnScale(a, a2, progress),
	}

	halfSize := float64(size) / 2

	xOffi := int((xOff - halfSize) + cp.X())
	yOffi := int((yOff - halfSize) + cp.Y())

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if gen.Shape.In(i, j, size) {
				//fmt.Println("In", i, j)
				//render.ShinySet(buff, c, xOffi+i, yOffi+j)
				buff.Set(xOffi+i, yOffi+j, c)
			}
		}
	}

	//fmt.Println("Drawing particle", size, xOffi, yOffi)
}

// GetLayer returns baseParticle GetLayer. This is a safety check against auto-generated
// code which would not contain the nil check here
func (cp *ColorParticle) GetLayer() int {
	if cp == nil {
		return render.Undraw
	}
	return cp.baseParticle.GetLayer()
}

// GetPos returns the middle of a color particle
func (cp *ColorParticle) GetPos() physics.Vector {
	return physics.NewVector(cp.X()-cp.size/2, cp.Y()-cp.size/2)
}

// GetDims returns the color particle's size, twice
func (cp *ColorParticle) GetDims() (int, int) {
	progress := cp.Life / cp.totalLife
	if progress < 0 {
		progress = 0
	}

	size := int(((progress) * cp.size) + ((1 - progress) * cp.endSize))
	//fmt.Println("Dim size", size)
	return size, size
}
