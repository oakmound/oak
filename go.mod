module github.com/oakmound/oak/v2

require (
	github.com/200sc/go-dist v1.0.0
	github.com/200sc/klangsynthese v0.2.0
	github.com/BurntSushi/toml v0.3.1
	github.com/akavel/polyclip-go v0.0.0-20160111220610-2cfdb71461bd
	github.com/davecgh/go-spew v1.1.1
	github.com/disintegration/gift v1.2.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/go-mp3 v0.2.1 // indirect
	github.com/oakmound/libudev v0.2.1
	github.com/oakmound/oak v2.0.0+incompatible
	github.com/oakmound/shiny v0.4.1-0.20191119013337-fdd972eb9250
	github.com/oakmound/w32 v2.1.0+incompatible
	github.com/oov/directsound-go v0.0.0-20141101201356-e53e59c700bf // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/image v0.0.0-20190227222117-0694c2d4d067
	golang.org/x/mobile v0.0.0-20190415191353-3e0bab5405d6
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6
)

replace github.com/oakmound/oak => ../oak

go 1.13
