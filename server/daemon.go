package main

import "time"

type Daemon struct {
	queue      *Queue
	monitor    *Monitor
	tickrate   int
	scalarFunc func(int) float64 // this allows us to have cyclic behavior, so we can simulate traffic over time
	ticks      int
}

func NewDaemon(queue *Queue, monitor *Monitor, tickrate int, scalarFunc func(int) float64) *Daemon {
	return &Daemon{queue: queue, monitor: monitor, tickrate: tickrate, scalarFunc: scalarFunc, ticks: 0}
}

func (d *Daemon) Run() {
	for {
		time.Sleep(time.Duration(1000/d.tickrate) * time.Millisecond)
		_, cpuUsage, memoryUsage, ioUsage := d.queue.Tick(d.scalarFunc(d.ticks))
		d.monitor.Update(cpuUsage, memoryUsage, ioUsage)
		d.ticks++
	}
}
