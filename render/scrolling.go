package render

import (
	"errors"
	"fmt"
	"image/draw"
	"strconv"
	"time"
)

//Rate of scrool
//Sprite to draw to
//List of things to scroll

//Needs have a start/stop on the scrolling : Pause
//needs bool to track

//scrollrate = some unit of scrolling px per ms  takes that and transforms to duration
//nextscroll = time

//This will only scroll x for now

type Scrolling struct {
	*Sprite
	Rs                       []Renderable
	nextScrollX, nextScrollY time.Time
	scrollRateX, scrollRateY time.Duration
	View, reappear           Point
	dirX, dirY               float64

	paused bool
}

func NewScrolling(rs []Renderable, milliPerPixelX, milliPerPixelY, width, height int) *Scrolling {
	s := new(Scrolling)
	s.Rs = rs
	s.View = Point{float64(width), float64(height)}

	s.SetScrollRate(milliPerPixelX, milliPerPixelY)
	s.reappear = Point{-1 * s.dirX * s.View.X, -1 * s.dirY * s.View.Y}

	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
	s.Sprite = NewEmptySprite(0, 0, width, height)
	fmt.Println("Made a scrolling ")

	s.drawRenderables()
	return s
}

func (s *Scrolling) DrawOffset(buff draw.Image, xOff, yOff float64) {
	s.update()
	s.Sprite.DrawOffset(buff, xOff, yOff)

}
func (s *Scrolling) Draw(buff draw.Image) {
	s.DrawOffset(buff, 0, 0)
}

func (s *Scrolling) update() {
	updatedFlag := false
	if s.paused {
		return
	}
	if s.dirX != 0 && time.Now().After(s.nextScrollX) {
		pixelsMovedX := int64(time.Now().Sub(s.nextScrollX))/int64(s.scrollRateX) + 1
		fmt.Println("Scrolled by ", pixelsMovedX)
		s.nextScrollX = time.Now().Add(s.scrollRateX)

		newS := NewEmptySprite(s.Sprite.X, s.Sprite.Y, int(s.View.X), int(s.View.Y))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftX(-1 * s.dirX * float64(pixelsMovedX))
			if (s.dirX == 1 && m.GetX() <= s.reappear.X) || (s.dirX == -1 && m.GetX() >= s.reappear.X) {
				m.ShiftX(-1 * s.reappear.X) //Hope that delta is not higher than reappear...
			}

		}
		*s.Sprite = *newS
		updatedFlag = true
	}
	if s.dirY != 0 && time.Now().After(s.nextScrollY) {
		pixelsMovedY := int64(time.Now().Sub(s.nextScrollY))/int64(s.scrollRateY) + 1
		s.nextScrollY = time.Now().Add(s.scrollRateY)

		newS := NewEmptySprite(s.Sprite.X, s.Sprite.Y, int(s.View.X), int(s.View.Y))
		newS.SetLayer(s.Sprite.GetLayer())
		for _, m := range s.Rs {
			m.ShiftY(-1 * s.dirY * float64(pixelsMovedY))
			if (s.dirY == 1 && m.GetY() <= s.reappear.Y) || (s.dirY == -1 && m.GetY() >= s.reappear.Y) {
				m.ShiftY(-1 * s.reappear.Y) //Hope that delta is not higher than reappear...
			}
		}
		*s.Sprite = *newS
		updatedFlag = true
	}
	if updatedFlag {
		s.drawRenderables()
	}
}
func (s *Scrolling) Pause() {
	s.paused = true
}
func (s *Scrolling) Unpause() {
	s.paused = false
	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
}

func (s *Scrolling) SetReappearPos(x, y float64) error {
	s.reappear = Point{x, y}
	if x*s.dirX > 0 {
		return errors.New("Scrolling will not loop with direction.X: " + strconv.Itoa(int(s.dirX)) + " and reappear.X: " + strconv.Itoa(int(x)))
	}
	if y*s.dirY > 0 {
		return errors.New("Scrolling will not loop with direction.Y: " + strconv.Itoa(int(s.dirY)) + " and reappear.X: " + strconv.Itoa(int(y)))
	}
	return nil
}

func (s *Scrolling) SetScrollRate(milliPerPixelX, milliPerPixelY int) {
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

func (s *Scrolling) AddRenderable(rs ...Renderable) {
	for _, r := range rs {
		switch r.(type) {
		case *Text, *IntText:
			r.SetPos(r.GetX()*-1, r.GetY()*-1)
		}
		s.Rs = append(s.Rs, r)
	}
	s.drawRenderables()
}

func (s *Scrolling) drawRenderables() {
	for _, r := range s.Rs {
		r.DrawOffset(s.GetRGBA(), -2*r.GetX(), -2*r.GetY())
		if s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX(), -2*r.GetY()+s.reappear.Y)
		}
		if s.scrollRateX != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX()+s.reappear.X, -2*r.GetY())
		}
		if s.scrollRateX != 0 && s.scrollRateY != 0 {
			r.DrawOffset(s.GetRGBA(), -2*r.GetX()+s.reappear.X, -2*r.GetY()+s.reappear.Y)
		}
	}
}
