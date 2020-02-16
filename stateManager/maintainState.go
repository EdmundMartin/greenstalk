package stateManager

import (
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"time"
)

func ManageState(changes chan *HeapValue, manager protocol.Storage) {
	heap := NewJobHeap()
	lastReset := time.After(time.Second * 180)
	for {
		select {
		case msg := <- changes:
			updateHeap(heap, msg)
		case <- time.After(time.Second * 5):
			fmt.Println("We should check the heap here")
			handleStateCheck(heap, manager)
		case <- lastReset:
			fmt.Println("We should remove uncollected changes")
			manager.Reset()
			lastReset = time.After(time.Second * 180)
		}
	}
}

func updateHeap(hp *JobHeap, msg *HeapValue) {
	jID := msg.JobID
	if msg.Status == "DELETED" {
		_, ok := hp.inHeap[jID]
		if ok {
			hp.deletedJobs[jID] = true
		}
		return
	} else {
		hp.Insert(msg)
		hp.inHeap[jID] = true
	}
}


func handleStateCheck(hp *JobHeap, storage protocol.Storage) {
	for {
		peeked := hp.Peek()
		if peeked == nil {
			return
		}
		_, ok := hp.deletedJobs[peeked.JobID]
		if ok {
			hp.Remove()
		} else if peeked.UnixStamp < time.Now().Unix() {
			fmt.Println(peeked.UnixStamp)
			storage.UpdateJob(peeked.JobID, peeked.Status)
			fmt.Println("Updated reserved")
			hp.Remove()
		} else {
			return
		}
	}
}