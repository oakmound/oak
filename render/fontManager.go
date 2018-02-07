package render

import "github.com/oakmound/oak/oakerr"

// A FontManager is just a map for fonts that contains a default font
type FontManager map[string]*Font

// NewFontManager returns a FontManager where 'def' is the default font
func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)["def"] = (&FontGenerator{}).Generate()
	return fm
}

// NewFont adds to the font manager and potentially returns if the key
// was already defined in the map
func (fm *FontManager) NewFont(name string, fg FontGenerator) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = oakerr.ExistingElement{
			InputName:   name,
			InputType:   "font",
			Overwritten: true,
		}
	}
	manager[name] = (&fg).Generate()
	return err

}

// Get retrieves a font from a manager
func (fm *FontManager) Get(name string) *Font {
	manager := (*fm)
	font := manager[name]
	return font
}
