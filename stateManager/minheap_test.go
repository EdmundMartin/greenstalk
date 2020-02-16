package stateManager

import "testing"

func TestJobHeap_Insert(t *testing.T) {
	h := NewJobHeap()
	hv := &HeapValue{
		JobID:     1,
		UnixStamp: 100,
		Status:    "",
	}
	hv2 := &HeapValue{
		JobID:     0,
		UnixStamp: 75,
		Status:    "",
	}
	h.Insert(hv)
	h.Insert(hv2)
	result := h.Peek()
	if result.UnixStamp != 75 {
		t.Errorf("incorrect heapValue returned")
	}
}

func TestJobHeap_Remove(t *testing.T) {
	h := NewJobHeap()
	hv := &HeapValue{
		JobID:     0,
		UnixStamp: 75,
		Status:    "",
	}
	hv2 := &HeapValue{
		JobID:     0,
		UnixStamp: 100,
		Status:    "",
	}
	hv3 := &HeapValue{
		JobID:     0,
		UnixStamp: 125,
		Status:    "",
	}
	h.Insert(hv)
	h.Insert(hv2)
	h.Insert(hv3)
	removed := h.Remove()
	if removed != hv {
		t.Errorf("incorrect value removed from heap")
	}
	peeked := h.Peek()
	if peeked.UnixStamp != 100 {
		t.Errorf("heap peeked value incorrect")
	}
}
