package cmd

var funcCodeDir string

var showCheckBasicsResults string

type functionFields struct {
	Description     string
	Runtime         string
	Handler         string
	Timeout         int32
	Memory          int32
	OSSBucketName   string
	OSSObjectName   string
	VMInstanceType  string
	VMInstanceCount int32
}

type listCondition struct {
	Prefix    string
	StartKey  string
	NextToken string
	Limit     int32
}
