# scene

The scene package handles distinct re-loadable sections of application or game flow.

## Start

A scene itself has two fields: `Start` and `End`, functions called when a scene starts and ends respectively. Start is provided a `Context` object, which both serves as a standard Go context in persisting cancellations and arbitrary data, and ensures the operations within the scene operate on the correct data structures for the current logic and display.

## Context

A `Context` has fields to enable accessing the `window`, `render`, `key`, `collision`, `mouse`, and `event` packages directly and succinctly. These helpers are embedded into the context when possible, enabling calls like:

```go
render.Draw(...) // implied to run on the draw stack relevant for the current scene, fails in a multi-window context
// vs 
ctx.Draw(...) // explicitly runs on the draw stack for the current scene


mouse.DefaultTree.Hits(...) // ^
//vs
ctx.MouseTree.Hits(...) // ^

oak.SetFullScreen(true) 
// vs
ctx.Window.SetFullScreen(true)
```

The `Window` field on a `Context` will have different functionality based on the platform being compiled; e.g. window position manipulation is not possible when targeting Android or JS. These methods will fail to exist at compile time if used, preventing runtime errors. 

## End

When `Context.Window.GoToScene` or `Context.Window.NextScene` is called, the `End` function will be triggered. `End` enables a scene to define which other scene follow it and how it should transition to the next scene, if applicable.

- When `GoToScene` is used, the next scene output from `End` is ignored.
- It is valid to have a scene end and return to itself.
- It is valid to never define any `End` functions and solely rely on `GoToScene`.
- With the exception of persistent event bindings and the structure of the draw stack (not the elements on the draw stack), every entity is destroyed, undrawn, and unbound on scene end.

## Helpers

Scene has other utilities:

- `ctx.DoAfter` and `ctx.DoAfterContext` will safely register a time based callback which will not be called if the scene ends early
- `ctx.DrawForTime` will render an element for a specified duration (using `DoAfter`)
- The `Map` type serves as the underlying data structure for the currently registered scenes, available via oak.Window.SceneMap
- `GoTo` and `GoToPtr` offer shorthand for End functions.
