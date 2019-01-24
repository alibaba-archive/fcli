package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/fcli/util"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/aliyun/fc-go-sdk"
)

var client util.IClient

func getClient() (*fc.Client, error) {
	clientTemp, err := fc.NewClient(
		gConfig.Endpoint,
		gConfig.APIVersion,
		gConfig.AccessKeyID,
		gConfig.AccessKeySecret,
		fc.WithSecurityToken(gConfig.SecurityToken),
	)
	return clientTemp, err
}

// ResponseError :
type ResponseError struct {
	HTTPStatus   int32  `json:"HttpStatus"`
	RequestID    string `json:"RequestId"`
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	ErrorType    string `json:"errorType"`
	StackTrace   string `json:"stackTrace"`
}

func printStruct(content interface{}) (string, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// make error msg for user friendly
func wrapResponseError(err error) error {
	var errRes ResponseError
	errJSON := json.Unmarshal([]byte(err.Error()), &errRes)
	if errJSON != nil {
		return err
	}
	switch errRes.HTTPStatus {
	case http.StatusBadRequest:
		errRes.ErrorMessage += ",please check your region of log project, make sure the sls_endpoint " +
			"in the config match log's region"
		return fmt.Errorf("%v", errRes.ErrorMessage)
	case http.StatusForbidden:
		errRes.ErrorMessage += ",please check your role's policy has the access for these resources"
		return fmt.Errorf("%v", errRes.ErrorMessage)
	default:
		return err
	}
}

func prettyPrint(content interface{}, err error) {
	if content != nil {
		typeInp := reflect.TypeOf(content)
		if typeInp.Kind() == reflect.Struct {
			result, err := printStruct(content)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Fprintln(os.Stdout, result)
		} else {
			fmt.Fprintln(os.Stdout, content)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func prepareCommon() error {
	return nil
}

/** TRIGGER DECORATE HELPER **/

type triggerCliOutputDecorate struct {
	HTTPTriggerURL *string
}

func decorateTriggerOutput(triggerType *string, output *triggerCliOutputDecorate) {
	if *triggerType == fc.TRIGGER_TYPE_HTTP {
		temp := strings.Join([]string{
			gConfig.Endpoint,
			gConfig.APIVersion,
			"proxy",
			serviceName,
			functionName,
		}, "/")
		output.HTTPTriggerURL = &temp
	}
}
