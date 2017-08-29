package render

import (
	"errors"
	"image"
	"image/draw"
	"time"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/timing"
)

// Sheet is a 2D array of image rgbas
type Sheet [][]*image.RGBA

//SubSprite gets a sprite from a sheet at the given location
func (sh *Sheet) SubSprite(x, y int) *Sprite {
	return NewSprite(0, 0, (*sh)[x][y])
}

//Animation takes a set of frames and provides a framework for animating them
type Animation struct {
	LayeredPoint
	pauseBool
	InterruptBool
	sheetPos   int
	frameTime  int64
	frames     [][]int
	sheet      *Sheet
	lastChange time.Time
	cID        event.CID
}

//NewAnimation creates an Animation
func NewAnimation(sheetP *Sheet, fps float64, frames []int) (*Animation, error) {

	if len(frames)%2 != 0 {
		return nil, errors.New("Uneven number of animation coordinates")
	}

	splitFrames := make([][]int, len(frames)/2)
	for i := 0; i < len(frames); i += 2 {
		splitFrames[i/2] = []int{frames[i], frames[i+1]}
	}

	animation := Animation{
		LayeredPoint: LayeredPoint{
			Vector: physics.NewVector(0, 0),
		},
		InterruptBool: InterruptBool{
			Interruptable: true,
		},
		pauseBool: pauseBool{
			playing: true,
		},
		sheetPos:   0,
		frameTime:  timing.FPSToNano(fps),
		frames:     splitFrames,
		sheet:      sheetP,
		lastChange: time.Now(),
	}

	return &animation, nil
}

//Copy creates a new Modifiable that is a copy of the current animation.
func (a *Animation) Copy() Modifiable {
	newA := new(Animation)
	newA.LayeredPoint = a.LayeredPoint.Copy()
	newA.sheetPos = a.sheetPos
	newA.frameTime = a.frameTime
	newA.frames = a.frames
	newA.lastChange = a.lastChange
	newA.playing = a.playing
	newA.Interruptable = a.Interruptable
	newA.cID = a.cID
	// Manual deep copy of pointers
	aSheet := *a.sheet
	newA.sheet = new(Sheet)
	newSheet := make(Sheet, len(aSheet))
	for x, col := range aSheet {
		newSheet[x] = make([]*image.RGBA, len(aSheet[x]))
		for y, val := range col {
			newRGBA := new(image.RGBA)
			*newRGBA = *val
			newSheet[x][y] = newRGBA
		}
	}
	*newA.sheet = newSheet
	return newA
}

//SetTriggerID sets animation's id to the passed in id.
func (a *Animation) SetTriggerID(id event.CID) {
	a.cID = id
}

func (a *Animation) update() {
	if a.playing && time.Since(a.lastChange).Nanoseconds() > a.frameTime {
		a.lastChange = time.Now()
		a.sheetPos = (a.sheetPos + 1) % len(a.frames)
		// Eventually, if an animation is cut off before
		// it finishes by another animation starting,
		// AnimationCancelled (or maybe AnimationShortCircuit)
		// should trigger instead of AnimationEnd
		if a.sheetPos == 0 && a.cID != 0 {
			a.cID.Trigger("AnimationEnd", nil)
		}
	}
}

//DrawOffset draws the animation with some x,y  offset from its logical location
func (a *Animation) DrawOffset(buff draw.Image, xOff, yOff float64) {
	a.update()
	img := a.GetRGBA()
	ShinyDraw(buff, img, int(a.X()+xOff), int(a.Y()+yOff))
}

//Draw draws the animation at its logical location
func (a *Animation) Draw(buff draw.Image) {
	a.update()
	img := a.GetRGBA()
	ShinyDraw(buff, img, int(a.X()), int(a.Y()))
}

//GetRGBA returns the current frames rgba
func (a *Animation) GetRGBA() *image.RGBA {
	return (*a.sheet)[a.frames[a.sheetPos][0]][a.frames[a.sheetPos][1]]
}

//GetDims returns the dimensions of the animation in terms of  x, y
func (a *Animation) GetDims() (int, int) {
	r := a.GetRGBA()
	return r.Bounds().Max.X, r.Bounds().Max.Y
}

//Modify applies a set of modifications to the Animation
func (a *Animation) Modify(ms ...Modification) Modifiable {
	sheet := *a.sheet
	for x, row := range sheet {
		for y, rgba := range row {
			for _, m := range ms {
				sheet[x][y] = m(rgba)
			}
		}
	}
	return a
}

// IsStatic returns false for animations
func (a *Animation) IsStatic() bool {
	return false
}
