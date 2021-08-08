package audio

import (
	"github.com/oakmound/oak/v3/audio/font"
	"github.com/oakmound/oak/v3/oakerr"
)

const defaultFontKey = "def"

// DefaultFont is the font used for default functions. It can be publicly
// modified to apply a default font to generated audios through def
// methods. If it is not modified, it is a font of zero filters.
var DefaultFont = font.New()

// A FontManager is a map of names to Fonts that has a built in
// default font at name 'def'.
type FontManager map[string]*font.Font

// NewFontManager returns a manager with a single 'def' font
func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)[defaultFontKey] = &font.Font{}
	return fm
}

// NewFont adds a font to a manger with the given keyed name.
// NewFont can return an error indicating if the name assigned
// was already in use.
func (fm *FontManager) NewFont(name string, f *font.Font) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = oakerr.ExistingElement{
			InputName:   name,
			InputType:   "font",
			Overwritten: true,
		}
	}
	manager[name] = f
	return err

}

// Get returns whatever is at name in font
func (fm *FontManager) Get(name string) *font.Font {
	manager := (*fm)
	font := manager[name]
	return font
}

func (fm *FontManager) GetDefault() *font.Font {
	return fm.Get(defaultFontKey)
}
