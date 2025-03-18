package main

import (
	"math"
	"time"
)

func scalarFunc(ticks int) float64 {
	scalar := math.Sin((float64(ticks)-500)*math.Pi/1000)/2 + 1.5
	return scalar
}

func main() {
	queries := getQueries(100)
	probs := getExecutionProbs(100)
	defaultDelay := 10                                        // 10 tick / 100ms default delay
	metricsUpdateFrequency := time.Second * time.Duration(30) // 30 second metrics update frequency
	tickrate := 100                                           // 100 ticks per second

	// Channels for event comms between components
	resourceUpdateChan := make(chan ResourceUpdate, 100)    // Handles the resource updates from daemon --> monitor
	queueEventChan := make(chan []*Execution, 100)          // Handles the queue updates from daemon --> monitor
	monitorEventChan := make(chan *QueuedQuery, 100)        // Handles the queue updates from monitor --> server
	getQueuedChan := make(chan GetQueuedRequest, 100)       // Handles the synchronous get queued requests from server --> daemon
	getResourcesChan := make(chan GetResourcesRequest, 100) // Handles the synchronous get resources requests from server --> monitor
	delayChan := make(chan DelayRequest, 100)               // Handles the synchronous delay requests from server --> queue

	// Create components
	queue := NewQueue(queries, probs, defaultDelay)
	daemon := NewDaemon(queue, resourceUpdateChan, queueEventChan, tickrate, scalarFunc)
	monitor := NewMonitor(metricsUpdateFrequency, tickrate, monitorEventChan)
	server := NewServer(monitorEventChan, getQueuedChan, getResourcesChan, delayChan)

	// Start components
	go monitor.Run(resourceUpdateChan, queueEventChan, getResourcesChan)
	go daemon.Run(getQueuedChan)
	server.Start(":8080")
}
