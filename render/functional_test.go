package render

import (
	"image/color"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPos(t *testing.T) {
	spr := NewColorBox(10,10,color.Black)
	fnc := Functional{F:func () Renderable {return spr}}

	fnc.ShiftX(6)
	asrt := assert.New(t)
	asrt.Equal(float64(6),spr.X(),"wrong X value")
	w, h := fnc.GetDims()
	asrt.Equal(10,w,"wrong width")
	asrt.Equal(10,h,"wrong height")

	fnc.SetLayer(2)
	asrt.Equal(2,fnc.GetLayer(),"wrong layer")
}
