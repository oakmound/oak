package span

import "image/color"

type linearColor struct {
	r, g, b, a Span[uint32]
}

// NewLinearColor returns a linear color distribution between min and maxColor
func NewLinearColor(minColor, maxColor color.Color) Span[color.Color] {
	r, g, b, a := minColor.RGBA()
	r2, g2, b2, a2 := maxColor.RGBA()
	return linearColor{
		NewLinear(r, r2),
		NewLinear(g, g2),
		NewLinear(b, b2),
		NewLinear(a, a2),
	}
}

func (l linearColor) Clamp(c color.Color) color.Color {
	r3, g3, b3, a3 := c.RGBA()
	r4 := l.r.Clamp(r3)
	g4 := l.g.Clamp(g3)
	b4 := l.b.Clamp(b3)
	a4 := l.a.Clamp(a3)
	return rgbaFromInts(r4, g4, b4, a4)
}

func (l linearColor) MulSpan(i float64) Span[color.Color] {
	return linearColor{
		l.r.MulSpan(i),
		l.g.MulSpan(i),
		l.b.MulSpan(i),
		l.a.MulSpan(i),
	}
}

func (l linearColor) Poll() color.Color {
	r3 := l.r.Poll()
	g3 := l.g.Poll()
	b3 := l.b.Poll()
	a3 := l.a.Poll()
	return rgbaFromInts(r3, g3, b3, a3)
}

func (l linearColor) Percentile(f float64) color.Color {
	r3 := l.r.Percentile(f)
	g3 := l.g.Percentile(f)
	b3 := l.b.Percentile(f)
	a3 := l.a.Percentile(f)
	return rgbaFromInts(r3, g3, b3, a3)
}

func rgbaFromInts(r, g, b, a uint32) color.RGBA {
	return color.RGBA{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}
