package audio

import (
	"bitbucket.org/oakmoundstudio/oak/oakerr"
	"github.com/200sc/klangsynthese/font"
)

// A FontManager is a map of names to Fonts that has a built in
// default font at name 'def'.
type FontManager map[string]*font.Font

// NewFontManager returns a manager with a single 'def' font
func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)["def"] = &font.Font{}
	return fm
}

// NewFont adds a font to a manger with the given keyed name.
// NewFont can return an error indicating if the name assigned
// was already in use.
func (fm *FontManager) NewFont(name string, f *font.Font) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = oakerr.ExistingFontError{}
	}
	manager[name] = f
	return err

}

// Get returns whatever is at name in font
func (fm *FontManager) Get(name string) *font.Font {
	manager := (*fm)
	font, _ := manager[name]
	return font
}
