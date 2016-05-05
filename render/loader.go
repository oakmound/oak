package render

import (
	"errors"
	"fmt"
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
	loadedImages = make(map[string]*screen.Buffer)
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

	fmt.Println(fileName, "buffer on load:", buff)

	return &buff
}

func LoadSprite(fileName string) Sprite {
	if _, ok := loadedImages[fileName]; !ok {
		loadedImages[fileName] = loadPNG(fileName)
	}
	return Sprite{buffer: loadedImages[fileName]}
}

func LoadSheet(fileName string, w, h, pad int) (*Sheet, error) {
	if _, ok := loadedImages[fileName]; !ok {
		loadedImages[fileName] = loadPNG(fileName)
	}
	buffer := loadedImages[fileName]
	bounds := (*buffer).Size()
	rgba := (*buffer).RGBA()

	sheetW := bounds.X / w
	remainderW := bounds.X % w
	sheetH := bounds.Y / h
	remainderH := bounds.Y % h

	widthBuffers := remainderW / pad
	heightBuffers := remainderH / pad

	if sheetW < 1 || sheetH < 1 ||
		widthBuffers != sheetW-1 ||
		heightBuffers != sheetH-1 {
		return nil, errors.New("Bad dimensions given to load sheet")
	}

	sheet := make(Sheet, sheetW)
	i := 0
	for x := 0; x < bounds.X; x += (w + pad) {
		sheet[i] = make([]*image.RGBA, sheetH)
		j := 0
		for y := 0; y < bounds.Y; y += (h + pad) {
			sheet[i][j] = subImage(rgba, x, y, w, h)
			j++
		}
		i++
	}

	//fmt.Println("Sheet[0][0]", sheet[0][0])

	return &sheet, nil

}

func subImage(rgba *image.RGBA, x, y, w, h int) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			out.Set(i, j, rgba.At(x+i, y+j))
		}
	}
	return out
}
