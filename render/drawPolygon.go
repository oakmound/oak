package render

import (
	clip "github.com/akavel/polyclip-go"
)

var (
	usingDrawPolygon = false
	drawPolygon      clip.Polygon
)

func SetDrawPolygon(p clip.Polygon) {
	usingDrawPolygon = true
	drawPolygon = p
}

func DrawPolygonDim() (int, int, int, int) {
	if !usingDrawPolygon {
		return 0, 0, 0, 0
	}
	mbr := drawPolygon.BoundingBox()
	return int(mbr.Min.X), int(mbr.Min.Y), int(mbr.Max.X), int(mbr.Max.Y)
}

func InDrawPolygon(xi, yi, x2i, y2i int) bool {
	if usingDrawPolygon {
		x := float64(xi)
		y := float64(yi)
		x2 := float64(x2i)
		y2 := float64(y2i)
		p2 := clip.Polygon{{{X: x, Y: y}, {X: x, Y: y2}, {X: x2, Y: y2}, {X: x2, Y: y}}}
		intsct := drawPolygon.Construct(clip.INTERSECTION, p2)
		return len(intsct) != 0
	}
	return true
}
