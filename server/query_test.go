package main

func (s *TestSuite) TestGetQueryCPU() {
	meanCpuUsage := 0
	meanMemoryUsage := 0
	meanIoUsage := 0
	for i := 0; i < 1000; i++ {
		query := getQuery(CPU)
		meanCpuUsage += query.cpuUsage
		meanMemoryUsage += query.memoryUsage
		meanIoUsage += query.ioUsage
	}
	meanCpuUsage /= 1000
	meanMemoryUsage /= 1000
	meanIoUsage /= 1000
	s.InDelta(meanCpuUsage, 55, 10)
	s.InDelta(meanMemoryUsage, 40, 10)
	s.InDelta(meanIoUsage, 40, 10)
}

func (s *TestSuite) TestGetQueryMemory() {
	meanCpuUsage := 0
	meanMemoryUsage := 0
	meanIoUsage := 0
	for i := 0; i < 1000; i++ {
		query := getQuery(Memory)
		meanCpuUsage += query.cpuUsage
		meanMemoryUsage += query.memoryUsage
		meanIoUsage += query.ioUsage
	}
	meanCpuUsage /= 1000
	meanMemoryUsage /= 1000
	meanIoUsage /= 1000
	s.InDelta(meanCpuUsage, 40, 10)
	s.InDelta(meanMemoryUsage, 55, 10)
	s.InDelta(meanIoUsage, 40, 10)
}

func (s *TestSuite) TestGetQueryIO() {
	meanCpuUsage := 0
	meanMemoryUsage := 0
	meanIoUsage := 0
	for i := 0; i < 1000; i++ {
		query := getQuery(IO)
		meanCpuUsage += query.cpuUsage
		meanMemoryUsage += query.memoryUsage
		meanIoUsage += query.ioUsage
	}
	meanCpuUsage /= 1000
	meanMemoryUsage /= 1000
	meanIoUsage /= 1000
	s.InDelta(meanCpuUsage, 40, 10)
	s.InDelta(meanMemoryUsage, 40, 10)
	s.InDelta(meanIoUsage, 55, 10)
}

func (s *TestSuite) TestGetQueries() {
	queries := getQueries(10)
	s.Equal(10, len(queries))
	totalCpuUsage := 0
	totalMemoryUsage := 0
	totalIoUsage := 0
	for _, query := range queries {
		totalCpuUsage += query.cpuUsage
		totalMemoryUsage += query.memoryUsage
		totalIoUsage += query.ioUsage
	}
	meanCpuUsage := totalCpuUsage / len(queries)
	meanMemoryUsage := totalMemoryUsage / len(queries)
	meanIoUsage := totalIoUsage / len(queries)
	s.InDelta(meanCpuUsage, 30, 5)
	s.InDelta(meanMemoryUsage, 30, 5)
	s.InDelta(meanIoUsage, 30, 5)
}
