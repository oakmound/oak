package btn

// Group links several btns together
type Group struct {
	members []Btn
	active  Btn
}

// GetActive returns the active btn from the group
func (g *Group) GetActive() Btn {
	return g.active
}
