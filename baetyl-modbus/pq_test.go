package main

import (
	"container/heap"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestPQ(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	t0:= time.Unix(100, 0)
	tasks := []*task {
		&task {
			runTime: t0.Add(time.Second * 3),
		},
		&task {
			runTime: t0.Add(time.Second),
		},
		&task {
			runTime: t0.Add(time.Second * 2),
		},
	}
	for _, ta := range tasks {
		heap.Push(&pq, ta)
	}
	expect := []*task {
		&task {
			runTime: t0.Add(time.Second),
		},
		&task {
			runTime: t0.Add(time.Second * 2),
		},
		&task {
			runTime: t0.Add(time.Second * 3),
		},
	}

	i := 0
	for pq.Len() > 0 {
		ta := heap.Pop(&pq).(*task)
		assert.Equal(t, ta.runTime, expect[i].runTime)
		i++
	}
}


