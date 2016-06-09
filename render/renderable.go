package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/draw"
	"log"
	//"image/color"
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
	//ShiftX(x int)
	//ShiftY(y int)
	Draw(buff screen.Buffer)

	// Squish the renderable into a geometry
	// FitTo(geometry) *Renderable
}

func SetScreen(s *screen.Screen) {
	gameScreen = s
}

func GetScreen() *screen.Screen {
	return gameScreen
}

func RGBAtoBuffer(img *image.RGBA) *screen.Buffer {
	buff, err := (*gameScreen).NewBuffer(img.Bounds().Max)
	if err != nil {
		log.Fatal(err)
	}
	draw.Draw(buff.RGBA(), img.Bounds(), img, image.Point{0, 0}, draw.Src)
	return &buff
}
