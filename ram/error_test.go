package ram

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestServiceError(t *testing.T) {
	suite.Run(t, new(ServiceErrorTestSuite))
}

type ServiceErrorTestSuite struct {
	suite.Suite
}

// TestError ..
func (s *ServiceErrorTestSuite) TestError() {
	serviceError := &ServiceError{
		HTTPStatus: 500,
		RequestID: "1",
		HostID: "1",
		ErrorCode: "500",
		ErrorMessage: "Error",
	}

	errMsg := serviceError.Error()
	s.NotNil(errMsg)
	s.Equal(`{
    "HttpStatus": 500,
    "RequestId": "1",
    "HostId": "1",
    "Code": "500",
    "Message": "Error"
}`, errMsg)
}