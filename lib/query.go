package lib

import (
	"github.com/google/uuid"
)

type Query struct {
	id          uuid.UUID
	cpuUsage    int
	memoryUsage int
	ioUsage     int
}

// Profile is what the query execution is bound by
type Profile int

const (
	CPU Profile = iota
	IO
	Memory
)

func getQuery(profile Profile) *Query {
	query := Query{
		id: uuid.New(),
	}
	switch profile {
	case CPU:
		query.cpuUsage = int(SkewNorm(70, 15, -10))
		query.memoryUsage = int(SkewNorm(30, 15, 10))
		query.ioUsage = int(SkewNorm(30, 15, 10))
	case IO:
		query.cpuUsage = int(SkewNorm(30, 15, 10))
		query.memoryUsage = int(SkewNorm(30, 15, 10))
		query.ioUsage = int(SkewNorm(70, 15, -10))
	case Memory:
		query.cpuUsage = int(SkewNorm(30, 15, 10))
		query.memoryUsage = int(SkewNorm(70, 15, -10))
		query.ioUsage = int(SkewNorm(30, 15, 10))
	}
	return &query
}

func getQueries(n int) []*Query {
	queries := make([]*Query, n)
	for i := 0; i < n; i++ {
		queries[i] = getQuery(Profile(i % 3))
	}
	return queries
}
