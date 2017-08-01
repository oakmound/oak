package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeScene(t *testing.T) {
	tree, err := NewTree(10, 20)
	assert.Nil(t, err)

	// Empty tree checks
	hits := tree.Hits(NewSpace(0, 0, 10, 10, 0))
	assert.Empty(t, hits)
	hit := tree.HitLabel(NewUnassignedSpace(0, 0, 10, 10))
	assert.Nil(t, hit)
	hits = tree.Hit(NewLabeledSpace(0, 0, 10, 10, 0))
	assert.Empty(t, hits)
}
