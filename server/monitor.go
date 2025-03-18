package main

import (
	"time"
)

type Monitor struct {
	cpuUsage         []int
	memoryUsage      []int
	ioUsage          []int
	lastUpdate       time.Time
	lastCpu          ResourceUsage
	lastMemory       ResourceUsage
	lastIo           ResourceUsage
	updateFrequency  time.Duration
	tickrate         int
	monitorEventChan chan<- *QueuedQuery
}

type ResourceUsage struct {
	Average int `json:"average"`
	Min     int `json:"min"`
	Max     int `json:"max"`
}

type ResourceUpdate struct {
	CPU    int
	Memory int
	IO     int
}

func NewMonitor(updateFrequency time.Duration, tickrate int, monitorEventChan chan<- *QueuedQuery) *Monitor {
	return &Monitor{
		cpuUsage:         make([]int, 0),
		memoryUsage:      make([]int, 0),
		ioUsage:          make([]int, 0),
		lastUpdate:       time.Now(),
		lastCpu:          ResourceUsage{0, 0, 0},
		lastMemory:       ResourceUsage{0, 0, 0},
		lastIo:           ResourceUsage{0, 0, 0},
		updateFrequency:  updateFrequency,
		tickrate:         tickrate,
		monitorEventChan: monitorEventChan,
	}
}

func (m *Monitor) Run(resourceUpdateChan <-chan ResourceUpdate, queueEventChan <-chan []*Execution, getResourcesChan <-chan GetResourcesRequest) {
	ticker := time.NewTicker(m.updateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.aggregate()
		case update := <-resourceUpdateChan:
			m.update(update.CPU, update.Memory, update.IO)
		case queueUpdate := <-queueEventChan:
			m.queueEvent(queueUpdate)
		case req := <-getResourcesChan:
			lastUpdate, lastCpu, lastMemory, lastIo := m.getResources()
			req.ResponseChan <- &ResourceMetrics{
				Timestamp: lastUpdate.UnixMilli(),
				CPU:       lastCpu,
				Memory:    lastMemory,
				IO:        lastIo,
			}
		}
	}
}

func (m *Monitor) aggregate() {
	m.lastCpu = getResourceStats(m.cpuUsage)
	m.lastMemory = getResourceStats(m.memoryUsage)
	m.lastIo = getResourceStats(m.ioUsage)
	m.lastUpdate = time.Now()
	m.cpuUsage = make([]int, 0)
	m.memoryUsage = make([]int, 0)
	m.ioUsage = make([]int, 0)
}

func (m *Monitor) update(cpuUsage int, memoryUsage int, ioUsage int) {
	m.cpuUsage = append(m.cpuUsage, cpuUsage)
	m.memoryUsage = append(m.memoryUsage, memoryUsage)
	m.ioUsage = append(m.ioUsage, ioUsage)
}

func (m *Monitor) getResources() (time.Time, ResourceUsage, ResourceUsage, ResourceUsage) {
	return m.lastUpdate, m.lastCpu, m.lastMemory, m.lastIo
}

func (m *Monitor) queueEvent(queueUpdate []*Execution) {
	for _, execution := range queueUpdate {
		offset := time.Duration(float64(execution.delay)/float64(m.tickrate)) * time.Second
		m.monitorEventChan <- &QueuedQuery{
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
	}
}

func getResourceStats(usage []int) ResourceUsage {
	if len(usage) == 0 {
		return ResourceUsage{0, 0, 0}
	}

	min := usage[0]
	max := usage[0]
	sum := 0

	for _, val := range usage {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
		sum += val
	}

	avg := sum / len(usage)
	return ResourceUsage{avg, min, max}
}
