package ram

import "encoding/json"

// ServiceError :
type ServiceError struct {
	HTTPStatus   int    `json:"HttpStatus"`
	RequestID    string `json:"RequestId"`
	HostID       string `json:"HostId"`
	ErrorCode    string `json:"Code"`
	ErrorMessage string `json:"Message"`
}

// String :
func (e ServiceError) String() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

// Error :
func (e ServiceError) Error() string {
	return e.String()
}
