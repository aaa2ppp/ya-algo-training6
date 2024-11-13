package main

import (
	"container/heap"
	"log"
)

type Item struct {
	id         int
	a, b       int
	idxA, idxB int
}

type QueueA []*Item

func (pq QueueA) Len() int { return len(pq) }

func (pq QueueA) Less(i, j int) bool {
	return pq[i].a > pq[j].a ||
		pq[i].a == pq[j].a && (pq[i].b > pq[j].b ||
			pq[i].b == pq[j].b && pq[i].id < pq[j].id)
}

func (pq QueueA) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idxA = i
	pq[j].idxA = j
}

func (pq *QueueA) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.idxA = n
	*pq = append(*pq, item)
}

func (pq *QueueA) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // don't stop the GC from reclaiming the item eventually
	item.idxA = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type QueueB []*Item

func (pq QueueB) Len() int { return len(pq) }

func (pq QueueB) Less(i, j int) bool {
	return pq[i].b > pq[j].b ||
		pq[i].b == pq[j].b && (pq[i].a > pq[j].a ||
			pq[i].a == pq[j].a && pq[i].id < pq[j].id)
}

func (pq QueueB) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idxB = i
	pq[j].idxB = j
}

func (pq *QueueB) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.idxB = n
	*pq = append(*pq, item)
}

func (pq *QueueB) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // don't stop the GC from reclaiming the item eventually
	item.idxB = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func solve1(aa, bb, pp []int) []int {
	n := len(aa)

	items := make([]Item, n)
	for i := 0; i < n; i++ {
		items[i] = Item{
			id:   i + 1,
			a:    aa[i],
			b:    bb[i],
			idxA: i,
			idxB: i,
		}
	}

	qa := make(QueueA, n)
	qb := make(QueueB, n)

	for i := 0; i < n; i++ {
		qa[i] = &items[i]
		qb[i] = &items[i]
	}

	heap.Init(&qa)
	heap.Init(&qb)

	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if debugEnable {
			log.Println(items)
		}

		switch pp[i] {
		case 0:
			it := heap.Pop(&qa).(*Item)
			res = append(res, it.id)
			heap.Remove(&qb, it.idxB)
		case 1:
			it := heap.Pop(&qb).(*Item)
			res = append(res, it.id)
			heap.Remove(&qa, it.idxA)
		}
	}

	if debugEnable {
		log.Println(items)
	}

	return res
}
