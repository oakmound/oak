package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/shape"
	"github.com/stretchr/testify/assert"
)

func TestGradientParticle(t *testing.T) {
	g := NewGradientGenerator(
		Rotation(floatrange.Constant(1)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Color2(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(intrange.Constant(5)),
		EndSize(intrange.Constant(10)),
		Shape(shape.Heart),
		Progress(render.HorizontalProgress),
		And(
			NewPerFrame(floatrange.Constant(20)),
		),
		Pos(20, 20),
		LifeSpan(floatrange.Constant(10)),
		Angle(floatrange.Constant(0)),
		Speed(floatrange.Constant(0)),
		Spread(10, 10),
		Duration(intrange.Constant(10)),
		Gravity(10, 10),
		SpeedDecay(1, 1),
		End(func(_ Particle) {}),
		Layer(func(_ physics.Vector) int { return 0 }),
	)
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*GradientParticle)

	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))
	assert.Equal(t, 0, p.GetLayer())

	p.Life = -1
	sz, _ := p.GetDims()
	assert.Equal(t, float64(sz), p.endSize)
	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))

	_, _, ok := g.GetParticleSize()
	assert.True(t, ok)
}
