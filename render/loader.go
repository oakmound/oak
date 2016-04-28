package render

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

var (
	// Form ...main/core.go/../assets/images,
	// the image directory.
	wd, _ = os.Getwd()
	dir   = filepath.Join(
		filepath.Dir(wd),
		"assets",
		"images")
	loadedSprites = make(map[string]*screen.Buffer)
)

func loadPNG(fileName string) *screen.Buffer {

	s := *GetScreen()

	imgFile, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	buff, err := s.NewBuffer(img.Bounds().Max)
	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(buff.RGBA(), img.Bounds(), img, image.Point{0, 0}, draw.Src)

	return &buff
}

func LoadSprite(fileName string) Sprite {
	if _, ok := loadedSprites[fileName]; !ok {
		loadedSprites[fileName] = loadPNG(fileName)
	}
	return Sprite{buffer: loadedSprites[fileName]}
}

// For loading a sheet as an array of images

//func loadSheet(fileName s, frameWidth int, frameHeight int) {
//
//}

// I don't know if we need this.
//func unloadFile(fileName s) {
//
//}
