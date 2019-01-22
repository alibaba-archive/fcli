package ram

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestModel(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}

type ModelTestSuite struct {
	suite.Suite
}
