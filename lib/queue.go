package lib

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

// QueuedOperation represents a query in the queue
type QueuedOperation struct {
	Query     QueuedQuery     `json:"query"`
	Execution QueuedExecution `json:"execution"`
}

type QueuedQuery struct {
	ID string `json:"id"`
}

type QueuedExecution struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

func newQueue(queries []*Query, probs *[]float64, defaultDelay int) *Queue {
	return &Queue{
		queued:       make(map[uuid.UUID]*Execution),
		queries:      queries,
		probs:        probs,
		defaultDelay: defaultDelay,
	}
}

func (q *Queue) getQueued() []*Execution {
	queued := make([]*Execution, 0, len(q.queued))
	for _, execution := range q.queued {
		queued = append(queued, execution)
	}
	return queued
}

func (q *Queue) tick(scalar float64) ([]*Execution, []*Execution) {
	newQueries := selectExecutedQueries(q.probs, q.queries, scalar, q.defaultDelay)
	for _, query := range newQueries {
		q.queued[query.id] = query
	}

	executed := make([]*Execution, 0)
	for id, execution := range q.queued {
		execution.delay -= 1
		if execution.delay <= 0 {
			executed = append(executed, execution)
			delete(q.queued, id)
		}
	}

	return newQueries, executed
}

func (q *Queue) delay(id uuid.UUID, delay int) error {
	if _, ok := q.queued[id]; !ok {
		return fmt.Errorf("execution not found")
	}
	q.queued[id].delay += delay
	return nil
}
