package show

import (
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

type StaticSlide struct {
	Rs          *render.Composite
	ContinueKey string
	PrevKey     string
	Transition  scene.Transition
	cont        bool
	prev        bool
}

func (ss *StaticSlide) Init() {
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

func (ss *StaticSlide) Continue() bool {
	return !ss.cont && !ss.prev
}

func (ss *StaticSlide) Prev() bool {
	ret := ss.prev
	ss.prev = false
	ss.cont = false
	return ret
}

func (ss *StaticSlide) Result() *scene.Result {
	if ss.Transition == nil {
		return nil
	}
	return &scene.Result{
		Transition: ss.Transition,
	}
}

func NewStaticSlide(continueKey, prevKey string, rs ...render.Modifiable) *StaticSlide {
	return &StaticSlide{
		Rs:          render.NewComposite(rs...),
		ContinueKey: continueKey,
		PrevKey:     prevKey,
	}
}

func (ss *StaticSlide) WithTransition(trans scene.Transition) *StaticSlide {
	ss.Transition = trans
	return ss
}
