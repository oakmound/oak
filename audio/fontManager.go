//+build windows

package audio

import "bitbucket.org/oakmoundstudio/oak/oakerr"

// A FontManager is a map of names to Fonts that has a built in
// default font at name 'def'.
type FontManager map[string]*Font

// NewFontManager returns a manager with a single 'def' font
func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)["def"] = &Font{}
	return fm
}

// NewFont adds a font to a manger with the given keyed name.
// NewFont can return an error indicating if the name assigned
// was already in use.
func (fm *FontManager) NewFont(name string, f *Font) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = oakerr.ExistingFontError{}
	}
	manager[name] = f
	return err

}

// Get returns whatever is at name in font
func (fm *FontManager) Get(name string) *Font {
	manager := (*fm)
	font, _ := manager[name]
	return font
}
