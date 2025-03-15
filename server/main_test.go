package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
