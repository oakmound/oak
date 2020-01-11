package mouse

import (
	"testing"

	"github.com/oakmound/oak/v2/collision"
	"github.com/stretchr/testify/assert"
)

func TestEventConversions(t *testing.T) {
	me := NewZeroEvent(1.0, 1.0)
	s := me.ToSpace()
	Add(collision.NewUnassignedSpace(1.0, 1.0, .1, .1))
	assert.NotEmpty(t, Hits(s))
}
