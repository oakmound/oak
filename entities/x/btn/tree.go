package btn

// AddChildren adds a generator to create a child btn
func AddChildren(cg ...Generator) Option {
	return func(g Generator) Generator {
		g.Children = append(g.Children, cg...)
		return g
	}
}
