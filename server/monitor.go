package main

import (
	"time"
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
}

type ResourceUsage struct {
	Average int `json:"average"`
	Min     int `json:"min"`
	Max     int `json:"max"`
}

func NewMonitor(updateFrequency time.Duration) *Monitor {
	return &Monitor{
		cpuUsage:        make([]int, 0),
		memoryUsage:     make([]int, 0),
		ioUsage:         make([]int, 0),
		lastUpdate:      time.Now(),
		lastCpu:         ResourceUsage{0, 0, 0},
		lastMemory:      ResourceUsage{0, 0, 0},
		lastIo:          ResourceUsage{0, 0, 0},
		updateFrequency: updateFrequency,
	}
}

func (m *Monitor) Update(cpuUsage int, memoryUsage int, ioUsage int) {
	m.cpuUsage = append(m.cpuUsage, cpuUsage)
	m.memoryUsage = append(m.memoryUsage, memoryUsage)
	m.ioUsage = append(m.ioUsage, ioUsage)
}

func (m *Monitor) GetResourceUsage() (time.Time, ResourceUsage, ResourceUsage, ResourceUsage) {
	if time.Since(m.lastUpdate) > m.updateFrequency {
		cpuUsage := getResourceStats(m.cpuUsage)
		memoryUsage := getResourceStats(m.memoryUsage)
		ioUsage := getResourceStats(m.ioUsage)
		m.lastUpdate = time.Now()
		m.lastCpu = cpuUsage
		m.lastMemory = memoryUsage
		m.lastIo = ioUsage
		m.cpuUsage = make([]int, 0)
		m.memoryUsage = make([]int, 0)
		m.ioUsage = make([]int, 0)
	}
	return m.lastUpdate, m.lastCpu, m.lastMemory, m.lastIo
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
