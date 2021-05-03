package intrange

// NewConstant returns a range which will always return the input constant
func NewConstant(i int) Range {
	return constant(i)
}

// constant implements Range as a poll
// which always returns the same integer.
type constant int

// Poll returns c cast to an int
func (c constant) Poll() int {
	return int(c)
}

// Mult returns this range scaled by i
func (c constant) Mult(i float64) Range {
	return constant(int(float64(int(c)) * i))
}

// EnforceRange on a constant must return the constant
func (c constant) EnforceRange(int) int {
	return int(c)
}

// Percentile can only return the constant itself
func (c constant) Percentile(float64) int {
	return int(c)
}
