package app

import (
	"container/heap"
	"fmt"
	"time"
)

type savingsList []saving

type saving struct {
	i, j  int
	value float64
}

func (h savingsList) Len() int           { return len(h) }
func (h savingsList) Less(i, j int) bool { return h[i].value > h[j].value }
func (h savingsList) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *savingsList) Push(x any) {
	*h = append(*h, x.(saving))
}

func (h *savingsList) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func (h *savingsList) popAll() {
	for h.Len() > 0 {
		e := heap.Pop(h).(saving)
		fmt.Printf("Popped: (%d, %d) = %f\n", e.i, e.j, e.value)
		time.Sleep(10 * time.Millisecond)
	}
}

func TestMaxHeap() {

	h := &savingsList{}
	heap.Init(h)

	heap.Push(h, saving{i: 0, j: 1, value: 10})
	heap.Push(h, saving{i: 2, j: 3, value: 25})
	heap.Push(h, saving{i: 1, j: 1, value: 5})
	h.popAll()

}
