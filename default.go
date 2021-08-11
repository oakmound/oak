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

var defaultWindow *Window

var initDefaultWindowOnce sync.Once

func initDefaultWindow() {
	initDefaultWindowOnce.Do(func() {
		defaultWindow = NewWindow()
	})
}

func Init(scene string, configOptions ...ConfigOption) error {
	initDefaultWindow()
	defaultWindow.DrawStack = render.GlobalDrawStack
	defaultWindow.logicHandler = event.DefaultBus
	return defaultWindow.Init(scene, configOptions...)
}

func AddScene(name string, sc scene.Scene) error {
	initDefaultWindow()
	return defaultWindow.AddScene(name, sc)
}

func IsDown(key string) bool {
	initDefaultWindow()
	return defaultWindow.IsDown(key)
}

func IsHeld(key string) (bool, time.Duration) {
	initDefaultWindow()
	return defaultWindow.IsHeld(key)
}

func SetUp(key string) {
	initDefaultWindow()
	defaultWindow.SetUp(key)
}

func SetDown(key string) {
	initDefaultWindow()
	defaultWindow.SetDown(key)
}

func SetViewportBounds(rect intgeom.Rect2) {
	initDefaultWindow()
	defaultWindow.SetViewportBounds(rect)
}

func ShiftScreen(x, y int) {
	initDefaultWindow()
	defaultWindow.ShiftScreen(x, y)
}

func SetScreen(x, y int) {
	initDefaultWindow()
	defaultWindow.SetScreen(x, y)
}

func MoveWindow(x, y, w, h int) error {
	initDefaultWindow()
	return defaultWindow.MoveWindow(x, y, w, h)
}

func UpdateViewSize(w, h int) error {
	initDefaultWindow()
	return defaultWindow.UpdateViewSize(w, h)
}

func SetFullScreen(fs bool) error {
	initDefaultWindow()
	return defaultWindow.SetFullScreen(fs)
}

func SetBorderless(bs bool) error {
	initDefaultWindow()
	return defaultWindow.SetBorderless(bs)
}

func SetTopMost(on bool) error {
	initDefaultWindow()
	return defaultWindow.SetTopMost(on)
}

func SetTitle(title string) error {
	initDefaultWindow()
	return defaultWindow.SetTitle(title)
}

func SetTrayIcon(icon string) error {
	initDefaultWindow()
	return defaultWindow.SetTrayIcon(icon)
}

func ShowNotification(title, msg string, icon bool) error {
	initDefaultWindow()
	return defaultWindow.ShowNotification(title, msg, icon)
}

func ScreenShot() *image.RGBA {
	initDefaultWindow()
	return defaultWindow.ScreenShot()
}

func SetLoadingRenderable(r render.Renderable) {
	initDefaultWindow()
	defaultWindow.SetLoadingRenderable(r)
}

func SetBackground(b Background) {
	initDefaultWindow()
	defaultWindow.SetBackground(b)
}

func SetColorBackground(img image.Image) {
	initDefaultWindow()
	defaultWindow.SetColorBackground(img)
}

func GetBackgroundImage() image.Image {
	initDefaultWindow()
	return defaultWindow.GetBackgroundImage()
}

func Width() int {
	initDefaultWindow()
	return defaultWindow.Width()
}

func Height() int {
	initDefaultWindow()
	return defaultWindow.Height()
}

func HideCursor() error {
	initDefaultWindow()
	return defaultWindow.HideCursor()
}

func GetCursorPosition() (x, y float64, err error) {
	initDefaultWindow()
	return defaultWindow.GetCursorPosition()
}
