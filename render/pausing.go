package render

// CanPause types have pause functions to start and stop animation
type CanPause interface {
	Pause()
	Unpause()
}

type pauseBool struct {
	playing bool
}

func (p *pauseBool) Pause() {
	p.playing = false
}

func (p *pauseBool) Unpause() {
	p.playing = true
}
