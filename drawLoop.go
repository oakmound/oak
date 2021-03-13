package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/timing"
	"github.com/oakmound/shiny/screen"
)

type Background interface {
	GetRGBA() *image.RGBA
}

var (
	// TODO V3: cleanup this interface
	// BackgroundColor is the starting background color for the draw loop. BackgroundImage or SetBackground will overwrite it.
	BackgroundColor image.Image = image.Black
	// BackgroundImage is the starting custom background for the draw loop. SetBackground will overwrite it.
	BackgroundImage Background
	// DrawTicker is the parallel to LogicTicker to set the draw framerate
	DrawTicker *timing.DynamicTicker

	setBackgroundCh = make(chan Background)

	bkgFn = func() image.Image {
		return BackgroundColor
	}
)

func SetBackground(b Background) {
	setBackgroundCh <- b
}

// DrawLoop
// Unless told to stop, the draw channel will repeatedly
// 1. draw the background color to a temporary buffer
// 2. draw all visible rendered elements onto the temporary buffer.
// 3. draw the buffer's data at the viewport's position to the screen.
// 4. publish the screen to display in window.
func drawLoop() {
	<-drawCh

	tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
	if err != nil {
		panic(err)
	}

	if BackgroundImage != nil {
		bkgFn = func() image.Image {
			return BackgroundImage.GetRGBA()
		}
	}

	draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), bkgFn(), zeroPoint, draw.Src)
	drawLoopPublish(tx)

	DrawTicker = timing.NewDynamicTicker()
	DrawTicker.SetTick(timing.FPSToDuration(DrawFrameRate))

	dlog.Verb("Draw Loop Start")
	for {
	drawSelect:
		select {
		case bkg := <-setBackgroundCh:
			bkgFn = func() image.Image {
				return bkg.GetRGBA()
			}
		case <-windowUpdateCh:
			<-windowUpdateCh
		case <-drawCh:
			dlog.Verb("Got something from draw channel")
			<-drawCh
			dlog.Verb("Starting loading")
			for {
				<-DrawTicker.C
				draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), bkgFn(), zeroPoint, draw.Src)
				if LoadingR != nil {
					LoadingR.Draw(winBuffer.RGBA())
				}
				drawLoopPublish(tx)

				select {
				case <-drawCh:
					break drawSelect
				case bkg := <-setBackgroundCh:
					bkgFn = func() image.Image {
						return bkg.GetRGBA()
					}
				case viewPoint := <-viewportCh:
					dlog.Verb("Got something from viewport channel (waiting on draw)")
					updateScreen(viewPoint[0], viewPoint[1])
				case viewPoint := <-viewportShiftCh:
					dlog.Verb("Got something from viewport shift channel (waiting on draw)")
					shiftViewPort(viewPoint[0], viewPoint[1])
				default:
				}
			}
		case viewPoint := <-viewportCh:
			dlog.Verb("Got something from viewport channel")
			updateScreen(viewPoint[0], viewPoint[1])
		case viewPoint := <-viewportShiftCh:
			dlog.Verb("Got something from viewport shift channel")
			shiftViewPort(viewPoint[0], viewPoint[1])
		case <-DrawTicker.C:
			draw.Draw(winBuffer.RGBA(), winBuffer.Bounds(), bkgFn(), zeroPoint, draw.Src)
			render.PreDraw()
			render.GlobalDrawStack.Draw(winBuffer.RGBA(), ViewPos, ScreenWidth, ScreenHeight)
			drawLoopPublish(tx)
		}
	}
}

var (
	drawLoopPublishDef = func(tx screen.Texture) {
		tx.Upload(zeroPoint, winBuffer, winBuffer.Bounds())
		windowControl.Scale(windowRect, tx, tx.Bounds(), draw.Src)
		windowControl.Publish()
	}
	drawLoopPublish = drawLoopPublishDef
)
