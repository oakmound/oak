package oak

import (
	"image"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var defaultController *Window

var initDefaultControllerOnce sync.Once

func initDefaultController() {
	initDefaultControllerOnce.Do(func() {
		defaultController = NewWindow()
	})
}

func Init(scene string, configOptions ...ConfigOption) error {
	initDefaultController()
	defaultController.DrawStack = render.GlobalDrawStack
	defaultController.logicHandler = event.DefaultBus
	return defaultController.Init(scene, configOptions...)
}

func AddScene(name string, sc scene.Scene) error {
	initDefaultController()
	return defaultController.AddScene(name, sc)
}

func IsDown(key string) bool {
	initDefaultController()
	return defaultController.IsDown(key)
}

func IsHeld(key string) (bool, time.Duration) {
	initDefaultController()
	return defaultController.IsHeld(key)
}

func SetUp(key string) {
	initDefaultController()
	defaultController.SetUp(key)
}

func SetDown(key string) {
	initDefaultController()
	defaultController.SetDown(key)
}

func SetViewportBounds(rect intgeom.Rect2) {
	initDefaultController()
	defaultController.SetViewportBounds(rect)
}

func ShiftScreen(x, y int) {
	initDefaultController()
	defaultController.ShiftScreen(x, y)
}

func SetScreen(x, y int) {
	initDefaultController()
	defaultController.SetScreen(x, y)
}

func SetFullScreen(fs bool) error {
	initDefaultController()
	return defaultController.SetFullScreen(fs)
}

func SetBorderless(bs bool) error {
	initDefaultController()
	return defaultController.SetBorderless(bs)
}

func ScreenShot() *image.RGBA {
	initDefaultController()
	return defaultController.ScreenShot()
}

func SetLoadingRenderable(r render.Renderable) {
	initDefaultController()
	defaultController.SetLoadingRenderable(r)
}

func SetBackground(b Background) {
	initDefaultController()
	defaultController.SetBackground(b)
}

func SetColorBackground(img image.Image) {
	initDefaultController()
	defaultController.SetColorBackground(img)
}

func GetBackgroundImage() image.Image {
	initDefaultController()
	return defaultController.GetBackgroundImage()
}

func Width() int {
	initDefaultController()
	return defaultController.Width()
}

func Height() int {
	initDefaultController()
	return defaultController.Height()
}

func HideCursor() error {
	initDefaultController()
	return defaultController.HideCursor()
}

func GetCursorPosition() (x, y float64, err error) {
	return defaultController.GetCursorPosition()
}
