package intgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtDistance(t *testing.T) {
	a := Point{0, 0}
	b := Point{3, 0}
	assert.Equal(t, a.Distance(b), b.Distance(a))
	assert.Equal(t, a.Distance(b), 3.0)
}
