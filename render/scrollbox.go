package render

import (
	"errors"
	"image/draw"
	"strconv"
	"time"

	"github.com/oakmound/oak/physics"
)

// A ScrollBox is a renderable that draws other renderables to itself in a scrolling fashion,
// for animating ticker tape feeds or rotating background animations.
type ScrollBox struct {
	*Sprite
	pauseBool
	Rs                       []Renderable
	nextScrollX, nextScrollY time.Time
	scrollRateX, scrollRateY time.Duration
	View, reappear           physics.Vector
	dirX, dirY               float64
}

// NewScrollBox returns a ScrollBox of the input renderables and the given dimensions.
// milliPerPixel represents the number of milliseconds it will take for the scroll box
// to move a horizontal or vertical pixel respectively. A negative value for milliPerPixel
// will move in a negative direction.
func NewScrollBox(rs []Renderable, milliPerPixelX, milliPerPixelY, width, height int) *ScrollBox {
	s := new(ScrollBox)
	s.pauseBool = pauseBool{playing: true}
	s.Rs = rs
	s.View = physics.NewVector(float64(width), float64(height))

	s.SetScrollRate(milliPerPixelX, milliPerPixelY)
	s.reappear = physics.NewVector(-1*s.dirX*s.View.X(), -1*s.dirY*s.View.Y())

	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
	s.Sprite = NewEmptySprite(0, 0, width, height)

	s.drawRenderables()
	return s
}

// DrawOffset draws this scroll box at +xOff, +yOff
func (s *ScrollBox) DrawOffset(buff draw.Image, xOff, yOff float64) {
	s.update()
	s.Sprite.DrawOffset(buff, xOff, yOff)

}

// Draw draws this scroll box to the input buffer
func (s *ScrollBox) Draw(buff draw.Image) {
	s.DrawOffset(buff, 0, 0)
}

func (s *ScrollBox) update() {
	if !s.playing {
		return
	}
	// ScrollBoxes update in a discontinuous fashion with Animation and Sequence
	// Both of the mentioned types will only ever advance one frame per update,
	// whereas ScrollBox will move however many pixels it should have moved in
	// the case of a long lag in draw calls.
	if s.dirX != 0 && time.Now().After(s.nextScrollX) {
		pixelsMovedX := int64(time.Since(s.nextScrollX))/int64(s.scrollRateX) + 1
		s.nextScrollX = time.Now().Add(s.scrollRateX)

		newS := NewEmptySprite(s.Sprite.X(), s.Sprite.Y(), int(s.View.X()), int(s.View.Y()))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftX(-1 * s.dirX * float64(pixelsMovedX))
			if s.shouldReappearX(m) {
				m.ShiftX(-1 * s.reappear.X()) //Hope that delta is not higher than reappear...
			}
		}
		*s.Sprite = *newS
		s.drawRenderables()
	}
	if s.dirY != 0 && time.Now().After(s.nextScrollY) {
		pixelsMovedY := int64(time.Since(s.nextScrollY))/int64(s.scrollRateY) + 1
		s.nextScrollY = time.Now().Add(s.scrollRateY)

		newS := NewEmptySprite(s.Sprite.X(), s.Sprite.Y(), int(s.View.X()), int(s.View.Y()))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftY(-1 * s.dirY * float64(pixelsMovedY))
			if s.shouldReappearY(m) {
				m.ShiftY(-1 * s.reappear.Y()) //Hope that delta is not higher than reappear...
			}
		}
		*s.Sprite = *newS
		s.drawRenderables()
	}
}

func (s *ScrollBox) shouldReappearY(m Renderable) bool {
	return (s.dirY == 1 && m.Y() <= s.reappear.Y()) || (s.dirY == -1 && m.Y() >= s.reappear.Y())
}

func (s *ScrollBox) shouldReappearX(m Renderable) bool {
	return (s.dirX == 1 && m.X() <= s.reappear.X()) || (s.dirX == -1 && m.X() >= s.reappear.X())
}

// Unpause resumes this scroll box's scrolling. Will delay the next scroll frame
// if already unpaused.
func (s *ScrollBox) Unpause() {
	s.pauseBool.Unpause()
	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
}

// SetReappearPos sets at what point renderables in this box should loop back on
// themselves to begin scrolling again
func (s *ScrollBox) SetReappearPos(x, y float64) error {
	s.reappear = physics.NewVector(x, y)
	if x*s.dirX > 0 {
		return errors.New("ScrollBox will not loop with direction.X: " + strconv.Itoa(int(s.dirX)) + " and reappear.X: " + strconv.Itoa(int(x)))
	}
	if y*s.dirY > 0 {
		return errors.New("ScrollBox will not loop with direction.Y: " + strconv.Itoa(int(s.dirY)) + " and reappear.X: " + strconv.Itoa(int(y)))
	}
	return nil
}

// SetScrollRate sets how fast this scroll box should rotate its x and y axes
// Maybe BUG, Consider: The next time that the box will scroll at is not updated
// immediately after this is called, only after the box is drawn.
func (s *ScrollBox) SetScrollRate(milliPerPixelX, milliPerPixelY int) {
	s.dirX = 1
	s.dirY = 1
	if milliPerPixelX < 0 {
		milliPerPixelX *= -1
		s.dirX = -1
	} else if milliPerPixelX == 0 {
		s.dirX = 0
	}
	if milliPerPixelY < 0 {
		milliPerPixelY *= -1
		s.dirY = -1
	} else if milliPerPixelY == 0 {
		s.dirY = 0
	}

	s.scrollRateX = time.Duration(milliPerPixelX) * time.Millisecond
	s.scrollRateY = time.Duration(milliPerPixelY) * time.Millisecond
}

// AddRenderable adds the inputs to this scrollbox.
func (s *ScrollBox) AddRenderable(rs ...Renderable) {
	for _, r := range rs {
		// We don't do this specific position swapping (which is to attempt
		// to do what we think users actually want) at initialization,
		// I suppose because we assume the inputs are at 0,0?
		switch r.(type) {
		case *Text:
			r.SetPos(r.X()*-1, r.Y()*-1)
		}
		s.Rs = append(s.Rs, r)
	}
	s.drawRenderables()
}

func (s *ScrollBox) drawRenderables() {
	for _, r := range s.Rs {
		// This might be the only place where we draw to a buffer that isn't
		// oak's main buffer.
		r.DrawOffset(s.GetRGBA(), -2*r.X(), -2*r.Y())
		if s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.X(), -2*r.Y()+s.reappear.Y())
		}
		if s.scrollRateX != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.X()+s.reappear.X(), -2*r.Y())
		}
		if s.scrollRateX != 0 && s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.X()+s.reappear.X(), -2*r.Y()+s.reappear.Y())
		}
	}
}
