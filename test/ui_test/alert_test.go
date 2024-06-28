package ui_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	test_utils "github.com/ApicaSystem/apicactl/test/utils"
	"github.com/ApicaSystem/apicactl/types"
	"github.com/ApicaSystem/apicactl/ui"
	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

func TestListAlertsCommand(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)

	defer tearDownTestcase(t)
	testCases := []test_utils.TestCase{
		{
			Name: "List All Alerts",
			Expected: []types.Alert{
				{
					Id:   1,
					Name: "table: column >= 150",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  150,
					},
					State:         "ok",
					LastTriggered: "2022-08-31T18:25:42.792Z",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
				{
					Id:   2,
					Name: "table: column >= 100",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  100,
					},
					State:         "unknown",
					LastTriggered: "",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
			},
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/list/success.json",
					Err:        nil,
				},
			},
		},
		{
			Name:     "List Alerts when no alerts are created",
			Expected: []types.Alert{},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/list/empty.json",
					Err:        nil,
				},
			},
			OutputFormat: test_utils.DefaultOutputFormat,
		},
		{
			Name:     "List Alerts when there is an error in api response",
			Expected: "Error: Internal server error",
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodGet,
					StatusCode: 500,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/list/500_error.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "List all alerts with epoch time format",
			Expected: []types.Alert{
				{
					Id:   1,
					Name: "table: column >= 150",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  150,
					},
					State:         "ok",
					LastTriggered: "1661970342",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
				{
					Id:   2,
					Name: "table: column >= 100",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  100,
					},
					State:         "unknown",
					LastTriggered: "",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
			},
			OutputFormat: map[string]string{
				"time-format": "epoch",
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/list/success.json",
					Err:        nil,
				},
			},
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test_utils.SetupOutputFormat(testCase.OutputFormat)
			test_utils.MockApiResponse(&testCase.MockResponses)

			alertsList, err := ui.ListAlerts()
			actual, _ := json.Marshal(&alertsList)
			if err != nil {
				expected := testCase.Expected
				assert.Equal(t, expected, err.Error(), testCase.Name)
			} else {
				expected, _ := json.Marshal(&testCase.Expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.Name)
			}
		})
	}
}

func TestGetAlertCommand(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)
	defer tearDownTestcase(t)

	// testcases
	testCases := []test_utils.TestCase{
		{
			Name: "Get Alert By Id",
			Input: map[string]string{
				"id": "1",
			},
			Expected: types.Alert{
				Id:   1,
				Name: "table: column >= 150",
				AlertOption: types.AlertOption{
					Column: "column",
					Op:     ">=",
					Value:  150,
				},
				State:         "ok",
				LastTriggered: "2022-08-31T18:25:42.792Z",
				Rearm:         3,
				Query: types.Query{
					Id:                1,
					LatestQueryDataId: 1,
					Name:              "table",
					Description:       "",
					Query:             "select column from table;",
					DataSourceId:      3,
					Version:           1,
					Tags:              []string{},
					QueryOptions: types.QueryOptions{
						Parameters: []map[string]interface{}{},
					},
				},
				QueryId: 1,
			},
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts/1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/get/success.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Get Alert By Id for which doest not exist",
			Input: map[string]string{
				"id": "1",
			},
			Expected:     "Error: Alert does not exist",
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts/1",
					HttpMethod: http.MethodGet,
					StatusCode: 404,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/get/404_error.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Get Alert By Id for which the api fails",
			Input: map[string]string{
				"id": "1",
			},
			Expected:     "Error: Internal server error",
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts/1",
					HttpMethod: http.MethodGet,
					StatusCode: 500,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/get/500_error.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Get alert by id with epcoh time format",
			Input: map[string]string{
				"id": "1",
			},
			Expected: types.Alert{
				Id:   1,
				Name: "table: column >= 150",
				AlertOption: types.AlertOption{
					Column: "column",
					Op:     ">=",
					Value:  150,
				},
				State:         "ok",
				LastTriggered: "1661970342",
				Rearm:         3,
				Query: types.Query{
					Id:                1,
					LatestQueryDataId: 1,
					Name:              "table",
					Description:       "",
					Query:             "select column from table;",
					DataSourceId:      3,
					Version:           1,
					Tags:              []string{},
					QueryOptions: types.QueryOptions{
						Parameters: []map[string]interface{}{},
					},
				},
				QueryId: 1,
			},
			OutputFormat: map[string]string{
				"time-format": "epoch",
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts/1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/get/success.json",
					Err:        nil,
				},
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test_utils.SetupOutputFormat(testCase.OutputFormat)
			id := testCase.Input.(map[string]string)["id"]
			test_utils.MockApiResponse(&testCase.MockResponses)

			alert, err := ui.GetAlert(id)
			actual, _ := json.Marshal(&alert)
			if err != nil {
				expected := testCase.Expected
				assert.Equal(t, expected, err.Error(), testCase.Name)
			} else {
				expected, _ := json.Marshal(&testCase.Expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.Name)
			}
		})
	}
}

func TestCreateAlertCommand(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)
	defer tearDownTestcase(t)

	// testcases
	testCases := []test_utils.TestCase{
		{
			Name: "Create Alert",
			Input: []map[string]interface{}{
				{
					"options": map[string]interface{}{
						"column": "impression_count",
						"op":     "<=",
						"value":  150,
					},
					"name":     "table: column >= 150",
					"query_id": 15,
				},
				{
					"options": map[string]interface{}{
						"column": "impression_count",
						"op":     "<=",
						"value":  150,
					},
					"name":     "table: column >= 150",
					"query_id": 15,
				},
			},
			Expected: []types.Alert{
				{
					Id:   1,
					Name: "table: column >= 150",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  150,
					},
					State:         "unknown",
					LastTriggered: "",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
				{
					Id:   1,
					Name: "table: column >= 150",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  150,
					},
					State:         "unknown",
					LastTriggered: "",
					Rearm:         3,
					Query: types.Query{
						Id:                1,
						LatestQueryDataId: 1,
						Name:              "table",
						Description:       "",
						Query:             "select column from table;",
						DataSourceId:      3,
						Version:           1,
						Tags:              []string{},
						QueryOptions: types.QueryOptions{
							Parameters: []map[string]interface{}{},
						},
					},
					QueryId: 1,
				},
			},
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/create/success.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create alert with invalid payload",
			Input: []map[string]interface{}{
				{
					"options": map[string]interface{}{
						"column": "impression_count",
						"op":     "<=",
						"value":  20,
					},
					"name": "Go - Google Ads: impression_count <= 20",
				},
			},
			Expected:     "Error: Missing query id",
			OutputFormat: test_utils.DefaultOutputFormat,
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodPost,
					StatusCode: 400,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/create/error.json",
					Err:        nil,
				},
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test_utils.SetupOutputFormat(testCase.OutputFormat)
			test_utils.MockApiResponse(&testCase.MockResponses)
			payload, err := json.Marshal(testCase.Input)
			if err != nil {
				assert.Error(t, fmt.Errorf("unable to create payload"), testCase.Name)
			}
			actual, err := ui.CreateAlert(string(payload))
			if err != nil {
				assert.Equal(t, testCase.Expected, err.Error(), testCase.Name)
			} else {
				expected, _ := json.Marshal(&testCase.Expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.Name)
			}
		})
	}
}
