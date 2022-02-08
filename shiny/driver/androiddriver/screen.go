//go:build android
// +build android

package androiddriver

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

var _ screen.Screen = &screenImpl{}

type screenImpl struct {
	event.Deque

	app   app.App
	glctx gl.Context

	images       *glutil.Images
	activeImages []*imageImpl

	lastSz size.Event
}

func (s *screenImpl) NewImage(size image.Point) (screen.Image, error) {
	img := &imageImpl{
		screen: s,
		size:   size,
		img:    s.images.NewImage(size.X, size.Y),
	}
	s.activeImages = append(s.activeImages, img)
	return img, nil
}

func (s *screenImpl) NewTexture(size image.Point) (screen.Texture, error) {
	return NewTexture(s, size), nil
}

var _ screen.Window = &screenImpl{}

func (s *screenImpl) NewWindow(opts screen.WindowGenerator) (screen.Window, error) {
	// android does not support multiple windows
	return s, nil
}

func (w *screenImpl) Publish() screen.PublishResult {
	return screen.PublishResult{}
}

func (w *screenImpl) Release()                                                                      {}
func (w *screenImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle)                   {}
func (w *screenImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)                          {}
func (w *screenImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op)     {}
func (w *screenImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {}
func (w *screenImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op)       {}
func (w *screenImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	t := src.(*textureImpl)
	t.img.img.Draw(
		w.lastSz,
		geom.Point{},
		geom.Point{X: w.lastSz.WidthPt},
		geom.Point{Y: w.lastSz.HeightPt},
		t.img.Bounds(),
	)
	t.img.img.Upload()
	w.app.Publish()
}
