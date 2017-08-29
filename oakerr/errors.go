package oakerr

import "strconv"

// NotLoaded is returned when something is queried that is not yet loaded.
type NotLoaded struct{}

func (nle NotLoaded) Error() string {
	return "File not loaded"
}

// ExistingFont is returned when a font is overwritten in a font manager
type ExistingFont struct{}

func (efe ExistingFont) Error() string {
	return "Font name already used, overwriting old font"
}

// InsufficientInputs is returned when something requires at least some number
// of inputs in a variadic argument, but that minimum was not supplied.
type InsufficientInputs struct {
	AtLeast   int
	InputName string
}

func (ii InsufficientInputs) Error() string {
	return "Must supply at least " + strconv.Itoa(ii.AtLeast) + " " + ii.InputName
}
