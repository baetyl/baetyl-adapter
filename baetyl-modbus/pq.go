package main

// A PriorityQueue implements heap.Interface and holds task.
type PriorityQueue []*task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].runTime.Sub(pq[j].runTime) < 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	t := x.(*task)
	t.index = n
	*pq = append(*pq, t)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // avoid memory leak
	task.index = -1 // for safety
	*pq = old[0 : n-1]
	return task
}
