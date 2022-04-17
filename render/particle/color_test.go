package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/v3/alg/span"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/shape"
)

func TestColorParticle(t *testing.T) {
	g := NewColorGenerator(
		Rotation(span.NewConstant(1.0)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(span.NewConstant(5)),
		EndSize(span.NewConstant(10)),
		Shape(shape.Heart),
	)
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*ColorParticle)

	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)), 0, 0)
	if p.GetLayer() != 0 {
		t.Fatalf("expected 0 layer, got %v", p.GetLayer())
	}

	p.Life = -1
	sz, _ := p.GetDims()
	if sz != int(p.endSize) {
		t.Fatalf("expected size %v at end of particle's life, got %v", p.endSize, sz)
	}
	p.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)), 0, 0)

	var cp2 *ColorParticle
	if cp2.GetLayer() != render.Undraw {
		t.Fatalf("uninitialized particle was not set to the undraw layer")
	}

	_, _, ok := g.GetParticleSize()
	if !ok {
		t.Fatalf("get particle size not particle-specified")
	}
}
