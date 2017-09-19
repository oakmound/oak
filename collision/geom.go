// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import (
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
)

// minDist computes the square of the distance from a point to a rectangle.
// If the point is contained in the rectangle then the distance is zero.
//
// Implemented per Definition 2 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func minDist(p floatgeom.Point3, r floatgeom.Rect3) float64 {
	sum := 0.0
	for i, pi := range p {
		if pi < r.Min[i] {
			d := pi - r.Min[i]
			sum += d * d
		} else if pi > r.Max[i] {
			d := pi - r.Max[i]
			sum += d * d
		} else {
			sum += 0
		}
	}
	return sum
}

// minMaxDist computes the minimum of the maximum distances from p to points
// on r.  If r is the bounding box of some geometric objects, then there is
// at least one object contained in r within minMaxDist(p, r) of p.
//
// Implemented per Definition 4 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func minMaxDist(p floatgeom.Point3, r floatgeom.Rect3) float64 {
	// by definition, MinMaxDist(p, r) =
	// min{1<=k<=n}(|pk - rmk|^2 + sum{1<=i<=n, i != k}(|pi - rMi|^2))
	// where rmk and rMk are defined as follows:

	rm := func(k int) float64 {
		if p[k] <= (r.Min[k]+r.Max[k])/2 {
			return r.Min[k]
		}
		return r.Max[k]
	}

	rM := func(k int) float64 {
		if p[k] >= (r.Min[k]+r.Max[k])/2 {
			return r.Min[k]
		}
		return r.Max[k]
	}

	// This formula can be computed in linear time by precomputing
	// S = sum{1<=i<=n}(|pi - rMi|^2).

	S := 0.0
	for i := range p {
		d := p[i] - rM(i)
		S += d * d
	}

	// Compute MinMaxDist using the precomputed S.
	min := math.MaxFloat64
	for k := range p {
		d1 := p[k] - rM(k)
		d2 := p[k] - rm(k)
		d := S - d1*d1 + d2*d2
		if d < min {
			min = d
		}
	}

	return min
}

// boundingBox constructs the smallest rectangle containing both r1 and r2.
func boundingBox(r1, r2 floatgeom.Rect3) floatgeom.Rect3 {
	return r1.GreaterOf(r2)
}

// boundingBoxN constructs the smallest rectangle containing all of r...
func boundingBoxN(rects ...floatgeom.Rect3) (bb floatgeom.Rect3) {
	if len(rects) == 1 {
		return rects[0]
	}
	bb = boundingBox(rects[0], rects[1])
	for _, rect := range rects[2:] {
		bb = boundingBox(bb, rect)
	}
	return
}
