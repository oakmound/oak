package oak

import (
	"image"
	"image/draw"

	"github.com/oakmound/oak/v2/timing"
)

type Background interface {
	GetRGBA() *image.RGBA
}

func (c *Controller) SetBackground(b Background) {
	c.setBackgroundCh <- b
}

// DrawLoop
// Unless told to stop, the draw channel will repeatedly
// 1. draw the background color to a temporary buffer
// 2. draw all visible rendered elements onto the temporary buffer.
// 3. draw the buffer's data at the viewport's position to the screen.
// 4. publish the screen to display in window.
func (c *Controller) drawLoop() {
	<-c.drawCh

	tx, err := c.screenControl.NewTexture(c.winBuffer.Bounds().Max)
	if err != nil {
		panic(err)
	}

	if c.BackgroundImage != nil {
		c.bkgFn = func() image.Image {
			return c.BackgroundImage.GetRGBA()
		}
	}

	draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
	c.drawLoopPublish(c, tx)

	c.DrawTicker = timing.NewDynamicTicker()
	c.DrawTicker.SetTick(timing.FPSToDuration(c.DrawFrameRate))

	for {
	drawSelect:
		select {
		case bkg := <-c.setBackgroundCh:
			c.bkgFn = func() image.Image {
				return bkg.GetRGBA()
			}
		case <-c.windowUpdateCh:
			<-c.windowUpdateCh
		case <-c.drawCh:
			<-c.drawCh
			for {
				<-c.DrawTicker.C
				draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
				if c.LoadingR != nil {
					c.LoadingR.Draw(c.winBuffer.RGBA(), 0, 0)
				}
				c.drawLoopPublish(c, tx)

				select {
				case <-c.drawCh:
					break drawSelect
				case viewPoint := <-c.viewportCh:
					c.updateScreen(viewPoint)
				case viewPoint := <-c.viewportShiftCh:
					c.shiftViewPort(viewPoint)
				default:
				}
			}
		case viewPoint := <-c.viewportCh:
			c.updateScreen(viewPoint)
		case viewPoint := <-c.viewportShiftCh:
			c.shiftViewPort(viewPoint)
		case <-c.DrawTicker.C:
			draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
			c.DrawStack.PreDraw()
			c.DrawStack.DrawToScreen(c.winBuffer.RGBA(), c.ViewPos, c.ScreenWidth, c.ScreenHeight)
			c.drawLoopPublish(c, tx)
		}
	}
}
