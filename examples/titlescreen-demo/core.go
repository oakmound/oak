package main

import (
	"image/color"
	"os"

	oak "github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/entities"
	"github.com/oakmound/oak/v2/key"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

// Axes are the plural of axis
type Axes uint8 

// This is an enum for what axes to center around
const (
	X Axes = iota
	Y
	Both
)

func center(obj render.Renderable, ax Axes) {
	objWidth, objHeight := obj.GetDims()

	switch ax {
	case Both:
		obj.SetPos(float64(oak.ScreenWidth/2-objWidth/2),
			float64(oak.ScreenHeight-objHeight)/2) //distributive property
	case X:
		obj.SetPos(float64(oak.ScreenWidth-objWidth)/2, obj.Y())
	case Y:
		obj.SetPos(obj.X(), float64(oak.ScreenHeight-objHeight)/2)
	}
}

func main() {
	//make the scene for the titlescreen
	oak.Add("titlescreen", func(string, interface{}) { //scene start

		//create text saying titlescreen in placeholder position
		titleText := render.NewStrText("titlescreen", 0, 0)

		//center text along both axes
		center(titleText, Both)

		//tell the draw loop to draw titleText
		render.Draw(titleText)

		//do the same for the text with button instuctions, but this time Y position is not a placeholder (X still is)
		instructionText := render.NewStrText("press Enter to start, or press Q to quit", 0, float64(oak.ScreenHeight*3/4))
		//this time we only center the X axis, otherwise it would overlap titleText
		center(instructionText, X)
		render.Draw(instructionText)
	}, func() bool { //scene loop
		//if the enter key is pressed, go to the next scene
		if oak.IsDown(key.Enter) {
			return false
		}
		//exit the program if the q key is pressed
		if oak.IsDown(key.Q) {
			os.Exit(0)
		}
		return true
	}, func() (string, *scene.Result) { //scene end
		return "game", nil //set the next scene to "game"
	})

	//we declare this here so it can be accesed by the scene start and scene loop
	var player *entities.Moving

	//define the "game" (it's just a square that can be moved with WASD)
	oak.Add("game", func(string, interface{}) {
		//create the player, a blue 32x32 square at 100,100
		player = entities.NewMoving(100, 100, 32, 32,
			render.NewColorBox(32, 32, color.RGBA{0, 0, 255, 255}),
			nil, 0, 0)
		//because the player is more than visuals (it has a hitbox, even though we don't use it),
		//we have to get the visual part specificaly, and not the whole thing.
		render.Draw(player.R)

		controlsText := render.NewStrText("WASD to move, ESC to return to titlescreen", 5, 20)
		//we draw the text on layer 1 (instead of the default layer 0)
		//because we want it to show up above the player
		render.Draw(controlsText, 1)
	}, func() bool {
		//if escape is pressed, go to the next scene (titlescreen)
		if oak.IsDown(key.Escape) {
			return false
		}
		//controls
		if oak.IsDown(key.S) {
			//if S is pressed, set the player's vertical speed to 2 (positive == down)
			player.Delta.SetY(2)
		} else if oak.IsDown(key.W) {
			player.Delta.SetY(-2)
		} else {
			//if the now buttons are pressed for vertical movement, don't move verticaly
			player.Delta.SetY(0)
		}

		//do the same thing as before, but horizontaly
		if oak.IsDown(key.D) {
			player.Delta.SetX(2)
		} else if oak.IsDown(key.A) {
			player.Delta.SetX(-2)
		} else {
			player.Delta.SetX(0)
		}
		//apply the player's speed to their position
		player.ShiftPos(player.Delta.X(), player.Delta.Y())
		return true
	}, func() (string, *scene.Result) {
		return "titlescreen", nil //set the next scene to be titlescreen
	})
	//start the game on the titlescreen
	oak.Init("titlescreen")
}
