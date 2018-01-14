package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"

	"github.com/oakmound/oak"
)

var (
	Font              *render.Font
	text              *render.Text
	playerAlive       bool
	waitrand          float64
	i                 int
	exclamationPoints string
	end               = make(chan bool)
)

const (
	NONE = iota
	CLIFFRACER
)

type CliffRacer struct {
	entities.Moving
}

func (cr *CliffRacer) Init() event.CID {
	cid := event.NextID(cr)
	cr.CID = cid
	return cid
}

func NewCliffRacer(y float64) *CliffRacer {
	cr := new(CliffRacer)
	sp, err := render.LoadSprite("raw/cliffracer.png")
	if err != nil {
		dlog.Error(err)
		return nil
	}
	cr.Moving = entities.NewMoving(640, y, 80, 80, sp, nil, cr.Init(), 0)
	cr.Speed = physics.NewVector(rand.Float64()*10+3, rand.Float64()*4-2)
	render.Draw(cr.R, 100)
	cr.Space = collision.NewLabeledSpace(cr.X(), cr.Y(), 80, 80, CLIFFRACER)
	collision.Add(cr.Space)
	cr.CID.Bind(moveCliffRacer, "EnterFrame")
	cr.CID.Bind(func(id int, nothing interface{}) int {
		event.GetEntity(id).(*CliffRacer).Destroy()
		return 0
	}, "PlayerHit")
	return cr
}

func moveCliffRacer(id int, nothing interface{}) int {
	cr := event.GetEntity(id).(*CliffRacer)
	cr.ShiftX(-cr.Speed.X())
	cr.ShiftY(cr.Speed.Y())
	if cr.X() < -100 {
		cr.Destroy()
	}
	return 0
}

type Player struct {
	entities.Solid
}

func (p *Player) Init() event.CID {
	cid := event.NextID(p)
	p.CID = cid
	return cid
}

func NewPlayer() {
	p := new(Player)
	p.Solid = entities.NewSolid(50, 100, 10, 10, render.NewColorBox(10, 10, color.RGBA{255, 0, 0, 255}), nil, p.Init())
	render.Draw(p.R, 80)
	collision.Add(p.Space)
	p.CID.Bind(playerEnter, "EnterFrame")
}

func playerEnter(id int, nothing interface{}) int {
	p := event.GetEntity(id).(*Player)
	if oak.IsDown("W") {
		p.ShiftY(-5)
	} else if oak.IsDown("S") {
		p.ShiftY(5)
	}
	if oak.IsDown("A") {
		p.ShiftX(-5)
	} else if oak.IsDown("D") {
		p.ShiftX(5)
	}
	if p.X() < 0 {
		p.ShiftX(-1 * p.X())
	} else if p.X() > float64(oak.ScreenWidth-10) {
		p.ShiftX(-1 * (p.X() - float64(oak.ScreenWidth-10)))
	}
	if p.Y() < 0 {
		p.ShiftY(-1 * p.Y())
	} else if p.Y() > float64(oak.ScreenHeight-10) {
		p.ShiftY(-1 * (p.Y() - float64(oak.ScreenHeight-10)))
	}

	if collision.HitLabel(p.Space, CLIFFRACER) != nil {
		playerAlive = false
	}
	return 0
}

func main() {
	oak.LoadConf("oak.config")

	oak.Add("cliffRacers",
		func(prevScene string, data interface{}) {
			playerAlive = true
			bkg, err := render.LoadSprite("raw/background.png")
			if err != nil {
				dlog.Error(err)
				return
			}
			render.Draw(bkg, 1)
			text = render.DefFont().NewStrText("Dodge the Cliff Racers!", 70.0, 70.0)
			render.Draw(text, 60000)
			NewPlayer()
			waitrand = 5000.0
			i = 1
			exclamationPoints = ""
			go func() {
				for {
					select {
					case <-time.After(((time.Duration(rand.Intn(int(waitrand)))) * time.Millisecond) + 50*time.Millisecond):
						NewCliffRacer(float64(rand.Intn(200) + 50))
						i++
						if i%10 == 0 && waitrand > 400 {
							exclamationPoints += "!"
							text.SetString("Next Level" + exclamationPoints)
							waitrand *= .7
						}
					case <-end:
						return
					}
				}
			}()
		}, func() bool {
			return playerAlive
		}, func() (string, *scene.Result) {
			end <- true
			return "cliffRacers", nil
		})
	oak.Init("cliffRacers")
}
