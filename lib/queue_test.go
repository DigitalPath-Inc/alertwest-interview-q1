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

	// Create a queue with a default delay of 10 ticks
	queue := newQueue(queries, probs, 10)

	// Test initial state
	s.Empty(queue.queued)
	s.Equal(100, len(queue.queries))
	s.Equal(100, len(*queue.probs))
	s.Equal(10, queue.defaultDelay)

	// Track resource usage over time
	cpuCounts := make(map[int]int)
	memoryCounts := make(map[int]int)
	ioCounts := make(map[int]int)

	numTicks := 1000

	type QueryExecution struct {
		QueryID string `json:"query_id"`
		CPU     int    `json:"cpu"`
		Memory  int    `json:"memory"`
		IO      int    `json:"io"`
	}

	type TickLog struct {
		Tick      int              `json:"tick"`
		NumQueued int              `json:"num_queued"`
		Executed  []QueryExecution `json:"executed"`
		TotalCPU  int              `json:"total_cpu"`
		TotalMem  int              `json:"total_memory"`
		TotalIO   int              `json:"total_io"`
	}

	executionLog := make([]TickLog, 0, numTicks)

	// Run for 1000 ticks
	for i := 0; i < numTicks; i++ {
		scalar := 1.0

		// Tick the queue
		queued, executed := queue.tick(scalar)

		summed := sumResources(executed)

		// Log tick details
		tickLog := TickLog{
			Tick:      i,
			NumQueued: len(queued),
			Executed:  make([]QueryExecution, 0, len(executed)),
			TotalCPU:  summed.CPU,
			TotalMem:  summed.Memory,
			TotalIO:   summed.IO,
		}

		for _, exec := range executed {
			tickLog.Executed = append(tickLog.Executed, QueryExecution{
				QueryID: exec.query.id.String(),
				CPU:     exec.query.cpuUsage,
				Memory:  exec.query.memoryUsage,
				IO:      exec.query.ioUsage,
			})
		}

		executionLog = append(executionLog, tickLog)

		s.InDelta(len(queued), 1, 5) // 0-6 queries queued per tick

		// Record resource usage
		cpuCounts[summed.CPU]++
		memoryCounts[summed.Memory]++
		ioCounts[summed.IO]++
	}

	// Save detailed execution log
	execLogJson, err := json.MarshalIndent(executionLog, "", "  ")
	if err == nil {
		err = os.WriteFile("queue_execution_log.json", execLogJson, 0644)
		s.NoError(err, "Failed to write execution log to file")
	} else {
		s.Fail("Failed to marshal execution log to JSON")
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
