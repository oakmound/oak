package main

import (
	"fmt"
	"path/filepath"

	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/oak/shape"
	"golang.org/x/image/colornames"
)

const (
	width  = 1920
	height = 1080
)

var (
	Express28  = show.FontSize(28)(show.Express)
	Gnuolane28 = show.FontSize(28)(show.Gnuolane)
	Libel28    = show.FontSize(28)(show.Libel)

	Express44  = show.FontSize(44)(show.Express)
	Gnuolane44 = show.FontSize(44)(show.Gnuolane)
	Libel44    = show.FontSize(44)(show.Libel)

	Express72  = show.FontSize(72)(show.Express)
	Gnuolane72 = show.FontSize(72)(show.Gnuolane)
	Libel72    = show.FontSize(72)(show.Libel)
)

func main() {

	bz1, _ := shape.BezierCurve(
		width/15, height/5,
		width/15, height/15,
		width/5, height/15)

	bz2, _ := shape.BezierCurve(
		width-(width/15), height/5,
		width-(width/15), height/15,
		width-(width/5), height/15)

	bz3, _ := shape.BezierCurve(
		width/15, height-(height/5),
		width/15, height-(height/15),
		width/5, height-(height/15))

	bz4, _ := shape.BezierCurve(
		width-(width/15), height-(height/5),
		width-(width/15), height-(height/15),
		width-(width/5), height-(height/15))

	bkg := render.NewComposite(
		render.NewColorBox(width, height, colornames.Seagreen),
		render.BezierThickLine(bz1, colornames.White, 1),
		render.BezierThickLine(bz2, colornames.White, 1),
		render.BezierThickLine(bz3, colornames.White, 1),
		render.BezierThickLine(bz4, colornames.White, 1),
	)

	sslides := static.NewSlideSet(12,
		static.Background(bkg),
		static.Transition(scene.Fade(4, 12)),
	)

	intro := 0

	// Intro: three slides
	// Title
	sslides[intro].Append(
		TxtAt(Gnuolane72, "Applying Go to Game Programming", .5, .4),
		TxtAt(Gnuolane44, "Patrick Stephen", .5, .6),
	)
	// Me
	sslides[intro+1].Append(TxtAt(Gnuolane72, "Who Am I", .5, .3))
	sslides[intro+1].Append(
		TxtSetAt(Gnuolane44, .5, .63, 0, .07,
			"Graduate Student at University of Minnesota",
			"Maintainer / Programmer of Oak",
			"github.com/200sc  github.com/oakmound/oak",
			"patrick.d.stephen@gmail.com",
		)...,
	)

	sslides[intro+2].Append(TxtAt(Gnuolane72, "Games I Made", .5, .3))
	sslides[intro+2].Append(
		Image("botanist.PNG", .1, .5).Modify(mod.Scale(.5, .5)),
		Image("agent.PNG", .1, .11).Modify(mod.Scale(.75, .75)),
		Image("dyscrasia.PNG", .33, .65).Modify(mod.Scale(.5, .5)),
		Image("esque.PNG", .4, .4).Modify(mod.Scale(.5, .5)),
		Image("fantastic.PNG", .5, .65).Modify(mod.Scale(.5, .5)),
		Image("flower.PNG", .7, .4).Modify(mod.Scale(.75, .75)),
		Image("jeremy.PNG", .7, .1).Modify(mod.Scale(.5, .5)),
		Image("wolf.PNG", .7, .7).Modify(mod.Scale(.5, .5)),
	) // screenshots

	// What I'm going to talk about
	sslides[intro+3].Append(TxtAt(Gnuolane72, "Topics", .5, .2))
	sslides[intro+3].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Why Go",
			"- Design Philosophy",
			"- Particles",
			"- Modifying Renderables",
			"- AI with Interfaces",
			"- Level Building with Interfaces",
		)...,
	)

	whyGo := 4
	// Topic: Why Go
	sslides[whyGo].Append(TxtAt(Gnuolane72, "Why Go", .5, .4))
	sslides[whyGo+1].Append(TxtAt(Gnuolane72, "Why Go", .5, .2))
	sslides[whyGo+1].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Concurrency",
			"- Fast Development",
			"- Scales Well",
		)...,
	)
	sslides[whyGo+2].Append(TxtAt(Gnuolane72, "Why Not Go", .5, .2))
	sslides[whyGo+2].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Difficult to use Graphics Cards",
			"- Difficult to vectorize instructions",
			"- C is Unavoidable",
		)...,
	)

	// Philosophy, engine discussion
	philosophy := 7
	sslides[philosophy].Append(TxtAt(Gnuolane72, "Design Philosophy", .5, .4))
	sslides[philosophy+1].Append(TxtAt(Gnuolane72, "Design Philosophy", .5, .2))
	sslides[philosophy+1].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- No non-Go dependencies",
			"- Ease of API",
			"- If it's useful and generic, put it in the engine",
		)...,
	)
	sslides[philosophy+2].Append(TxtAt(Gnuolane72, "Update Loops and Functions", .5, .2))
	sslides[philosophy+2].Append(
		Image("updateCode1.PNG", .3, .4),
		Image("updateCode3.PNG", .6, .4),
	)
	//
	// Some game engines model their exposed API as a loop--
	// stick all your logic inside update()
	//
	// In larger projects, this leads directly to an explicit splitting up of that
	// loop into at least two parts-- update all entities, then
	// draw all entities.
	//
	// The combining of these elements into one loop causes
	// a major problem-- tying the rate at which entities update themselves
	// to the rate at which entities are drawn. This leads to inflexible
	// engines, and in large projects you'll have to do something to work around
	// this, or if you hard lock your draw rate modders will post funny videos
	// of your physics breaking when they try to fix your frame rate.
	//
	// Oak handles this loop for you, and splits it into two loops, one for
	// drawing elements and one for logical frame updating.
	//
	sslides[philosophy+3].Append(TxtAt(Gnuolane72, "Update Loops and Functions", .5, .2))
	sslides[philosophy+3].Append(
		Image("updateCode2.PNG", .3, .5),
		Image("updateCode3.PNG", .6, .5),
	)
	//
	// Another pattern used, in parallel with the Update Loop,
	// is the Update Function. Give every entity in your game the
	// Upate() function, and then your game logic is handled by calling Update()
	// on everything. At a glance, this works very well in Go because your entities
	// all fit into this single-function interface, but in games with a lot of
	// entities you'll end up with a lot of entities that don't need to do
	// anything on each frame.
	//
	// The engine needs to provide a way to handle game objects that don't
	// need to be updated as well as those that do, and separating these into
	// two groups explicitly makes the engine less extensible. Oak uses an
	// event handler for this instead, where each entity that wants to use
	// an update function binds that function to their entity id once.
	//
	sslides[philosophy+4].Append(TxtAt(Gnuolane72, "Useful Packages", .5, .2))

	oak.SetupConfig.Screen = oak.Screen{
		Width:  width,
		Height: height,
	}

	slides := make([]show.Slide, len(sslides))
	for i, s := range sslides {
		slides[i] = s
	}
	show.AddNumberShortcuts(len(slides))
	show.Start(slides...)
}

func TxtSetAt(f *render.Font, xpos, ypos, xadv, yadv float64, txts ...string) []render.Renderable {
	rs := make([]render.Renderable, len(txts))
	for i, txt := range txts {
		rs[i] = TxtAt(f, txt, xpos, ypos)
		xpos += xadv
		ypos += yadv
	}
	return rs
}

func TxtSetFrom(f *render.Font, xpos, ypos, xadv, yadv float64, txts ...string) []render.Renderable {
	rs := make([]render.Renderable, len(txts))
	for i, txt := range txts {
		rs[i] = TxtFrom(f, txt, xpos, ypos)
		xpos += xadv
		ypos += yadv
	}
	return rs
}

func TxtAt(f *render.Font, txt string, xpos, ypos float64) render.Renderable {
	return Pos(f.NewStrText(txt, 0, 0), xpos, ypos)
}

func Pos(r render.Renderable, xpos, ypos float64) render.Renderable {
	XPos(r, xpos)
	YPos(r, ypos)
	return r
}

func XPos(r render.Renderable, pos float64) render.Renderable {
	w, _ := r.GetDims()
	r.SetPos(width*pos-float64(w/2), r.Y())
	return r
}

func YPos(r render.Renderable, pos float64) render.Renderable {
	_, h := r.GetDims()
	r.SetPos(r.X(), height*pos-float64(h/2))
	return r
}

func TxtFrom(f *render.Font, txt string, xpos, ypos float64) render.Renderable {
	return f.NewStrText(txt, width*xpos, height*ypos)
}

func Image(file string, xpos, ypos float64) render.Modifiable {
	s, err := render.LoadSprite(filepath.Join("raw", file))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	s.SetPos(width*xpos, height*ypos)
	return s
}
