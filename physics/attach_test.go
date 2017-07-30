package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttach(t *testing.T) {
	// Attach
	v := NewVector(0, 0)
	v2 := NewVector(0, 0)
	v3 := NewVector(0, 0)
	v2 = v2.Attach(v)
	v3 = v3.Attach(v, 10, 10)
	v.SetPos(10, 10)
	assert.Equal(t, v.X(), v2.X())
	assert.Equal(t, v.X(), v3.X()-10)
	assert.Equal(t, v.Y(), v2.Y())
	assert.Equal(t, v.Y(), v3.Y()-10)

	// AttachX, AttachY
	v4 := NewVector(0, 0)
	v4 = v4.AttachX(v, 5)
	assert.Equal(t, v.X(), v4.X()-5)
	assert.NotEqual(t, v.Y(), v4.Y())

	v5 := NewVector(0, 0)
	v5 = v5.AttachY(v, 5)
	assert.Equal(t, v.Y(), v5.Y()-5)
	assert.NotEqual(t, v.X(), v5.X())

	// Detach
	v6 := v5.Detach()
	v.SetPos(100, 100)
	assert.NotEqual(t, v.Y(), v6.Y()-5)
}
