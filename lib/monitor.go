package lib

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type Monitor struct {
	cpuUsage        []int
	memoryUsage     []int
	ioUsage         []int
	lastUpdate      time.Time
	lastCpu         ResourceUsage
	lastMemory      ResourceUsage
	lastIo          ResourceUsage
	updateFrequency time.Duration
	tickrate        int
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

// ResourceMetrics represents the resource utilization metrics
type ResourceMetrics struct {
	CPU       ResourceUsage `json:"cpu"`
	IO        ResourceUsage `json:"io"`
	Memory    ResourceUsage `json:"memory"`
	Timestamp int64         `json:"timestamp"`
}

func (r ResourceMetrics) MarshalZerologObject(log *zerolog.Event) {
	log.Str("CPU", fmt.Sprintf("avg: %d, min: %d, max: %d", r.CPU.Average, r.CPU.Min, r.CPU.Max))
	log.Str("IO", fmt.Sprintf("avg: %d, min: %d, max: %d", r.IO.Average, r.IO.Min, r.IO.Max))
	log.Str("Memory", fmt.Sprintf("avg: %d, min: %d, max: %d", r.Memory.Average, r.Memory.Min, r.Memory.Max))
	log.Time("Timestamp", time.UnixMilli(r.Timestamp))
}

func newMonitor(updateFrequency time.Duration, tickrate int) *Monitor {
	return &Monitor{
		cpuUsage:        make([]int, 0),
		memoryUsage:     make([]int, 0),
		ioUsage:         make([]int, 0),
		lastUpdate:      time.Now(),
		lastCpu:         ResourceUsage{0, 0, 0},
		lastMemory:      ResourceUsage{0, 0, 0},
		lastIo:          ResourceUsage{0, 0, 0},
		updateFrequency: updateFrequency,
		tickrate:        tickrate,
	}
}

func (m *Monitor) run(resourceUpdateChan <-chan ResourceUpdate) {
	ticker := time.NewTicker(m.updateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.aggregate()
		case update := <-resourceUpdateChan:
			m.update(update.CPU, update.Memory, update.IO)
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

func (m *Monitor) getResources() *ResourceMetrics {
	return &ResourceMetrics{
		Timestamp: m.lastUpdate.UnixMilli(),
		CPU:       m.lastCpu,
		Memory:    m.lastMemory,
		IO:        m.lastIo,
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
