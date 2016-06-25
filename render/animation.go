package render

import (
	"errors"
	//"fmt"
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/dlog"
	"golang.org/x/exp/shiny/screen"
	"image"
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

func (a Animation) GetRGBA() *image.RGBA {
	return (*a.sheet)[a.frames[a.sheetPos][0]][a.frames[a.sheetPos][1]]
}

func (a *Animation) GetLayer() int {
	return a.layer
}

func (a *Animation) SetLayer(l int) {
	a.layer = l
}

// Creates a new sheet and then sets the animation's sheet to be the new sheet
// func (a Animation) Scale()
