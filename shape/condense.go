package shape

import "github.com/oakmound/oak/v2/alg/intgeom"

// Condense finds a set of rectangles that covers the shape.
// Used to return a minimal set of rectangles in an appropriate time.
func Condense(sh Shape, w, h int) []intgeom.Rect2 {
	condensed := []intgeom.Rect2{}
	remainingSpaces := make(map[intgeom.Point2]struct{})

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if sh.In(x, y, w, h) {
				remainingSpaces[intgeom.Point2{x, y}] = struct{}{}
			}
		}
	}

	for k := range remainingSpaces {
		topLeft := k
		w := 0
		h := 0
		right := true
		left := true
		up := true
		down := true
		xIncrement := intgeom.Point2{1, 0}
		yIncrement := intgeom.Point2{0, 1}
		xDecrement := intgeom.Point2{-1, 0}
		yDecrement := intgeom.Point2{0, -1}
		for right || left || up || down {
			var toCheck intgeom.Point2
			if right {
				toCheck = topLeft.Add(intgeom.Point2{w + 1, 0})
				for i := 0; i <= h; i++ {
					if _, ok := remainingSpaces[toCheck]; !ok {
						right = false
						break
					}
					toCheck = toCheck.Add(yIncrement)
				}
				if right {
					w++
				}
			}
			if left {
				toCheck = topLeft.Add(intgeom.Point2{-1, 0})
				for i := 0; i <= h; i++ {
					if _, ok := remainingSpaces[toCheck]; !ok {
						left = false
						break
					}
					toCheck = toCheck.Add(yIncrement)
				}
				if left {
					w++
					topLeft = topLeft.Add(xDecrement)
				}
			}
			if up {
				toCheck = topLeft.Add(intgeom.Point2{0, -1})
				for i := 0; i <= w; i++ {
					if _, ok := remainingSpaces[toCheck]; !ok {
						up = false
						break
					}
					toCheck = toCheck.Add(xIncrement)
				}
				if up {
					h++
					topLeft = topLeft.Add(yDecrement)
				}
			}
			if down {
				toCheck = topLeft.Add(intgeom.Point2{0, h + 1})
				for i := 0; i <= w; i++ {
					if _, ok := remainingSpaces[toCheck]; !ok {
						down = false
						break
					}
					toCheck = toCheck.Add(xIncrement)
				}
				if down {
					h++
				}
			}

		}
		condensed = append(condensed, intgeom.NewRect2WH(topLeft.X(), topLeft.Y(), w, h))
		for x := topLeft.X(); x <= topLeft.X()+w; x++ {
			for y := topLeft.Y(); y <= topLeft.Y()+h; y++ {
				delete(remainingSpaces, intgeom.Point2{x, y})
			}
		}
	}
	return condensed
}
