package screen

// Todo: flesh out options as a generator style constructor
// with (optional) variadic arguments

// A WindowGenerator can generate windows based on various new window settings.
type WindowGenerator struct {
	// Width and Height specify the dimensions of the new window. If Width
	// or Height are zero, a driver-dependent default will be used for each
	// zero value dimension.
	Width, Height int

	// Title specifies the window title.
	Title string

	// Fullscreen determines whether the new window will be fullscreen or not.
	Fullscreen bool

	// Borderless determines whether the new window will have borders or not
	Borderless bool

	// TopMost determines whether the new window will stay on top of other windows
	// even when out of focus.
	TopMost bool

	// NoScaling determines whether the new window will have scaling allowed.
	// With a zero value of false, scaling is allowed.
	NoScaling bool

	// X and Y determine the location the new window should be created at. If
	// either are zero, a driver-dependant default will be used for each zero
	// value. If Fullscreen is true, these values will be ignored.
	X, Y int32
}

// A WindowOption is any function that sets up a WindowGenerator.
type WindowOption func(*WindowGenerator)

// Title sets a sanitized form of the input string. In particular, its length will
// not exceed 4096, and it may be further truncated so that it is valid UTF-8
// and will not contain the NUL byte.
func Title(s string) WindowOption {
	return func(g *WindowGenerator) {
		g.Title = sanitizeUTF8(s, 4096)
	}
}

// Dimensions sets the width and height of new windows
func Dimensions(w, h int) WindowOption {
	return func(g *WindowGenerator) {
		g.Width = w
		g.Height = h
	}
}

// Position sets the starting position of the new window
func Position(x, y int32) WindowOption {
	return func(g *WindowGenerator) {
		g.X = x
		g.Y = y
	}
}

// Fullscreen sets the starting fullscreen boolean of the new window
func Fullscreen(on bool) WindowOption {
	return func(g *WindowGenerator) {
		g.Fullscreen = on
	}
}

// Borderless sets the starting borderless boolean of the new window
// defaults to borders on.
func Borderless(on bool) WindowOption {
	return func(g *WindowGenerator) {
		g.Borderless = on
	}
}

// TopMost sets the starting topmost boolean of the new window, determining
// whether the window should appear above other windows even when unfocused.
func TopMost(on bool) WindowOption {
	return func(g *WindowGenerator) {
		g.TopMost = on
	}
}

// NewWindowGenerator creates a window generator with zero values,
// then calls all options passed in on it.
func NewWindowGenerator(opts ...WindowOption) WindowGenerator {
	wg := &WindowGenerator{}
	for _, o := range opts {
		o(wg)
	}
	return *wg
}
