package lib

import (
	"math/rand"

	"github.com/google/uuid"
)

type Execution struct {
	query *Query
	id    uuid.UUID // id is the unique identifier for the execution
	delay int       // delay is the number of ticks before the query is executed
}

// getExecutionProbs returns the probability of execution at a given tick for each query
func getExecutionProbs(n int) *[]float64 {
	probs := make([]float64, n)
	for i := 0; i < n; i++ {
		probs[i] = rand.ExpFloat64() / float64(n)
	}
	return &probs
}

// selectExecutedIdx returns the indices of the queries that are executed at a given tick
// scalar is the scalar by which the probabilities are multiplied, to provide the opportunity
// to control the number of queries executed at a given tick
func selectExecutedIdx(probs *[]float64, scalar float64) []int {
	executed := make([]int, 0)
	for i := 0; i < len(*probs); i++ {
		if rand.Float64() < (*probs)[i]*scalar {
			executed = append(executed, i)
		}
	}
	return executed
}

func selectExecutedQueries(probs *[]float64, queries []*Query, scalar float64, delay int) []*Execution {
	executed := selectExecutedIdx(probs, scalar)
	executedQueries := make([]*Execution, len(executed))
	for i, idx := range executed {
		executedQueries[i] = &Execution{
			query: queries[idx],
			id:    uuid.New(),
			delay: delay,
		}
	}
	return executedQueries
}
