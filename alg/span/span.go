package span

// A Span represents some enumerable range.
type Span[T any] interface {
	// Poll returns a pseudorandom value within this span.
	Poll() T
	// Clamp, if v lies within the boundary of this span, returns v.
	// Otherwise, CLamp returns a modified version of v that is rounded to the closest value
	// that does lie within the boundary of this span.
	Clamp(v T) T
	// Percentile returns the value along this span that is at the provided percentile through the span,
	// e.g. providing .5 will return the middle of the span, providing 1 will return the maximum value in
	// the span. Providing a value less than 0 or greater than 1 may extend the span by where it would theoretically
	// progress, but should not be relied upon unless a given implementation specifies what it will do. If this span
	// represents multiple degrees of freedom, this will pin all those degrees to the single provided percent.
	Percentile(float64) T
	// MulSpan returns this span with its entire range multiplied by the given constant.
	MulSpan(float64) Span[T]
}
