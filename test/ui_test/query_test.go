package ui_test

import (
	"encoding/json"
	"fmt"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	test_utils "github.com/logiqai/logiqctl/test/utils"
)

func TestCreateQueryResult(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		{
			Input:    test_utils.BASE_TEST_DATA_DIR + "/data/query/create_query_result/input-1.json",
			Expected: test_utils.BASE_TEST_DATA_DIR + "/data/query/create_query_result/output-1.json",
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/query_results",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/query_results/create/200.json",
					Err:        nil,
				},
				{
					Url:        "api/jobs/job_1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/jobs/get/200.json",
					Err:        nil,
				},
				{
					Url:        "api/query_results/1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/query_results/get/200_response.json",
					Err:        nil,
				},
			},
		},
	}
	httpmock.Activate()
	defer httpmock.Deactivate()
	for _, testcase := range testcases {
		payload, err := ioutil.ReadFile(testcase.Input.(string))
		test_utils.MockApiResponse(&testcase.MockResponses)
		if err != nil {
			fmt.Println("Error Reading input file")
			os.Exit(-1)
		}
		queryResult := types.QueryResult{}
		err = json.Unmarshal(payload, &queryResult)
		if err != nil {
			fmt.Printf("Unable to convert payload to query result %s \n", err.Error())
			os.Exit(-1)
		}
		actual, err := ui.ExecuteQuery(queryResult)
		if err != nil {
			assert.Equal(t, testcase.Expected.(string), err.Error(), testcase.Name)
		} else {
			expected, err := ioutil.ReadFile(testcase.Expected.(string))
			if err != nil {
				fmt.Println("Error Reading Expected output file")
				os.Exit(-1)
			}
			assert.JSONEq(t, string(expected), actual, testcase.Name)
		}
	}
}
