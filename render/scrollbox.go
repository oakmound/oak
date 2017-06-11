package render

import (
	"errors"
	"image/draw"
	"strconv"
	"time"

	"bitbucket.org/oakmoundstudio/oak/physics"
)

//Needs have a start/stop on the ScrollBox : Pause
//needs bool to track

//scrollrate = some unit of ScrollBox px per ms  takes that and transforms to duration
//nextscroll = time

type ScrollBox struct {
	*Sprite
	Rs                       []Renderable
	nextScrollX, nextScrollY time.Time
	scrollRateX, scrollRateY time.Duration
	View, reappear           physics.Vector
	dirX, dirY               float64

	paused bool
}

func NewScrollBox(rs []Renderable, milliPerPixelX, milliPerPixelY, width, height int) *ScrollBox {
	s := new(ScrollBox)
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

func (s *ScrollBox) DrawOffset(buff draw.Image, xOff, yOff float64) {
	s.update()
	s.Sprite.DrawOffset(buff, xOff, yOff)

}
func (s *ScrollBox) Draw(buff draw.Image) {
	s.DrawOffset(buff, 0, 0)
}

func (s *ScrollBox) update() {
	updatedFlag := false
	if s.paused {
		return
	}
	if s.dirX != 0 && time.Now().After(s.nextScrollX) {
		pixelsMovedX := int64(time.Since(s.nextScrollX))/int64(s.scrollRateX) + 1
		s.nextScrollX = time.Now().Add(s.scrollRateX)

		newS := NewEmptySprite(s.Sprite.X(), s.Sprite.Y(), int(s.View.X()), int(s.View.Y()))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftX(-1 * s.dirX * float64(pixelsMovedX))
			if (s.dirX == 1 && m.GetX() <= s.reappear.X()) || (s.dirX == -1 && m.GetX() >= s.reappear.X()) {
				m.ShiftX(-1 * s.reappear.X()) //Hope that delta is not higher than reappear...
			}

		}
		*s.Sprite = *newS
		updatedFlag = true
	}
	if s.dirY != 0 && time.Now().After(s.nextScrollY) {
		pixelsMovedY := int64(time.Since(s.nextScrollY))/int64(s.scrollRateY) + 1
		s.nextScrollY = time.Now().Add(s.scrollRateY)

		newS := NewEmptySprite(s.Sprite.X(), s.Sprite.Y(), int(s.View.X()), int(s.View.Y()))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftY(-1 * s.dirY * float64(pixelsMovedY))
			if (s.dirY == 1 && m.GetY() <= s.reappear.Y()) || (s.dirY == -1 && m.GetY() >= s.reappear.Y()) {
				m.ShiftY(-1 * s.reappear.Y()) //Hope that delta is not higher than reappear...
			}
		}
		*s.Sprite = *newS
		updatedFlag = true
	}
	if updatedFlag {
		s.drawRenderables()
	}
}
func (s *ScrollBox) Pause() {
	s.paused = true
}
func (s *ScrollBox) Unpause() {
	s.paused = false
	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
}

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

func (s *ScrollBox) AddRenderable(rs ...Renderable) {
	for _, r := range rs {
		switch r.(type) {
		case *Text:
			r.SetPos(r.GetX()*-1, r.GetY()*-1)
		}
		s.Rs = append(s.Rs, r)
	}
	s.drawRenderables()
}

func (s *ScrollBox) drawRenderables() {
	for _, r := range s.Rs {
		r.DrawOffset(s.GetRGBA(), -2*r.GetX(), -2*r.GetY())
		if s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX(), -2*r.GetY()+s.reappear.Y())
		}
		if s.scrollRateX != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX()+s.reappear.X(), -2*r.GetY())
		}
		if s.scrollRateX != 0 && s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX()+s.reappear.X(), -2*r.GetY()+s.reappear.Y())
		}
	}
}
