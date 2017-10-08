package main

import (
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

	sslides := static.NewSlideSet(6,
		static.Background(bkg),
		static.Transition(scene.Fade(5, 10)),
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

	sslides[intro+2].Append(TxtAt(Gnuolane72, "Things I Made", .5, .3))
	sslides[intro+2].Append() // screenshots

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
			"-",
		)...,
	)

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
