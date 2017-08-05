package oakerr

//Todo 2.0: remove "Error" from the end of these types

// NotLoadedError is returned when something is queried that is not yet loaded.
type NotLoadedError struct{}

func (nle NotLoadedError) Error() string {
	return "File not loaded"
}

// ExistingFont is returned when a font is overwritten in a font manager
type ExistingFontError struct{}

func (efe ExistingFontError) Error() string {
	return "Font name already used, overwriting old font"
}
