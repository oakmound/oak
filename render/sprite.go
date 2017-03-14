package render

import (
	"bitbucket.org/oakmoundstudio/oak/physics"
	"image"
	"image/color"
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
		LayeredPoint: LayeredPoint{
			Vector: physics.Vector{
				X: x,
				Y: y,
			},
		},
		r: r,
	}
}

func (s *Sprite) GetRGBA() *image.RGBA {
	return s.r
}

func (s *Sprite) SetRGBA(r *image.RGBA) {
	s.r = r
}

func (s *Sprite) DrawOffset(buff draw.Image, xOff, yOff float64) {
	ShinyDraw(buff, s.r, int(s.X+xOff), int(s.Y+yOff))
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

func (s *Sprite) ApplyColor(c color.Color) Modifiable {
	s.r = ApplyColor(s.r, c)
	return s
}

func (s *Sprite) FillMask(img image.RGBA) Modifiable {
	s.r = FillMask(s.r, img)
	return s
}

func (s *Sprite) ApplyMask(img image.RGBA) Modifiable {
	s.r = ApplyMask(s.r, img)
	return s
}

func (s *Sprite) Rotate(degrees int) Modifiable {
	s.r = Rotate(s.r, degrees)
	return s
}
func (s *Sprite) Scale(xRatio float64, yRatio float64) Modifiable {
	s.r = Scale(s.r, xRatio, yRatio)
	return s
}
func (s *Sprite) FlipX() Modifiable {
	s.r = FlipX(s.r)
	return s
}
func (s *Sprite) FlipY() Modifiable {
	s.r = FlipY(s.r)
	return s
}
func (s *Sprite) Fade(alpha int) Modifiable {
	s.r = Fade(s.r, alpha)
	return s
}
func OverlaySprites(sps []Sprite) *Sprite {
	tmpSprite := sps[len(sps)-1].Copy().(*Sprite)
	for i := len(sps) - 1; i > 0; i-- {
		tmpSprite.r = FillMask(tmpSprite.r, *sps[i-1].r)
	}
	return tmpSprite
}

func ParseSubSprite(s string, x, y, w, h, pad int) *Sprite {
	sh, _ := LoadSheet(dir, s, w, h, pad)
	return sh.SubSprite(x, y)
}
