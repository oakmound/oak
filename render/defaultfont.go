package render

import (
	"fmt"
)

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
