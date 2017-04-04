package render

import (
	"image"
	"image/color"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/physics"
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
	s.SetRGBA(ApplyColor(s.GetRGBA(), c))
	return s
}

func (s *Sprite) FillMask(img image.RGBA) Modifiable {
	s.SetRGBA(FillMask(s.GetRGBA(), img))
	return s
}

func (s *Sprite) ApplyMask(img image.RGBA) Modifiable {
	s.SetRGBA(ApplyMask(s.GetRGBA(), img))
	return s
}

func (s *Sprite) Rotate(degrees int) Modifiable {
	s.SetRGBA(Rotate(s.GetRGBA(), degrees))
	return s
}
func (s *Sprite) Scale(xRatio float64, yRatio float64) Modifiable {
	s.SetRGBA(Scale(s.GetRGBA(), xRatio, yRatio))
	return s
}
func (s *Sprite) FlipX() Modifiable {
	s.SetRGBA(FlipX(s.GetRGBA()))
	return s
}
func (s *Sprite) FlipY() Modifiable {
	s.SetRGBA(FlipY(s.GetRGBA()))
	return s
}
func (s *Sprite) Fade(alpha int) Modifiable {
	s.SetRGBA(Fade(s.GetRGBA(), alpha))
	return s
}
func OverlaySprites(sps []Sprite) *Sprite {
	tmpSprite := sps[len(sps)-1].Copy().(*Sprite)
	for i := len(sps) - 1; i > 0; i-- {
		tmpSprite.SetRGBA(FillMask(tmpSprite.GetRGBA(), *sps[i-1].GetRGBA()))
	}
	return tmpSprite
}

func ParseSubSprite(s string, x, y, w, h, pad int) *Sprite {
	sh, _ := LoadSheet(dir, s, w, h, pad)
	return sh.SubSprite(x, y)
}
