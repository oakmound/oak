package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
)

var (
	gameScreen *screen.Screen
	// spriteNames = []string{"Empty":"textures/tile1.png", "textures/tile2.png"}

)

type Renderable interface {
	GetRGBA() *image.RGBA
	//ApplyColor(*image.color) *Renderable
	//ApplyMask(??) *Renderable
	//Rotate(degrees int) *Renderable
	//Scale(xRatio int, yRatio int) *Renderable
	ShiftX(x float64)
	ShiftY(y float64)
	Draw(buff screen.Buffer)
	GetLayer() int
	SetLayer(l int)
	UnDraw()
	SetPos(x, y float64)
	// Squish the renderable into a geometry
	// FitTo(geometry) *Renderable
}

func SetScreen(s *screen.Screen) {
	gameScreen = s
}

func GetScreen() *screen.Screen {
	return gameScreen
}
