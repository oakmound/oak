package oak

import (
	"image"
	"image/draw"
)

type Background interface {
	GetRGBA() *image.RGBA
}

// DrawLoop
// Unless told to stop, the draw channel will repeatedly
// 1. draw the background color to a temporary buffer
// 2. draw all visible rendered elements onto the temporary buffer.
// 3. draw the buffer's data at the viewport's position to the screen.
// 4. publish the screen to display in window.
func (c *Controller) drawLoop() {
	<-c.drawCh

	draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
	c.publish()

	for {
	drawSelect:
		select {
		case <-c.drawCh:
			<-c.drawCh
			for {
				select {
				case <-c.drawCh:
					break drawSelect
				case <-c.DrawTicker.C:
					draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
					if c.LoadingR != nil {
						c.LoadingR.Draw(c.winBuffer.RGBA(), 0, 0)
					}
					c.publish()
				}
			}
		case <-c.DrawTicker.C:
			draw.Draw(c.winBuffer.RGBA(), c.winBuffer.Bounds(), c.bkgFn(), zeroPoint, draw.Src)
			c.DrawStack.PreDraw()
			c.DrawStack.DrawToScreen(c.winBuffer.RGBA(), &c.viewPos, c.ScreenWidth, c.ScreenHeight)
			c.publish()
		}
	}
}

func (c *Controller) publish() {
	c.prePublish(c, c.windowTexture)
	c.windowTexture.Upload(zeroPoint, c.winBuffer, c.winBuffer.Bounds())
	c.windowControl.Scale(c.windowRect, c.windowTexture, c.windowTexture.Bounds(), draw.Src)
	c.windowControl.Publish()
}
