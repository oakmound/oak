package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPoint(t *testing.T) {
	p := NewPoint(nil, 10, 10)
	assert.Equal(t, 10.0, p.X())
	assert.Equal(t, 10.0, p.Y())
}
