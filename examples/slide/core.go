package main

import (
	"fmt"
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/render/particle"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/examples/slide/show"
	"github.com/oakmound/oak/examples/slide/show/static"
	"github.com/oakmound/oak/render"
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

	show.SetDims(width, height)
	show.SetTitleFont(Gnuolane72)

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

	oak.LoadingR = bkg

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
		//static.Transition(scene.Fade(4, 12)),
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
	oak.SetupConfig.FrameRate = 30
	oak.SetupConfig.DrawFrameRate = 30

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
		show.Title("Applying Go to Game Programming"),
		show.TxtAt(Gnuolane44, "Patrick Stephen", .5, .6),
	)
	// Thanks everybody for coming to this talk. I'm going to be talking about
	// design patterns, philosophies, and generally useful tricks for
	// developing video games in Go.

	sslides[i+1].Append(show.Header("Who Am I"))
	sslides[i+1].Append(
		show.TxtSetAt(Gnuolane44, 0.5, 0.63, 0.0, 0.07,
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

	sslides[i+2].Append(show.Header("Games I Made"))
	sslides[i+2].Append(show.TxtAt(Gnuolane28, "White = Me, Blue = Oakmound", .5, .24))
	sslides[i+2].Append(
		show.ImageCaption("botanist.PNG", .67, .1, .5, Libel28, "Space Botanist"),
		show.ImageCaption("agent.PNG", .1, .11, .85, RLibel28, "Agent Blue"),
		show.ImageCaption("dyscrasia.PNG", .5, .65, .5, RLibel28, "Dyscrasia"),
		show.ImageCaption("esque.PNG", .4, .37, .5, RLibel28, "Esque"),
		show.ImageCaption("fantastic.PNG", .33, .65, .5, RLibel28, "A Fantastic Doctor"),
		show.ImageCaption("flower.PNG", .7, .41, .75, Libel28, "Flower Son"),
		show.ImageCaption("jeremy.PNG", .07, .5, .66, Libel28, "Jeremy The Clam"),
		show.ImageCaption("wolf.PNG", .68, .71, .5, Libel28, "The Wolf Comes Out At 18:00"),
	)
	// These are games that I've made in the past, most being made
	// for game jams, built in somewhere between 2 days and 2 weeks.
	//
	// We'll mostly be focusing on these three games, which are those
	// that we've been working on in Go-- Agent Blue, Jeremy the Clam,
	// and A Fantastic Doctor.

	sslides[i+3].Append(show.Header("This Talk is Not About..."))
	sslides[i+3].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
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

	sslides[i+4].Append(show.Header("Topics"))
	sslides[i+4].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Why Go",
			"- Design Philosophy",
			"- Particles",
			"- AI with Interfaces",
			"- Level Building with Interfaces",
			"- Other Examples",
		)...,
	)

	// What we will talk about, is why Go is particularly useful for
	// developing games, the philosophy behind our engine and development
	// strategy, and then some interesting use cases for applying
	// design patterns that Go makes easy with particle generation,
	// artificial intelligence, and level construction.
}

var (
	why = slideSetup{
		addWhy,
		3,
	}
)

func addWhy(i int, sslides []*static.Slide) {
	sslides[i].Append(show.Title("Why Go"))
	sslides[i+1].Append(show.Header("Why Go"))
	sslides[i+1].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Concurrency",
			"- Fast Development",
			"- Scales Well",
			"- Multiplatform Support",
		)...,
	)

	// So Go is particularly nice for building games one the one hand
	// for its speed-- If you're used to building games with javascript
	// or pygame, you'll have way more cpu cycles than you know how to
	// deal with, especially if you use concurrency well on machines with
	// multiple CPUs, which is going to be most of your audience.
	//
	// More importantly, Go is just as fast to develop with as those slower
	// languages but it scales so much better. A little effort into decoupling
	// your components with interfaces, and your code becomes far easier to read
	// and increment on.

	sslides[i+2].Append(show.Header("Why Not Go"))
	sslides[i+2].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Execution Speed",
			"- Difficult to use Graphics Cards",
			"- Difficult to vectorize instructions",
			"- C is Unavoidable",
		)...,
	)

	// But what I've said so far mostly applies to game jam style games--
	// how do you make a quick and dirty game in a few days without your
	// code falling all over itself. If you're interested in doing something
	// with heavy performance requirements, Go isn't the language to use.
	// Go's speed isn't good enough for AAA games because it doesn't have easy
	// access to things like OpenGL, Vulkan, or SIMD CPU instructions.
	// What Go can do with these things is call out to C to do the work for it,
	// but every C call in Go has overhead, and that overhead adds up if you're
	// calling out to it thousands of times per second.
	//
	// There's other practical issues if you want to develop in Go even if you
	// don't have high performance requirements-- depending on your platform
	// you may need to install audio dependencies, usb dependencies, and so on,
	// and for all of Go's benefits in cross compilation these dependencies
	// completely break the hope of your game working the same on multiple
	// platforms without you going in and testing it manually.
}

var (
	philo = slideSetup{
		addPhilo,
		7,
	}
)

func addPhilo(i int, sslides []*static.Slide) {
	// Philosophy, engine discussion
	sslides[i].Append(show.Title("Design Philosophy"))
	sslides[i+1].Append(show.Header("Design Philosophy"))
	sslides[i+1].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- No non-Go dependencies",
			"- Ease / Terseness of API",
			"- If it's useful and generic, put it in the engine",
		)...,
	)

	// That brings us to our design philosphy in Oak.
	// First, if we have a non-Go dependency, we also have an issue to
	// replace that non-Go dependency ASAP. Right now we have just one.
	//
	// The motivation for having few dependencies isn't just so we can
	// feel confident that all of our platforms are supported, but also
	// making the engine easy to use. Most libraries in Go can be used
	// with 'go get', and we want the same thing here-- a developer
	// should be able to go get oak and immediately start working with it.
	//
	// After that, we want our API to be easy to use and small. Part of our
	// motivation to start building Oak was that other game engines at the
	// time took 500 lines to draw a cube or 400 lines to write Pong. Maybe
	// at their core, those problems do take that many lines, but a lot of that
	// code we can write for you (and also for us, so we don't have to keep
	// re-writing it).
	//
	// In line with this, we follow a rule where if we have to rewrite something
	// more than once for a game or for a package of the engine, that probably
	// means that should be its own package and feature the engine provides.
	// This does go against the go proverbs-- we do not follow the idea that
	// a little copying is better than a little dependency, so long as we
	// treat that dependency as part of the larger, engine dependency.

	sslides[i+2].Append(show.Header("Update Loops and Functions"))
	sslides[i+2].Append(
		show.Image("updateCode1.PNG", .27, .4),
		show.Image("updateCode3.PNG", .57, .4),
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
	sslides[i+3].Append(show.Header("Update Loops and Functions"))
	sslides[i+3].Append(
		show.Image("updateCode2.PNG", .27, .4),
		show.Image("updateCode3.PNG", .57, .4),
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
	sslides[i+4].Append(show.Header("Useful Packages"))
	sslides[i+4].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- oak/alg/intgeom, oak/alg/floatgeom",
			"- oak/alg",
			"- oak/physics",
			"- oak/render/particle",
		)...,
	)
	//
	// These are some of the less obvious useful packages we've taken
	// from games or sub-packages and built into their own package--
	//
	// intgeom and floatgeom should be self explanatory-- we and every
	// other Go package continually redefine X,Y and X,Y,Z points of
	// integers and floats, and we needed to stop redoing that work.
	//
	sslides[i+5].Append(show.Header("oak/alg"))
	sslides[i+5].Append(show.ImageAt("degToRad.PNG", .3, .5, mod.Scale(1.25, 1.25)))
	sslides[i+5].Append(show.ImageAt("chooseX.PNG", .6, .5, mod.Scale(1.25, 1.25)))
	//
	// in alg, we store things like rounding and selection algorithms.
	// We found that we really needed to pick a random element from
	// a list of weighted floats a lot, so we split it off here.
	//
	sslides[i+6].Append(show.Header("oak/physics"))
	sslides[i+6].Append(show.ImageAt("push.PNG", .3, .5, mod.Scale(1.25, 1.25)))
	sslides[i+6].Append(show.ImageAt("attachable.PNG", .7, .5, mod.Scale(1.25, 1.25)))
	//
	// Physics was built to store some physics primitives for handling
	// propagation of forces, mass, friction, but was mostly built so
	// we could attach entities to each other and stop having to move
	// every sub-component in an entity when we moved the entity.
	//
	// And lastly, particle, where we figured being able to generate
	// a lot of small images or colors in patterns was something that could easily
	// spice up most games.
}

var (
	particles = slideSetup{
		addParticles,
		5,
	}
)

func addParticles(i int, sslides []*static.Slide) {
	sslides[i].Append(show.Title("Particles"))
	sslides[i].OnClick = func() {
		go particle.NewColorGenerator(
			particle.Size(intrange.Constant(4)),
			particle.EndSize(intrange.Constant(7)),
			particle.Angle(floatrange.NewLinear(0, 359)),
			particle.Pos(width/2, height/2),
			particle.Speed(floatrange.NewSpread(5, 2)),
			particle.NewPerFrame(floatrange.NewSpread(5, 5)),
			particle.Color(
				color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0},
				color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 0},
			),
		).Generate(0)
	}
	//
	// Speaking of particles, that's our first example
	// of applying some techniques Go provides for making this API something I
	// would consider to be really special.
	//
	// A particle generator is something like what's showing on screen now--
	// a source of a bunch of colors or effects or images, and they're complex
	// to implement only because of the vast number of options you can take in
	// to a particle emitter.
	sslides[i+1].Append(show.Header("Particles in Other Engines"))
	sslides[i+1].Append(
		show.ImageCaption("craftyParticle.PNG", .17, .3, 1.25, Libel28, "CraftyJS"),
		show.ImageCaption("phaserParticle.PNG", .40, .3, 1.25, Libel28, "PhaserJS"),
	)
	//
	// For context, we'll look at how some other engines do their
	// particle APIs. Before starting Oak we worked with CraftyJS,
	// which has the nice feature that these giant blocks of settings
	// can be stored and reused for new particles, but then you get
	// giant settings blocks.
	//
	// Phaser uses the reverse approach-- you can't keep particle settings
	// around but you don't need to set a bunch of settings you don't need.
	//
	// These examples aren't making the same particle emitter, by the way,
	// they're just the first examples I found from the respective engine's
	// documentation.

	sslides[i+2].Append(show.Header("Particle Generators in Oak"))
	sslides[i+2].Append(show.Image("AndPt.PNG", .13, .59).Modify(mod.Scale(1.25, 1.25)))
	sslides[i+2].Append(show.Image("colorGen.PNG", .13, .29).Modify(mod.Scale(1.25, 1.25)))
	sslides[i+2].Append(show.Image("particleOpts.PNG", .53, .29).Modify(mod.Scale(1.25, 1.25)))
	//
	// We wanted to apply what crafty did with saving settings, but we wanted
	// settings to not all be mandatory, so our functional pattern starts by
	// setting a bunch of defaults, then applying all the options that are passed in.
	// Because the Option type is an exported type, users can define their own settings,
	// and one of the settings I like to define is the And helper, shown here.
	//
	sslides[i+3].Append(show.Header("Particle Generators in Oak"))
	sslides[i+3].Append(show.Image("oakParticle.PNG", .27, .4).Modify(mod.Scale(1.25, 1.25)))
	//
	// This is what this code looks like in practice-- often
	// particle effects will get thrown off to their own package in our games
	// so we can use a shorter import for particle.
	//
	sslides[i+4].Append(show.Header("Aside: Filtering Audio with Functional Options"))
	sslides[i+4].Append(show.Image("filter.PNG", .13, .29).Modify(mod.Scale(1.25, 1.25)))
	sslides[i+4].Append(show.Image("audio.PNG", .53, .29).Modify(mod.Scale(1.25, 1.25)))
	sslides[i+4].Append(show.Image("dataFilter.PNG", .13, .59).Modify(mod.Scale(1.25, 1.25)))
	//
	// On the implementation side, though, if you have multiple types of
	// particle generators, it's really frustrating to have to define interfaces
	// for them accepting a whole bunch of different kinds of settings or not.
	// While we haven't refactored particles to use this approach yet, our
	// audio library fixes this by defining all filters on audio as functions
	// that define their own Apply function--  so the logic for whether or not
	// a particle type supports a setting can be confined to the type of filter.
}

var (
	ai = slideSetup{
		addAI,
		5,
	}
)

func addAI(i int, sslides []*static.Slide) {
	sslides[i].Append(show.Title("Building AI with Interfaces"))

	sslides[i+1].Append(show.Header("When Your Interface is Massive"))
	sslides[i+1].Append(show.ImageAt("agentAI.PNG", .4, .5))
	sslides[i+1].Append(show.ImageAt("agentCharacter.PNG", .7, .5))

	// But now that you've stored all of your enemy types as themselves,
	// if you've got a bunch of procedures that run on your AI for
	// pathing, targetting, or attacking, you'll run into this issue
	// where the interface that defines your AI needs to know a lot
	// of different infromation for each of these different behaviors.

	sslides[i+2].Append(show.Header("Condensing Massive Interfaces"))
	sslides[i+2].Append(show.ImageAt("agentEnemy.PNG", .5, .5))

	// The solution to this is to implement this sort of interface,
	// where you compose all of your entities with a struct that has
	// a function to return itself (as a pointer). Define an interface
	// of just that function and...

	sslides[i+3].Append(show.Header("Reusable AI"))
	sslides[i+3].Append(show.ImageAt("doctorEntity.PNG", .39, .5, mod.Scale(1.25, 1.25)))
	sslides[i+3].Append(show.ImageAt("doctorHasE.PNG", .15, .5, mod.Scale(1.25, 1.25)))
	sslides[i+3].Append(show.ImageAt("doctorBounce.PNG", .75, .5, mod.Scale(1.25, 1.25)))

	// ... now you can store all of the things
	// any AI entity needs in one embedded struct and run all of your
	// entities on any AI procedure you have.

	sslides[i+4].Append(show.Header("Aside: Composition for Private Features"))
	sslides[i+4].Append(show.ImageAt("phase.PNG", .3, .5))
	sslides[i+4].Append(show.ImageAt("phaseCollision.PNG", .7, .5))
}

var (
	levels = slideSetup{
		addLevels,
		8,
	}
)

func addLevels(i int, sslides []*static.Slide) {
	sslides[i].Append(show.Title("Designing Levels with Interfaces"))
	sslides[i+1].Append(show.Header("A Poor Approach"))
	sslides[i+1].Append(show.ImageCaption("agentLevelGenerate.PNG", .3, .1, 1, Libel28, "Level Generation in Agent Blue"))
	sslides[i+1].Append(show.ImageAt("agentRoom.PNG", .6, .5, mod.Scale(1.25, 1.25)))
	//
	// Our first approach to building levels didn't use interfaces.
	// We're going to go through why this was a terrible idea.
	// Agent Blue was the first game we started making in Go and so
	// it also has all of our worst patterns in its code.
	// ...
	sslides[i+2].Append(show.Header("A Poor Approach"))
	sslides[i+2].Append(show.ImageCaption("agentTutorial.png", .2, .4, 5, Libel28, "Agent Blue Tutorial Map"))
	sslides[i+2].Append(show.ImageCaption("agentLevelSelect.PNG", .5, .3, 1.25, Libel28, "Agent Blue Level Select Room (demo)"))
	// Level vs LevelSelect
	sslides[i+3].Append(show.Header("A Poor Approach"))
	sslides[i+3].Append(show.ImageAt("agentDoodad.PNG", .5, .5, mod.Scale(1.25, 1.25)))
	//
	//	Doodads
	//
	sslides[i+4].Append(show.Header("Modular Tile Enumeration"))
	sslides[i+4].Append(show.ImageAt("jeremyLevel.PNG", .4, .5))
	sslides[i+4].Append(show.ImageAt("jeremyTileTypes.PNG", .6, .5))
	//
	// So we were restricted because our tiles had too limited functionality.
	// In Jeremy the Clam I tried to adapt this out by giving tiles explicit
	// types based on their value, and each tile built itself using a Place
	// function. So where before we had a tile interface, now tiles are just
	// integers, making shared functionality a lot easier.
	//
	sslides[i+5].Append(show.Header("Modular Tile Enumeration"))
	sslides[i+5].Append(show.ImageAt("jeremyTilePlace.PNG", .5, .5))
	//
	// However, the immediate problem caused from this was that we no longer
	// could place multiple tiles in the same tile position. Before we could
	// stack floor tiles below wall tiles or doodad tiles, and now we need
	// tiles that are labeled as non-floors to place floors below them
	// when they get inserted during the start of the level.
	//
	sslides[i+6].Append(show.Header("Level Interfaces"))
	sslides[i+6].Append(show.ImageAt("doctorLevelPlace.PNG", .3, .5))
	sslides[i+6].Append(show.ImageAt("doctorLevel.PNG", .7, .5))
	//
	// We addressed this in A Fantastic Doctor by moving entity creation
	// out of levels themselves, but also by abstracting the concept of a
	// level in the first place. In this case, all a level (or in this
	// game, an Organ) needs to provide is a Place() function to initialize
	// all of its components when it is entered.
	//
	//
	// So while we didn't do this, that means that we can extend organ
	// functionality by making organs with layers of tiles instead of
	// just one 2d layer.
	sslides[i+7].Append(show.Header("Level Files"))
	sslides[i+7].Append(show.ImageCaption("jeremyFile.PNG", .2, .3, 1.0, Libel28, "A Jeremy Level File"))
	sslides[i+7].Append(show.ImageCaption("doctorFile.PNG", .6, .3, 2.0, Libel28, "A Fantastic Doctor Level File"))
}

var (
	conclusion = slideSetup{
		addConclusion,
		3,
	}
)

func addConclusion(i int, sslides []*static.Slide) {
	sslides[i].Append(show.Header("Thanks To"))
	sslides[i].Append(
		show.TxtSetFrom(Gnuolane44, .25, .35, 0, .07,
			"- Nate Fudenberg, John Ficklin",
			"- Contributors on Github",
			"- You, Audience",
		)...,
	)

	// And I'll end by thanking the other people in Oakmound
	// for working with me on our engine, those who've tried
	// out the engine and raised issues or PRs, and all of you
	// for listening.
	sslides[i+1].Append(show.Header("Other GUI Programs"))
	sslides[i+1].Append(show.ImageCaption("slide.PNG", .2, .4, 1.25, Libel28, "This Slideshow"))

	// As a final note, to demonstrate some kind of versatility,
	// this slide show was written in Oak as well.

	sslides[i+2].Append(show.Title("Questions"))

	// ?
}
