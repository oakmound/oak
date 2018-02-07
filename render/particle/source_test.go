package particle

import (
	"image/color"
	"testing"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/shape"
	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	g := NewGradientGenerator(
		Rotation(floatrange.Constant(1)),
		Color(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Color2(color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 0, 0, 255}),
		Size(intrange.Constant(5)),
		EndSize(intrange.Constant(10)),
		Shape(shape.Heart),
		Progress(render.HorizontalProgress),
		And(
			NewPerFrame(floatrange.Constant(200)),
		),
		Pos(20, 20),
		LifeSpan(floatrange.Constant(10)),
		Angle(floatrange.Constant(0)),
		Speed(floatrange.Constant(0)),
		Spread(10, 10),
		Duration(intrange.Constant(10)),
		Gravity(10, 10),
		SpeedDecay(1, 1),
		End(func(_ Particle) {}),
		Layer(func(_ physics.Vector) int { return 0 }),
	)

	src := g.Generate(0)

	ended := false

	src.EndFunc = func() {
		ended = true
	}

	for i := 0; i < 1000; i++ {
		rotateParticles(int(src.CID), nil)
	}
	for clearParticles(int(src.CID), nil) != event.UnbindEvent {
	}

	assert.True(t, ended)

	src.Pause()
	assert.True(t, src.paused)
	src.UnPause()
	assert.False(t, src.paused)
	x, y := src.Generator.GetPos()
	src.ShiftX(10)
	src.ShiftY(10)
	x2, y2 := src.Generator.GetPos()
	assert.Equal(t, x+10, x2)
	assert.Equal(t, y+10, y2)
	src.SetPos(-20, -30)
	x2, y2 = src.Generator.GetPos()
	assert.Equal(t, -20.0, x2)
	assert.Equal(t, -30.0, y2)

	var src2 *Source
	src2.Stop()
}
