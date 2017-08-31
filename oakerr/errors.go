package oakerr

import "strconv"

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
