package ram

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSignature(t *testing.T) {
	suite.Run(t, new(SignatureTestSuite))
}

type SignatureTestSuite struct {
	suite.Suite
}

// TestCreateQueryStr ..
func (s *SignatureTestSuite) TestCreateQueryStr() {
	args := &map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	queryStr := CreateQueryStr(*args)
	s.NotNil(queryStr)
	s.Equal(`key1=value1&key2=value2`, queryStr)
}