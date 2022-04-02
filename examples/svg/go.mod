module github.com/oakmound/oak/examples/svg

go 1.18

require (
	github.com/oakmound/oak/v3 v3.0.0-alpha.1
	github.com/srwiley/oksvg v0.0.0-20210320200257-875f767ac39a
	github.com/srwiley/rasterx v0.0.0-20200120212402-85cb7272f5e9
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
)

replace github.com/oakmound/oak/v3 => ../..
