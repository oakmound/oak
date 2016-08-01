package render

import (
	"image"
	"image/color"
	"image/draw"
)

type Sprite struct {
	Point
	Layered
	r *image.RGBA
}

func (s *Sprite) GetRGBA() *image.RGBA {
	return s.r
}

func (s *Sprite) Draw(buff draw.Image) {
	ShinyDraw(buff, s.r, int(s.X), int(s.Y))
}

func (s *Sprite) Copy() Modifiable {
	newS := new(Sprite)
	*newS = *s
	return newS
}

func (s *Sprite) IsNil() bool {
	return s.r == nil
}

func (s *Sprite) ApplyColor(c color.Color) {
	s.r = ApplyColor(s.r, c)
}

func (s *Sprite) FillMask(img image.RGBA) {
	s.r = FillMask(s.r, img)
}

func (s *Sprite) ApplyMask(img image.RGBA) {
	s.r = ApplyMask(s.r, img)
}

func (s *Sprite) Rotate(degrees int) {
	s.r = Rotate(s.r, degrees)
}
func (s *Sprite) Scale(xRatio float64, yRatio float64) {
	s.r = Scale(s.r, xRatio, yRatio)
}
func (s *Sprite) FlipX() {
	s.r = FlipX(s.r)
}
func (s *Sprite) FlipY() {
	s.r = FlipY(s.r)
}

func ParseSubSprite(s string, x, y, w, h, pad int) *Sprite {
	sh, _ := LoadSheet(dir, s, w, h, pad)
	return sh.SubSprite(x, y)
}
