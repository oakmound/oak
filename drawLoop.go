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

	draw.Draw(w.winBuffers[w.bufferIdx].RGBA(), w.winBuffers[w.bufferIdx].Bounds(), w.bkgFn(), zeroPoint, draw.Src)
	w.publish()

	drawFrame := func() {
		buff := w.winBuffers[w.bufferIdx]
		if buff.RGBA() != nil {
			// Publish what was drawn last frame to screen, then work on preparing the next frame.
			w.publish()
			draw.Draw(buff.RGBA(), buff.Bounds(), w.bkgFn(), zeroPoint, draw.Src)
			w.DrawStack.PreDraw()
			p := w.viewPos
			w.DrawStack.DrawToScreen(buff.RGBA(), &p, w.ScreenWidth, w.ScreenHeight)
		}
	}

	drawLoadingFrame := func() {
		buff := w.winBuffers[w.bufferIdx]
		w.publish()
		draw.Draw(buff.RGBA(), w.winBuffers[w.bufferIdx].Bounds(), w.bkgFn(), zeroPoint, draw.Src)
		if w.LoadingR != nil {
			w.LoadingR.Draw(w.winBuffers[w.bufferIdx].RGBA(), 0, 0)
		}
	}

	if w.config.UnlimitedDrawFrameRate {
		// this code is duplicated as an optimization: it's much faster
		// to have a 'default' case than to flood a channel, and we can't conditionally
		// add a default case to a select.
		for {
			select {
			case <-w.quitCh:
				return
			case <-w.drawCh:
				<-w.drawCh
			loadingSelectUnlimited:
				for {
					select {
					case <-w.ParentContext.Done():
						return
					case <-w.quitCh:
						return
					case <-w.drawCh:
						break loadingSelectUnlimited
					case <-w.animationFrame:
						drawLoadingFrame()
					case <-w.DrawTicker.C:
						drawLoadingFrame()
					default:
						drawLoadingFrame()
					}
				}
			case f := <-w.betweenDrawCh:
				f()
			case <-w.animationFrame:
				drawFrame()
			case <-w.DrawTicker.C:
				drawFrame()
			default:
				drawFrame()
			}
		}
	}

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
					drawLoadingFrame()
				case <-w.DrawTicker.C:
					drawLoadingFrame()
				}
			}
		case f := <-w.betweenDrawCh:
			f()
		case <-w.animationFrame:
			drawFrame()
		case <-w.DrawTicker.C:
			drawFrame()
		}
	}
}

func (w *Window) publish() {
	w.prePublish(w, w.windowTextures[w.bufferIdx])
	w.windowTextures[w.bufferIdx].Upload(zeroPoint, w.winBuffers[w.bufferIdx], w.winBuffers[w.bufferIdx].Bounds())
	w.windowControl.Scale(w.windowRect, w.windowTextures[w.bufferIdx], w.windowTextures[w.bufferIdx].Bounds(), draw.Src)
	w.windowControl.Publish()
	// every frame, swap buffers. This enables drivers which might hold on to the rgba buffers we publish as if they
	// were immutable.
	w.bufferIdx = (w.bufferIdx + 1) % bufferCount
}

// DoBetweenDraws will execute the given function in-between draw frames
func (w *Window) DoBetweenDraws(f func()) {
	go func() {
		w.betweenDrawCh <- f
	}()
}
