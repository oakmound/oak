package btn

import (
	"fmt"

	"github.com/oakmound/oak/v2/entities"

	"github.com/oakmound/oak/v2/render"
)

// NewText creates some uitext
func NewText(f *render.Font, str string, x, y float64, layers ...int) *entities.Doodad {
	d := entities.NewDoodad(x, y, f.NewStrText(str, x, y), 0)
	render.Draw(d.R, layers...)
	return d
}

// NewIntText creates some uitext from an integer
func NewIntText(f *render.Font, str *int, x, y float64, layers ...int) *entities.Doodad {
	d := entities.NewDoodad(x, y, f.NewIntText(str, x, y), 0)
	render.Draw(d.R, layers...)
	return d
}

// NewRawText creates some uitext from a stringer
func NewRawText(f *render.Font, str fmt.Stringer, x, y float64, layers ...int) *entities.Doodad {
	d := entities.NewDoodad(x, y, f.NewText(str, x, y), 0)
	render.Draw(d.R, layers...)
	return d
}
