package stateManager

import (
	"fmt"
	"github.com/EdmundMartin/greenstalk/protocol"
	"time"
)

func ManageState(changes chan *HeapValue, manager protocol.Storage) {
	heap := NewJobHeap()
	for {
		select {
		case msg := <- changes:
			fmt.Println(msg)
			heap.Insert(msg)
		case <- time.After(time.Second * 5):
			fmt.Println("We should check the heap here")
		case <- time.After(time.Second * 180):
			fmt.Println("We should remove uncollected changes")
		}
	}
}


func handleStateCheck() {

}