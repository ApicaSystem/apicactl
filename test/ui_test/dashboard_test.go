package ui_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	test_utils "github.com/logiqai/logiqctl/test/utils"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
	"github.com/stretchr/testify/assert"
)

func TestCreatePublishDashboardSpec(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		{
			Name: "Create Dashboard by reusing datasources and query",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": types.Datasource{
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
				Widgets: []types.Widget{
					{
						Options: map[string]interface{}{
							"isHidden": false,
							"parameterMappings": map[string]map[string]string{
								"duration": {
									"mapTo": "duration",
									"name":  "duration",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
								"step": {
									"mapTo": "step",
									"name":  "step",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
							},
							"position": map[string]interface{}{
								"autoHeight": false,
								"col":        3,
								"maxSizeX":   6,
								"maxSizeY":   1000,
								"minSizeX":   1,
								"minSizeY":   5,
								"row":        40,
								"sizeX":      3,
								"sizeY":      8,
							},
						},
						Width: 1,
						Visualization: types.Visualization{
							Name: "Chart",
							Type: "CHART",
							Options: map[string]interface{}{
								"columnMapping": map[string]string{
									"timestamp": "x",
									"value":     "y",
								},
								"customCode":     "// Available variables are x, ys, element, and Plotly\n// Type console.log(x, ys); for more info about x and ys\n// To plot your graph call Plotly.plot(element, ...)\n// Plotly examples and docs: https://plot.ly/javascript/",
								"dateTimeFormat": "DD/MM/YY HH:mm",
								"defaultColumns": 3,
								"defaultRows":    8,
								"error_y": map[string]interface{}{
									"type":    "data",
									"visible": true,
								},
								"globalSeriesType": "line",
								"legend": map[string]interface{}{
									"enabled": true,
								},
								"minColumns":    1,
								"minRows":       5,
								"numberFormat":  "0,0[.]00000",
								"percentFormat": "0[.]00%",
								"series": map[string]interface{}{
									"error_y": map[string]interface{}{
										"type":    "data",
										"visible": true,
									},
									"stacking": nil,
								},
								"seriesOptions": map[string]map[string]interface{}{
									"value": {
										"index":  0,
										"type":   "line",
										"yAxis":  0,
										"zIndex": 0,
									},
								},
								"showDataLabels": false,
								"sortX":          true,
								"textFormat":     "",
								"xAxis": map[string]interface{}{
									"labels": map[string]bool{
										"enabled": true,
									},
									"type": "-",
								},
								"yAxis": []map[string]interface{}{
									{
										"type": "linear",
									},
									{
										"opposite": true,
										"type":     "linear",
									},
								},
							},
							Query: types.Query{
								DataSourceId: 1,
								Description:  "",
								Name:         "Query 1",
								QueryOptions: types.QueryOptions{
									Parameters: []interface{}{
										map[string]interface{}{
											"enumOptions": "1h\n2h\n3h\n1d\n2d\n3d",
											"locals":      []string{},
											"name":        "duration",
											"title":       "Duration",
											"type":        "enum",
											"value":       "2h",
										},
										map[string]interface{}{
											"enumOptions": "10s\n30s\n60s\n1m\n5m",
											"locals":      []string{},
											"name":        "step",
											"title":       "Step",
											"type":        "enum",
											"value":       "1m",
										},
									},
								},
								Query: "query=rate(redis_keyspace_hits_total[5m])\u0026 duration={{ duration }}\u0026step={{ step }}",
							},
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/datasource/get/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/queries",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/queries/list/200_success.json",
					Err:        nil,
				},
				{
					Url:        "api/widgets",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/widget/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/visualizations",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/visualization/create/201_success.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create Dashboard by creating datasources and query",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
				Widgets: []types.Widget{
					{
						Options: map[string]interface{}{
							"isHidden": false,
							"parameterMappings": map[string]map[string]string{
								"duration": {
									"mapTo": "duration",
									"name":  "duration",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
								"step": {
									"mapTo": "step",
									"name":  "step",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
							},
							"position": map[string]interface{}{
								"autoHeight": false,
								"col":        3,
								"maxSizeX":   6,
								"maxSizeY":   1000,
								"minSizeX":   1,
								"minSizeY":   5,
								"row":        40,
								"sizeX":      3,
								"sizeY":      8,
							},
						},
						Width: 1,
						Visualization: types.Visualization{
							Name: "Chart",
							Type: "CHART",
							Options: map[string]interface{}{
								"columnMapping": map[string]string{
									"timestamp": "x",
									"value":     "y",
								},
								"customCode":     "// Available variables are x, ys, element, and Plotly\n// Type console.log(x, ys); for more info about x and ys\n// To plot your graph call Plotly.plot(element, ...)\n// Plotly examples and docs: https://plot.ly/javascript/",
								"dateTimeFormat": "DD/MM/YY HH:mm",
								"defaultColumns": 3,
								"defaultRows":    8,
								"error_y": map[string]interface{}{
									"type":    "data",
									"visible": true,
								},
								"globalSeriesType": "line",
								"legend": map[string]interface{}{
									"enabled": true,
								},
								"minColumns":    1,
								"minRows":       5,
								"numberFormat":  "0,0[.]00000",
								"percentFormat": "0[.]00%",
								"series": map[string]interface{}{
									"error_y": map[string]interface{}{
										"type":    "data",
										"visible": true,
									},
									"stacking": nil,
								},
								"seriesOptions": map[string]map[string]interface{}{
									"value": {
										"index":  0,
										"type":   "line",
										"yAxis":  0,
										"zIndex": 0,
									},
								},
								"showDataLabels": false,
								"sortX":          true,
								"textFormat":     "",
								"xAxis": map[string]interface{}{
									"labels": map[string]bool{
										"enabled": true,
									},
									"type": "-",
								},
								"yAxis": []map[string]interface{}{
									{
										"type": "linear",
									},
									{
										"opposite": true,
										"type":     "linear",
									},
								},
							},
							Query: types.Query{
								DataSourceId: 1,
								Description:  "",
								Name:         "Query 1",
								QueryOptions: types.QueryOptions{
									Parameters: []interface{}{
										map[string]interface{}{
											"enumOptions": "1h\n2h\n3h\n1d\n2d\n3d",
											"locals":      []string{},
											"name":        "duration",
											"title":       "Duration",
											"type":        "enum",
											"value":       "2h",
										},
										map[string]interface{}{
											"enumOptions": "10s\n30s\n60s\n1m\n5m",
											"locals":      []string{},
											"name":        "step",
											"title":       "Step",
											"type":        "enum",
											"value":       "1m",
										},
									},
								},
								Query: "query=rate(redis_keyspace_hits_total[5m])\u0026 duration={{ duration }}\u0026step={{ step }}",
							},
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/datasource/get/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/queries",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/queries/list/200_success.json",
					Err:        nil,
				},
				{
					Url:        "api/queries",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/queries/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/widgets",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/widget/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/visualizations",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/visualization/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/datasource/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/alerts/create/success.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create Dashboard without widgets",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/datasource/get/200_response.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create Dashboard with name which already exists",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard exist",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
			},
			Expected: "Dashboard with name \"Dashboard exist\" already exists",
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       "../mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/datasource/get/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       "../mock_response/dashboards/list/200_response.json",
					Err:        nil,
				},
			},
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			test_utils.MockApiResponse(testcase.MockResponses)
			dashboardSpec := testcase.Input.(types.DashboardSpec)
			payload, err := json.Marshal(dashboardSpec)
			actual, err := ui.CreateAndPublishDashboardSpec(string(payload))
			if err != nil {
				assert.Equal(t, testcase.Expected, err.Error(), testcase.Name)
			}
			actualDashboardSpec := types.DashboardSpec{}
			err = json.Unmarshal([]byte(actual), &actualDashboardSpec)
			if err != nil {
				assert.Fail(t, err.Error(), testcase.Name)
			} else {
				assert.Greater(t, actualDashboardSpec.Dashboard.Id, 0, "dashboard is not created")
				assert.Greater(t, len(actualDashboardSpec.Datasources), 0, "Datasource is empty")
				if len(dashboardSpec.Alerts) > 0 {
					assert.Equal(t, len(actualDashboardSpec.Alerts), 2, "Alerts is empty")
					for _, alert := range dashboardSpec.Alerts {
						assert.Greater(t, alert.Id, 0, testcase.Name)
					}
				}
				for _, datasource := range actualDashboardSpec.Datasources {
					assert.Greater(t, datasource.Id, 0, testcase.Name)
				}
				for _, alert := range actualDashboardSpec.Alerts {
					assert.Greater(t, alert.Id, 0, testcase.Name)
				}
				for _, widget := range actualDashboardSpec.Widgets {
					assert.Greater(t, widget.Id, 0, "widget is not created")
					assert.Greater(t, widget.Visualization.Id, 0, "visualization is not created")
					assert.Greater(t, widget.Visualization.Query.Id, 0, "Query is not created")
				}
			}
		})
	}
}
