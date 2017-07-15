# Oak 
### A pure Go game engine

----

## Usage
1. run `go get -u github.com/OakmoundStudio/oak`
1. Get one or all of the examples
1. Go to the chosen example's folder
1. run `go run core.go / demo.go`

## Reason
The initial version of oak was made to support Oakmound Studio's game:
[Agent Blue](http://oakmound.blogspot.com/) and was developed in parallel.
Oak supports Windows with no dependencies and Linux with limited audio dependencies.
 We hope that users will be able to make great pure Go games with oak and potentially improve oak.



## Features
1. Window Rendering
    - Windows and key events through [shiny](https://github.com/golang/exp/tree/master/shiny)
1. Asset Management
    - Loading
    - Batch Loading
    - Manipulation
        - Recoloring
        - Transforming
        - Moving
        - Shading
        - Copying
1. Mouse Handling 
    - Click Collision
    - Drag Handling
1. Audio Support
    - From [klangsynthese](https://github.com/200sc/klangsynthese)
1. Collision
    - Collision rTrees from [rtreego](https://github.com/dhconnelly/rtreego)
    - 2D Raycasting
    - Collision Spaces
        - Attachable to Objects
        - React to collisions with events
        - Start/Stop collision with targeted objects
1. Physics System
    - Vectors
        - Attachable to Objects / Renderables
        - Momentum
        - Friction
        - Force / Pushing
1. Event Driven System
    - Entities can bind and trigger events
1. Timing system
    - Accurate time tracking
    - FPS conversion
    - Manipulable Time Tick Rate (speed up slow down timers tick rate)
1. Shaping 
    - Shapes from x->y equations
    - Convert shapes into containment checks
    - Convert shapes into outlines
1. Console Commands
    - Supports the easy addition of new console commands
1. Logging (Probably going to not roll our own soon!)
    - Controlled config
    - Filterable


## More Examples

Oak contains a few snippet examples, but a number of examples exist as external packages.

See examples/README.md

## Contributions
See CONTRIBUTING.md