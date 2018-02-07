package particle

import (
	"testing"

	"github.com/oakmound/oak/render"
	"github.com/stretchr/testify/assert"
)

func TestParticle(t *testing.T) {
	var bp *baseParticle
	w, h := bp.GetDims()
	assert.Equal(t, 0, w)
	assert.Equal(t, 0, h)

	assert.Equal(t, render.Undraw, bp.GetLayer())

	bp = new(baseParticle)
	bp.setPID(100)
	assert.Equal(t, 100, bp.pID)
}
