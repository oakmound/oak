package oak

import (
	"testing"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

func TestDefaultFunctions(t *testing.T) {
	t.Run("SuperficialCoverage", func(t *testing.T) {
		IsDown(key.A)
		IsHeld(key.A)
		AddScene("test", scene.Scene{
			Start: func(ctx *scene.Context) {
				ScreenShot()
				ctx.Window.Quit()
			},
		})
		SetViewportBounds(intgeom.NewRect2(0, 0, 1, 1))
		SetViewport(intgeom.Point2{})
		ShiftViewport(intgeom.Point2{})
		UpdateViewSize(10, 10)
		Bounds()
		SetLoadingRenderable(render.EmptyRenderable())
		SetColorBackground(nil)
		SetBackground(render.EmptyRenderable())
		Init("test")
	})
}
