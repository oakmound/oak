package scene

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := NewMap()
	_, ok := m.Get("badScene")
	assert.False(t, ok)
	assert.Nil(t, m.Add("test", nil, nil, nil))
	assert.NotNil(t, m.Add("test", nil, nil, nil))
	_, ok = m.Get("test")
	assert.True(t, ok)
}
