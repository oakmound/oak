package particle

import (
	"image"
	"image/color"
	"testing"

	"github.com/oakmound/oak/v4/alg/span"
	"github.com/oakmound/oak/v4/physics"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/shape"
)

func TestGradientParticle(t *testing.T) {
	g := NewGradientGenerator(
		Rotation(span.NewConstant(1.0)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Color2(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(span.NewConstant(5)),
		EndSize(span.NewConstant(10)),
		Shape(shape.Heart),
		Progress(render.HorizontalProgress),
		And(
			NewPerFrame(span.NewConstant(20.0)),
		),
		Pos(20, 20),
		LifeSpan(span.NewConstant(10.0)),
		Angle(span.NewConstant(0.0)),
		Speed(span.NewConstant(0.0)),
		Spread(10, 10),
		Duration(span.NewConstant(10)),
		Gravity(10, 10),
		SpeedDecay(1, 1),
		End(func(_ Particle) {}),
		Layer(func(_ physics.Vector) int { return 0 }),
	)
	src := g.Generate(0)
	src.addParticles()

	p := src.particles[0].(*GradientParticle)

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

	_, _, ok := g.GetParticleSize()
	if !ok {
		t.Fatalf("get particle size not particle-specified")
	}
}
