package androiddriver

import (
	"image"
	"image/color"
	"sync"

	"github.com/oakmound/oak/v3/shiny/driver/internal/event"
	"github.com/oakmound/oak/v3/shiny/driver/internal/lifecycler"
	"github.com/oakmound/oak/v3/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var _ screen.Screen = &screenImpl{}

type screenImpl struct {
	event.Deque
	glctx        gl.Context
	images       *glutil.Images
	activeImages []*imageImpl
	activeImage  *imageImpl
	texture      struct {
		program gl.Program
		pos     gl.Attrib
		mvp     gl.Uniform
		uvp     gl.Uniform
		inUV    gl.Attrib
		sample  gl.Uniform
		quad    gl.Buffer
	}
	fill struct {
		program gl.Program
		pos     gl.Attrib
		mvp     gl.Uniform
		color   gl.Uniform
		quad    gl.Buffer
	}

	lifecycler lifecycler.State

	// glctxMu is a mutex that enforces the atomicity of methods like
	// Texture.Upload or Window.Draw that are conceptually one operation
	// but are implemented by multiple OpenGL calls. OpenGL is a stateful
	// API, so interleaving OpenGL calls from separate higher-level
	// operations causes inconsistencies.
	glctxMu sync.Mutex
	worker  gl.Worker

	// szMu protects only sz. If you need to hold both glctxMu and szMu, the
	// lock ordering is to lock glctxMu first (and unlock it last).
	szMu sync.Mutex
	sz   size.Event
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
	// gl.Flush is a lightweight (on modern GL drivers) blocking call
	// that ensures all GL functions pending in the gl package have
	// been passed onto the GL driver before the app package attempts
	// to swap the screen buffer.
	//
	// This enforces that the final receive (for this paint cycle) on
	// gl.WorkAvailable happens before the send on publish.
	w.glctxMu.Lock()
	w.glctx.Flush()
	w.glctxMu.Unlock()

	return screen.PublishResult{}
}

func (w *screenImpl) Release() {
	// TODO
}

func (w *screenImpl) Upload(dp image.Point, src screen.Image, sr image.Rectangle)                   {}
func (w *screenImpl) Fill(dr image.Rectangle, src color.Color, op draw.Op)                          {}
func (w *screenImpl) Draw(src2dst f64.Aff3, src screen.Texture, sr image.Rectangle, op draw.Op)     {}
func (w *screenImpl) DrawUniform(src2dst f64.Aff3, src color.Color, sr image.Rectangle, op draw.Op) {}
func (w *screenImpl) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op)       {}
func (w *screenImpl) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
	w.glctxMu.Lock()
	defer w.glctxMu.Unlock()

	t := src.(*textureImpl)
	w.activeImage = t.img
}
