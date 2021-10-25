package oak

import (
	"image"
	"image/draw"
)

// A Background can be used as a background draw layer. Backgrounds will be drawn as the first
// element in each frame, and are expected to cover up data drawn on the previous frame.
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
		select {
		case <-w.quitCh:
			return
		case <-w.drawCh:
			<-w.drawCh
		loadingSelect:
			for {
				select {
				case <-w.ParentContext.Done():
					return
				case <-w.quitCh:
					return
				case <-w.drawCh:
					break loadingSelect
				case <-w.animationFrame:
					draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
					if w.LoadingR != nil {
						w.LoadingR.Draw(w.winBuffer.RGBA(), 0, 0)
					}
					w.publish()
				case <-w.DrawTicker.C:
					draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
					if w.LoadingR != nil {
						w.LoadingR.Draw(w.winBuffer.RGBA(), 0, 0)
					}
					w.publish()
				}
			}
		case f := <-w.betweenDrawCh:
			f()
		case <-w.animationFrame:
			draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
			w.DrawStack.PreDraw()
			p := w.viewPos
			w.DrawStack.DrawToScreen(w.winBuffer.RGBA(), &p, w.ScreenWidth, w.ScreenHeight)
			w.publish()
		case <-w.DrawTicker.C:
			draw.Draw(w.winBuffer.RGBA(), w.winBuffer.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
			w.DrawStack.PreDraw()
			p := w.viewPos
			w.DrawStack.DrawToScreen(w.winBuffer.RGBA(), &p, w.ScreenWidth, w.ScreenHeight)
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

// DoBetweenDraws will execute the given function in-between draw frames
func (w *Window) DoBetweenDraws(f func()) {
	go func() {
		w.betweenDrawCh <- f
	}()
}
