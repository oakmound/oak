package show

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/render"
)

var (
	width, height float64
)

func SetDims(w, h float64) {
	width = w
	height = h
}

var (
	titleFont *render.Font
)

func SetTitleFont(f *render.Font) {
	titleFont = f
}

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
	return TxtAt(titleFont, str, .5, .4)
}

func Header(str string) render.Renderable {
	return TxtAt(titleFont, str, .5, .2)
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

func ImageAt(file string, xpos, ypos float64, mods ...mod.Mod) render.Modifiable {
	m := Image(file, xpos, ypos)
	m.Modify(mods...)
	w, h := m.GetDims()
	m.ShiftX(float64(-w / 2))
	m.ShiftY(float64(-h / 2))
	return m
}

func ImageCaptionSize(file string, xpos, ypos float64, w, h float64, f *render.Font, cap string) render.Renderable {
	r, err := render.LoadSprite(filepath.Join("raw", file))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	w2, h2 := r.GetDims()
	w3 := float64(w2) / width
	h3 := float64(h2) / height
	wScale := w / w3
	hScale := h / h3
	if wScale > hScale {
		wScale = hScale
	} else {
		hScale = wScale
	}
	r.Modify(mod.Scale(wScale, hScale))
	w4, h4 := r.GetDims()
	r.SetPos(width*xpos, height*ypos)
	r.ShiftX(float64(-w4 / 2))
	r.ShiftY(float64(-h4 / 2))

	x := r.X() + float64(w4)/2
	y := r.Y() + float64(h4) + 42

	caps := strings.Split(cap, "\n")
	for i := 1; i < len(caps); i++ {
		// remove whitespace
		caps[i] = strings.TrimSpace(caps[i])
	}
	s := TxtSetAt(f, float64(x)/width, float64(y)/height, 0, .04, caps...)

	return render.NewCompositeR(append(s, r)...)
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
