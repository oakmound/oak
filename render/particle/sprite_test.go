package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/stretchr/testify/assert"

	"github.com/oakmound/oak/render"
)

func TestSpriteParticle(t *testing.T) {
	s := render.NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	g := NewSpriteGenerator(Sprite(s), Rotation(floatrange.Constant(1)), SpriteRotation(floatrange.Constant(1)))
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*SpriteParticle)

	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))

	w, h, ok := g.GetParticleSize()
	assert.Equal(t, 10.0, w)
	assert.Equal(t, 10.0, h)
	assert.False(t, ok)
}
