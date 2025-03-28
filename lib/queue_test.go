package lib

import (
	"encoding/json"
	"os"
	"sort"
)

func (s *TestSuite) TestQueue() {
	// Create test queries
	queries := getQueries(100)
	probs := getExecutionProbs(100)

	// Create a queue with a tickrate of 10 and default delay of 5 ticks
	queue := newQueue(queries, probs, 5)

	// Test initial state
	s.Empty(queue.queued)
	s.Equal(100, len(queue.queries))
	s.Equal(100, len(*queue.probs))
	s.Equal(5, queue.defaultDelay)

	// Track resource usage over time
	cpuCounts := make(map[int]int)
	memoryCounts := make(map[int]int)
	ioCounts := make(map[int]int)

	numTicks := 1000
	usageHistory := make(map[string][]int)
	usageHistory["cpu"] = make([]int, numTicks)
	usageHistory["memory"] = make([]int, numTicks)
	usageHistory["io"] = make([]int, numTicks)

	// Run for 1000 ticks
	for i := 0; i < numTicks; i++ {
		scalar := 1.0

		// Tick the queue
		queued, executed := queue.tick(scalar)

		s.InDelta(len(queued), 1, 5) // 0-6 queries queued per tick

		summed := sumResources(executed)

		usageHistory["cpu"][i] = summed.CPU
		usageHistory["memory"][i] = summed.Memory
		usageHistory["io"][i] = summed.IO

		// Record resource usage
		cpuCounts[summed.CPU]++
		memoryCounts[summed.Memory]++
		ioCounts[summed.IO]++
	}

	// Save usage history to disk for analysis
	file, err := json.MarshalIndent(usageHistory, "", "  ")
	if err == nil {
		err = os.WriteFile("queue_usage_history.json", file, 0644)
		s.NoError(err, "Failed to write usage history to file")
	} else {
		s.Fail("Failed to marshal usage history to JSON")
	}

	meanCpuUsage := calculateMean(cpuCounts, 1000)
	meanMemoryUsage := calculateMean(memoryCounts, 1000)
	meanIoUsage := calculateMean(ioCounts, 1000)

	// Calculate p95 (95th percentile)
	p95CpuUsage := calculateP95(cpuCounts, 1000)
	p95MemoryUsage := calculateP95(memoryCounts, 1000)
	p95IoUsage := calculateP95(ioCounts, 1000)

	// Verify resource usage is within expected ranges
	// These are more relaxed bounds since queue behavior is different from tick
	s.InDelta(meanCpuUsage, 50, 25)
	s.InDelta(meanMemoryUsage, 50, 25)
	s.InDelta(meanIoUsage, 50, 25)
	s.InDelta(p95CpuUsage, 150, 50)
	s.InDelta(p95MemoryUsage, 150, 50)
	s.InDelta(p95IoUsage, 150, 50)
}

// Helper function to calculate mean
func calculateMean(counts map[int]int, totalSamples int) int {
	total := 0
	for usage, count := range counts {
		total += usage * count
	}
	return total / totalSamples
}

// Helper function to calculate 95th percentile
func calculateP95(counts map[int]int, totalSamples int) int {
	values := make([]int, 0, len(counts))
	for usage := range counts {
		values = append(values, usage)
	}
	sort.Ints(values)

	total := 0
	p95Threshold := int(float64(totalSamples) * 0.95)

	for _, usage := range values {
		total += counts[usage]
		if total >= p95Threshold {
			return usage
		}
	}

	return 0
}
