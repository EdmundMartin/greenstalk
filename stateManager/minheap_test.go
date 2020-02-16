package stateManager

import "testing"

func TestJobHeap_Insert(t *testing.T) {
	h := NewJobHeap()
	hv := &HeapValue{
		jobID:     1,
		unixStamp: 100,
		status:    "",
	}
	hv2 := &HeapValue{
		jobID:     0,
		unixStamp: 75,
		status:    "",
	}
	h.Insert(hv)
	h.Insert(hv2)
	result := h.Peek()
	if result.unixStamp != 75 {
		t.Errorf("incorrect heapValue returned")
	}
}

func TestJobHeap_Remove(t *testing.T) {
	h := NewJobHeap()
	hv := &HeapValue{
		jobID:     0,
		unixStamp: 75,
		status:    "",
	}
	hv2 := &HeapValue{
		jobID:     0,
		unixStamp: 100,
		status:    "",
	}
	hv3 := &HeapValue{
		jobID:     0,
		unixStamp: 125,
		status:    "",
	}
	h.Insert(hv)
	h.Insert(hv2)
	h.Insert(hv3)
	removed := h.Remove()
	if removed != hv {
		t.Errorf("incorrect value removed from heap")
	}
	peeked := h.Peek()
	if peeked.unixStamp != 100 {
		t.Errorf("heap peeked value incorrect")
	}
}
