package main

import (
	"context"
	"fmt"
	"image/color"
	"image/draw"
	"os"
	"sync"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/audio/klang"
	"github.com/oakmound/oak/v3/audio/pcm"
	"github.com/oakmound/oak/v3/audio/synth"
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

func newKey(note synth.Pitch, c keyColor, k string) *entities.Solid {
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
	s.Bind(key.Down+k, func(c event.CID, i interface{}) int {
		if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
			return 0
		}
		playPitch(note)
		sw.Set("down")
		return 0
	})
	s.Bind(key.Up+k, func(c event.CID, i interface{}) int {
		if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
			return 0
		}
		releasePitch(note)
		sw.Set("up")
		return 0
	})
	s.Bind(mouse.PressOn, func(c event.CID, i interface{}) int {
		playPitch(note)
		me := i.(*mouse.Event)
		me.StopPropagation = true
		sw.Set("down")
		return 0
	})
	s.Bind(mouse.Release, func(c event.CID, i interface{}) int {
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

var keycharOrder = []string{
	"Z", "S", "X", "D", "C",
	"V", "G", "B", "H", "N", "J", "M",
	key.Comma, "L", key.Period, key.Semicolon, key.Slash,
	"Q", "2", "W", "3", "E", "4", "R",
	"T", "6", "Y", "7", "U",
	"I", "9", "O", "0", "P", key.HyphenMinus, key.LeftSquareBracket,
}

var playLock sync.Mutex
var cancelFuncs = map[synth.Pitch]func(){}

var synthKind func(...synth.Option) (pcm.Reader, error)

func playPitch(pitch synth.Pitch) {
	playLock.Lock()
	defer playLock.Unlock()
	if cancel, ok := cancelFuncs[pitch]; ok {
		cancel()
	}
	a, _ := synthKind(synth.AtPitch(pitch))
	toPlay := pcm.LoopReader(a)
	format := toPlay.PCMFormat()
	speaker, err := pcm.NewWriter(format)
	if err != nil {
		fmt.Println("writer failed:", err)
		return
	}
	monitor := newPCMMonitor(speaker)
	monitor.SetPos(0, 0)
	render.Draw(monitor)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = pcm.Play(ctx, monitor, toPlay)
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
	err := pcm.Init()
	if err != nil {
		fmt.Println("init failed:", err)
		os.Exit(1)
	}

	oak.AddScene("piano", scene.Scene{
		Start: func(ctx *scene.Context) {
			src := synth.Int16
			src.Format = klang.Format{
				SampleRate: 40000,
				Channels:   2,
				Bits:       16,
			}
			synthKind = src.SawPCM
			pitch := synth.C3
			kc := keyColorWhite
			x := 20.0
			y := 200.0
			i := 0
			for i < len(keycharOrder) && x+kc.Width() < float64(ctx.Window.Width()-10) {
				ky := newKey(pitch, kc, keycharOrder[i])
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
			event.GlobalBind(key.Down+key.S, func(c event.CID, i interface{}) int {
				if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
					synthKind = src.SinPCM
				}
				return 0
			})
			event.GlobalBind(key.Down+key.W, func(c event.CID, i interface{}) int {
				if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
					synthKind = src.SawPCM
				}
				return 0
			})
			event.GlobalBind(key.Down+key.T, func(c event.CID, i interface{}) int {
				if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
					synthKind = src.TrianglePCM
				}
				return 0
			})
			event.GlobalBind(key.Down+key.P, func(c event.CID, i interface{}) int {
				if oak.IsDown(key.LeftShift) || oak.IsDown(key.RightShift) {
					synthKind = src.PulsePCM(2)
				}
				return 0
			})
			help1 := render.NewText("Shift+([S]in/[T]ri/[P]ulse/sa[W]) to change wave style", 10, 500)
			help2 := render.NewText("Keyboard / mouse to play", 10, 520)
			render.Draw(help1)
			render.Draw(help2)

			event.GlobalBind(mouse.ScrollDown, func(c event.CID, i interface{}) int {
				mag := globalMagnification - 0.05
				if mag < 1 {
					mag = 1
				}
				globalMagnification = mag
				return 0
			})
			event.GlobalBind(mouse.ScrollUp, func(c event.CID, i interface{}) int {
				globalMagnification += 0.05
				return 0
			})
		},
	})
	oak.Init("piano", func(c oak.Config) (oak.Config, error) {
		c.Screen.Height = 600
		c.Title = "Piano Example"
		return c, nil
	})
}

type pcmMonitor struct {
	event.CID
	render.LayeredPoint
	pcm.Writer
	written []byte
	at      int
}

var globalMagnification float64 = 1

func newPCMMonitor(w pcm.Writer) *pcmMonitor {
	fmt := w.PCMFormat()
	pm := &pcmMonitor{
		Writer:       w,
		LayeredPoint: render.NewLayeredPoint(0, 0, 0),
		written:      make([]byte, fmt.BytesPerSecond()*pcm.WriterBufferLengthInSeconds),
	}
	pm.Init()
	return pm
}

func (pm *pcmMonitor) Init() event.CID {
	pm.CID = event.NextID(pm)
	return pm.CID
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

		// look for pair for this int16
		var val int16
		switch wIndex % 4 {
		case 0, 2:
			val = int16(pm.written[wIndex+1])<<8 +
				int16(pm.written[wIndex])
		case 1, 3:
			val = int16(pm.written[wIndex])<<8 +
				int16(pm.written[wIndex-1])
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
