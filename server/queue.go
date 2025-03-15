package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Queue struct {
	queued       map[uuid.UUID]*Execution // queued is a map of executed queries to their remaining time in the queue in ticks
	queries      []*Query                 // queries is the list of possible queries
	probs        *[]float64               // probs is the list of probabilities that a given query is selected
	defaultDelay int                      // defaultDelay is the default delay of a query in ticks
}

func NewQueue(queries []*Query, probs *[]float64, defaultDelay int) *Queue {
	return &Queue{
		queued:       make(map[uuid.UUID]*Execution),
		queries:      queries,
		probs:        probs,
		defaultDelay: defaultDelay,
	}
}

func (q *Queue) GetQueued() []*Execution {
	queued := make([]*Execution, 0, len(q.queued))
	for _, execution := range q.queued {
		queued = append(queued, execution)
	}
	return queued
}

func (q *Queue) Tick(scalar float64) ([]*Execution, int, int, int) {
	newQueries := selectExecutedQueries(q.probs, q.queries, scalar, q.defaultDelay)
	for _, query := range newQueries {
		q.queued[query.id] = query
	}

	executed := make([]*Execution, 0)
	cpuUsage := 0
	memoryUsage := 0
	ioUsage := 0
	for id, execution := range q.queued {
		execution.delay -= 1
		if execution.delay <= 0 {
			executed = append(executed, execution)
			delete(q.queued, id)
			cpuUsage += execution.query.cpuUsage
			memoryUsage += execution.query.memoryUsage
			ioUsage += execution.query.ioUsage
		}
	}

	return executed, cpuUsage, memoryUsage, ioUsage
}

func (q *Queue) Delay(id uuid.UUID, delay int) error {
	if _, ok := q.queued[id]; !ok {
		return fmt.Errorf("execution not found")
	}
	q.queued[id].delay += delay
	return nil
}
