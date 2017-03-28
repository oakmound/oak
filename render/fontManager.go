package render

import "errors"

type FontManager map[string]*Font

func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)["def"] = (&FontGenerator{}).Generate()
	return fm
}

func (fm *FontManager) NewFont(name string, fg FontGenerator) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = errors.New("Font already existed, overwriting it now.")
	}
	manager[name] = (&fg).Generate()
	return err

}

func (fm *FontManager) Get(elementName string) *Font {
	manager := (*fm)
	font, _ := manager[elementName]
	return font
}
