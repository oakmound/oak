package oak

import (
	"image"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/intgeom"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var defaultWindow *Window

var initDefaultWindowOnce sync.Once

func initDefaultWindow() {
	initDefaultWindowOnce.Do(func() {
		defaultWindow = NewWindow()
	})
}

// Init calls Init on the default window. The default window
// will be set to use render.GlobalDrawStack and event.DefaultBus.
func Init(scene string, configOptions ...ConfigOption) error {
	initDefaultWindow()
	defaultWindow.DrawStack = render.GlobalDrawStack
	defaultWindow.eventHandler = event.DefaultBus
	return defaultWindow.Init(scene, configOptions...)
}

// AddScene calls AddScene on the default window.
func AddScene(name string, sc scene.Scene) error {
	initDefaultWindow()
	return defaultWindow.AddScene(name, sc)
}

// IsDown calls IsDown on the default window.
func IsDown(k key.Code) bool {
	initDefaultWindow()
	return defaultWindow.IsDown(k)
}

// IsHeld calls IsHeld on the default window.
func IsHeld(k key.Code) (bool, time.Duration) {
	initDefaultWindow()
	return defaultWindow.IsHeld(k)
}

// SetViewportBounds calls SetViewportBounds on the default window.
func SetViewportBounds(rect intgeom.Rect2) {
	initDefaultWindow()
	defaultWindow.SetViewportBounds(rect)
}

// ShiftViewport calls ShiftViewport on the default window.
func ShiftViewport(pt intgeom.Point2) {
	initDefaultWindow()
	defaultWindow.ShiftViewport(pt)
}

// SetViewport calls SetViewport on the default window.
func SetViewport(pt intgeom.Point2) {
	initDefaultWindow()
	defaultWindow.SetViewport(pt)
}

// UpdateViewSize calls UpdateViewSize on the default window.
func UpdateViewSize(w, h int) error {
	initDefaultWindow()
	return defaultWindow.UpdateViewSize(w, h)
}

// ScreenShot calls ScreenShot on the default window.
func ScreenShot() *image.RGBA {
	initDefaultWindow()
	return defaultWindow.ScreenShot()
}

// SetLoadingRenderable calls SetLoadingRenderable on the default window.
func SetLoadingRenderable(r render.Renderable) {
	initDefaultWindow()
	defaultWindow.SetLoadingRenderable(r)
}

// SetBackground calls SetBackground on the default window.
func SetBackground(b Background) {
	initDefaultWindow()
	defaultWindow.SetBackground(b)
}

// SetColorBackground calls SetColorBackground on the default window.
func SetColorBackground(img image.Image) {
	initDefaultWindow()
	defaultWindow.SetColorBackground(img)
}

// Bounds returns the default window's boundary.
func Bounds() intgeom.Point2 {
	initDefaultWindow()
	return defaultWindow.Bounds()
}
