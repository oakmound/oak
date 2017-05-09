package audio

import "bitbucket.org/oakmoundstudio/oak/oakerr"

type FontManager map[string]*Font

func NewFontManager() *FontManager {
	fm := &FontManager{}
	(*fm)["def"] = &Font{}
	return fm
}

func (fm *FontManager) NewFont(name string, f *Font) error {
	manager := (*fm)
	var err error
	if _, ok := manager[name]; ok {
		err = oakerr.ExistingFontError{}
	}
	manager[name] = f
	return err

}

func (fm *FontManager) Get(elementName string) *Font {
	manager := (*fm)
	font, _ := manager[elementName]
	return font
}
