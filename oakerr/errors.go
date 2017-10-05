package oakerr

import "strconv"

//Todo 2.0: remove "Error" from the end of these types

// NotLoadedError is returned when something is queried that is not yet loaded.
type NotLoadedError struct{}

func (nle NotLoadedError) Error() string {
	return "File not loaded"
}

// ExistingFontError is returned when a font is overwritten in a font manager
type ExistingFontError struct{}

func (efe ExistingFontError) Error() string {
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
