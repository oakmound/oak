//go:build js

package audio

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"syscall/js"

	"github.com/oakmound/oak/v4/oakerr"
)

func initOS() error {
	return nil
}

var processorIndex int32

func newWriter(f Format) (Writer, error) {
	if f.Bits != 32 {
		return nil, oakerr.InvalidInput{
			InputName: "f.Bits",
		}
	}
	window := js.Global()
	actxConstruct := window.Get("AudioContext")
	if actxConstruct.IsUndefined() || actxConstruct.IsNull() {
		actxConstruct = window.Get("webkitAudioContext")
	}
	audioCtx := actxConstruct.New(map[string]interface{}{
		"latencyHint": "interactive",
		"sampleRate":  f.SampleRate,
	})

	processorName := "oakPCM" + strconv.Itoa(int(atomic.Add(&processorIndex, 1)))
	window.Call("registerProcessor", processorName, "js class?")
	audioCtx.Get("audioWorklet").Call("addModule", processorName)

	audioBuffer := audioCtx.Call("createBuffer", f.Channels, f.SampleRate*WriterBufferLengthInSeconds, f.SampleRate)
	source := audioCtx.Call("createBufferSource")

	channelData := make([]js.Value, f.Channels)
	for i := 0; i < int(f.Channels); i++ {
		channelData[i] = audioBuffer.Call("getChannelData", i)

	}

	return &jsWriter{
		Format:      f,
		bufferSize:  f.BytesPerSecond() * WriterBufferLengthInSeconds,
		audioCtx:    audioCtx,
		buffer:      audioBuffer,
		channelData: channelData,
		source:      source,
	}, nil
}

type jsWriter struct {
	sync.Mutex
	Format
	buffer       js.Value
	channelData  []js.Value // Float32Array
	source       js.Value
	audioCtx     js.Value
	lockedOffset uint32
	bufferSize   uint32
	writeChannel int
	writeOffset  int
	playing      bool
}

func (jsw *jsWriter) Close() error {
	jsw.Lock()
	defer jsw.Unlock()

	// we can't release this object?
	if jsw.playing {
		jsw.source.Call("stop")
	}
	return nil
}

func (jsw *jsWriter) Reset() error {
	jsw.Lock()
	defer jsw.Unlock()
	// emptyBuff := make([]byte, jsw.bufferSize)
	// a, b, err := jsw.buff.LockBytes(0, jsw.bufferSize, 0)
	// if err != nil {
	// 	return err
	// }
	// copy(a, emptyBuff)
	// if len(b) != 0 {
	// 	copy(b, emptyBuff)
	// }
	// err = jsw.buff.UnlockBytes(a, b)
	// jsw.Seek(0, io.SeekStart)

	//jsw.audioBuffer.Call("copyToChannel")
	// make it a []float32 array somehow // then a byte array? or just convert
	// from byte to float32 adaptively
	return nil
}

func (jsw *jsWriter) WritePCM(data []byte) (n int, err error) {
	jsw.Lock()
	defer jsw.Unlock()

	// we cannot write less than four bytes -- float32
	readAt := 0
	for len(data[readAt:]) >= 4 {
		u32 := uint32(data[readAt]) +
			uint32(data[readAt+1])<<8 +
			uint32(data[readAt+2])<<16 +
			uint32(data[readAt+3])<<24
		f32 := float32(u32) / float32(math.MaxInt32)

		jsw.channelData[jsw.writeChannel].SetIndex(jsw.writeOffset, f32)

		readAt += 4
		jsw.writeChannel++
		jsw.writeChannel %= int(jsw.Channels)
		if jsw.writeChannel == 0 {
			jsw.writeOffset++
			if jsw.writeOffset >= int((jsw.bufferSize/4)/uint32(jsw.Channels)) {
				jsw.writeOffset = 0
			}
		}
	}

	jsw.source.Set("buffer", jsw.buffer)
	if !jsw.playing {
		fmt.Println("start playing")
		jsw.playing = true
		jsw.source.Set("loop", true)
		jsw.source.Call("connect", jsw.audioCtx.Get("destination"))
		jsw.source.Call("start")
	}

	return readAt, nil
}
