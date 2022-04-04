package oak

import (
	"image"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
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

// ShiftScreen calls ShiftScreen on the default window.
func ShiftScreen(x, y int) {
	initDefaultWindow()
	defaultWindow.ShiftScreen(x, y)
}

// SetScreen calls SetScreen on the default window.
func SetScreen(x, y int) {
	initDefaultWindow()
	defaultWindow.SetScreen(x, y)
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

// GetBackgroundImage calls GetBackgroundImage on the default window.
func GetBackgroundImage() image.Image {
	initDefaultWindow()
	return defaultWindow.GetBackgroundImage()
}

// Width calls Width on the default window.
func Width() int {
	initDefaultWindow()
	return defaultWindow.Width()
}

// Height calls Height on the default window.
func Height() int {
	initDefaultWindow()
	return defaultWindow.Height()
}
