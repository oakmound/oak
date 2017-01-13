package oak

import (
	"fmt"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/draw"
	"testing"
)

// func BenchmarkMakeNewBuffer(b *testing.B) {
// 	driver.Main(func(s screen.Screen) {
// 		s.NewWindow(&screen.NewWindowOptions{640, 480})
// 		for n := 0; n < b.N; n++ {
// 			s.NewBuffer(image.Point{10000, 10000})
// 		}
// 		//defer w.Release()
// 	})
// }

func BenchmarkFillBuffer(b *testing.B) {
	driver.Main(func(s screen.Screen) {

		//s.NewWindow(&screen.NewWindowOptions{640, 480})
		b2, err := s.NewBuffer(image.Point{10000, 10000})
		fmt.Println(b2, err)
		for n := 0; n < b.N; n++ {
			draw.Draw(b2.RGBA(), b2.Bounds(), image.Black, image.Point{0, 0}, screen.Src)
		}
		//defer w.Release()
	})
}
