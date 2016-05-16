package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/draw"
	//"image/color"
)

var (
	gameScreen *screen.Screen
	// spriteNames = []string{"Empty":"textures/tile1.png", "textures/tile2.png"}
	spriteNames = map[string]string{
		"Empty": "textures/tile1.png",
		"Wall":  "textures/wall.png",
		"Floor": "textures/floor.png"}
)

type Renderable interface {
	GetRGBA() *image.RGBA
	//ApplyColor(*image.color) *Renderable
	//ApplyMask(??) *Renderable
	//Rotate(degrees int) *Renderable
	//Scale(xRatio int, yRatio int) *Renderable
	//ShiftX(x int)
	//ShiftY(y int)
	Draw(buff screen.Buffer)

	// Squish the renderable into a geometry
	// FitTo(geometry) *Renderable
}

func ParseSprite(s string) Sprite {
	return LoadSprite(spriteNames[s])
}

func ParseSubSprite(s string, x, y, w, h, pad int) Sprite {
	sh, _ := LoadSheet(spriteNames[s], w, h, pad)
	b, _ := (*GetScreen()).NewBuffer(image.Point{w, h})
	draw.Draw(b.RGBA(), b.Bounds(), (*sh)[x][y], image.Point{0, 0}, draw.Src)
	return Sprite{
		x:      0,
		y:      0,
		buffer: &b,
	}
}

func SetScreen(s *screen.Screen) {
	gameScreen = s
}

func GetScreen() *screen.Screen {
	return gameScreen
}

type Sprite struct {
	x, y   float64
	buffer *screen.Buffer
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

//func (s *Sprite) ApplyColor(*image.color) *Renderable {

///}

//func (s *Sprite) ApplyMask(??) *Renderable

//func (s *Sprite) Rotate(degrees int) *Renderable {

//}
//func (s *Sprite) Scale(xRatio int, yRatio int) *Renderable {

//}
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
