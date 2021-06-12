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
func (w *Window) drawLoop() {
	<-w.drawCh

	draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
	w.publish()

	for {
	drawSelect:
		select {
		case <-w.quitCh:
			return
		case <-w.drawCh:
			<-w.drawCh
			for {
				select {
				case <-w.ParentContext.Done():
					return
				case <-w.quitCh:
					return
				case <-w.drawCh:
					break drawSelect
				case <-w.DrawTicker.C:
					draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
					if w.LoadingR != nil {
						w.LoadingR.Draw(w.winBuffer.RGBA(), 0, 0)
					}
					w.publish()
				}
			}
		case <-w.DrawTicker.C:
			draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
			w.DrawStack.PreDraw()
			w.DrawStack.DrawToScreen(w.winBuffer.RGBA(), &w.viewPos, w.ScreenWidth, w.ScreenHeight)
			w.publish()
		}
	}
}

func (w *Window) publish() {
	w.prePublish(w, w.windowTexture)
	w.windowTexture.Upload(zeroPoint, w.winBuffer, w.winBuffer.Bounds())
	w.windowControl.Scale(w.windowRect, w.windowTexture, w.windowTexture.Bounds(), draw.Src)
	w.windowControl.Publish()
}
