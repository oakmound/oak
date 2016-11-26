package render

import (
	"bitbucket.org/oakmoundstudio/plasticpiston/plastic/event"
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
	LayeredPoint
	sheetPos   int
	frameTime  int64
	frames     [][]int
	sheet      *Sheet
	lastChange time.Time
	playing    bool
	cID        event.CID
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
		LayeredPoint: LayeredPoint{
			Point: Point{
				X: 0.0,
				Y: 0.0,
			},
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

func (a *Animation) SetTriggerID(id event.CID) {
	a.cID = id
}

func (a *Animation) updateAnimation() {
	if a.playing && time.Since(a.lastChange).Nanoseconds() > a.frameTime {
		a.lastChange = time.Now()
		a.sheetPos = (a.sheetPos + 1) % len(a.frames)
		if a.sheetPos == 0 && a.cID != 0 {
			a.cID.Trigger("AnimationEnd", nil)
		}
	}
}

func (a *Animation) DrawOffset(buff draw.Image, xOff, yOff float64) {
	a.updateAnimation()
	img := a.GetRGBA()
	ShinyDraw(buff, img, int(a.X+xOff), int(a.Y+yOff))
}

func (a *Animation) Draw(buff draw.Image) {
	a.updateAnimation()
	img := a.GetRGBA()
	ShinyDraw(buff, img, int(a.X), int(a.Y))
}

func (a_p *Animation) GetRGBA() *image.RGBA {
	return (*a_p.sheet)[a_p.frames[a_p.sheetPos][0]][a_p.frames[a_p.sheetPos][1]]
}

func (a *Animation) ApplyColor(c color.Color) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = ApplyColor(rgba, c)
		}
	}
	return a
}

func (a *Animation) FillMask(img image.RGBA) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FillMask(rgba, img)
		}
	}
	return a
}

func (a *Animation) ApplyMask(img image.RGBA) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = ApplyMask(rgba, img)
		}
	}
	return a
}

func (a *Animation) Rotate(degrees int) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = Rotate(rgba, degrees)
		}
	}
	return a
}

func (a *Animation) Scale(xRatio float64, yRatio float64) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = Scale(rgba, xRatio, yRatio)
		}
	}
	return a
}

func (a *Animation) FlipX() Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FlipX(rgba)
		}
	}
	return a
}

func (a *Animation) FlipY() Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = FlipY(rgba)
		}
	}
	return a
}
func (a *Animation) Fade(alpha int) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			sheet[x][y] = Fade(rgba, alpha)
		}
	}
	return a
}

func (a *Animation) Pause() {
	a.playing = false
}

func (a *Animation) Unpause() {
	a.playing = true
}
