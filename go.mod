module github.com/oakmound/oak/v3

go 1.18

require (
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20201218220906-28db891af037 // osx, shiny
	github.com/BurntSushi/xgb v0.0.0-20210121224620-deaf085860bc // linux, shiny
	github.com/BurntSushi/xgbutil v0.0.0-20190907113008-ad855c713046 // linux, shiny
	github.com/disintegration/gift v1.2.0 // render
	github.com/eaburns/flac v0.0.0-20171003200620-9a6fb92396d1
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20210410170116-ea3d685f79fb // osx, shiny
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/go-mp3 v0.3.1
	github.com/jfreymuth/pulse v0.1.0 // linux, audio
	github.com/oakmound/alsa v0.0.2 // linux, audio
	github.com/oakmound/libudev v0.2.1 // linux, joystick
	github.com/oakmound/w32 v2.1.0+incompatible // windows, shiny
	github.com/oov/directsound-go v0.0.0-20141101201356-e53e59c700bf // windows, audio
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6
	golang.org/x/mobile v0.0.0-20220112015953-858099ff7816
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220403205710-6acee93ad0eb
)

require (
	github.com/eaburns/bit v0.0.0-20131029213740-7bd5cd37375d // indirect
	golang.org/x/exp v0.0.0-20220328175248-053ad81199eb // indirect
)
