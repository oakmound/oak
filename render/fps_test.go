package render

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawFPS(t *testing.T) {

	initTestFont()

	dfps := NewDrawFPS()

	dfps.PreDraw()

	assert.Nil(t, dfps.Add(nil))

	dfps.Replace(nil, nil, 0)
	assert.NotNil(t, dfps.Copy())
	dfps.draw(image.NewRGBA(image.Rect(0, 0, 100, 100)), image.Point{0, 0}, 10, 10)
}
