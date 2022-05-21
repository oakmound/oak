// Package grid provides structures for aligning grids of buttons
package grid

import (
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/entities/x/btn"
	"github.com/oakmound/oak/v4/scene"
)

// A Grid is a 2D slice of entities
type Grid [][]*entities.Entity

// A Generator defines the variables used to create grids from optional arguments
type Generator struct {
	Content    [][]btn.Option
	Defaults   btn.Option
	XGap, YGap float64
}

var (
	// A number of these fields could be removed, because they are the zero
	// value, but are left for documentation
	defaultGenerator = Generator{
		Content: [][]btn.Option{
			{
				nil,
			},
		},
		Defaults: nil,
		XGap:     0,
		YGap:     0,
	}
)

// Generate creates a Grid from a Generator
func (g *Generator) Generate(ctx *scene.Context) Grid {
	grid := make([][]*entities.Entity, len(g.Content))
	for x := 0; x < len(g.Content); x++ {
		grid[x] = make([]*entities.Entity, len(g.Content[x]))
		for y := 0; y < len(g.Content[x]); y++ {
			grid[x][y] = btn.New(ctx,
				g.Defaults,
				g.Content[x][y],
				btn.Offset(float64(x)*g.XGap, float64(y)*g.YGap),
			)
		}
	}
	return grid
}

// New creates a grid of buttons from a set of options
func New(ctx *scene.Context, opts ...Option) Grid {
	g := defaultGenerator
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		g = opt(g)
	}
	return g.Generate(ctx)
}
