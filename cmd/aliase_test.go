package cmd

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAlias(t *testing.T) {
	suite.Run(t, new(AliasTestSuite))
}

type AliasTestSuite struct {
	suite.Suite
}
