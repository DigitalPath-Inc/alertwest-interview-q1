package lib

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type DB struct {
	queue              *Queue
	daemon             *Daemon
	monitor            *Monitor
	resourceUpdateChan <-chan ResourceUpdate
	QueueEventChan     <-chan *QueuedQuery
}

func scalarFunc(ticks int) float64 {
	scalar := math.Sin((float64(ticks)-500)*math.Pi/10000)/2 + 1.5
	return scalar
}

func NewDB() *DB {
	queries := getQueries(100)
	probs := getExecutionProbs(100)
	defaultDelay := 10                                        // 10 tick / 100ms default delay
	metricsUpdateFrequency := time.Second * time.Duration(30) // 30 second metrics update frequency
	tickrate := 100                                           // 100 ticks per second

	// Channels for event comms between components
	resourceUpdateChan := make(chan ResourceUpdate, 100) // Handles the resource updates from daemon --> monitor
	queueEventChan := make(chan *QueuedQuery, 100)       // Handles the queue updates from daemon --> ???

	// Create components
	queue := newQueue(queries, probs, defaultDelay)
	daemon := newDaemon(queue, resourceUpdateChan, queueEventChan, tickrate, scalarFunc)
	monitor := newMonitor(metricsUpdateFrequency, tickrate)

	return &DB{
		queue:              queue,
		daemon:             daemon,
		monitor:            monitor,
		resourceUpdateChan: resourceUpdateChan,
		QueueEventChan:     queueEventChan,
	}
}

func (d *DB) Run() {
	// Start components
	go d.monitor.run(d.resourceUpdateChan)
	go d.daemon.run()
}

func (d *DB) GetQueued() []*QueuedQuery {
	return d.daemon.getQueued()
}

func (d *DB) GetResources() *ResourceMetrics {
	return d.monitor.getResources()
}

func (d *DB) Delay(id uuid.UUID, delay int) {
	d.queue.delay(id, delay)
}
