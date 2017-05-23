package render

import (
	"image"
	"image/draw"
)

type Sprite struct {
	LayeredPoint
	r *image.RGBA
}

func NewEmptySprite(x, y float64, w, h int) *Sprite {
	r := image.NewRGBA(image.Rect(0, 0, w, h))
	return NewSprite(x, y, r)
}

func NewSprite(x, y float64, r *image.RGBA) *Sprite {
	return &Sprite{
		LayeredPoint: NewLayeredPoint(x, y, 0),
		r:            r,
	}
}

func (s *Sprite) GetRGBA() *image.RGBA {
	return s.r
}
func (s *Sprite) GetDims() (int, int) {
	rgba := s.r
	if rgba == nil {
		return 6, 6
	}
	return rgba.Bounds().Max.X, rgba.Bounds().Max.Y
}

func (s *Sprite) SetRGBA(r *image.RGBA) {
	s.r = r
}

func (s *Sprite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	ShinyDraw(buff, s.r, int(s.X()+xOff), int(s.Y()+yOff))
}

func (s *Sprite) Draw(buff draw.Image) {
	ShinyDraw(buff, s.r, int(s.X()), int(s.Y()))
}

func (s *Sprite) Copy() Modifiable {
	newS := new(Sprite)
	if s.r != nil {
		newS.r = new(image.RGBA)
		*newS.r = *s.r
	}
	newS.LayeredPoint = s.LayeredPoint.Copy()
	return newS
}

func (s *Sprite) IsNil() bool {
	return s.r == nil
}

func (s *Sprite) Modify(ms ...Modification) Modifiable {
	for _, m := range ms {
		s.r = m(s.GetRGBA())
	}
	return s
}

func OverlaySprites(sps []Sprite) *Sprite {
	tmpSprite := sps[len(sps)-1].Copy().(*Sprite)
	for i := len(sps) - 1; i > 0; i-- {
		tmpSprite.SetRGBA(FillMask(*sps[i-1].GetRGBA())(tmpSprite.GetRGBA()))
	}
	return tmpSprite
}

func ParseSubSprite(s string, x, y, w, h, pad int) *Sprite {
	sh, _ := LoadSheet(dir, s, w, h, pad)
	return sh.SubSprite(x, y)
}
