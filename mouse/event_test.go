package mouse

import (
	"testing"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"
	"github.com/stretchr/testify/assert"
)

func TestEventConversions(t *testing.T) {
	me := Event{1.0, 1.0, "", ""}
	v := me.ToVector()
	assert.Equal(t, v, physics.NewVector(1.0, 1.0))
	s := me.ToSpace()
	Add(collision.NewUnassignedSpace(1.0, 1.0, .1, .1))
	assert.NotEmpty(t, Hits(s))
}
