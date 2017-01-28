package render

import (
	"image/draw"
	"time"
	"fmt"
)


//Rate of scrool
//Sprite to draw to
//List of things to scroll

//Needs have a start/stop on the scrolling : Pause
//needs bool to track


//scrollrate = some unit of scrolling px per ms  takes that and transforms to duration
//nextscroll = time

//This will only scroll x for now





type  Scrolling  struct{
	*Sprite
	Ms []Modifiable
	nextScrollX, nextScrollY time.Time
	scrollRateX, scrollRateY time.Duration
	View, Reappear Point
	dirX, dirY float64

	paused bool
}






func NewScrolling(ms []Modifiable,  milliPerPixelX, milliPerPixelY, width, height int)  *Scrolling{
	s := new(Scrolling)
	s.Ms =  ms
	s.View = Point{float64(width), float64(height)}
	s.Reappear = Point{s.View.X, s.View.Y}

	s.SetScrollRate(milliPerPixelX, milliPerPixelY)

	s.nextScrollX = time.Now().Add(s.scrollRateX)
	s.nextScrollY = time.Now().Add(s.scrollRateY)
	s.Sprite = NewEmptySprite(0, 0, width, height)
	fmt.Println("Made a scrolling ")

	for _, m := range s.Ms {
		m.DrawOffset(s.Sprite.GetRGBA(), -2*m.GetX(), -2*m.GetY())
		m.DrawOffset(s.Sprite.GetRGBA(), -2*m.GetX()-s.Reappear.X, -2*m.GetY())
	}
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
	if s.paused{
		return
	}
	if  time.Now().After(s.nextScrollX) {
		pixelsMovedX := int64(time.Now().Sub(s.nextScrollX)) / int64(s.scrollRateX) + 1
		fmt.Println("Scrolled by " , pixelsMovedX)
		s.nextScrollX = time.Now().Add(s.scrollRateX)

		newS := NewEmptySprite(s.Sprite.X, s.Sprite.Y, int(s.View.X), int(s.View.Y))

		for _, m := range s.Ms {
			m.ShiftX(-1 * s.dirX *  float64(pixelsMovedX))
			if m.GetX() <= -1*s.Reappear.X {
				m.ShiftX(s.Reappear.X) //Hope that delta is not higher than reappear...
			}
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX(), -2*m.GetY())
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX(), -2*m.GetY()-s.Reappear.Y)
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX()-s.Reappear.X, -2*m.GetY())
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX()-s.Reappear.X, -2*m.GetY()-s.Reappear.Y)
		}
		s.Sprite = newS
	}
	if time.Now().After(s.nextScrollY) {
		pixelsMovedY := int64(time.Now().Sub(s.nextScrollY)) / int64(s.scrollRateY) + 1
		s.nextScrollY = time.Now().Add(s.scrollRateY)

		newS := NewEmptySprite(s.Sprite.X, s.Sprite.Y, int(s.View.X), int(s.View.Y))

		for _, m := range s.Ms {
			m.ShiftY(-1 * s.dirY *  float64(pixelsMovedY))
			if m.GetY() <= -1*s.Reappear.Y {
				m.ShiftY(s.Reappear.Y) //Hope that delta is not higher than reappear...
			}
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX(), -2*m.GetY())
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX(), -2*m.GetY()-s.Reappear.Y)
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX()-s.Reappear.X, -2*m.GetY())
			m.DrawOffset(newS.GetRGBA(), -2*m.GetX()-s.Reappear.X, -2*m.GetY()-s.Reappear.Y)
		}
		s.Sprite = newS
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

func (s * Scrolling) SetScrollRate(milliPerPixelX, milliPerPixelY  int){
	s.dirX = 1
	s.dirY = 1
	if milliPerPixelX < 0{
		milliPerPixelX *=-1
		s.dirX = -1
	}
	if milliPerPixelY < 0{
		milliPerPixelY *=-1
		s.dirY = -1
	}
	s.scrollRateX = time.Duration(milliPerPixelX) * time.Millisecond
	s.scrollRateY = time.Duration(milliPerPixelY) * time.Millisecond
}



