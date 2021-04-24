package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/dlog"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	oak.AddScene("blank", scene.Scene{
		Start: func(ctx *scene.Context) {
			ctx.DrawStack.Draw(render.NewDrawFPS(0, nil, 10, 10))
			ctx.DrawStack.Draw(render.NewLogicFPS(0, nil, 10, 20))
		},
	})
	oak.SetupConfig.Debug.Level = dlog.INFO.String()
	oak.Init("blank")
}
