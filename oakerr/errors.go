package oakerr

var (
	_ error = NotFound{}
	_ error = ExistingElement{}
	_ error = InsufficientInputs{}
	_ error = UnsupportedFormat{}
	_ error = NilInput{}
	_ error = IndivisibleInput{}
	_ error = InvalidInput{}
	_ error = UnsupportedPlatform{}
)

// NotFound is returned when some input was queried but not found.
type NotFound struct {
	InputName string
}

func (nf NotFound) Error() string {
	return errorString(codeNotFound, nf.InputName)
}

// ExistingElement is an alternative to ExistingFont, where in this case the
// existing element is -not- overwritten.
type ExistingElement struct {
	InputName   string
	InputType   string
	Overwritten bool
}

func (ee ExistingElement) Error() string {
	if ee.Overwritten {
		return errorString(codeExistingElementOverwritten, ee.InputName, ee.InputType)
	} else {
		return errorString(codeExistingElement, ee.InputName, ee.InputType)
	}
}

// InsufficientInputs is returned when something requires at least some number
// of inputs in a variadic argument, but that minimum was not supplied.
type InsufficientInputs struct {
	AtLeast   int
	InputName string
}

func (ii InsufficientInputs) Error() string {
	return errorString(codeInsufficientInputs, ii.AtLeast, ii.InputName)
}

// UnsupportedFormat is returned by functions expecting formatted data or
// files which received a format they can't use.
type UnsupportedFormat struct {
	Format string
}

func (uf UnsupportedFormat) Error() string {
	return errorString(codeUnsupportedFormat, uf.Format)
}

// NilInput is returned from functions expecting a non-nil pointer which
// receive a nil pointer.
type NilInput struct {
	InputName string
}

func (ni NilInput) Error() string {
	return errorString(codeNilInput, ni.InputName)
}

// IndivisibleInput is returned from functions expecting a count of inputs
// in a slice or variadic argument divisible by some integer, or an integer
// value divisible by some integer. IsList represents which input type was expected.
type IndivisibleInput struct {
	InputName    string
	MustDivideBy int
}

func (ii IndivisibleInput) Error() string {
	return errorString(codeIndivisibleInput, ii.InputName, ii.MustDivideBy)
}

// InvalidInput is a generic struct returned for otherwise invalid input.
type InvalidInput struct {
	InputName string
}

func (ii InvalidInput) Error() string {
	return errorString(codeInvalidInput, ii.InputName)
}

// UnsupportedPlatform is returned when functionality isn't supported
// on the hardware or operating system used.
type UnsupportedPlatform struct {
	Operation string
}

func (up UnsupportedPlatform) Error() string {
	return errorString(codeUnsupportedPlatform, up.Operation)
}
