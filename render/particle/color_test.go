package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/shape"
	"github.com/stretchr/testify/assert"
)

func TestColorParticle(t *testing.T) {
	g := NewColorGenerator(
		Rotation(floatrange.Constant(1)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(intrange.Constant(5)),
		EndSize(intrange.Constant(10)),
		Shape(shape.Heart),
	)
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*ColorParticle)

	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))
	assert.Equal(t, 0, p.GetLayer())

	p.Life = -1
	sz, _ := p.GetDims()
	assert.Equal(t, float64(sz), p.endSize)
	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))

	var cp2 *ColorParticle
	assert.Equal(t, render.Undraw, cp2.GetLayer())

	_, _, ok := g.GetParticleSize()
	assert.True(t, ok)
}
