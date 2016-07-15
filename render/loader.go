package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"errors"
	// "fmt"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexpSingleNumber, _ = regexp.Compile("^\\d+$")
	regexpTwoNumbers, _   = regexp.Compile("^\\d+x\\d+$")
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
	loadedSheets = make(map[string]*Sheet)
	// move to some batch load settings
	defaultPad = 0
)

func loadPNG(directory, fileName string) *screen.Buffer {

	if _, ok := loadedImages[fileName]; ok {
		return loadedImages[fileName]
	}

	s := *GetScreen()

	imgFile, err := os.Open(filepath.Join(directory, fileName))
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

	loadedImages[fileName] = &buff

	dlog.Verb("Loaded filename: ", fileName)

	return &buff
}

func LoadSprite(fileName string) *Sprite {
	return &Sprite{buffer: loadPNG(dir, fileName)}
}

func GetSheet(fileName string) [][]*Sprite {
	sprites := make([][]*Sprite, 0)
	dlog.Verb(loadedSheets, fileName, loadedSheets[fileName])

	sheet, _ := LoadSheet(dir, fileName, 0, 0, 0)
	for x, row := range *sheet {
		sprites = append(sprites, make([]*Sprite, 0))
		for y := range row {
			sprites[x] = append(sprites[x], sheet.SubSprite(x, y))
		}
	}

	return sprites
}

func LoadSheet(directory, fileName string, w, h, pad int) (*Sheet, error) {
	if _, ok := loadedImages[fileName]; !ok {
		dlog.Verb("Missing file in loaded images: ", fileName)
		loadedImages[fileName] = loadPNG(directory, fileName)
	}
	if sheet_p, ok := loadedSheets[fileName]; ok {
		return sheet_p, nil
	}
	buffer := loadedImages[fileName]
	bounds := (*buffer).Size()
	rgba := (*buffer).RGBA()

	sheetW := bounds.X / w
	remainderW := bounds.X % w
	sheetH := bounds.Y / h
	remainderH := bounds.Y % h

	var widthBuffers, heightBuffers int
	if pad != 0 {
		widthBuffers = remainderW / pad
		heightBuffers = remainderH / pad
	} else {
		widthBuffers = sheetW - 1
		heightBuffers = sheetH - 1
	}

	if sheetW < 1 || sheetH < 1 ||
		widthBuffers != sheetW-1 ||
		heightBuffers != sheetH-1 {
		dlog.Error("Bad dimensions given to load sheet")
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

	dlog.Verb("Loaded sheet into map")
	sheet_p := &sheet
	loadedSheets[fileName] = sheet_p

	return loadedSheets[fileName], nil
}

func LoadSheetAnimation(fileName string, w, h, pad int, fps float64, frames []int) (*Animation, error) {
	sheet, err := LoadSheet(dir, fileName, w, h, pad)
	if err != nil {
		return nil, err
	}
	return LoadAnimation(sheet, w, h, pad, fps, frames)
}

func LoadAnimation(sheet *Sheet, w, h, pad int, fps float64, frames []int) (*Animation, error) {
	animation, err := NewAnimation(sheet, fps, frames)
	if err != nil {
		return nil, err
	}
	return animation, nil
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

func BatchLoad(baseFolder string) error {

	// dir2 := filepath.Join(dir, "textures")
	folders, _ := ioutil.ReadDir(baseFolder)

	for i, folder := range folders {
		//Pull in from alias file here to grab things correectly for size

		dlog.Verb("folder ", i, folder.Name())
		if folder.IsDir() {

			var frameW int
			var frameH int

			if folder.Name() == "raw" {
				frameW = 0
				frameH = 0
			} else if result := regexpTwoNumbers.Find([]byte(folder.Name())); result != nil {
				vals := strings.Split(string(result), "x")
				dlog.Verb("Extracted dimensions: ", vals)
				frameW, _ = strconv.Atoi(vals[0])
				frameH, _ = strconv.Atoi(vals[1])
			} else if result := regexpSingleNumber.Find([]byte(folder.Name())); result != nil {
				val, _ := strconv.Atoi(string(result))
				frameW = val
				frameH = val
			} else {
				return errors.New("Aliases for folders are not implemented yet")
			}

			files, _ := ioutil.ReadDir(filepath.Join(baseFolder, folder.Name()))
			for _, file := range files {
				if !file.IsDir() {
					n := file.Name()
					switch n[len(n)-4:] {
					case ".png":
						dlog.Verb("loading file ", n)
						buff := loadPNG(baseFolder, filepath.Join(folder.Name(), n))

						w := (*buff).Size().X
						h := (*buff).Size().Y

						dlog.Verb("buffer: ", w, h, " frame: ", frameW, frameH)

						if frameW == 0 || frameH == 0 {
							continue
						} else if w < frameW || h < frameH {
							dlog.Error("File ", n, " in folder", folder.Name(), " is too small for these folder dimensions")
							return errors.New("File in folder is too small for these folder dimensions")

							// Load this as a sheet if it is greater
							// than the folder size's frame size
						} else if w != frameW || h != frameH {
							dlog.Verb("Loading as sprite sheet")
							LoadSheet(baseFolder, filepath.Join(folder.Name(), n), frameW, frameH, defaultPad)
						}
					default:
						dlog.Error("Unsupported file ending for batchLoad: ", n)
					}
				}
			}
		} else {
			dlog.Verb("Not Folder")
		}

	}
	return nil
}
