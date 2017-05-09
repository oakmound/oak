package oakerr

// NotLoadedError is returned when something is queried that is not yet loaded.
type NotLoadedError struct{}

func (nle NotLoadedError) Error() string {
	return "File not loaded"
}

// ExistingFontError is returned when a font is overwritten in a font manager
type ExistingFontError struct{}

func (efe ExistingFontError) Error() string {
	return "Font already existed, overwriting it now"
}
