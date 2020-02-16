package stateManager


type HeapValue struct {
	jobID int
	unixStamp int
	status string
}

type JobHeap struct {
	heap []*HeapValue
}

func NewJobHeap() *JobHeap {
	return &JobHeap{}
}

func (jh *JobHeap) isLeaf(idx int) bool {
	if idx >= len(jh.heap) / 2 && idx <= len(jh.heap) {
		return true
	}
	return false
}

func (jh *JobHeap) parent(idx int) int {
	return (idx - 1) / 2
}

func (jh *JobHeap) leftChild(idx int) int {
	return (2 * idx) + 1
}

func (jh *JobHeap) rightChild(idx int) int {
	return (2 * idx) + 2
}

func (jh *JobHeap) Insert(hv *HeapValue) {
	jh.heap = append(jh.heap, hv)
	jh.upHeapify(len(jh.heap)-1)
}

func (jh *JobHeap) swap(first, second int) {
	temp := jh.heap[first]
	jh.heap[first] = jh.heap[second]
	jh.heap[second] = temp
}

func (jh *JobHeap) upHeapify(idx int) {
	for jh.heap[idx].unixStamp < jh.heap[jh.parent(idx)].unixStamp {
		jh.swap(idx, jh.parent(idx))
	}
}

func (jh *JobHeap) downHeapify(current int) {
	if jh.isLeaf(current) {
		return
	}
	smallest := current
	leftChild := jh.leftChild(current)
	rightChild := jh.rightChild(current)
	currentSize := len(jh.heap)
	if leftChild < currentSize && jh.heap[leftChild].unixStamp < jh.heap[smallest].unixStamp {
		smallest = leftChild
	}
	if rightChild < currentSize && jh.heap[rightChild].unixStamp < jh.heap[smallest].unixStamp {
		smallest = rightChild
	}
	if smallest != current {
		jh.swap(current, smallest)
		jh.downHeapify(smallest)
	}
	return
}

func (jh *JobHeap) Remove() *HeapValue {
	top := jh.heap[0]
	jh.heap[0] = jh.heap[len(jh.heap)-1]
	jh.heap = jh.heap[:len(jh.heap)-1]
	jh.downHeapify(0)
	return top
}

func (jh *JobHeap) Peek() *HeapValue {
	return jh.heap[0]
}