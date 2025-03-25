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
}

func scalarFunc(ticks int) float64 {
	scalar := math.Sin((float64(ticks)-500)*math.Pi/10000)/2 + 1.5
	return scalar
}

func NewDB() *DB {
	queries := getQueries(100)
	probs := getExecutionProbs(100)
	defaultDelay := 1                                        // 1 tick / 100ms default delay
	metricsUpdateFrequency := time.Second * time.Duration(5) // 5 second metrics update frequency
	tickrate := 10                                           // 10 ticks per second

	// Channels for event comms between components
	resourceUpdateChan := make(chan ResourceUpdate, 100) // Handles the resource updates from daemon --> monitor

	// Create components
	queue := newQueue(queries, probs, defaultDelay)
	daemon := newDaemon(queue, resourceUpdateChan, tickrate, scalarFunc)
	monitor := newMonitor(metricsUpdateFrequency, tickrate)

	return &DB{
		queue:              queue,
		daemon:             daemon,
		monitor:            monitor,
		resourceUpdateChan: resourceUpdateChan,
	}
}

func (d *DB) Run() {
	// Start components
	go d.monitor.run(d.resourceUpdateChan)
	go d.daemon.run()
}

func (d *DB) AddQueueListener(listener chan *QueuedOperation) {
	d.daemon.addQueueListener(listener)
}

func (d *DB) GetQueued() []*QueuedOperation {
	return d.daemon.getQueued()
}

func (d *DB) GetResources() *ResourceMetrics {
	return d.monitor.getResources()
}

func (d *DB) Delay(id uuid.UUID, delay int) error {
	return d.queue.delay(id, delay)
}
