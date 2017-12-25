package particle

import (
	"image"
	"testing"

	"github.com/oakmound/oak/collision"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 1, w)
	assert.Equal(t, 1, h)
	cp.Draw(image.NewRGBA(image.Rect(0, 0, 20, 20)))
	cp.Cycle(g)
	collision.Add(collision.NewLabeledSpace(-20, -20, 40, 40, 1))
	cp.Cycle(g)

	_, _, ok := g.GetParticleSize()
	assert.True(t, ok)
}
