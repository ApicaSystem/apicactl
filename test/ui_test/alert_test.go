package ui_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/logiqai/logiqctl/test/test_double"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	m "github.com/stretchr/testify/mock"
)

type mockResponse struct {
	statusCode int32
	body       string
	err        error
}

type testCase struct {
	name     string
	input    interface{}
	expected interface{}
	mockResponse
	outputFormat map[string]string
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	viper.Set("cluster", "localhost")
	return func(t *testing.T) {}
}

func setupOutputFormat(outputFormat map[string]string) {
	for key, value := range outputFormat {
		if key == "time-format" {
			utils.FlagTimeFormat = value
		}
	}
}

var defaultOutputFormat map[string]string = map[string]string{
	"time-format": "RFC3339",
}

func TestListAlertsCommand(t *testing.T) {
	tearDownTestcase := setupTestCase(t)

	defer tearDownTestcase(t)
	testCases := []testCase{
		{
			name: "List All Alerts",
			expected: []types.Alert{
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
					},
				},
				{
					Id:   2,
					Name: "table: column >= 100",
					AlertOption: types.AlertOption{
						Column: "column",
						Op:     ">=",
						Value:  100,
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
					},
				},
			},
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 200,
				body: `[{
					"id": 1,
					"name": "table: column >= 150",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 150
					},
					"state": "ok",
					"last_triggered_at": "2022-08-31T18:25:42.792Z",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				},
				{
					"id": 2,
					"name": "table: column >= 100",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 100
					},
					"state": "ok",
					"last_triggered_at": "2022-08-31T18:25:42.792Z",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				}]`,
				err: nil,
			},
		},
		{
			name:     "List Alerts when no alerts are created",
			expected: []types.Alert{},
			mockResponse: mockResponse{
				statusCode: 200,
				body:       `[]`,
				err:        nil,
			},
			outputFormat: defaultOutputFormat,
		},
		{
			name:     "List Alerts when there is an error in api response",
			expected: "Error: Unable to fetch alerts",
			mockResponse: mockResponse{
				statusCode: 500,
				body:       fmt.Errorf("%s", "Unable to fetch alerts").Error(),
				err:        nil,
			},
		},
		{
			name: "List all alerts with epoch time format",
			expected: []types.Alert{
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
					},
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
					},
				},
			},
			outputFormat: map[string]string{
				"time-format": "epoch",
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body: `[{
					"id": 1,
					"name": "table: column >= 150",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 150
					},
					"state": "ok",
					"last_triggered_at": "2022-08-31T18:25:42.792Z",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				},
				{
					"id": 2,
					"name": "table: column >= 100",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 100
					},
					"state": "unknown",
					"last_triggered_at": "",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				}]`,
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setupOutputFormat(testCase.outputFormat)
			client := &test_double.HttpMock{}
			client.On("MakeApiCall", http.MethodGet, fmt.Sprintf("https://%s/api/alerts", viper.GetString("cluster")), m.Anything).Return(
				&http.Response{StatusCode: int(testCase.statusCode), Body: ioutil.NopCloser(bytes.NewBufferString(testCase.mockResponse.body))},
				testCase.mockResponse.err).Once()

			alertsList, err := ui.ListAlerts(client)
			actual, _ := json.Marshal(&alertsList)
			if err != nil {
				expected := testCase.expected
				assert.Equal(t, expected, err.Error(), testCase.name)
			} else {
				expected, _ := json.Marshal(&testCase.expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.name)
			}
		})
	}
}

func TestGetAlertCommand(t *testing.T) {
	tearDownTestcase := setupTestCase(t)
	defer tearDownTestcase(t)

	// testcases
	testCases := []testCase{
		{
			name: "Get Alert By Id",
			input: map[string]string{
				"id": "1",
			},
			expected: types.Alert{
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
				},
			},
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 200,
				body: `{
					"id": 1,
					"name": "table: column >= 150",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 150
					},
					"state": "ok",
					"last_triggered_at": "2022-08-31T18:25:42.792Z",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				}`,
				err: nil,
			},
		},
		{
			name: "Get Alert By Id for which doest not exist",
			input: map[string]string{
				"id": "1",
			},
			expected:     "Error: Alert does not exist",
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 404,
				body:       ``,
				err:        nil,
			},
		},
		{
			name: "Get Alert By Id for which the api fails",
			input: map[string]string{
				"id": "1",
			},
			expected:     "Error: Unable to fetch alert",
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 500,
				body:       ``,
				err:        nil,
			},
		},
		{
			name: "Get alert by id with epcoh time format",
			input: map[string]string{
				"id": "1",
			},
			expected: types.Alert{
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
				},
			},
			outputFormat: map[string]string{
				"time-format": "epoch",
			},
			mockResponse: mockResponse{
				statusCode: 200,
				body: `{
					"id": 1,
					"name": "table: column >= 150",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 150
					},
					"state": "ok",
					"last_triggered_at": "2022-08-31T18:25:42.792Z",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						},
						"last_modified_by": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/user_1?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				}`,
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setupOutputFormat(testCase.outputFormat)
			client := &test_double.HttpMock{}
			id := testCase.input.(map[string]string)["id"]
			client.On("MakeApiCall", http.MethodGet, fmt.Sprintf("https://%s/api/alerts/%s", viper.GetString("cluster"), id), m.Anything).Return(
				&http.Response{StatusCode: int(testCase.statusCode), Body: ioutil.NopCloser(bytes.NewBufferString(testCase.mockResponse.body))},
				testCase.mockResponse.err).Once()

			alert, err := ui.GetAlert(client, id)
			actual, _ := json.Marshal(&alert)
			if err != nil {
				expected := testCase.expected
				assert.Equal(t, expected, err.Error(), testCase.name)
			} else {
				expected, _ := json.Marshal(&testCase.expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.name)
			}
		})
	}
}

func TestCreateAlertCommand(t *testing.T) {
	tearDownTestcase := setupTestCase(t)
	defer tearDownTestcase(t)

	// testcases
	testCases := []testCase{
		{
			name: "Create Alert",
			input: map[string]interface{}{
				"options": map[string]interface{}{
					"column": "impression_count",
					"op":     "<=",
					"value":  20,
				},
				"name":     "Go - Google Ads: impression_count <= 20",
				"query_id": 15,
			},
			expected: types.Alert{
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
				},
			},
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 200,
				body: `{
					"id": 1,
					"name": "table: column >= 150",
					"options": {
						"column": "column",
						"op": ">=",
						"value": 150
					},
					"state": "unknown",
					"last_triggered_at": "",
					"updated_at": "2022-08-31T18:25:42.787Z",
					"created_at": "2022-08-31T18:11:13.470Z",
					"rearm": 3,
					"query": {
						"id": 1,
						"latest_query_data_id": 1,
						"name": "table",
						"description": null,
						"query": "select column from table;",
						"query_hash": "1234567890",
						"schedule": null,
						"api_key": "api_secret",
						"is_archived": false,
						"is_draft": false,
						"updated_at": "2022-08-31T18:37:31.569Z",
						"created_at": "2022-08-31T18:08:05.365Z",
						"data_source_id": 3,
						"options": {
							"parameters": []
						},
						"version": 1,
						"tags": [],
						"is_safe": true,
						"user": {
							"id": 1,
							"name": "flash-admin@foo.com",
							"email": "flash-admin@foo.com",
							"profile_image_url": "https://www.gravatar.com/avatar/b'66dbe9c206b1758bfbb5f2bb06130358'?s=40&d=identicon",
							"groups": [1, 2],
							"updated_at": "2022-09-03T09:18:44.421Z",
							"created_at": "2022-08-23T19:15:16.850Z",
							"disabled_at": null,
							"is_disabled": false,
							"active_at": "2022-09-03T09:18:28Z",
							"is_invitation_pending": false,
							"is_email_verified": true,
							"auth_type": "password"
						}
					}
				}`,
				err: nil,
			},
		},
		{
			name: "Create alert with invalid payload",
			input: map[string]interface{}{
				"options": map[string]interface{}{
					"column": "impression_count",
					"op":     "<=",
					"value":  20,
				},
				"name": "Go - Google Ads: impression_count <= 20",
			},
			expected:     "Error: Missing query id",
			outputFormat: defaultOutputFormat,
			mockResponse: mockResponse{
				statusCode: 400,
				body:       `Missing query id`,
				err:        nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			setupOutputFormat(testCase.outputFormat)
			client := &test_double.HttpMock{}
			client.On("MakeApiCall", http.MethodPost, fmt.Sprintf("https://%s/api/alerts", viper.GetString("cluster")), m.Anything).Return(
				&http.Response{StatusCode: int(testCase.statusCode), Body: ioutil.NopCloser(bytes.NewBufferString(testCase.mockResponse.body))},
				testCase.mockResponse.err).Once()
			payload, err := json.Marshal(testCase.input)
			if err != nil {
				assert.Error(t, fmt.Errorf("unable to create payload"), testCase.name)
			}
			actual, err := ui.CreateAlert(client, string(payload))
			if err != nil {
				assert.Equal(t, testCase.expected, err.Error(), testCase.name)
			} else {
				expected, _ := json.Marshal(&testCase.expected)

				assert.JSONEq(t, string(expected), string(actual), testCase.name)
			}
		})
	}
}
