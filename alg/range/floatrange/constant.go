package floatrange

// constant is a range that represents some constant float
type constant float64

// NewConstant returns a range that will always poll to return f
func NewConstant(f float64) Range {
	return constant(f)
}

// Poll returns the float behind the constant
func (c constant) Poll() float64 {
	return float64(c)
}

// Mult scales the constant by f
func (c constant) Mult(f float64) Range {
	c = constant(float64(c) * f)
	return c
}

// EnforceRange returns the float behind the constant
func (c constant) EnforceRange(float64) float64 {
	return float64(c)
}

// Percentile returns the float behind the constant
func (c constant) Percentile(float64) float64 {
	return float64(c)
}
