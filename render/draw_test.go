package render

import (
	"container/heap"
	"testing"
)

type HeapNode struct {
	r    Renderable
	next *HeapNode
}

func (hn *HeapNode) Push(r Renderable) {
	if hn.r.GetLayer() >= r.GetLayer() {
		if hn.next == nil || hn.next.r.GetLayer() < r.GetLayer() {
			nhn := new(HeapNode)
			nhn.next = hn.next
			nhn.r = r
			hn.next = nhn
		} else {
			hn.next.Push(r)
		}
	} else {
		nhn := new(HeapNode)
		nhn.next = hn
		nhn.r = r
	}
}

func BenchmarkHeapAddOne(b *testing.B) {
	rh := new(RenderableHeap)
	for n := 0; n < b.N; n++ {
		heap.Push(rh, new(Sprite))
	}
}

func BenchmarkNodeAddOne(b *testing.B) {
	hn := new(HeapNode)
	hn.r = new(Sprite)
	for n := 0; n < b.N; n++ {
		hn.Push(new(Sprite))
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

func BenchmarkNodeAdd100(b *testing.B) {
	hn := new(HeapNode)
	hn.r = new(Sprite)
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			hn.Push(new(Sprite))
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

func BenchmarkNodePullAll(b *testing.B) {
	hn := new(HeapNode)
	hn.r = new(Sprite)
	for i := 100; i > 0; i-- {
		hn.Push(new(Sprite))
	}
	for n := 0; n < b.N; n++ {
		for hn != nil {
			hn = hn.next
		}
	}
}
