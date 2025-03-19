package lib

import (
	"time"
)

type Daemon struct {
	queue              *Queue
	resourceUpdateChan chan<- ResourceUpdate
	queueEventChan     chan *QueuedQuery
	tickrate           int
	scalarFunc         func(int) float64 // this allows us to have cyclic behavior, so we can simulate traffic over time
	ticks              int
}

func newDaemon(queue *Queue, resourceUpdateChan chan<- ResourceUpdate, queueEventChan chan *QueuedQuery, tickrate int, scalarFunc func(int) float64) *Daemon {
	return &Daemon{
		queue:              queue,
		resourceUpdateChan: resourceUpdateChan,
		queueEventChan:     queueEventChan,
		tickrate:           tickrate,
		scalarFunc:         scalarFunc,
		ticks:              0,
	}
}

func (d *Daemon) run() {
	ticker := time.NewTicker(time.Duration(1000/d.tickrate) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			queued, executed := d.queue.tick(d.scalarFunc(d.ticks))
			d.queueEvent(queued)
			d.resourceUpdateChan <- sumResources(executed)
			d.ticks++
		}
	}
}

func (d *Daemon) getQueued() []*QueuedQuery {
	queued := d.queue.getQueued()
	res := make([]*QueuedQuery, 0, len(queued))
	for _, q := range queued {
		res = append(res, &QueuedQuery{
			Query: struct {
				ID string `json:"id"`
			}{
				ID: q.query.id.String(),
			},
			Execution: struct {
				ID        string `json:"id"`
				Timestamp int64  `json:"timestamp"`
			}{
				ID:        q.id.String(),
				Timestamp: time.Now().Add(time.Duration(q.delay) * time.Millisecond).UnixMilli(),
			},
		})
	}
	return res
}

func (d *Daemon) queueEvent(queueUpdate []*Execution) {
	for _, execution := range queueUpdate {
		offset := time.Duration(float64(execution.delay)/float64(d.tickrate)) * time.Second

		// Create the QueuedQuery object
		queuedQuery := &QueuedQuery{
			Query: struct {
				ID string `json:"id"`
			}{
				ID: execution.query.id.String(),
			},
			Execution: struct {
				ID        string `json:"id"`
				Timestamp int64  `json:"timestamp"`
			}{
				ID:        execution.id.String(),
				Timestamp: time.Now().Add(offset).UnixMilli(),
			},
		}

		// Try to send to channel, but if it's full, remove oldest item first
		select {
		case d.queueEventChan <- queuedQuery:
			// Successfully sent
		default:
			// Channel is full, remove oldest item and then send
			<-d.queueEventChan
			d.queueEventChan <- queuedQuery
		}
	}
}

func sumResources(executions []*Execution) ResourceUpdate {
	cpu := 0
	memory := 0
	io := 0
	for _, execution := range executions {
		cpu += execution.query.cpuUsage
		memory += execution.query.memoryUsage
		io += execution.query.ioUsage
	}
	return ResourceUpdate{cpu, memory, io}
}
