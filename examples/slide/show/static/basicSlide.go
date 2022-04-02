package static

import (
	"fmt"
	"os"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

type Slide struct {
	Rs          *render.CompositeR
	ContinueKey key.Code
	PrevKey     key.Code
	transition  scene.Transition
	cont        bool
	prev        bool
	OnClick     func()
}

func (ss *Slide) Init(ctx *scene.Context) {
	oak.SetFullScreen(true)
	render.Draw(ss.Rs, 0)

	event.GlobalBind(ctx, key.Up(ss.ContinueKey), func(key.Event) event.Response {

		fmt.Println("continue key pressed")
		ss.cont = true
		return 0
	})

	event.GlobalBind(ctx, key.Up(ss.PrevKey), func(key.Event) event.Response {
		fmt.Println("prev key pressed")
		ss.prev = true
		return 0
	})

	event.GlobalBind(ctx, key.Up(key.Escape), func(key.Event) event.Response {
		os.Exit(0)
		return 0
	})
	if ss.OnClick != nil {
		event.GlobalBind(ctx, mouse.Press, func(*mouse.Event) event.Response {
			ss.OnClick()
			return 0
		})
	}
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

func (ss *Slide) Append(rs ...render.Renderable) {
	for _, r := range rs {
		ss.Rs.Append(r)
	}
}

func (ss *Slide) Transition() scene.Transition {
	return ss.transition
}

func NewSlide(rs ...render.Renderable) *Slide {
	return &Slide{
		Rs:          render.NewCompositeR(rs...),
		ContinueKey: key.RightArrow,
		PrevKey:     key.LeftArrow,
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

func ControlKeys(cont, prev key.Code) SlideOption {
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
