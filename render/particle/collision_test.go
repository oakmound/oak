package particle

import (
	"image"
	"testing"

	"github.com/oakmound/oak/v2/collision"
)

func TestCollisionParticle(t *testing.T) {
	hm := map[collision.Label]collision.OnHit{
		1: func(_, _ *collision.Space) {},
	}
	g := NewCollisionGenerator(NewColorGenerator(), HitMap(hm), Fragile(true)).(*CollisionGenerator)
	src := g.Generate(0)
	src.addParticles()
	cp := src.particles[0].(*CollisionParticle)
	w, h := cp.GetDims()
	if w != 1 {
		t.Fatalf("expected 1 width, got %v", w)
	}
	if h != 1 {
		t.Fatalf("expected 1 height, got %v", h)
	}
	cp.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))
	cp.Cycle(g)
	collision.Add(collision.NewLabeledSpace(-20, -20, 40, 40, 1))
	cp.Cycle(g)

	_, _, ok := g.GetParticleSize()
	if !ok {
		t.Fatalf("get particle size not particle-specified")
	}
}
