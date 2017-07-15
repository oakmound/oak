#Oak Engine
###An engine in golang for golang games

----

##Usage
1. Go get the oak engine
1. Get one or all of the examples
1. Go to the chosen example's folder
1. Type: go run core.go
1. Hit enter

##Reason
The initial version of Oak was made to support Oakmound Studio's game:
[Agent Blue](http://oakmound.blogspot.com/) and was developed in parallel.
Oak supports both Windows and Linux without outside dependencies like a C compiler.
 OakmoundStudio hopes that users will be able to make great pure Go games with Oak and potentially improve the oak engine.



##Features
1. Asset Management
    1. Loading
    1. Batch Loading
    1. Drawing
    1. Manipulation
        1. Recoloring
        1. Transforming
        1. Moving
        1. Shading
        1. Copying
1. Mouse Handling 
    1. Click Collision
    1. Drag Handling
1. Audio Support
    1. From [klangsynthese](https://github.com/200sc/klangsynthese)
1. Collision
    1. Collision rTrees
        1. Multiple disjoint trees
        1. Object Collsion trees
        1. Mouse Collision trees
1. Physics System
    1. Attachable to any object with a logical location
    1. Vectors
        1. Momentum
        1. Friction
        1. Force
1. Event Driven System
    1. Entities can bind and trigger events
1. Timing system
    1. Accurate time tracking
    1. FPS conversion
    1. Manipulable Time Tick Rate (speed up slow down timers tick rate)
1. Console Commands
    1. Supports the easy addition of new console commands
1. Logging
    1. Controlled config
    1. Filterable


## Contributions
Contributions to this project are welcome, though please send mail before
starting work on anything major.


