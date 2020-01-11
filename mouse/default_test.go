package mouse

import (
	"testing"

	"github.com/oakmound/oak/v2/collision"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFunctions(t *testing.T) {
	Clear()
	s := collision.NewUnassignedSpace(0, 0, 10, 10)
	Add(s)
	Remove(s)
	assert.Empty(t, Hits(collision.NewUnassignedSpace(1, 1, 1, 1)))

	Add(s)
	assert.Nil(t, ShiftSpace(3, 3, s))
	assert.Empty(t, Hits(collision.NewUnassignedSpace(1, 1, 1, 1)))

	assert.Nil(t, UpdateSpace(0, 0, 10, 10, s))
	assert.NotEmpty(t, Hits(collision.NewUnassignedSpace(1, 1, 1, 1)))

	Clear()
	assert.Empty(t, Hits(collision.NewUnassignedSpace(1, 1, 1, 1)))

	s = collision.NewLabeledSpace(0, 0, 10, 10, collision.Label(2))
	Add(s)
	assert.NotEmpty(t, HitLabel(collision.NewUnassignedSpace(1, 1, 1, 1), collision.Label(2)))
}
