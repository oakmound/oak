package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

type Sheet [][]*image.RGBA

func (sh *Sheet) SubSprite(x, y int) *Sprite {
	return &Sprite{
		r: (*sh)[x][y],
	}
}

type Animation struct {
	Point
	Layered
	sheetPos   int
	frameTime  int64
	frames     [][]int
	sheet      *Sheet
	lastChange time.Time
	playing    bool
}

func NewAnimation(sheet_p *Sheet, fps float64, frames []int) (*Animation, error) {

	if len(frames)%2 != 0 {
		return nil, errors.New("Uneven number of animation coordinates")
	}

	frameTime := math.Pow(10, 9) / fps
	splitFrames := make([][]int, len(frames)/2)
	for i := 0; i < len(frames); i += 2 {
		splitFrames[i/2] = []int{frames[i], frames[i+1]}
	}

	animation := Animation{
		Point: Point{
			X: 0.0,
			Y: 0.0,
		},
		sheetPos:   0,
		frameTime:  int64(frameTime),
		frames:     splitFrames,
		sheet:      sheet_p,
		lastChange: time.Now(),
		playing:    true,
	}

	return &animation, nil
}

func (a_p *Animation) Copy() Modifiable {
	newA := new(Animation)
	*newA = *a_p
	// Manual deep copy of pointers
	aSheet := *a_p.sheet
	sheetPointer := new(Sheet)
	newSheet := make(Sheet, len(aSheet))
	for x, col := range aSheet {
		newSheet[x] = make([]*image.RGBA, len(aSheet[x]))
		for y, val := range col {
			newRGBA := new(image.RGBA)
			*newRGBA = *val
			newSheet[x][y] = newRGBA
		}
	}
	*sheetPointer = newSheet
	newA.sheet = sheetPointer
	return newA
}

func (a_p *Animation) updateAnimation() {
	if a_p.playing && time.Since(a_p.lastChange).Nanoseconds() > a_p.frameTime {
		dlog.Verb("Increment sheetPos")
		a_p.lastChange = time.Now()
		a_p.sheetPos = (a_p.sheetPos + 1) % len(a_p.frames)
	}
}

func (a_p *Animation) Draw(buff draw.Image) {
	a_p.updateAnimation()
	img := a_p.GetRGBA()
	ShinyDraw(buff, img, int(a_p.X), int(a_p.Y))
}

func (a_p *Animation) GetRGBA() *image.RGBA {
	return (*a_p.sheet)[a_p.frames[a_p.sheetPos][0]][a_p.frames[a_p.sheetPos][1]]
}

func (a *Animation) ApplyColor(c color.Color) {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = ApplyColor(rgba, c)
		}
	}
}

func (a *Animation) FillMask(img image.RGBA) {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FillMask(rgba, img)
		}
	}
}

func (a *Animation) ApplyMask(img image.RGBA) {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = ApplyMask(rgba, img)
		}
	}
}

func (a *Animation) Rotate(degrees int) {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = Rotate(rgba, degrees)
		}
	}
}

func (a *Animation) Scale(xRatio float64, yRatio float64) {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = Scale(rgba, xRatio, yRatio)
		}
	}
}

func (a *Animation) FlipX() {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FlipX(rgba)
		}
	}
}

func (a *Animation) FlipY() {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FlipY(rgba)
		}
	}
}

func (a *Animation) Pause() {
	a.playing = false
}

func (a *Animation) Unpause() {
	a.playing = true
}
