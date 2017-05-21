package collision

import (
	"bitbucket.org/oakmoundstudio/oak/event"
	higgs "github.com/Sythe2o0/rtreego"
	//"github.com/dhconnelly/rtreego"
	"math/rand"
	"testing"
	"time"
)

// DHconnely impl.
// func BenchmarkRTree(b *testing.B) {
// 	curSeed := time.Now().UTC().UnixNano()
// 	rand.Seed(curSeed)
// 	rt := rtreego.NewTree(2, 20, 40)
// 	var j event.CID = 0
// 	for i := 0; i < 500; i++ {

// 		loc := NewSpace(100*rand.Float64(), 100*rand.Float64(), rand.Float64(), rand.Float64(), j)
// 		rt.Insert(loc)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {

// 		r, _ := rtreego.NewRect(rtreego.Point{100 * rand.Float64(), 100 * rand.Float64()}, []float64{10, 10})
// 		rt.SearchIntersect(r)
// 	}
// }

// patrick-higgins / our fork
func Benchmark2DRTree(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)
	rt2 := higgs.NewTree(20, 40)
	var j event.CID
	for i := 0; i < 500; i++ {

		loc := NewSpace(100*rand.Float64(), 100*rand.Float64(), rand.Float64(), rand.Float64(), j)
		rt2.Insert(loc)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		r, _ := higgs.NewRect(higgs.Point{100 * rand.Float64(), 100 * rand.Float64()}, [3]float64{10, 10, 0})
		rt2.SearchIntersect(&r)
	}
}
