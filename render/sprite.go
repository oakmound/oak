package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"math"
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
	r1, g1, b1, _ := c.RGBA()
	bounds := (*s.buffer).Size()
	rgba := (*s.buffer).RGBA()
	w := bounds.X
	h := bounds.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()
			tmp := color.RGBA{uint8(r1 * r2 / 255),
				uint8(g1 * g2 / 255),
				uint8(b1 * b2 / 255),
				uint8(a2)}
			newRgba.Set(x, y, tmp)
		}
	}
	out := RGBAtoBuffer(newRgba)
	s.buffer = out
}

func (s *Sprite) ApplyMask(img image.RGBA) {
	// Instead of static color it just two buffers melding
	bounds := (*s.buffer).Size()
	rgba := (*s.buffer).RGBA()
	w := bounds.X
	h := bounds.Y
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r2, g2, b2, a2 := rgba.At(x, y).RGBA()
			tmp := color.RGBA{uint8(r1 * r2 / 255),
				uint8(g1 * g2 / 255),
				uint8(b1 * b2 / 255),
				uint8(a2)}
			newRgba.Set(x, y, tmp)
		}
	}
	out := RGBAtoBuffer(newRgba)
	s.buffer = out
}

func (s *Sprite) Rotate(degrees int) {
	//otates clockwise by the given degrees
	// Will shear any pixels that land outside the given buffer
	angle := float64(degrees) / 180 * math.Pi
	bounds := (*s.buffer).Size()
	rgba := (*s.buffer).RGBA()
	w := bounds.X
	h := bounds.Y
	centerX := float64(w / 2)
	centerY := float64(h / 2)

	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xf := float64(x)
			yf := float64(y)
			newX := int(math.Cos(angle)*(xf-centerX) - math.Sin(angle)*(yf-centerY) + centerX)
			newY := int(math.Sin(angle)*(xf-centerX) + math.Cos(angle)*(yf-centerY) + centerY)
			newRgba.Set(newX, newY, rgba.At(x, y))
		}
	}
	out := RGBAtoBuffer(newRgba)
	s.buffer = out
}
func (s *Sprite) Scale(xRatio float64, yRatio float64) {
	bounds := (*s.buffer).Size()
	rgba := (*s.buffer).RGBA()
	w := int(math.Floor(float64(bounds.X) * xRatio))
	h := int(math.Floor(float64(bounds.Y) * yRatio))
	newRgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			newRgba.Set(x, y, rgba.At(int(math.Floor(float64(x)/xRatio)), int(math.Floor(float64(y)/yRatio))))
		}
	}
	out := RGBAtoBuffer(newRgba)
	s.buffer = out
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
