package static

import (
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

type Slide struct {
	Rs          *render.Composite
	ContinueKey string
	PrevKey     string
	transition  scene.Transition
	cont        bool
	prev        bool
}

func (ss *Slide) Init() {
	render.Draw(ss.Rs, 0)
	event.GlobalBind(func(int, interface{}) int {
		ss.cont = true
		return 0
	}, "KeyUp"+ss.ContinueKey)
	event.GlobalBind(func(int, interface{}) int {
		ss.prev = true
		return 0
	}, "KeyUp"+ss.PrevKey)
}

func (ss *Slide) Continue() bool {
	return !ss.cont && !ss.prev
}

func (ss *Slide) Prev() bool {
	ret := ss.prev
	ss.prev = false
	ss.cont = false
	return ret
}

func (ss *Slide) Transition() scene.Transition {
	return ss.transition
}

func NewSlide(rs ...render.Modifiable) *Slide {
	return &Slide{
		Rs:          render.NewComposite(rs...),
		ContinueKey: "RightArrow",
		PrevKey:     "LeftArrow",
	}
}

func Transition(trans scene.Transition) SlideOption {
	return func(s *Slide) *Slide {
		s.transition = trans
		return s
	}
}

func Background(r render.Modifiable) SlideOption {
	return func(s *Slide) *Slide {
		s.Rs.Prepend(r)
		return s
	}
}

func ControlKeys(cont, prev string) SlideOption {
	return func(s *Slide) *Slide {
		s.ContinueKey = cont
		s.PrevKey = prev
		return s
	}
}

type SlideOption func(*Slide) *Slide

func NewSlideSet(n int, opts ...SlideOption) []*Slide {
	slides := make([]*Slide, n)
	for i := range slides {
		slides[i] = NewSlide()
		for _, opt := range opts {
			slides[i] = opt(slides[i])
		}
	}
	return slides
}
