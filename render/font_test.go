package render

import (
	"sync"
	"testing"
)

func TestFont_UnsafeCopy(t *testing.T) {
	// Still thinking about if this behavior is correct
	initTestFont()
	f := DefaultFont()
	f.Unsafe = true
	f2 := f.Copy()
	if f2 != f {
		t.Fatalf("unsafe should have prevented the copy from actually copying")
	}
}

var initTestFontOnce sync.Once

func initTestFont() {
	initTestFontOnce.Do(func() {
		DefFontGenerator = FontGenerator{RawFile: luxisrTTF}
		SetFontDefaults("", "", "", "", "white", "", 10, 10)
	})
}
