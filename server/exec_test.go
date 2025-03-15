package main

func (s *TestSuite) TestGetExecutionProbs() {
	probs := getExecutionProbs(100)
	mean := 0.0
	for _, prob := range *probs {
		mean += prob
	}
	mean /= float64(len(*probs))
	s.InDelta(mean, 0.01, 0.01)
}

func (s *TestSuite) TestSelectExecutedIdx() {
	probs := getExecutionProbs(100)
	executed := selectExecutedIdx(probs, 1)
	s.InDelta(len(executed), 1, 2)
}
