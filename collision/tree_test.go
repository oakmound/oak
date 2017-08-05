package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeScene(t *testing.T) {
	tree, err := NewTree(20, 10)
	assert.Nil(t, tree)
	assert.NotNil(t, err)
	tree, err = NewTree(10, 20)
	assert.Nil(t, err)

	// Empty tree checks
	assert.Empty(t, tree.Hits(NewSpace(0, 0, 10, 10, 0)))
	assert.Nil(t, tree.HitLabel(NewUnassignedSpace(0, 0, 10, 10)))
	assert.Empty(t, tree.Hit(NewLabeledSpace(0, 0, 10, 10, 0), WithLabels(1)))

	// Positive hit checks
	s1 := NewFullSpace(0, 0, 10, 10, 1, 3)
	s2 := NewFullSpace(10, 10, 20, 20, 2, 4)
	tree.Add(s1, s2)
	assert.Equal(t, len(tree.Hits(NewSpace(5, 5, 1, 1, 0))), 1)
	assert.NotNil(t, tree.HitLabel(NewSpace(15, 15, 1, 1, 0), 2))
	// Self-hit should not happen
	assert.Empty(t, tree.Hits(s1))

	// Filters
	assert.NotEmpty(t, tree.Hit(NewSpace(0, 0, 100, 100, 0), WithoutLabels(2)))
	assert.NotEmpty(t, tree.Hit(NewSpace(0, 0, 100, 100, 0), WithLabels(2)))
	assert.NotEmpty(t, tree.Hit(NewSpace(0, 0, 100, 100, 0), WithoutCIDs(3)))
	assert.NotEmpty(t, tree.Hit(NewSpace(0, 0, 100, 100, 0), FirstLabel(1)))
	assert.Empty(t, tree.Hit(NewSpace(0, 0, 100, 100, 0), FirstLabel(5)))

	// Update functions
	assert.NotNil(t, tree.ShiftSpace(0, 0, nil))
	assert.NotNil(t, tree.UpdateSpace(0, 0, 0, 0, nil))

	assert.Nil(t, tree.ShiftSpace(1, 1, s1))
	assert.Empty(t, tree.Hits(NewSpace(0, 0, 1, 1, 0)))
	assert.Nil(t, tree.UpdateSpace(20, 20, 5, 5, s2))
	assert.NotNil(t, tree.HitLabel(NewSpace(21, 21, 20, 20, 0), 2))

	// Removal, Clear
	assert.Equal(t, 1, tree.Remove(s2))
	assert.Nil(t, tree.HitLabel(NewSpace(21, 21, 20, 20, 0), 2))

	tree.Clear()

	assert.Empty(t, tree.Hits(NewSpace(0, 0, 100, 100, 0)))
}
