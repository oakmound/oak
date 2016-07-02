package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
)

var (
	spriteNames = map[string]string{
		"Empty": "textures/tile1.png",
		"Wall":  "textures/wall.png",
		"Floor": "textures/floor.png"}
)

type Sprite struct {
	x, y   float64
	buffer *screen.Buffer
	layer  int
}

func ParseSprite(s string) *Sprite {
	return LoadSprite(spriteNames[s])
}

func ParseSubSprite(s string, x, y, w, h, pad int) *Sprite {
	sh, _ := LoadSheet(spriteNames[s], w, h, pad)
	b, _ := (*GetScreen()).NewBuffer(image.Point{w, h})
	draw.Draw(b.RGBA(), b.Bounds(), (*sh)[x][y], image.Point{0, 0}, draw.Src)
	return &Sprite{buffer: &b}
}

func (s Sprite) GetRGBA() *image.RGBA {
	return (*s.buffer).RGBA()
}

func (s Sprite) HasBuffer() bool {
	if s.buffer != nil {
		return true
	}
	return false
}

func (s *Sprite) ApplyColor(c color.Color) {
	out := ApplyColor((*s.buffer).RGBA(), c)
	s.buffer = RGBAtoBuffer(out)
}

func (s *Sprite) ApplyMask(img image.RGBA) {
	out := ApplyMask((*s.buffer).RGBA(), img)
	s.buffer = RGBAtoBuffer(out)
}

func (s *Sprite) Rotate(degrees int) {
	out := Rotate((*s.buffer).RGBA(), degrees)
	s.buffer = RGBAtoBuffer(out)
}
func (s *Sprite) Scale(xRatio float64, yRatio float64) {
	out := Scale((*s.buffer).RGBA(), xRatio, yRatio)
	s.buffer = RGBAtoBuffer(out)
}
func (s *Sprite) FlipX() {
	out := FlipX((*s.buffer).RGBA())
	s.buffer = RGBAtoBuffer(out)
}
func (s *Sprite) FlipY() {
	out := FlipY((*s.buffer).RGBA())
	s.buffer = RGBAtoBuffer(out)
}

func (s_p *Sprite) ShiftX(x float64) {
	s_p.x += x
}
func (s_p *Sprite) ShiftY(y float64) {
	s_p.y += y
}

func (s Sprite) Draw(buff screen.Buffer) {
	// s := *s_p
	img := (&s).GetRGBA()
	draw.Draw(buff.RGBA(), buff.Bounds(),
		img, image.Point{int((&s).x),
			int((&s).y)}, draw.Over)
}

func (s *Sprite) GetLayer() int {
	return s.layer
}

func (s *Sprite) SetLayer(l int) {
	s.layer = l
}

func (s *Sprite) UnDraw() {
	s.layer = -1
}
