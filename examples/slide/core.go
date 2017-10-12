package main

import (
	"fmt"

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

	RLibel28 = show.FontColor(colornames.Blue)(Libel28)

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

	setups := []slideSetup{
		intro,
		why,
		philo,
		particles,
		ai,
		levels,
		conclusion,
	}

	total := 0

	for _, setup := range setups {
		total += setup.len
	}

	fmt.Println("Total slides", total)

	sslides := static.NewSlideSet(total,
		static.Background(bkg),
		static.Transition(scene.Fade(4, 12)),
	)

	nextStart := 0

	for _, setup := range setups {
		setup.add(nextStart, sslides)
		nextStart += setup.len
	}

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

type slideSetup struct {
	add func(int, []*static.Slide)
	len int
}

var (
	intro = slideSetup{
		addIntro,
		5,
	}
)

func addIntro(i int, sslides []*static.Slide) {
	// Intro: three slides
	// Title
	sslides[i].Append(
		Title("Applying Go to Game Programming"),
		TxtAt(Gnuolane44, "Patrick Stephen", .5, .6),
	)
	// Thanks everybody for coming to this talk. I'm going to be talking about
	// design patterns, philosophies, and generally useful tricks for
	// developing video games in Go.

	// Me
	sslides[i+1].Append(Header("Who Am I"))
	sslides[i+1].Append(
		TxtSetAt(Gnuolane44, 0.5, 0.63, 0.0, 0.07,
			"Graduate Student at University of Minnesota",
			"Maintainer / Programmer of Oak",
			"github.com/200sc  github.com/oakmound/oak",
			"patrick.d.stephen@gmail.com",
			"oakmoundstudio@gmail.com",
		)...,
	)
	// My name is Patrick Stephen, I'm currently a Master's student at
	// the University of Minnesota. I'm one of two primary maintainers
	// of oak's source code, Oak being the game engine that we built
	// to make our games with.
	// If you have any questions that don't get answered in or after
	// this talk, feel free to send those questions either to me
	// personally or to our team's email, or if it applies, feel free
	// to raise an issue on the repository.

	sslides[i+2].Append(Header("Games I Made"))
	sslides[i+2].Append(TxtAt(Gnuolane28, "White = Me, Blue = Oakmound", .5, .24))
	sslides[i+2].Append(
		ImageCaption("botanist.PNG", .67, .1, .5, Libel28, "Space Botanist"),
		ImageCaption("agent.PNG", .1, .11, .85, RLibel28, "Agent Blue"),
		ImageCaption("dyscrasia.PNG", .5, .65, .5, RLibel28, "Dyscrasia"),
		ImageCaption("esque.PNG", .4, .37, .5, RLibel28, "Esque"),
		ImageCaption("fantastic.PNG", .33, .65, .5, RLibel28, "A Fantastic Doctor"),
		ImageCaption("flower.PNG", .7, .41, .75, Libel28, "Flower Son"),
		ImageCaption("jeremy.PNG", .07, .5, .66, Libel28, "Jeremy The Clam"),
		ImageCaption("wolf.PNG", .68, .71, .5, Libel28, "The Wolf Comes Out At 18:00"),
	)
	// These are games that I've made in the past, most being made
	// for game jams, built in somewhere between 2 days and 2 weeks.
	//
	// We'll mostly be focusing on these three games, which are those
	// that we've been working on in Go-- Agent Blue, Jeremy the Clam,
	// and A Fantastic Doctor.

	sslides[i+3].Append(Header("This Talk is Not About..."))
	sslides[i+3].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Optimizing Go",
			"- 3D Graphics in Go",
			"- Mobile Games in Go",
		)...,
	)

	// And just to get this out of the way, as you will probably
	// note from the games I just showed, we aren't going to be
	// talking about 3D games here or really performance intensive
	// games, or games for non-desktop platforms, just because,
	// while we haven't ignored these things I don't have
	// any revolutionary breakthroughs to share about them right now.

	sslides[i+4].Append(Header("Topics"))
	sslides[i+4].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Why Go",
			"- Design Philosophy",
			"- Particles",
			"- AI with Interfaces",
			"- Level Building with Interfaces",
		)...,
	)

	// What we will talk about, is why Go is particular useful for
	// developing games, the philosophy behind our engine and development
	// strategy, and then some interesting use cases for applying
	// design patterns that Go makes easy with particle generation,
	// artificial intelligence, and constructing levels.
}

var (
	why = slideSetup{
		addWhy,
		3,
	}
)

func addWhy(i int, sslides []*static.Slide) {
	// Topic: Why Go
	sslides[i].Append(Title("Why Go"))
	sslides[i+1].Append(Header("Why Go"))
	sslides[i+1].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Concurrency",
			"- Fast Development",
			"- Scales Well",
		)...,
	)
	sslides[i+2].Append(Header("Why Not Go"))
	sslides[i+2].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Difficult to use Graphics Cards",
			"- Difficult to vectorize instructions",
			"- C is Unavoidable",
		)...,
	)
}

var (
	philo = slideSetup{
		addPhilo,
		5,
	}
)

func addPhilo(i int, sslides []*static.Slide) {
	// Philosophy, engine discussion
	sslides[i].Append(Title("Design Philosophy"))
	sslides[i+1].Append(Header("Design Philosophy"))
	sslides[i+1].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- No non-Go dependencies",
			"- Ease of API",
			"- If it's useful and generic, put it in the engine",
		)...,
	)
	sslides[i+2].Append(Header("Update Loops and Functions"))
	sslides[i+2].Append(
		Image("updateCode1.PNG", .27, .4),
		Image("updateCode3.PNG", .57, .4),
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
	sslides[i+3].Append(Header("Update Loops and Functions"))
	sslides[i+3].Append(
		Image("updateCode2.PNG", .27, .4),
		Image("updateCode3.PNG", .57, .4),
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
	sslides[i+4].Append(Header("Useful Packages"))
	sslides[i+4].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- oak/alg",
			"- oak/intgeom, oak/floatgeom",
			"- oak/physics",
			"- oak/render/particle",
		)...,
	)
}

var (
	particles = slideSetup{
		addParticles,
		5,
	}
)

func addParticles(i int, sslides []*static.Slide) {
	sslides[i].Append(Title("Particles"))
	sslides[i+1].Append(Header("Particles in CraftyJS"))
	sslides[i+2].Append(Header("Variadic Functional Options"))
	sslides[i+3].Append(Header("Particle Generators in Oak"))
	sslides[i+4].Append(Header("Particle Generators in Oak"))
}

var (
	ai = slideSetup{
		addAI,
		7,
	}
)

func addAI(i int, sslides []*static.Slide) {
	sslides[i].Append(Title("Building AI with Interfaces"))
	sslides[i+1].Append(Header("Building AI with Interfaces"))
	sslides[i+2].Append(Header("Storing Small Interface Types"))
	sslides[i+3].Append(Header("Storing Small Interface Types"))
	sslides[i+4].Append(Header("When Your Interface is Massive"))
	sslides[i+5].Append(Header("Condensing Massive Interfaces"))
	sslides[i+6].Append(Header("... And you've got reusable AI"))
}

var (
	levels = slideSetup{
		addLevels,
		8,
	}
)

func addLevels(i int, sslides []*static.Slide) {
	sslides[i].Append(Title("Designing Levels with Interfaces"))
	sslides[i+1].Append(Header("Components of a Level"))
	sslides[i+2].Append(Header("An Approach without Interfaces"))
	sslides[i+3].Append(Header("An Approach without Interfaces"))
	sslides[i+4].Append(Header("Level Interfaces: Attempt 1"))
	sslides[i+5].Append(Header("Level Interfaces: Attempt 1"))
	sslides[i+6].Append(Header("Level Interfaces: Attempt 2"))
	sslides[i+7].Append(Header("Level Interfaces: Attempt 2"))
}

var (
	conclusion = slideSetup{
		addConclusion,
		2,
	}
)

func addConclusion(i int, sslides []*static.Slide) {
	sslides[i].Append(Header("Thanks To"))
	sslides[i].Append(
		TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Nate Fudenberg, John Ficklin",
			"- Contributors on Github",
			"- You, Audience",
		)...,
	)
	sslides[i+1].Append(Title("Questions"))
}
