package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

//go:embed assets/images/TestShapes.svg
var testShapes []byte

func main() {
	oak.AddScene("svg", scene.Scene{
		Start: func(*scene.Context) {
			// load svg
			// svg from oksvg testdata
			icon, err := oksvg.ReadIconStream(bytes.NewBuffer(testShapes))
			if err != nil {
				fmt.Println(err)
			}
			// put it in the thing

			inputW, inputH := icon.ViewBox.W, icon.ViewBox.H
			iconAspect := inputW / inputH
			const width = 640
			const height = 480

			buff := image.NewRGBA(image.Rect(0, 0, width, height))

			viewAspect := float64(width) / float64(height)
			outputW, outputH := width, height
			if viewAspect > iconAspect {
				outputW = int(float64(height) * iconAspect)
			} else if viewAspect < iconAspect {
				outputH = int(float64(width) / iconAspect)
			}
			scanner := rasterx.NewScannerGV(int(inputW), int(inputH), buff, image.Rect(0, 0, width, height))
			scanner.SetBounds(10000, 10000)
			dasher := rasterx.NewDasher(width, height, scanner)
			icon.SetTarget(0, 0, float64(outputW), float64(outputH))
			icon.Draw(dasher, 1)

			render.Draw(render.NewSprite(0, 0, buff))
		},
	})

	oak.Init("svg")
}
