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
	defaultDelay := 10                                        // 100ms default delay
	metricsUpdateFrequency := time.Second * time.Duration(30) // 30 second metrics update frequency

	queue := NewQueue(queries, probs, defaultDelay)
	monitor := NewMonitor(metricsUpdateFrequency)
	tickrate := 10 // 10 ticks per second

	daemon := NewDaemon(queue, monitor, tickrate, scalarFunc)
	go daemon.Run()
	server := NewServer(queue, monitor, tickrate)
	server.Start(":8080") // Removed 'go' to prevent main from exiting
}
