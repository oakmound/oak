package main

import (
	"fmt"

	gokey "golang.org/x/mobile/event/key"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/entities/x/btn"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"

	"github.com/atotto/clipboard"
)

func main() {
	oak.AddScene("clipboard-test", scene.Scene{
		Start: func(ctx *scene.Context) {
			newClipboardCopyText("click-me-to-copy", 20, 20)
			newClipboardCopyText("click-to-copy-me-too", 20, 50)
			newClipboardPaster("click-or-ctrl+v-to-paste-here", 20, 200)
		},
	})
	oak.Init("clipboard-test")
}

func newClipboardCopyText(text string, x, y float64) {
	btn.New(
		btn.Font(render.DefaultFont()),
		btn.Text(text),
		btn.Pos(x, y),
		btn.Height(20),
		btn.FitText(20),
		btn.Click(func(event.CallerID, interface{}) int {
			err := clipboard.WriteAll(text)
			if err != nil {
				fmt.Println(err)
			}
			return 0
		}),
	)
}

func newClipboardPaster(placeholder string, x, y float64) {
	textPtr := new(string)
	*textPtr = placeholder
	btn.New(
		btn.Font(render.DefaultFont()),
		btn.TextPtr(textPtr),
		btn.Pos(x, y),
		btn.Height(20),
		btn.FitText(20),
		btn.Binding(key.Down+key.V, func(_ event.CallerID, payload interface{}) int {
			kv := payload.(key.Event)
			if kv.Modifiers&gokey.ModControl == gokey.ModControl {
				got, err := clipboard.ReadAll()
				if err != nil {
					fmt.Println(err)
					return 0
				}
				*textPtr = got
			}
			return 0
		}),
		btn.Click(func(event.CallerID, interface{}) int {
			got, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println(err)
				return 0
			}
			*textPtr = got
			return 0
		}),
	)
}
