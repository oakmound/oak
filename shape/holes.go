package shape

import (
	"github.com/oakmound/oak/v2/alg/intgeom"
)

// GetHoles finds sets of points which are not In this shape that
// are adjacent.
func GetHoles(sh Shape, w, h int) [][]intgeom.Point2 {
	return getHoles(sh, w, h, false)
}

// GetBorderHoles finds sets of points which are not In this shape that
// are adjacent in addition to the space around the shape
// (ie points that border the shape)
func GetBorderHoles(sh Shape, w, h int) [][]intgeom.Point2 {
	return getHoles(sh, w, h, true)
}

// getHoles is an internal function that finds sets of points which are not In this shape that
// are adjacent.
func getHoles(sh Shape, w, h int, includeBorder bool) [][]intgeom.Point2 {
	flooding := make(map[intgeom.Point2]bool)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if !sh.In(x, y) {
				flooding[intgeom.Point2{x, y}] = true
			}
		}
	}
	if !includeBorder {
		border := borderPoints(w, h)
		for _, p := range border {
			if !sh.In(p.X(), p.Y()) {
				bfsFlood(flooding, p)
			}
		}
		// flooding is now a map of holes, points which are false
		// but not on the border.
	}

	out := make([][]intgeom.Point2, 0)

	for len(flooding) > 0 {
		for k := range flooding {
			out = append(out, bfsFlood(flooding, k))
		}
	}

	return out
}

func borderPoints(w, h int) []intgeom.Point2 {
	out := make([]intgeom.Point2, (w*2+h*2)-4)
	i := 0
	for x := 0; x < w; x++ {
		out[i] = intgeom.Point2{x, 0}
		out[i+1] = intgeom.Point2{x, h - 1}
		i += 2
	}
	for y := 1; y < h-1; y++ {
		out[i] = intgeom.Point2{0, y}
		out[i+1] = intgeom.Point2{w - 1, y}
		i += 2
	}
	return out
}

func bfsFlood(m map[intgeom.Point2]bool, start intgeom.Point2) []intgeom.Point2 {
	visited := []intgeom.Point2{}
	toVisit := map[intgeom.Point2]bool{start: true}

	for len(toVisit) > 0 {
		for next := range toVisit {
			delete(m, next)
			delete(toVisit, next)
			visited = append(visited, next)
			// literally adjacent points for adjacency
			for x := -1; x <= 1; x++ {
				for y := -1; y <= 1; y++ {
					p := intgeom.Point2{x + next.X(), y + next.Y()}
					if _, ok := m[p]; ok {
						toVisit[p] = true
					}
				}
			}
		}
	}

	return visited
}
