package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/jarcoal/httpmock"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
)

type MockResponse struct {
	Url         string
	StatusCode  int32
	Body        string
	Err         error
	HttpMethod  string
	FromPayload bool
}

type TestCase struct {
	Name          string
	Input         interface{}
	Expected      interface{}
	MockResponses []MockResponse
	OutputFormat  map[string]string
}

func SetupTestCase(t *testing.T, mockResponseList *[]MockResponse) func(t *testing.T) {
	viper.Set("cluster", "dummyhost")
	viper.Set("uitoken", "dummy_token")
	httpmock.Activate()
	if mockResponseList != nil {
		MockApiResponse(mockResponseList)
	}
	return func(t *testing.T) {
		httpmock.Deactivate()
	}
}

func SetupOutputFormat(outputFormat map[string]string) {
	for key, value := range outputFormat {
		if key == "time-format" {
			utils.FlagTimeFormat = value
		}
	}
}

var DefaultOutputFormat = map[string]string{
	"time-format": "RFC3339",
}

func MockApiResponse(mockResponseList *[]MockResponse) {
	for _, mockResponse := range *mockResponseList {
		mockResponse.mock()
	}
}

func (m *MockResponse) mock() {
	body, _ := ioutil.ReadFile(m.Body)
	url := fmt.Sprintf("http://%s/%s", viper.GetString(utils.KeyCluster), m.Url)
	response := http.Response{
		StatusCode: int(m.StatusCode),
		Body:       httpmock.NewRespBodyFromString(string(body)),
		Header:     make(http.Header),
	}
	response.Header.Add("Content-Type", "application/json;utf-8")
	if !m.FromPayload || m.HttpMethod == http.MethodGet {
		httpmock.RegisterResponder(m.HttpMethod, url, httpmock.ResponderFromResponse(&response))
	} else {
		httpmock.RegisterResponder(m.HttpMethod, url,
			func(req *http.Request) (*http.Response, error) {
				payload, _ := ioutil.ReadAll(req.Body)
				defer req.Body.Close()
				respBody, _ := jsonpatch.MergePatch(payload, body)
				response.Body = httpmock.NewRespBodyFromString(string(respBody))
				return &response, nil
			},
		)
	}
}

var (
	BASE_TEST_DATA_DIR = "../../test_data"
)
