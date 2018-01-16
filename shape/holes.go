package shape

import (
	"github.com/oakmound/oak/alg/intgeom"
)

func GetHoles(sh Shape, w, h int) [][]intgeom.Point2 {

	flooding := make(map[intgeom.Point2]bool)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if !sh.In(x, y) {
				flooding[intgeom.Point2{x, y}] = true
			}
		}
	}

	border := BorderPoints(w, h)

	for _, p := range border {
		if !sh.In(p.X(), p.Y()) {
			BFSFlood(flooding, p)
		}
	}

	// flooding is now a map of holes, points which are false
	// but not on the border.

	out := make([][]intgeom.Point2, 0)

	for len(flooding) > 0 {
		for k := range flooding {
			out = append(out, BFSFlood(flooding, k))
		}
	}

	return out
}

func BorderPoints(w, h int) []intgeom.Point2 {
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

func BFSFlood(m map[intgeom.Point2]bool, start intgeom.Point2) []intgeom.Point2 {
	visited := []intgeom.Point2{}
	toVisit := []intgeom.Point2{start}
	for len(toVisit) > 0 {
		next := toVisit[0]
		delete(m, next)
		toVisit = toVisit[1:]
		visited = append(visited, next)

		// literally adjacent points for adjacency
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				p := intgeom.Point2{x + next.X(), y + next.Y()}
				if _, ok := m[p]; ok {
					toVisit = append(toVisit, p)
				}
			}
		}
	}

	return visited
}
