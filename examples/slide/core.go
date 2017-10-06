package main

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
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
)

func main() {

	bkg := render.NewColorBox(width, height, colornames.Firebrick)

	sslides := static.NewSlideSet(3,
		static.Background(bkg),
		static.Transition(scene.Fade(5, 10)),
	)

	sslides[0].Rs.Append(Express44.NewStrText("Slide 1", 100, 100).ToSprite())
	sslides[1].Rs.Append(Express44.NewStrText("Slide 2", 100, 100).ToSprite())
	sslides[2].Rs.Append(Express44.NewStrText("Slide 3", 100, 100).ToSprite())

	oak.SetupConfig.Screen = oak.Screen{
		Width:  width,
		Height: height,
	}

	slides := make([]show.Slide, len(sslides))
	for i, s := range sslides {
		slides[i] = s
	}

	show.AddSlides(slides...)
	show.Start()
}
