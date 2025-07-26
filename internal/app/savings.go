package app

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
