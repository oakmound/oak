package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/v3/alg/span"
	"github.com/oakmound/oak/v3/render"
)

func TestSpriteParticle(t *testing.T) {
	s := render.NewColorBox(10, 10, color.RGBA{255, 0, 0, 255})
	g := NewSpriteGenerator(
		Sprite(s),
		Rotation(span.NewConstant(1.0)),
		SpriteRotation(span.NewConstant(1.0)),
	)
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*SpriteParticle)

	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)), 0, 0)

	w, h, ok := g.GetParticleSize()
	if w != 10 {
		t.Fatalf("expected 10 x, got %v", w)
	}
	if h != 10 {
		t.Fatalf("expected 10 y, got %v", h)
	}
	if ok {
		t.Fatalf("sprite particle generator should not have particle-specific sizes")
	}
}
