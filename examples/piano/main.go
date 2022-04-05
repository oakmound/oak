package main

import (
	"context"
	"fmt"
	"image/color"
	"image/draw"
	"math"
	"os"
	"sync"
	"time"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/audio"
	"github.com/oakmound/oak/v3/audio/pcm"
	"github.com/oakmound/oak/v3/audio/synth"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

const (
	whiteKeyWidth  = 26
	whiteKeyHeight = 200
	blackKeyWidth  = 13
	blackKeyHeight = 140

	whiteBlackOverlap = 5

	labelWhiteKey = 0
	labelBlackKey = 1
)

type keyColor int

const keyColorWhite keyColor = 0
const keyColorBlack keyColor = 1

func (kc keyColor) Width() float64 {
	if kc == keyColorBlack {
		return blackKeyWidth
	}
	return whiteKeyWidth
}

func (kc keyColor) Height() float64 {
	if kc == keyColorBlack {
		return blackKeyHeight
	}
	return whiteKeyHeight
}

func (kc keyColor) Color() color.RGBA {
	if kc == keyColorBlack {
		return color.RGBA{60, 60, 60, 255}
	}
	return color.RGBA{255, 255, 255, 255}
}

func newKey(ctx *scene.Context, note synth.Pitch, c keyColor, k key.Code) *entities.Solid {
	w := c.Width()
	h := c.Height()
	clr := c.Color()
	downClr := clr
	downClr.R -= 60
	downClr.B -= 60
	downClr.G -= 60
	sw := render.NewSwitch("up", map[string]render.Modifiable{
		"up": render.NewCompositeM(
			render.NewColorBox(int(w), int(h), clr),
			render.NewLine(0, 0, 0, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(0, h, w, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, h, w, 0, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, 0, 0, 0, color.RGBA{0, 0, 0, 255}),
		).ToSprite(),
		"down": render.NewCompositeM(
			render.NewColorBox(int(w), int(h), downClr),
			render.NewLine(0, 0, 0, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(0, h, w, h, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, h, w, 0, color.RGBA{0, 0, 0, 255}),
			render.NewLine(w, 0, 0, 0, color.RGBA{0, 0, 0, 255}),
		).ToSprite(),
	})
	s := entities.NewSolid(0, 0, w, h, sw, mouse.DefaultTree, 0)
	if c == keyColorBlack {
		s.Space.SetZLayer(1)
		s.Space.Label = labelBlackKey
	} else {
		s.Space.SetZLayer(2)
		s.Space.Label = labelWhiteKey
	}
	mouse.UpdateSpace(s.X(), s.Y(), s.W, s.H, s.Space)
	event.GlobalBind(ctx, key.Down(k), func(ev key.Event) event.Response {
		// TODO: add helper function for this?
		if ev.Modifiers&key.ModShift == key.ModShift {
			return 0
		}
		playPitch(ctx, note)
		sw.Set("down")
		return 0
	})
	event.GlobalBind(ctx, key.Up(k), func(ev key.Event) event.Response {
		if ev.Modifiers&key.ModShift == key.ModShift {
			return 0
		}
		releasePitch(note)
		sw.Set("up")
		return 0
	})
	event.Bind(ctx, mouse.PressOn, s, func(_ *entities.Solid, me *mouse.Event) event.Response {
		playPitch(ctx, note)
		me.StopPropagation = true
		sw.Set("down")
		return 0
	})
	event.Bind(ctx, mouse.Release, s, func(_ *entities.Solid, me *mouse.Event) event.Response {
		releasePitch(note)
		sw.Set("up")
		return 0
	})
	return s
}

type keyDef struct {
	color keyColor
	pitch synth.Pitch
	x     float64
}

var keycharOrder = []key.Code{
	key.Z, key.S, key.X, key.D, key.C,
	key.V, key.G, key.B, key.H, key.N, key.J, key.M,
	key.Comma, key.L, key.FullStop, key.Semicolon, key.Slash,
	key.Q, key.Num2, key.W, key.Num3, key.E, key.Num4, key.R,
	key.T, key.Num6, key.Y, key.Num7, key.U,
	key.I, key.Num9, key.O, key.Num0, key.P, key.HyphenMinus, key.LeftSquareBracket,
}

var playLock sync.Mutex
var cancelFuncs = map[synth.Pitch]func(){}

var synthKind synth.Wave

func playPitch(ctx *scene.Context, pitch synth.Pitch) {
	playLock.Lock()
	defer playLock.Unlock()
	if cancel, ok := cancelFuncs[pitch]; ok {
		cancel()
	}
	a := synthKind(synth.AtPitch(pitch))
	toPlay := audio.LoopReader(a)
	format := toPlay.PCMFormat()
	speaker, err := audio.NewWriter(format)
	if err != nil {
		fmt.Println("new writer failed:", err)
		return
	}
	monitor := newPCMMonitor(ctx, speaker)
	monitor.SetPos(0, 0)
	render.Draw(monitor)
	gctx, cancel := context.WithCancel(ctx)
	go func() {
		fadeIn := audio.FadeIn(100*time.Millisecond, toPlay)
		err = audio.Play(gctx, monitor, fadeIn)
		if err != nil {
			fmt.Println("play error:", err)
		}
		speaker.Close()
		monitor.Undraw()
	}()
	cancelFuncs[pitch] = cancel
}

func releasePitch(pitch synth.Pitch) {
	playLock.Lock()
	defer playLock.Unlock()
	if cancel, ok := cancelFuncs[pitch]; ok {
		cancel()
		delete(cancelFuncs, pitch)
	}
}

func main() {
	err := audio.InitDefault()
	if err != nil {
		fmt.Println("init failed:", err)
		os.Exit(1)
	}

	oak.AddScene("piano", scene.Scene{
		Start: func(ctx *scene.Context) {
			src := synth.Int16
			src.Format = pcm.Format{
				SampleRate: 80000,
				Channels:   2,
				Bits:       32,
			}
			synthKind = src.Sin
			pitch := synth.C3
			kc := keyColorWhite
			x := 20.0
			y := 200.0
			i := 0
			for i < len(keycharOrder) && x+kc.Width() < float64(ctx.Window.Width()-10) {
				ky := newKey(ctx, pitch, kc, keycharOrder[i])
				ky.SetPos(x, y)
				layer := 0
				if kc == keyColorBlack {
					layer = 1
				}
				render.Draw(ky.R, layer)
				x += kc.Width()
				pitch = pitch.Up(synth.HalfStep)
				if pitch.IsAccidental() {
					x -= whiteBlackOverlap
					kc = keyColorBlack
				} else if kc != keyColorWhite {
					x -= whiteBlackOverlap
					kc = keyColorWhite
				}
				i++
			}
			// Consider: Adding volume control
			codeKinds := map[key.Code]func(src synth.Source) synth.Wave{
				key.S: func(src synth.Source) synth.Wave { return src.Sin },
				key.W: func(src synth.Source) synth.Wave { return src.Saw },
				key.T: func(src synth.Source) synth.Wave { return src.Triangle },
				key.P: func(src synth.Source) synth.Wave { return src.Pulse(2) },
			}
			for kc, synfn := range codeKinds {
				kc := kc
				synfn := synfn
				event.GlobalBind(ctx, key.Down(kc), func(ev key.Event) event.Response {
					if ev.Modifiers&key.ModShift == key.ModShift {
						synthKind = synfn(src)
					}
					return 0
				})
			}

			help1 := render.NewText("Shift+([S]in/[T]ri/[P]ulse/sa[W]) to change wave style", 10, 500)
			help2 := render.NewText("Keyboard / mouse to play", 10, 520)
			render.Draw(help1)
			render.Draw(help2)

			event.GlobalBind(ctx, mouse.ScrollDown, func(_ *mouse.Event) event.Response {
				mag := globalMagnification - 0.05
				if mag < 1 {
					mag = 1
				}
				globalMagnification = mag
				return 0
			})
			event.GlobalBind(ctx, mouse.ScrollUp, func(_ *mouse.Event) event.Response {
				globalMagnification += 0.05
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad0), func(_ key.Event) event.Response {
				// TODO: synth all sound like pulse waves at 8 bit
				src.Bits = 8
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad1), func(_ key.Event) event.Response {
				src.Bits = 16
				return 0
			})
			event.GlobalBind(ctx, key.Down(key.Keypad2), func(_ key.Event) event.Response {
				src.Bits = 32
				return 0
			})
		},
	})
	oak.Init("piano", func(c oak.Config) (oak.Config, error) {
		c.Screen.Height = 600
		c.Title = "Piano Example"
		c.Debug.Level = dlog.INFO.String()
		return c, nil
	})
}

type pcmMonitor struct {
	event.CallerID
	render.LayeredPoint
	pcm.Writer
	pcm.Format
	written []byte
	at      int
}

var globalMagnification float64 = 1

func newPCMMonitor(ctx *scene.Context, w pcm.Writer) *pcmMonitor {
	fmt := w.PCMFormat()
	pm := &pcmMonitor{
		Writer:       w,
		Format:       w.PCMFormat(),
		LayeredPoint: render.NewLayeredPoint(0, 0, 0),
		written:      make([]byte, int(float64(fmt.BytesPerSecond())*audio.WriterBufferLengthInSeconds)),
	}
	return pm
}

func (pm *pcmMonitor) CID() event.CallerID {
	return pm.CallerID
}

func (pm *pcmMonitor) PCMFormat() pcm.Format {
	return pm.Format
}

func (pm *pcmMonitor) WritePCM(b []byte) (n int, err error) {
	copy(pm.written[pm.at:], b)
	if len(b) > len(pm.written[pm.at:]) {
		copy(pm.written[0:], b[len(pm.written[pm.at:]):])
	}
	pm.at += len(b)
	pm.at %= len(pm.written)
	return pm.Writer.WritePCM(b)
}

func (pm *pcmMonitor) Draw(buf draw.Image, xOff, yOff float64) {
	const width = 640
	const height = 200.0
	xJump := len(pm.written) / width
	xJump = int(float64(xJump) / globalMagnification)
	c := color.RGBA{255, 255, 255, 255}
	for x := 0.0; x < width; x++ {
		wIndex := int(x) * xJump

		var val int16
		switch pm.Format.Bits {
		case 8:
			val8 := pm.written[wIndex]
			val = int16(val8) << 8
		case 16:
			wIndex -= wIndex % 2
			val = int16(pm.written[wIndex+1])<<8 +
				int16(pm.written[wIndex])
		case 32:
			wIndex = wIndex - wIndex%4
			val32 := int32(pm.written[wIndex+3])<<24 +
				int32(pm.written[wIndex+2])<<16 +
				int32(pm.written[wIndex+1])<<8 +
				int32(pm.written[wIndex])
			val = int16(val32 / int32(math.Pow(2, 16)))
		}

		// -32768 -> 200
		// 0 -> 100
		// 32768 -> 0
		var y float64
		if val < 0 {
			y = height/2 + float64(val)*float64(height/2/-32768.0)
		} else {
			y = height/2 + -(float64(val) * float64(height/2/32768.0))
		}
		buf.Set(int(x+xOff+pm.X()), int(y+yOff+pm.Y()), c)
	}
}
