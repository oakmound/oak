package render

import (
	// embed is used here to embed our default font file.
	_ "embed"
	"fmt"
)

// Oak ships with a free font to enable text display without needing to set up
// a font on the user's machine to import. This is that font. It is embedded into
// the Go code to ensure it is not stripped from the code by vendoring, for example.
// The file is called luxisr.ttf.

//go:embed luxisr.ttf
var luxisrTTF []byte

// Functions in this file operate on the default font, and are equivalent to
// DefaultFont().Call. DefaultFont() does perform work to generate the default font,
// so storing the result and calling these functions on the stored Font is
// recommended in cases where performance is a concern.

// NewStringerText creates a text element using the default font and a stringer.
func NewStringerText(str fmt.Stringer, x, y float64) *Text {
	return DefaultFont().NewStringerText(str, x, y)
}

// NewIntText wraps the given int pointer in a stringer interface and creates
// a text renderable that will diplay the underlying int value.
func NewIntText(str *int, x, y float64) *Text {
	return DefaultFont().NewIntText(str, x, y)
}

// NewText is a helper to create a text element with the default font and a string.
func NewText(str string, x, y float64) *Text {
	return DefaultFont().NewText(str, x, y)
}

// NewStrPtrText is a helper to take in a string pointer for NewText
func NewStrPtrText(str *string, x, y float64) *Text {
	return DefaultFont().NewStrPtrText(str, x, y)
}
