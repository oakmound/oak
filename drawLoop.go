package oak

import (
	"image"
	"image/draw"

	"runtime/debug"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
	"github.com/oakmound/shiny/screen"
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
	<-drawCh

	debug.SetPanicOnFault(true)

	tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
	if err != nil {
		panic(err)
	}

	draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, draw.Src)
	drawLoopPublish(tx)

	DrawTicker = timing.NewDynamicTicker()
	DrawTicker.SetTick(timing.FPSToDuration(DrawFrameRate))

	dlog.Verb("Draw Loop Start")
	for {
	drawSelect:
		select {
		case <-windowUpdateCh:
			<-windowUpdateCh
		case <-drawCh:
			dlog.Verb("Got something from draw channel")
			<-drawCh
			dlog.Verb("Starting loading")
			for {
				<-DrawTicker.C
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, draw.Src)
				if LoadingR != nil {
					LoadingR.Draw(winBuffer.RGBA())
				}
				drawLoopPublish(tx)

				select {
				case <-drawCh:
					break drawSelect
				case viewPoint := <-viewportCh:
					dlog.Verb("Got something from viewport channel (waiting on draw)")
					updateScreen(viewPoint[0], viewPoint[1])
				default:
				}
			}
		case viewPoint := <-viewportCh:
			dlog.Verb("Got something from viewport channel")
			updateScreen(viewPoint[0], viewPoint[1])
		case <-DrawTicker.C:
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), imageBlack, zeroPoint, draw.Src)
			render.PreDraw()
			render.GlobalDrawStack.Draw(winBuffer.RGBA(), ViewPos, ScreenWidth, ScreenHeight)
			drawLoopPublish(tx)
		}
	}
}

var (
	drawLoopPublishDef = func(tx screen.Texture) {
		tx.Upload(zeroPoint, winBuffer, winBuffer.Bounds())
		windowControl.Scale(windowRect, tx, tx.Bounds(), draw.Src, nil)
		windowControl.Publish()
	}
	drawLoopPublish = drawLoopPublishDef
)
