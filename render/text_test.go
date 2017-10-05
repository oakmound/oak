package render

import (
	"image"
	"path/filepath"
	"testing"

	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render/internal/testdata/fonts"
	"github.com/stretchr/testify/assert"
)

func TestTextFns(t *testing.T) {
	initTestFont()

	txt := DefFont().NewStrText("Test", 0, 0)

	fg := FontGenerator{
		File:  filepath.Join("default_assets", "font", "luxisr.ttf"),
		Color: image.Black,
	}

	f := fg.Generate()

	txt.SetFont(f)
	assert.Equal(t, f, txt.d)

	txt.SetString("Test2")
	assert.Equal(t, "Test2", txt.text.String())

	n := 100
	txt.SetIntP(&n)

	n = 50
	assert.Equal(t, "50", txt.text.String())

	txt.SetInt(n + 1)
	assert.Equal(t, "51", txt.text.String())

	txt.SetText(dummyStringer{})
	assert.Equal(t, "Dummy", txt.text.String())
	assert.Equal(t, "Text[Dummy]", txt.String())

	txts := txt.Wrap(1, 10)
	assert.Equal(t, 5, len(txts))

	for i, wrap := range txts {
		assert.Equal(t, float64(i*10), wrap.Y())
	}

	txts = txt.Wrap(3, 10)
	assert.Equal(t, 2, len(txts))

	assert.NotNil(t, txt.ToSprite())

	txt.Center()
	assert.Equal(t, float64(-2), txt.X())

}

type dummyStringer struct{}

func (d dummyStringer) String() string {
	return "Dummy"
}

// Todo: move this to font_test.go, once we have font_test.go
func initTestFont() {
	DefFontGenerator = FontGenerator{File: filepath.Join("default_assets", "font", "luxisr.ttf")}
	fileutil.BindataDir = fonts.AssetDir
	fileutil.BindataFn = fonts.Asset
	SetFontDefaults("", "", "", "", "white", "", 10, 10)
}
