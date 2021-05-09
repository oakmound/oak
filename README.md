# Oak

## A Pure Go game engine

[![Go Reference](https://pkg.go.dev/badge/github.com/oakmound/oak/v2.svg)](https://pkg.go.dev/github.com/oakmound/oak/v2)
[![Code Coverage](https://codecov.io/gh/oakmound/oak/branch/develop/graph/badge.svg)](https://codecov.io/gh/oakmound/oak)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

## Table of Contents

1. [Installation](#installation)

1. [Motivation](#motivation)

1. [Features](#features)

1. [Support](#support)

1. [Quick Start](#quick-start)

1. [Implementation and Examples](#)

1. [Finished Games](#finished-games)

***

## Installation <a name="installation"/>

`go get -u github.com/oakmound/oak/v2/`

## Motivation <a name="motivation"/>

The initial version of oak was made to support Oakmound Studio's game,
[Agent Blue](https://github.com/OakmoundStudio/AgentRelease), and was developed in parallel.

Because Oak wants to have as few non-Go dependencies as possible, Oak does not by default use OpenGL or [GLFW](https://github.com/go-gl/glfw).

### On Pure Go

Oak has recently brought in dependencies that include C code, but we still describe the engine as a Pure Go engine, which at face value seems contradictory. Oak's goal is that, by default, a user can pull down the engine and create a fully functional game or GUI application on a machine with no C compiler installed, so when we say Pure Go we mean that, by default, the library is configured so no C compilation is required, and that no major features are locked behind C compliation.  

We anticipate in the immediate future needing to introduce alternate drivers that include C dependencies for performance improvements in some scenarios, and currently we have no OSX solution that lacks objective C code.

## Features <a name="features"></a>

1. Window Rendering
    - Windows and key events forked from [shiny](https://github.com/oakmound/oak/v2/shiny)
    - Logical frame rate distinct from Draw rate
    - Fullscreen, Window Positioning support
    - Auto-scaling for screen size changes
1. [Image Management](https://godoc.org/github.com/oakmound/oak/render)
    - `render.Renderable` interface
    - Sprite Sheet Batch Loading at startup
    - Manipulation
        - `render.Modifiable` interface
        - Built in Transformations and Filters
        - Some built-ins via [gift](https://github.com/disintegration/gift)
        - Extensible Modification syntax `func(image.Image) *image.RGBA`
    - Built in `Renderable` types covering common use cases
        - `Sprite`, `Sequence`, `Switch`, `Composite`
        - Primitive builders, `ColorBox`, `Line`, `Bezier`
        - History-tracking `Reverting`
    - Primarily 2D
1. [Particle System](https://godoc.org/github.com/oakmound/oak/render/particle)
    - <details>
      <summary>Click to see gif captured in examples/particle-demo</summary>

        ![particles!](examples\particle-demo\overviewExample.gif)
    </details>
1. [Mouse Handling](https://godoc.org/github.com/oakmound/oak/mouse)
    - Click Collision
    - MouseEnter / MouseExit reaction events
    - Drag Handling
1. [Joystick Support](https://godoc.org/github.com/oakmound/oak/joystick)
    - <details>
      <summary>Click to see gif captured in examples/joystick-viz</summary>

        ![joysticks!](examples\joystick-viz\example.gif)
    </details>
1. [Audio Support](https://godoc.org/github.com/oakmound/oak/audio)
    - Positional filters to pan and scale audio based on a listening position
1. [Collision](https://godoc.org/github.com/oakmound/oak/collision)
    - Collision R-Tree forked from [rtreego](https://github.com/dhconnelly/rtreego)
    - [2D Raycasting](https://godoc.org/github.com/oakmound/oak/collision/ray)
    - Collision Spaces
        - Attachable to Objects
        - Auto React to collisions through events
        - OnHit bindings `func(s1,s2 *collision.Space)`
        - Start/Stop collision with targeted objects
1. [2D Physics System](https://godoc.org/github.com/oakmound/oak/physics)
    - Vectors
        - Attachable to Objects / Renderables
        - Friction
1. [Event Handler, Bus](https://godoc.org/github.com/oakmound/oak/event)
    - PubSub system: `event.CID` can `Bind(fn,eventName)` and `Trigger(eventName)` events
1. [Shaping](https://godoc.org/github.com/oakmound/oak/shape)
    - Convert shapes into:
        - Containment checks
        - Outlines
        - 2D arrays
1. [Custom Console Commands](debugConsole.go)
1. [Logging](https://godoc.org/github.com/oakmound/oak/dlog)
    - Swappable with custom implementations
    - Default Implementation: 4 log levels, writes to file and stdout

## Support <a name="support"></a>

For discussions not significant enough to be an Issue or PR, see the #oak channel on the [gophers slack](https://invite.slack.golangbridge.org/).

## Quick Start <a name="quick-start"></a>

This is an example of the most basic oak program:

```go
package main

import (
    "github.com/oakmound/oak/v2"
    "github.com/oakmound/oak/v2/scene"
)

func main() {
    oak.AddScene("firstScene", scene.Scene{
        Start: func(*scene.Context) {
            // ... draw entities, bind callbacks ... 
        }, 
    })
    oak.Init("firstScene")
}
```

See the [examples](examples) folder for longer demos, [godoc](https://godoc.org/github.com/oakmound/oak) for reference documentation, and the [wiki](https://github.com/oakmound/oak/wiki) for more guided feature sets, tutorials and walkthroughs.

## Implementation and Examples <a name="examples"></a>

### Platformer

![Platformer](examples/platformer-tutorial/6-complete/example.gif)

Build up to a simple platforming game step by step in the guided walkthrough. // TODO Link wiki

```go
char := entities.NewMoving(100, 100, 16, 32,
    render.NewColorBox(16, 32, color.RGBA{255, 0, 0, 255}),
nil, 0, 0)

char.Bind(func(id event.CID, nothing interface{}) int {
    char := ie.E().(*entities.Moving)

    // Move left and right with A and D
    if oak.IsDown(key.A) {
        char.Delta.SetX(-char.Speed.X())
    } else if oak.IsDown(key.D) {
        char.Delta.SetX(char.Speed.X())
    } else {
        char.Delta.SetX(0)
    }
    char.ShiftPos(char.Delta.X(), char.Delta.Y())
    return 0
}, event.Enter)
```

### Top Down Shooter

Learn to use the collision library and move the viewport as characters move in the guided walkthrough. // TODO link wiki  

![Shoota](examples/top-down-shooter-tutorial/6-performance/example.gif)

### Radar

Often times you might want to create a minimap or a radar for a game, check out this example for a barebones implementation

![Radar](examples/radar-demo/example.gif)

### Slideshow

A different way to use the oak engine.

![Slideshow](examples/slide/example.gif)

## Examples of Finished Games <a name="finished-games"/>

[Agent Blue](https://oakmound.itch.io/agent-blue)

![AgentBlue](https://img.itch.zone/aW1hZ2UvMTk4MjIxLzkyNzUyOC5wbmc=/original/aRusLc.png)

[Fantastic Doctor](https://github.com/oakmound/lowrez17)

![Fantastic Overview](https://img.itch.zone/aW1hZ2UvMTY4NDk1Lzc4MDk1Mi5wbmc=/original/hIjzFm.png)

![Fantastic Overview 2](https://img.itch.zone/aW1hZ2UvMTY4NDk1LzI0MjMxNTEuZ2lm/original/1zpD6g.gif)

[Jeremy The Clam](https://github.com/200sc/jeremy)

![Clammy](https://img.itch.zone/aW1hZ2UvMTYzNjgyLzc1NDkxOS5wbmc=/original/%2BwvZ7j.png)
