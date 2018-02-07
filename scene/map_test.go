package scene

import (
	"image"
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
	m.CurrentScene = "test"
	_, ok = m.GetCurrent()
	assert.True(t, ok)
}

func TestTransition(t *testing.T) {
	fadeFn := Fade(1, 10)
	assert.False(t, fadeFn(nil, 11))
	assert.True(t, fadeFn(image.NewRGBA(image.Rect(0, 0, 50, 50)), 2))
	zoomFn := Zoom(.5, .5, 10, .1)
	assert.False(t, zoomFn(nil, 11))
	assert.True(t, zoomFn(image.NewRGBA(image.Rect(0, 0, 50, 50)), 2))
}
