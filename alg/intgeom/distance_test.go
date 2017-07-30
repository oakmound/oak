package intgeom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtDistance(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(3, 0)
	assert.Equal(t, a.Distance(b), b.Distance(a))
	assert.Equal(t, a.Distance(b), 3.0)

	c := a.Add(b)

	assert.Equal(t, c.Distance(b), 0.0)
}
