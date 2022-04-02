package main

import (
	"image/color"
	"os"

	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

// Axes are the plural of axis
type Axes uint8

// This is an enum for what axes to center around
const (
	X Axes = iota
	Y
	Both
)

func center(ctx *scene.Context, obj render.Renderable, ax Axes) {
	objWidth, objHeight := obj.GetDims()

	switch ax {
	case Both:
		obj.SetPos(float64(ctx.Window.Width()/2-objWidth/2),
			float64(ctx.Window.Height()-objHeight)/2) //distributive property
	case X:
		obj.SetPos(float64(ctx.Window.Width()-objWidth)/2, obj.Y())
	case Y:
		obj.SetPos(obj.X(), float64(ctx.Window.Height()-objHeight)/2)
	}
}

func main() {
	//make the scene for the titlescreen
	oak.AddScene("titlescreen", scene.Scene{Start: func(ctx *scene.Context) {

		//create text saying titlescreen in placeholder position
		titleText := render.NewText("titlescreen", 0, 0)

		//center text along both axes
		center(ctx, titleText, Both)

		//tell the draw loop to draw titleText
		render.Draw(titleText)

		//do the same for the text with button instuctions, but this time Y position is not a placeholder (X still is)
		instructionText := render.NewText("press Enter to start, or press Q to quit", 0, float64(ctx.Window.Height()*3/4))
		//this time we only center the X axis, otherwise it would overlap titleText
		center(ctx, instructionText, X)
		render.Draw(instructionText)
		event.GlobalBind(ctx, key.Down(key.Enter), func(key.Event) event.Response {
			// Go to the next scene if enter is pressed. Next scene is the game
			ctx.Window.NextScene()
			return 0
		})
		event.GlobalBind(ctx, key.Down(key.Q), func(key.Event) event.Response {
			// exit the game if q is pressed
			os.Exit(0)
			return 0
		})

	}, End: func() (string, *scene.Result) {
		return "game", nil //set the next scene to "game"
	}})

	//we declare this here so it can be accesed by the scene start and scene loop
	var player *entities.Moving

	//define the "game" (it's just a square that can be moved with WASD)
	oak.AddScene("game", scene.Scene{Start: func(ctx *scene.Context) {
		//create the player, a blue 32x32 square at 100,100
		player = entities.NewMoving(100, 100, 32, 32,
			render.NewColorBox(32, 32, color.RGBA{0, 0, 255, 255}),
			nil, 0, 0)
		//because the player is more than visuals (it has a hitbox, even though we don't use it),
		//we have to get the visual part specificaly, and not the whole thing.
		render.Draw(player.R)

		controlsText := render.NewText("WASD to move, ESC to return to titlescreen", 5, 20)
		//we draw the text on layer 1 (instead of the default layer 0)
		//because we want it to show up above the player
		render.Draw(controlsText, 1)
		event.GlobalBind(ctx, key.Down(key.Escape), func(key.Event) event.Response {
			// Go to the next scene if escape is pressed. Next scene is titlescreen
			ctx.Window.NextScene()
			return 0
		})
		event.GlobalBind(ctx, event.Enter, func(event.EnterPayload) event.Response {
			if oak.IsDown(key.SStr) {
				//if S is pressed, set the player's vertical speed to 2 (positive == down)
				player.Delta.SetY(2)
			} else if oak.IsDown(key.WStr) {
				player.Delta.SetY(-2)
			} else {
				//if the now buttons are pressed for vertical movement, don't move verticaly
				player.Delta.SetY(0)
			}

			//do the same thing as before, but horizontaly
			if oak.IsDown(key.DStr) {
				player.Delta.SetX(2)
			} else if oak.IsDown(key.AStr) {
				player.Delta.SetX(-2)
			} else {
				player.Delta.SetX(0)
			}
			//apply the player's speed to their position
			player.ShiftPos(player.Delta.X(), player.Delta.Y())
			return 0
		})
	}, End: func() (string, *scene.Result) {
		return "titlescreen", nil //set the next scene to be titlescreen
	}})
	//start the game on the titlescreen
	oak.Init("titlescreen")
}
