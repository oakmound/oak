package main

import (
	"fmt"
	"path/filepath"

	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/render"
)

func TxtSetAt(f *render.Font, xpos, ypos, xadv, yadv float64, txts ...string) []render.Renderable {
	rs := make([]render.Renderable, len(txts))
	for i, txt := range txts {
		rs[i] = TxtAt(f, txt, xpos, ypos)
		xpos += xadv
		ypos += yadv
	}
	return rs
}

func TxtAt(f *render.Font, txt string, xpos, ypos float64) render.Renderable {
	return Pos(f.NewStrText(txt, 0, 0), xpos, ypos)
}

func Title(str string) render.Renderable {
	return TxtAt(Gnuolane72, str, .5, .4)
}

func Header(str string) render.Renderable {
	return TxtAt(Gnuolane72, str, .5, .2)
}

func TxtSetFrom(f *render.Font, xpos, ypos, xadv, yadv float64, txts ...string) []render.Renderable {
	rs := make([]render.Renderable, len(txts))
	for i, txt := range txts {
		rs[i] = TxtFrom(f, txt, xpos, ypos)
		xpos += xadv
		ypos += yadv
	}
	return rs
}

func TxtFrom(f *render.Font, txt string, xpos, ypos float64) render.Renderable {
	return f.NewStrText(txt, width*xpos, height*ypos)
}

func Pos(r render.Renderable, xpos, ypos float64) render.Renderable {
	XPos(r, xpos)
	YPos(r, ypos)
	return r
}

func XPos(r render.Renderable, pos float64) render.Renderable {
	w, _ := r.GetDims()
	r.SetPos(width*pos-float64(w/2), r.Y())
	return r
}

func YPos(r render.Renderable, pos float64) render.Renderable {
	_, h := r.GetDims()
	r.SetPos(r.X(), height*pos-float64(h/2))
	return r
}

func Image(file string, xpos, ypos float64) render.Modifiable {
	s, err := render.LoadSprite(filepath.Join("raw", file))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	s.SetPos(width*xpos, height*ypos)
	return s
}

func ImageCaption(file string, xpos, ypos float64, scale float64, f *render.Font, cap string) render.Renderable {
	r := Image(file, xpos, ypos)
	r.Modify(mod.Scale(scale, scale))
	w, h := r.GetDims()

	x := r.X() + float64(w)/2
	y := r.Y() + float64(h) + 28

	s := f.NewStrText(cap, x, y)
	s.Center()

	return render.NewCompositeR(r, s)
}
