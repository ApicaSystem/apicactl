package utils

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
)

type MockResponse struct {
	Url        string
	StatusCode int32
	Body       string
	Err        error
	HttpMethod string
}

type TestCase struct {
	Name          string
	Input         interface{}
	Expected      interface{}
	MockResponses []MockResponse
	OutputFormat  map[string]string
}

func SetupTestCase(t *testing.T) func(t *testing.T) {
	viper.Set("cluster", "localhost")
	viper.Set("uitoken", "dummy_token")
	return func(t *testing.T) {}
}

func SetupOutputFormat(outputFormat map[string]string) {
	for key, value := range outputFormat {
		if key == "time-format" {
			utils.FlagTimeFormat = value
		}
	}
}

var DefaultOutputFormat map[string]string = map[string]string{
	"time-format": "RFC3339",
}

func MockApiResponse(mockResponseList []MockResponse) {
	for _, mockResponse := range mockResponseList {
		body, _ := ioutil.ReadFile(mockResponse.Body)
		url := fmt.Sprintf("https://%s/%s", viper.GetString("cluster"), mockResponse.Url)
		httpmock.RegisterResponder(mockResponse.HttpMethod, url, httpmock.NewStringResponder(int(mockResponse.StatusCode), string(body)))
	}
}
