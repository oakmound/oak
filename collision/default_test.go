package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFns(t *testing.T) {
	Clear()
	s := NewUnassignedSpace(0, 0, 10, 10)
	Add(s)
	Remove(s)
	assert.Empty(t, Hits(NewUnassignedSpace(1, 1, 1, 1)))

	Add(s)
	assert.Nil(t, ShiftSpace(3, 3, s))
	assert.Empty(t, Hits(NewUnassignedSpace(1, 1, 1, 1)))

	assert.Nil(t, s.Update(0, 0, 10, 10))
	assert.NotEmpty(t, Hits(NewUnassignedSpace(1, 1, 1, 1)))

	Clear()
	assert.Empty(t, Hits(NewUnassignedSpace(1, 1, 1, 1)))

	s = NewLabeledSpace(0, 0, 2, 2, Label(2))
	Add(s)
	assert.Empty(t, HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)))
	assert.Nil(t, s.SetDim(10, 10))
	assert.NotEmpty(t, HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)))
	s.UpdateLabel(Label(1))
	assert.Empty(t, HitLabel(NewUnassignedSpace(5, 5, 1, 1), Label(2)))

}
