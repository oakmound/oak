package render

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogicFPS(t *testing.T) {

	initTestFont()

	lfps := NewLogicFPS()

	lfps.PreDraw()

	assert.Nil(t, lfps.Add(nil))

	lfps.Replace(nil, nil, 0)
	assert.NotNil(t, lfps.Copy())
	lfps.draw(image.NewRGBA(image.Rect(0, 0, 100, 100)), image.Point{0, 0}, 10, 10)

	logicFPSBind(int(lfps.CID), nil)
}
