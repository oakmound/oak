package oakerr

import "strconv"

// Todo: add language switches to initialization to change what the errors return

// The goal of putting structs here instead of returning errors.New(string)s
// is to be able to easily recognize error types through checks on the consuming
// side, and to be able to translate errors into other languages in a localized
// area.

// NotLoaded is returned when something is queried that is not yet loaded.
type NotLoaded struct{}

func (NotLoaded) Error() string {
	return "File not loaded"
}

// ExistingElement is an alternative to ExistingFont, where in this case the
// existing element is -not- overwritten.
type ExistingElement struct {
	InputName   string
	InputType   string
	Overwritten bool
}

func (ee ExistingElement) Error() string {
	s := ee.InputName + " " + ee.InputType + " already defined"
	if ee.Overwritten {
		s += ", old " + ee.InputType + " overwritten."
	} else {
		s += ", nothing overwritten."
	}
	return s
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

// UnsupportedFormat is returned by functions expecting formatted data or
// files which received a format they can't use.
type UnsupportedFormat struct {
	Format string
}

func (uf UnsupportedFormat) Error() string {
	return "Unsupported Format: " + uf.Format
}

// NilInput is returned from functions expecting a non-nil pointer which
// receive a nil pointer.
type NilInput struct {
	InputName string
}

func (ni NilInput) Error() string {
	return "Expected a non-nil pointer for input: " + ni.InputName
}

// IndivisibleInput is returned from functions expecting a count of inputs
// in a slice or variadic argument divisible by some integer, or an integer
// value divisible by some integer. IsList represents which input type was expected.
type IndivisibleInput struct {
	InputName    string
	MustDivideBy int
	IsList       bool
}

func (ii IndivisibleInput) Error() string {
	s := "Input " + ii.InputName
	if ii.IsList {
		s += " length"
	}
	return s + " was not divisible by " + strconv.Itoa(ii.MustDivideBy)
}

// Todo: compose InvalidInput into other invalid input esque structs, add
// constructors.

// InvalidInput is a generic struct returned for otherwise invalid input.
type InvalidInput struct {
	InputName string
}

func (ii InvalidInput) Error() string {
	return "Invalid input: " + ii.InputName
}

// InvalidLength is returned when some input has an explicit required length
// that was not provided.
type InvalidLength struct {
	InputName      string
	Length         int
	RequiredLength int
}

func (il InvalidLength) Error() string {
	return "Invalid input length for " + il.InputName +
		". Was " + strconv.Itoa(il.Length) +
		", expected " + strconv.Itoa(il.Length)
}

// ConsError is returned by specific functions that can coalesce errors
// over a series of inputs.
type ConsError struct {
	First, Second error
}

func (ce ConsError) Error() string {
	return ce.First.Error() + "; " + ce.Second.Error()
}

// UnsupportedPlatform is returned when functionality isn't supported
// on the hardware or operating system used.
type UnsupportedPlatform struct {
	Operation string
}

func (up UnsupportedPlatform) Error() string {
	return up.Operation + " is not supported on this platform/OS"
}
