package grid

import (
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/entities/x/btn"
)

// An Option modifies a generator prior to grid generation
type Option func(Generator) Generator

// Content sets the button option to create for each
// button at x,y coordinates on this grid. This should
// not be used with Width and Height. Options set here act
// like And() when used with Defaults.
func Content(content [][]btn.Option) Option {
	return func(g Generator) Generator {
		g.Content = content
		return g
	}
}

// ContentAt sets the button uption to create for a button
// at a given x,y coordinate on this grid. If the grid has
// already had content defined by Width, Height, or Content,
// and the x,y value given would not fall on the defined grid,
// the grids dimensions will be expanded so that it will.
// Negative x or y values will result in this option having
// no effect.
// ContentAt will overwrite options set by Content. It acts
// like And() when used with Defaults.
func ContentAt(x, y int, opts ...btn.Option) Option {
	return func(g Generator) Generator {
		opt := btn.And(opts...)
		if x < 0 {
			dlog.Error("ContentAt option created with <0 x")
			return g
		}
		if y < 0 {
			dlog.Error("ContentAt option created with <0 y")
			return g
		}
		if len(g.Content) <= x {
			delta := (len(g.Content) - x) + 1
			for i := 0; i < delta; i++ {
				g.Content = append(g.Content, make([]btn.Option, 0))
			}
		}
		if len(g.Content[x]) <= y {
			delta := (len(g.Content[x]) - y) + 1
			newOpts := make([]btn.Option, delta)
			g.Content[x] = append(g.Content[x], newOpts...)
		}
		g.Content[x][y] = opt
		return g
	}
}

// Height sets the number of buttons vertically that this grid will make.
func Height(h int) Option {
	return func(g Generator) Generator {
		for x := 0; x < len(g.Content); x++ {
			if len(g.Content[x]) <= h {
				delta := (len(g.Content[x]) - h) + 1
				opts := make([]btn.Option, delta)
				g.Content[x] = append(g.Content[x], opts...)
			}
		}
		return g
	}
}

// Width sets the number of buttons horizontally that this grid will make.
func Width(w int) Option {
	return func(g Generator) Generator {
		if len(g.Content) <= w {
			delta := (len(g.Content) - w) + 1
			height := 1
			if len(g.Content) > 0 {
				height = len(g.Content[0])
			}
			for i := 0; i < delta; i++ {
				g.Content = append(g.Content, make([]btn.Option, height))
			}
		}
		return g
	}
}

// Defaults sets the starting option used to create buttons in this grid.
func Defaults(defaults btn.Option) Option {
	return func(g Generator) Generator {
		g.Defaults = defaults
		return g
	}
}

// YGap sets the gap between buttons vertically on this grid.
func YGap(gap float64) Option {
	return func(g Generator) Generator {
		g.YGap = gap
		return g
	}
}

// XGap sets the gap between buttons horizontally on this grid.
func XGap(gap float64) Option {
	return func(g Generator) Generator {
		g.XGap = gap
		return g
	}
}

// Todo: combine grids?
// So could have a grid of color definitions,
// and it with a grid of something else...
// Todo also: row or grid permutations?
// Todo also: ContentAnd
