package mods

import "github.com/oakmound/oak/alg/intgeom"

type Dir intgeom.Point2

var (
	Up        = Dir(intgeom.Point2{0, -1})
	Down      = Dir(intgeom.Point2{0, 1})
	Left      = Dir(intgeom.Point2{-1, 0})
	Right     = Dir(intgeom.Point2{1, 0})
	UpRight   = Up.And(Right)
	DownRight = Down.And(Right)
	DownLeft  = Down.And(Left)
	UpLeft    = Up.And(Left)
)

func (d Dir) And(d2 Dir) Dir {
	return Dir(intgeom.Point2(d).Add(intgeom.Point2(d2)))
}

func (d Dir) X() int {
	return intgeom.Point2(d).X()
}

func (d Dir) Y() int {
	return intgeom.Point2(d).Y()
}
