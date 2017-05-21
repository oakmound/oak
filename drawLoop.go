package oak

import (
	"fmt"
	"image"
	"image/draw"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/render"
	"bitbucket.org/oakmoundstudio/oak/timing"
	"golang.org/x/exp/shiny/screen"
)

var (
	imageBlack = image.Black
	// DrawTicker is an unused parallel to LogicTicker to set the draw framerate
	DrawTicker *timing.DynamicTicker
)

// DrawLoop
// Unless told to stop, the draw channel will repeatedly
// 1. draw black to a temporary buffer
// 2. run any functions bound to precede drawing.
// 3. draw all elements onto the temporary buffer.
// 4. run any functions bound to follow drawing.
// 5. draw the buffer's data at the viewport's position to the screen.
// 6. publish the screen to display in window.
func drawLoop() {
	<-drawChannel

	tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
	if err != nil {
		fmt.Println(err)
	}

	draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
	drawLoopPublish(tx)

	//DrawTicker = timing.NewDynamicTicker()
	//DrawTicker.SetTick(timing.FPSToDuration(DrawFrameRate))

	for {
		//<-DrawTicker.C
		dlog.Verb("Draw Loop")
	drawSelect:
		select {
		case <-windowUpdateCH:
			<-windowUpdateCH
		case <-drawChannel:
			dlog.Verb("Got something from draw channel")
			<-drawChannel
			dlog.Verb("Starting loading")
			for {
				//<-DrawTicker.C
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
				if LoadingR != nil {
					LoadingR.Draw(winBuffer.RGBA())
				}
				drawLoopPublish(tx)

				select {
				case <-drawChannel:
					break drawSelect
				case viewPoint := <-viewportChannel:
					dlog.Verb("Got something from viewport channel (waiting on draw)")
					updateScreen(viewPoint[0], viewPoint[1])
				default:
				}
			}
		case viewPoint := <-viewportChannel:
			dlog.Verb("Got something from viewport channel")
			updateScreen(viewPoint[0], viewPoint[1])
		default:
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, screen.Src)
			render.PreDraw()
			render.GlobalDrawStack.Draw(winBuffer.RGBA(), ViewPos, ScreenWidth, ScreenHeight)
			drawLoopPublish(tx)
			//event.Trigger("PostDraw", nil)
		}
	}
}

func drawLoopPublish(tx screen.Texture) {
	tx.Upload(zeroPoint, winBuffer, winBuffer.Bounds())
	windowControl.Scale(windowRect, tx, tx.Bounds(), screen.Src, nil)
	windowControl.Publish()
}
