package shape

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/oakmound/oak/v3/oakerr"
)

const randTestCt = 100

func TestBezierCurve(t *testing.T) {
	t.Parallel()
	t.Run("EvenInputs", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			ct := (rand.Intn(50) * 2) + 2
			floats := make([]float64, ct)
			_, err := BezierCurve(floats...)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		}
	})
	t.Run("MatchesPoint", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			x := rand.Float64()
			y := rand.Float64()
			b, err := BezierCurve(x, y)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			bp, ok := b.(BezierPoint)
			if !ok {
				t.Fatalf("expected BezierPoint, got %T", b)
			}
			expectedBP := BezierPoint{x, y}
			if bp != expectedBP {
				t.Fatalf("expected point %+v, got %+v", expectedBP, bp)
			}
		}
	})
	t.Run("MatchesNode", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			x1 := rand.Float64()
			y1 := rand.Float64()
			x2 := rand.Float64()
			y2 := rand.Float64()
			b, err := BezierCurve(x1, y1, x2, y2)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			bn, ok := b.(BezierNode)
			if !ok {
				t.Fatalf("expected BezierNode, got %T", b)
			}
			expectedLeft := BezierPoint{x1, y1}
			left, ok := bn.Left.(BezierPoint)
			if !ok {
				t.Fatalf("expected left of bezier to be BezierNode, got %T", bn.Left)
			}
			if left != expectedLeft {
				t.Fatalf("expected left point %+v, got %+v", expectedLeft, left)
			}
			expectedRight := BezierPoint{x2, y2}
			right, ok := bn.Right.(BezierPoint)
			if !ok {
				t.Fatalf("expected right of bezier to be BezierNode, got %T", bn.Right)
			}
			if right != expectedRight {
				t.Fatalf("expected right point %+v, got %+v", expectedRight, right)
			}
		}
	})
}

func TestBezierPoint_Pos(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Parallel()
	for i := 0; i < randTestCt; i++ {
		x := rand.Float64()
		y := rand.Float64()
		bp := BezierPoint{x, y}
		gotX, gotY := bp.Pos(rand.Float64())
		if x != gotX {
			t.Fatalf("expected x value of %v, got %v", x, gotX)
		}
		if y != gotY {
			t.Fatalf("expected y value of %v, got %v", y, gotY)
		}
	}
}

func TestBezierNode_Pos(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Parallel()

	t.Run("VerticalLine", func(t *testing.T) {
		t.Parallel()
		bn := BezierNode{
			BezierPoint{0, 0},
			BezierPoint{0, 1},
		}
		for i := 0; i < randTestCt; i++ {
			r := rand.Float64()
			x, y := bn.Pos(r)
			if x != 0 {
				t.Fatalf("expected x value of %v, got %v", 0, x)
			}
			if y != r {
				t.Fatalf("expected y value of %v, got %v", r, y)
			}
		}
	})
	t.Run("HorizontalLine", func(t *testing.T) {
		t.Parallel()
		bn := BezierNode{
			BezierPoint{0, 0},
			BezierPoint{1, 0},
		}
		for i := 0; i < randTestCt; i++ {
			r := rand.Float64()
			x, y := bn.Pos(r)
			if x != r {
				t.Fatalf("expected x value of %v, got %v", r, x)
			}
			if y != 0 {
				t.Fatalf("expected y value of %v, got %v", 0, y)
			}
		}
	})
	t.Run("DiagonalLine", func(t *testing.T) {
		t.Parallel()
		bn := BezierNode{
			BezierPoint{0, 0},
			BezierPoint{1, 1},
		}
		for i := 0; i < randTestCt; i++ {
			r := rand.Float64()
			x, y := bn.Pos(r)
			if x != r {
				t.Fatalf("expected x value of %v, got %v", r, x)
			}
			if y != r {
				t.Fatalf("expected y value of %v, got %v", r, y)
			}
		}
	})
}

func TestBezierCurveErrors(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Parallel()
	t.Run("ZeroInputs", func(t *testing.T) {
		t.Parallel()
		_, err := BezierCurve()
		if err == nil {
			t.Fatalf("expected non-nil error")
		}
		insufficient := &oakerr.InsufficientInputs{}
		if !errors.As(err, insufficient) {
			t.Fatalf("expected insufficient error, got %T", err)
		}
		if insufficient.AtLeast != 2 {
			t.Fatalf("expected at least to be '2', got %v", insufficient.AtLeast)
		}
		if insufficient.InputName != "coords" {
			t.Fatalf("expected input name to be 'coords', got %v", insufficient.InputName)
		}
	})
	t.Run("UnevenInputs", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < randTestCt; i++ {
			ct := (rand.Intn(50) * 2) + 1
			floats := make([]float64, ct)
			_, err := BezierCurve(floats...)
			indivisible := &oakerr.IndivisibleInput{}
			if !errors.As(err, indivisible) {
				t.Fatalf("expected indivisible error, got %T", err)
			}
			if indivisible.MustDivideBy != 2 {
				t.Fatalf("expected must divide by to be '2', got %v", indivisible.MustDivideBy)
			}
			if indivisible.InputName != "coords" {
				t.Fatalf("expected input name to be 'coords', got %v", indivisible.InputName)
			}
		}
	})
}
