package priority

import (
	"container/heap"
	"english-util/domain"
	"fmt"
)

type WordQueue struct {
	PQ    PriorityQueue
	Index map[string]*Item
}

func BuildQueue(wlist []domain.Item) *WordQueue {
	wq := &WordQueue{
		PQ:    make(PriorityQueue, 0, len(wlist)),
		Index: make(map[string]*Item, len(wlist)),
	}

	for i := range wlist {
		item := &Item{
			Data:     &wlist[i],
			Priority: 0,
		}
		wq.PQ = append(wq.PQ, item)
		wq.Index[wlist[i].Word] = item
	}
	heap.Init(&wq.PQ)
	return wq
}

func (wq *WordQueue) Increase(word string) {
	item := wq.Index[word]
	if item == nil {
		fmt.Println("NO SUCH WORD:", word)
		return
	}
	item.Priority++
	heap.Fix(&wq.PQ, item.Index)
}

func (wq *WordQueue) Decrease(word string) {
	item := wq.Index[word]
	if item == nil {
		return
	}
	item.Priority--
	heap.Fix(&wq.PQ, item.Index)
}
