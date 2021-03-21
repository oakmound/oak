package render

import (
	"sync"
	"testing"
)

func TestFont_UnsafeCopy(t *testing.T) {
	// Still thinking about if this behavior is correct
	initTestFont()
	f := DefFont()
	f.Unsafe = true
	f2 := f.Copy()
	if f2 != f {
		t.Fatalf("unsafe should have prevented the copy from actually copying")
	}
}

var initTestFontOnce sync.Once

// Todo: move this to font_test.go, once we have font_test.go
func initTestFont() {
	initTestFontOnce.Do(func() {
		DefFontGenerator = FontGenerator{RawFile: luxisrTTF}
		SetFontDefaults("", "", "", "", "white", "", 10, 10)
	})
}
