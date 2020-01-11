package render

import (
	"image"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/oakerr"
)

// Sheet is a 2D array of image rgbas
type Sheet [][]*image.RGBA

// SubSprite gets a sprite from a sheet at the given location
func (sh *Sheet) SubSprite(x, y int) *Sprite {
	return NewSprite(0, 0, (*sh)[x][y])
}

// ToSprites returns this sheet as a 2D array of Sprites
func (sh *Sheet) ToSprites() [][]*Sprite {
	sprites := make([][]*Sprite, len(*sh))
	for x, row := range *sh {
		sprites[x] = make([]*Sprite, len(row))
		for y := range row {
			sprites[x][y] = sh.SubSprite(x, y)
		}
	}
	return sprites
}

// NewSheetSequence creates a Sequence from a sheet and a list of x,y frame coordinates.
// A sequence will be created by getting the sheet's [i][i+1]th elements incrementally
// from the input frames. If the number of input frames is uneven, an error is returned.
func NewSheetSequence(sheet *Sheet, fps float64, frames ...int) (*Sequence, error) {
	if len(frames)%2 != 0 {
		return nil, oakerr.IndivisibleInput{
			InputName:    "frames",
			IsList:       true,
			MustDivideBy: 2,
		}
	}

	sh := *sheet

	mods := make([]Modifiable, len(frames)/2)
	for i := 0; i < len(frames); i += 2 {
		if len(sh) <= frames[i] || len(sh[frames[i]]) <= frames[i+1] {
			dim2 := len(sh) - 1
			if len(sh) > frames[i] {
				dim2 = frames[i]
			}
			dlog.Error("Frame ", frames[i], frames[i+1], "requested but sheet has dimensions ", len(sh), len(sh[dim2]))
			return nil, oakerr.InvalidInput{InputName: "Frame requested does not exist "}
		}
		mods[i/2] = NewSprite(0, 0, sh[frames[i]][frames[i+1]])
	}
	return NewSequence(fps, mods...), nil
}
