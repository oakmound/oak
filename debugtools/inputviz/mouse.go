package inputviz

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v3/alg/floatgeom"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

type Mouse struct {
	Rect      floatgeom.Rect2
	BaseLayer int

	event.CallerID
	ctx *scene.Context

	rs map[mouse.Button]*render.Switch

	lastMousePos *posStringer
	posText      *render.Text

	stateIncLock sync.RWMutex
	stateInc     map[mouse.Button]int
}

func (m *Mouse) Init() event.CallerID {
	m.CID = m.ctx.CallerMap.NextID(m)
	return m.CID
}

func (m *Mouse) RenderAndListen(ctx *scene.Context, layer int) error {
	m.ctx = ctx
	m.Init()

	if m.Rect.W() == 0 || m.Rect.H() == 0 {
		m.Rect.Max = m.Rect.Min.Add(floatgeom.Point2{60, 100})
	}
	w, h := m.Rect.W(), m.Rect.H()

	m.stateInc = make(map[mouse.Button]int)
	m.rs = make(map[mouse.Button]*render.Switch)

	const background mouse.Button = -1
	m.rs[background] = render.NewSwitch("default", map[string]render.Modifiable{
		"default": render.NewColorBox(int(w), int(h), color.RGBA{100, 100, 100, 255}),
	})
	m.rs[background].SetLayer(layer)
	m.rs[mouse.ButtonLeft] = render.NewSwitch("released", map[string]render.Modifiable{
		"released": render.NewColorBox(int(w/2), int(h/2), color.RGBA{150, 150, 150, 255}),
		"pressed":  render.NewColorBox(int(w/2), int(h/2), color.RGBA{240, 240, 240, 255}),
	})
	m.rs[mouse.ButtonLeft].SetLayer(layer + 1)
	m.rs[mouse.ButtonRight] = render.NewSwitch("released", map[string]render.Modifiable{
		"released": render.NewColorBox(int(w/2), int(h/2), color.RGBA{150, 150, 150, 255}),
		"pressed":  render.NewColorBox(int(w/2), int(h/2), color.RGBA{240, 240, 240, 255}),
	})
	m.rs[mouse.ButtonRight].SetPos(w/2, 0)
	m.rs[mouse.ButtonRight].SetLayer(layer + 1)

	scrollDown := render.NewColorBox(int(w/5), int(h/6), color.RGBA{250, 250, 250, 255})
	scrollDown.SetPos(0, h/6)
	m.rs[mouse.ButtonMiddle] = render.NewSwitch("released", map[string]render.Modifiable{
		"released": render.NewColorBox(int(w/5), int(h/3), color.RGBA{160, 160, 160, 255}),
		"pressed":  render.NewColorBox(int(w/5), int(h/3), color.RGBA{250, 250, 250, 255}),
		"scrollup": render.NewCompositeM(
			render.NewColorBox(int(w/5), int(h/3), color.RGBA{160, 160, 160, 255}),
			render.NewColorBox(int(w/5), int(h/6), color.RGBA{250, 250, 250, 255}),
		),
		"scrolldown": render.NewCompositeM(
			render.NewColorBox(int(w/5), int(h/3), color.RGBA{160, 160, 160, 255}),
			scrollDown,
		),
	})
	m.rs[mouse.ButtonMiddle].SetPos(2*(w/5), h/3)
	m.rs[mouse.ButtonMiddle].SetLayer(layer + 2)
	m.lastMousePos = &posStringer{}
	m.posText = render.NewStringerText(m.lastMousePos, 0, 0)
	_, textH := m.posText.GetDims()
	m.posText.SetY(h - float64(textH+1))

	for _, r := range m.rs {
		r.ShiftPos(m.Rect.Min.X(), m.Rect.Min.Y())
		if m.BaseLayer == -1 {
			ctx.DrawStack.Draw(r)
		} else {
			ctx.DrawStack.Draw(r, m.BaseLayer)
		}
	}
	if m.BaseLayer == -1 {
		ctx.DrawStack.Draw(m.posText, layer+2)
	} else {
		ctx.DrawStack.Draw(m.posText, m.BaseLayer, layer+2)
	}

	m.Bind(mouse.Press, mouse.Binding(func(id event.CallerID, ev *mouse.Event) int {
		m, _ := m.ctx.CallerMap.GetEntity(id).(*Mouse)
		m.rs[ev.Button].Set("pressed")
		m.stateIncLock.Lock()
		m.stateInc[ev.Button]++
		m.stateIncLock.Unlock()
		return 0
	}))
	m.Bind(mouse.Release, mouse.Binding(func(id event.CallerID, ev *mouse.Event) int {
		m, _ := m.ctx.CallerMap.GetEntity(id).(*Mouse)
		m.rs[ev.Button].Set("released")
		m.stateIncLock.Lock()
		m.stateInc[ev.Button]++
		m.stateIncLock.Unlock()
		return 0
	}))
	m.Bind(mouse.ScrollDown, mouse.Binding(func(id event.CallerID, e *mouse.Event) int {
		m, _ := m.ctx.CallerMap.GetEntity(id).(*Mouse)
		m.rs[mouse.ButtonMiddle].Set("scrolldown")
		m.stateIncLock.Lock()
		m.stateInc[mouse.ButtonMiddle]++
		st := m.stateInc[mouse.ButtonMiddle]
		m.stateIncLock.Unlock()
		m.ctx.DoAfter(100*time.Millisecond, func() {
			m.stateIncLock.Lock()
			if m.stateInc[mouse.ButtonMiddle] == st {
				m.rs[mouse.ButtonMiddle].Set("released")
			}
			m.stateIncLock.Unlock()
		})
		return 0
	}))
	m.Bind(mouse.ScrollUp, mouse.Binding(func(id event.CallerID, e *mouse.Event) int {
		m, _ := m.ctx.CallerMap.GetEntity(id).(*Mouse)
		m.rs[mouse.ButtonMiddle].Set("scrollup")
		m.stateIncLock.Lock()
		m.stateInc[mouse.ButtonMiddle]++
		st := m.stateInc[mouse.ButtonMiddle]
		m.stateIncLock.Unlock()
		m.ctx.DoAfter(100*time.Millisecond, func() {
			m.stateIncLock.Lock()
			if m.stateInc[mouse.ButtonMiddle] == st {
				m.rs[mouse.ButtonMiddle].Set("released")
			}
			m.stateIncLock.Unlock()
		})
		return 0
	}))
	m.Bind(mouse.Drag, mouse.Binding(func(id event.CallerID, e *mouse.Event) int {
		m, _ := m.ctx.CallerMap.GetEntity(id).(*Mouse)
		m.lastMousePos.Point2 = e.Point2
		return 0
	}))

	return nil
}

type posStringer struct {
	floatgeom.Point2
}

func (ps *posStringer) String() string {
	return fmt.Sprintf("(%d,%d)", int(ps.X()), int(ps.Y()))
}

func (m *Mouse) Destroy() {
	m.UnbindAll()
	for _, r := range m.rs {
		r.Undraw()
	}
	m.posText.Undraw()
}
