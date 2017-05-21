package render

import (
	"container/heap"
	"testing"
)

func BenchmarkHeapAddOne(b *testing.B) {
	rh := new(RenderableHeap)
	for n := 0; n < b.N; n++ {
		heap.Push(rh, new(Sprite))
	}
}

func BenchmarkLambdaAddOne(b *testing.B) {
	lh := new(lambdaHeap)
	for n := 0; n < b.N; n++ {
		lh.Push(new(Sprite))
	}
}

func BenchmarkHeapAdd100(b *testing.B) {
	rh := new(RenderableHeap)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			heap.Push(rh, new(Sprite))
		}
	}
}

func BenchmarkLambdaAdd100(b *testing.B) {
	lh := new(lambdaHeap)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			lh.Push(new(Sprite))
		}
	}
}

func BenchmarkHeapPullAll(b *testing.B) {
	rh := new(RenderableHeap)
	for i := 100; i > 0; i-- {
		heap.Push(rh, new(Sprite))
	}
	for n := 0; n < b.N; n++ {
		for rh.Len() > 0 {
			heap.Pop(rh)
		}
		for i := 100; i > 0; i-- {
			heap.Push(rh, new(Sprite))
		}
	}
}

func BenchmarkLambdaPullAll(b *testing.B) {
	lh := new(lambdaHeap)
	for i := 100; i > 0; i-- {
		lh.Push(new(Sprite))
	}
	for n := 0; n < b.N; n++ {
		for len(lh.bh) > 0 {
			lh.Pop()
		}
		for i := 100; i > 0; i-- {
			lh.Push(new(Sprite))
		}
	}
}
