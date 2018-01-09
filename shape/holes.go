package shape

import (
	"github.com/oakmound/oak/alg/intgeom"
)

func GetHoles(sh Shape, w, h int) [][]intgeom.Point {

	flooding := make(map[intgeom.Point]bool)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if !sh.In(x, y) {
				flooding[intgeom.Point{x, y}] = true
			}
		}
	}

	border := BorderPoints(w, h)

	for _, p := range border {
		if !sh.In(p.X, p.Y) {
			BFSFlood(flooding, p)
		}
	}

	// flooding is now a map of holes, points which are false
	// but not on the border.

	out := make([][]intgeom.Point, 0)

	for len(flooding) > 0 {
		for k := range flooding {
			out = append(out, BFSFlood(flooding, k))
		}
	}

	return out
}

func BorderPoints(w, h int) []intgeom.Point {
	out := make([]intgeom.Point, (w*2+h*2)-4)
	i := 0
	for x := 0; x < w; x++ {
		out[i] = intgeom.Point{x, 0}
		out[i+1] = intgeom.Point{x, h - 1}
		i += 2
	}
	for y := 1; y < h-1; y++ {
		out[i] = intgeom.Point{0, y}
		out[i+1] = intgeom.Point{w - 1, y}
		i += 2
	}
	return out
}

func BFSFlood(m map[intgeom.Point]bool, start intgeom.Point) []intgeom.Point {
	visited := []intgeom.Point{}
	toVisit := []intgeom.Point{start}
	for len(toVisit) > 0 {
		next := toVisit[0]
		delete(m, next)
		toVisit = toVisit[1:]
		visited = append(visited, next)

		// literally adjacent points for adjacency
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				p := intgeom.Point{x + next.X, y + next.Y}
				if _, ok := m[p]; ok {
					toVisit = append(toVisit, p)
				}
			}
		}
	}

	return visited
}
