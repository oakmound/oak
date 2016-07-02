package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"errors"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

type Sheet [][]*image.RGBA

type Animation struct {
	x, y       float64
	sheetPos   int
	frameTime  int64
	frames     [][]int
	sheet      *Sheet
	lastChange time.Time
	playing    bool
	layer      int
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
		x:          0.0,
		y:          0.0,
		sheetPos:   0,
		frameTime:  int64(frameTime),
		frames:     splitFrames,
		sheet:      sheet_p,
		lastChange: time.Now(),
		playing:    false,
	}

	return &animation, nil
}

func (a_p *Animation) ShiftX(x float64) {
	a_p.x += x
}
func (a_p *Animation) ShiftY(y float64) {
	a_p.y += y
}

func (a Animation) SetPos(x, y float64) {
	(&a).x = x
	(&a).y = y
}

func (a_p *Animation) updateAnimation() {
	if time.Since(a_p.lastChange).Nanoseconds() > a_p.frameTime {
		dlog.Verb("Increment sheetPos")
		a_p.lastChange = time.Now()
		a_p.sheetPos = (a_p.sheetPos + 1) % len(a_p.frames)
	}
}

func (a_p *Animation) Draw(buff screen.Buffer) {
	a_p.playing = true
	a_p.updateAnimation()
	img := a_p.GetRGBA()
	draw.Draw(buff.RGBA(), buff.Bounds(),
		img, image.Point{int(a_p.x),
			int(a_p.y)}, draw.Over)
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

func (a *Animation) GetLayer() int {
	return a.layer
}

func (a *Animation) SetLayer(l int) {
	a.layer = l
}

func (a *Animation) UnDraw() {
	a.layer = -1
}
