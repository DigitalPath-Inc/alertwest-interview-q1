package lib

func (s *TestSuite) TestSkewNorm() {
	mu := 30.0
	sigma := 15.0
	lambda := 10.0
	numSamples := 10000000

	counts := make(map[int]int)
	for i := 0; i < numSamples; i++ {
		x := SkewNorm(mu, sigma, lambda)
		counts[int(x)]++
	}

	mean := 0.0
	total := 0
	median := 0.0
	mode := 0.0
	modeCount := 0

	for x, count := range counts {
		mean += float64(x) * float64(count) / float64(numSamples)
		total += count
		if total >= numSamples/2 && median == 0 {
			median = float64(x)
		}
		if count > modeCount {
			mode = float64(x)
			modeCount = count
		}
	}

	s.InDelta(mean, 40, 10)
	s.InDelta(median, 40, 10)
	s.InDelta(mode, mu, 3)
}
