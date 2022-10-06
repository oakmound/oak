# Oak

## A Pure Go game engine

[![Go Reference](https://pkg.go.dev/badge/github.com/oakmound/oak/v4.svg)](https://pkg.go.dev/github.com/oakmound/oak/v4)
[![Code Coverage](https://codecov.io/gh/oakmound/oak/branch/master/graph/badge.svg)](https://codecov.io/gh/oakmound/oak)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge-flat.svg)](https://github.com/avelino/awesome-go)

## Table of Contents

1. [Installation](#installation)

1. [Features](#features)

1. [Support](#support)

1. [Quick Start](#quick-start)

1. [Examples](#examples)

1. [Finished Games](#finished-games)

***

## Installation <a name="installation"/>

`go get -u github.com/oakmound/oak/v4`

## Features and Systems <a name="features"></a>

1. Window Management
    - Windows and key events forked from [shiny](https://pkg.go.dev/golang.org/x/exp/shiny)
    - Support for multiple windows running at the same time
1. [Image Rendering](https://pkg.go.dev/github.com/oakmound/oak/v4/render)
    - Manipulation
        - `render.Modifiable` interface
        - Integrated with optimized image manipulation via [gift](https://github.com/disintegration/gift)
    - Built in `Renderable` types covering common use cases
        - `Sprite`, `Sequence`, `Switch`, `Composite`
        - Primitive builders, `ColorBox`, `Line`, `Bezier`
        - History-tracking `Reverting`
    - Primarily 2D
1. [Particle System](https://pkg.go.dev/github.com/oakmound/oak/v4/render/particle)
1. [Mouse Handling](https://pkg.go.dev/github.com/oakmound/oak/v4/mouse)
1. [Joystick Support](https://pkg.go.dev/github.com/oakmound/oak/v4/joystick)
1. [Audio Support](https://pkg.go.dev/github.com/oakmound/oak/v4/audio)
1. [Collision](https://pkg.go.dev/github.com/oakmound/oak/v4/collision)
    - Collision R-Tree forked from [rtreego](https://github.com/dhconnelly/rtreego)
    - [2D Raycasting](https://pkg.go.dev/github.com/oakmound/oak/v4/collision/ray)
    - Collision Spaces
        - Attachable to Objects
        - Auto React to collisions through events
1. [2D Physics System](https://pkg.go.dev/github.com/oakmound/oak/v4/physics)
1. [Event Handler](https://pkg.go.dev/github.com/oakmound/oak/v4/event)

## Support <a name="support"></a>

For discussions not significant enough to be an Issue or PR, feel free to ping us in the #oak channel on the [gophers slack](https://invite.slack.golangbridge.org/). For insight into what is happening in oak see the [blog](https://200sc.dev/).

## Quick Start <a name="quick-start"></a>

This is an example of the most basic oak program:

```go
package main

import (
    "github.com/oakmound/oak/v4"
    "github.com/oakmound/oak/v4/scene"
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

See below or navigate to the [examples](examples) folder for demos. For more examples and documentation checkout  [godoc](https://pkg.go.dev/github.com/oakmound/oak/v4) for reference documentation, the [wiki](https://github.com/oakmound/oak/wiki), or our extended features in [grove](https://github.com/oakmound/grove). 

## Examples <a name="examples"></a>

| | | |
|:-------------------------:|:-------------------------:|:-------------------------:|
|<img width="200"  src="examples/platformer/example.gif" a=examples/platformer>  [Platformer](examples/platformer) |  <img width="200"  src="examples/top-down-shooter//example.gif"> [Top down shooter](examples/top-down-shooter)|<img width="200"  src="examples/flappy-bird//example.gif"> [Flappy Bird](examples/flappy-bird/)
|  <img width="200"  src="examples/bezier/example.PNG"> [Bezier Curves](examples/bezier) |<img width="200"  src="examples/joystick-viz/example.gif"> [Joysticks](examples/joystick-viz)|<img width="200"  src="examples/piano/example.gif"> [Piano](examples/piano)|
|<img width="200"  src="examples/screenopts/example.PNG"> [Screen Options](examples/screenopts)  |  <img width="200"  src="examples/multi-window/example.PNG"> [Multi Window](examples/multi-window) |<img width="200"  src="examples/particle-demo/overviewExample.gif"> [Particles](examples/particle-demo)|

## Games using Oak <a name="finished-games"></a>

To kick off a larger game project you can get started with [game-template](https://github.com/oakmound/game-template).

| | |
|:-------------------------:|:-------------------------:|
|<img width="200"  src="https://img.itch.zone/aW1hZ2UvMTk4MjIxLzkyNzUyOC5wbmc=/original/aRusLc.png" a=examples/platformer-tutorial>  [Agent Blue](https://oakmound.itch.io/agent-blue) |  <img width="200"  src="https://img.itch.zone/aW1hZ2UvMTY4NDk1Lzc4MDk1Mi5wbmc=/original/hIjzFm.png"> [Fantastic Doctor](https://github.com/oakmound/lowrez17)
|<img width="200"  src="https://img.itch.zone/aW1hZ2UvMzkwNjM5LzI2NzU0ODMucG5n/original/eaoFrd.png">  [Hiring Now: Looters](https://oakmound.itch.io/cheststacker) |  <img width="200"  src="https://img.itch.zone/aW1hZ2UvMTYzNjgyLzc1NDkxOS5wbmc=/original/%2BwvZ7j.png"> [Jeremy The Clam](https://github.com/200sc/jeremy)
|<img width="200"  src="https://img.itch.zone/aW1hZ2UvOTE0MjYzLzUxNjg3NDEucG5n/original/5btfEr.png">  [Diamond Deck Championship](https://oakmound.itch.io/diamond-deck-championship) |  <img width="200"  src="https://img.itch.zone/aW1nLzgzMDM5MjcucG5n/105x83%23/oA19CL.png">  [SokoPic](https://oakmound.itch.io/sokopic) 

## On Pure Go <a name="pure-go"/>

Oak has recently brought in dependencies that include C code, but we still describe the engine as a Pure Go engine, which at face value seems contradictory. Oak's goal is that, by default, a user can pull down the engine and create a fully functional game or GUI application on a machine with no C compiler installed, so when we say Pure Go we mean that, by default, the library is configured so no C compilation is required, and that no major features are locked behind C compliation.  

We anticipate in the immediate future needing to introduce alternate drivers that include C dependencies for performance improvements in some scasenarios, and currently we have no OSX solution that lacks objective C code.
