package render

import (
	"fmt"
)

// Functions in this file operate on the default font, and are equivalent to
// DefFont().Call. DefFont() does perform work to generate the default font,
// so storing the result and calling these functions on the stored Font is
// recommended in cases where performance is a concern.

// NewText creates a text element using the default font.
func NewText(str fmt.Stringer, x, y float64) *Text {
	return DefFont().NewText(str, x, y)
}

// NewIntText wraps the given int pointer in a stringer interface and creates
// a text renderable that will diplay the underlying int value.
func NewIntText(str *int, x, y float64) *Text {
	return DefFont().NewIntText(str, x, y)
}

// NewStrText is a helper to take in a string instead of a Stringer for NewText
func NewStrText(str string, x, y float64) *Text {
	return DefFont().NewStrText(str, x, y)
}
