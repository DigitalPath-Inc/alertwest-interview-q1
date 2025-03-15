package main

import (
	"math"
	"math/rand"
)

func SkewNorm(mu, sigma, lambda float64) float64 {
	delta := lambda / math.Sqrt(1+lambda*lambda)
	U := rand.NormFloat64()
	V := rand.NormFloat64()
	X := U*math.Sqrt(1-delta*delta) + delta*math.Abs(V)
	return mu + sigma*X
}
