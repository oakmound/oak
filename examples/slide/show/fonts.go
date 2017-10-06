package show

import (
	"path/filepath"

	"github.com/oakmound/oak/render"
)

var (
	Express = (&render.FontGenerator{
		File: fpFilter("expressway rg.ttf"),
	}).Generate()
	Gnuolane = (&render.FontGenerator{
		File: fpFilter("gnuolane rg.ttf"),
	}).Generate()
	Libel = (&render.FontGenerator{
		File: fpFilter("libel-suit-rg.ttf"),
	}).Generate()

	Express28 = func() *render.Font {
		e2 := Express.Copy()
		e2.Size = 28
		return e2.Generate()
	}()
)

// todo: we need to do this because some things
// haven't started in the engine yet (the engine
// doesn't know what our directories are for assets)
// Can we change this?
func fpFilter(file string) string {
	return filepath.Join("assets", "font", file)
}
