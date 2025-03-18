package main

import "time"

type Daemon struct {
	queue              *Queue
	resourceUpdateChan chan<- ResourceUpdate
	queueEventChan     chan<- []*Execution
	tickrate           int
	scalarFunc         func(int) float64 // this allows us to have cyclic behavior, so we can simulate traffic over time
	ticks              int
}

func NewDaemon(queue *Queue, resourceUpdateChan chan<- ResourceUpdate, queueEventChan chan<- []*Execution, tickrate int, scalarFunc func(int) float64) *Daemon {
	return &Daemon{
		queue:              queue,
		resourceUpdateChan: resourceUpdateChan,
		queueEventChan:     queueEventChan,
		tickrate:           tickrate,
		scalarFunc:         scalarFunc,
		ticks:              0,
	}
}

func (d *Daemon) Run(getQueuedChan <-chan GetQueuedRequest) {
	ticker := time.NewTicker(time.Duration(1000/d.tickrate) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			queued, executed := d.queue.Tick(d.scalarFunc(d.ticks))
			d.queueEventChan <- queued
			d.resourceUpdateChan <- sumResources(executed)
			d.ticks++
		case req := <-getQueuedChan:
			queued := d.queue.GetQueued()
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
			req.ResponseChan <- res
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
