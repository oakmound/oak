package main

import (
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/scene"
)

const (
	minX = 0
	minY = 0
	maxX = 578
	maxY = 416
)

var cache *render.Switch

func main() {
	oak.Add(
		"demo",
		func(string, interface{}) {
			layer := 0
			layerTxt := render.DefFont().NewIntText(&layer, 30, 20)
			layerTxt.SetLayer(100000000)
			render.Draw(layerTxt, 0)
			NewGopher(layer)
			layer++
			event.GlobalBind(func(int, interface{}) int {
				if oak.IsDown("K") {
					NewGopher(layer)
					layer++
				}
				return 0
			}, "EnterFrame")
			// Generate a rotation cache for comparison
			// Compare the use of the cache against the use of a reverting type below
			cache = render.NewSwitch("0", make(map[string]render.Modifiable))
			for i := 0; i < 360; i++ {
				s, err := render.LoadSprite(filepath.Join("raw", "gopher11.png"))
				if err != nil {
					dlog.Error(err)
					return
				}
				s.Modify(mod.Rotate(float32(i)))
				cache.Add(strconv.Itoa(i), s)
			}
		},
		func() bool {
			return true
		},
		func() (string, *scene.Result) {
			return "demo", nil
		},
	)

	render.SetDrawStack(
		render.NewHeap(false),
		render.NewDrawFPS(),
		render.NewLogicFPS(),
	)

	oak.SetupConfig.Screen.X = 1
	oak.SetupConfig.Screen.Y = 1
	oak.SetupConfig.LoadBuiltinCommands = true

	oak.Init("demo")
}

type Gopher struct {
	entities.Doodad
	deltaX, deltaY float64
	rotation       int
}

func (g *Gopher) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func NewGopher(layer int) {
	goph := Gopher{}
	goph.Doodad = entities.NewDoodad(
		rand.Float64()*576,
		rand.Float64()*416,
		render.NewSwitch("goph", map[string]render.Modifiable{"goph": render.EmptyRenderable()}),
		//render.NewReverting(render.LoadSprite(filepath.Join("raw", "gopher11.png"))),
		goph.Init())
	goph.R.SetLayer(layer)
	goph.Bind(gophEnter, "EnterFrame")
	goph.deltaX = 4 * float64(rand.Intn(2)*2-1)
	goph.deltaY = 4 * float64(rand.Intn(2)*2-1)
	goph.rotation = rand.Intn(360)
	render.Draw(goph.R, 0)
}

func gophEnter(cid int, nothing interface{}) int {
	goph := event.GetEntity(cid).(*Gopher)

	// Compare against this version of rotation
	// (also swap the comments on lines in goph.Doodad's renderable)
	//goph.R.(*render.Reverting).RevertAndModify(1, render.Rotate(goph.rotation))
	goph.R.(*render.Switch).Add("goph", render.NewSprite(0, 0, cache.GetSub(strconv.Itoa(goph.rotation)).GetRGBA()))
	if goph.X() < minX || goph.X() > maxX {
		goph.deltaX *= -1
	}
	if goph.Y() < minY || goph.Y() > maxY {
		goph.deltaY *= -1
	}
	goph.SetPos(goph.deltaX+goph.X(), goph.deltaY+goph.Y())
	goph.rotation = (goph.rotation + 1) % 360
	return 0
}
