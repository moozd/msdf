package msdf

type HeapItem struct {
	value    int
	priority float64
	index    int
}

type MaxHeap []*HeapItem

func (ca MaxHeap) Len() int { return len(ca) }

func (ca MaxHeap) Less(i, j int) bool {
	return ca[i].priority > ca[j].priority // Higher priority first
}

func (ca MaxHeap) Swap(i, j int) {
	ca[i], ca[j] = ca[j], ca[i]
	ca[i].index = i
	ca[j].index = j
}

func (ca *MaxHeap) Push(x interface{}) {
	n := len(*ca)
	item := x.(*HeapItem)
	item.index = n
	*ca = append(*ca, item)
}

func (ca *MaxHeap) Pop() interface{} {
	old := *ca
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*ca = old[0 : n-1]
	return item
}
