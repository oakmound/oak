package main

import (
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"golang.org/x/image/colornames"
)

func main() {
	slides := []show.Slide{
		show.NewStaticSlide(
			"RightArrow",
			"LeftArrow",
			render.NewColorBox(640, 480, colornames.Darkblue),
			show.Express28.NewStrText("Slide 1", 100, 100).ToSprite(),
		).WithTransition(scene.Fade(5, 10)),
		show.NewStaticSlide(
			"RightArrow",
			"LeftArrow",
			render.NewColorBox(640, 480, colornames.Darkblue),
			show.Express28.NewStrText("Slide 2", 100, 100).ToSprite(),
		).WithTransition(scene.Fade(5, 10)),
		show.NewStaticSlide(
			"RightArrow",
			"LeftArrow",
			render.NewColorBox(640, 480, colornames.Darkblue),
			show.Express28.NewStrText("Slide 3", 100, 100).ToSprite(),
		).WithTransition(scene.Fade(5, 10)),
	}
	show.AddSlides(slides...)
	show.Start()
}
