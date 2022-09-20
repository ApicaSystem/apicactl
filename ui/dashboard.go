package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/spf13/viper"

	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"net/http"
)

func NewListDashboardsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "logiqctl get dashboard|d <dashboard-slug>",
		Aliases: []string{"d"},
		Short:   "Get a dashboard",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing dashboard slug")
				os.Exit(-1)
			}
			exportDashboard(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get dashboard all",
		Short:   "List all the available dashboards",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDashboards()
		},
	})

	return cmd
}

func exportDashboard(args []string) {
	dashboardOut := map[string]interface{}{}

	if dashboardPtr, err := GetDashboard(args); err != nil {
		fmt.Println(err.Error())
	} else {
		dashboard := *dashboardPtr

		/* json dump exercise
		fmt.Println("dashboard=<", dashboard, ">")
		v, err := json.Marshal(dashboard)
		if (err!=nil) {
			fmt.Println("error=", err)
		} else {
			fmt.Println("no error found")
		}
		fmt.Println("v=<", string(v), ">")
		gg := map[string]interface{}{}
		_ = json.Unmarshal(v, &gg)
		fmt.Println("dashboard_json=<", gg, ">")
		*/

		dashboardParams := map[string]interface{}{}
		dashboardParams["name"] = dashboard["name"]
		dashboardParams["tags"] = dashboard["tags"]
		dashboardOut["dashboard"] = dashboardParams

		widgets := dashboard["widgets"].([]interface{})
		widgetOut := []interface{}{}
		dataSources := map[int]interface{}{}

		importWidget := true
		for _, w := range widgets {
			widget := w.(map[string]interface{})

			if _, ok := widget["visualization"]; ok {

				visualization := widget["visualization"].(map[string]interface{})
				query := visualization["query"].(map[string]interface{})

				dId := (int)(query["data_source_id"].(float64))
				if _, ok := dataSources[dId]; !ok {
					if dsPtr, dsErr := getDatasource([]string{fmt.Sprintf("%d", dId)}); dsErr == nil {
						datasource := *dsPtr
						dataSources[dId] = map[string]interface{}{
							"name":    datasource["name"],
							"options": datasource["options"],
							"type":    datasource["type"],
						}
					} else {
						fmt.Printf("Data source not found: %d", dId)
						fmt.Println("This widget will not be imported")
						importWidget = false
					}
				}
				if importWidget {
					wOutEntry := map[string]interface{}{}
					wOutEntry["options"] = widget["options"]
					wOutEntry["text"] = widget["text"]
					wOutEntry["width"] = widget["width"]
					wOutEntry["visualization"] = map[string]interface{}{
						"type":    visualization["type"],
						"name":    visualization["name"],
						"options": visualization["options"],
						"query": map[string]interface{}{
							"name":           query["name"],
							"options":        query["options"],
							"description":    query["description"],
							"data_source_id": query["data_source_id"],
							"query":          query["query"],
						},
					}
					widgetOut = append(widgetOut, wOutEntry)
					dashboardOut["widgets"] = widgetOut
					dashboardOut["datasources"] = dataSources
				}
			} else {
				if importWidget {
					wOutEntry := map[string]interface{}{}
					wOutEntry["options"] = widget["options"]
					wOutEntry["text"] = widget["text"]
					wOutEntry["width"] = widget["width"]
					widgetOut = append(widgetOut, wOutEntry)
					dashboardOut["widgets"] = widgetOut
				}
			}
		}

	}

	s, _ := json.MarshalIndent(dashboardOut, "", "    ")
	fmt.Println(string(s))
}

func GetDashboard(args []string) (*map[string]interface{}, error) {

	uri := GetUrlForResource(ResourceDashboardsGet, args...)
	client := getHttpClient()
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}

	req = utils.AddNetTrace(req)

	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode dashboard, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
	}
}

func createAndPublishDashboard(name string) (types.Dashboard, error) {
	dashboardParams := map[string]interface{}{
		"name": name,
	}

	if payloadBytes, jsonMarshallError := json.Marshal(dashboardParams); jsonMarshallError != nil {
		return types.Dashboard{}, jsonMarshallError
	} else {
		// Create dashboard
		uri := GetUrlForResource(ResourceDashboardsAll)
		client := &ApiClient{}
		resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))

		if err != nil {
			return types.Dashboard{}, err
		}

		if resp.StatusCode != http.StatusOK {
			return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard1, Http response error while creating dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		/*
			var v = map[string]interface{}{}
			if resp.StatusCode != http.StatusOK {
				return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard2, Http response error while creating dashboard, Error: %d", resp.StatusCode)
			}
		*/

		// Decode create response
		var dashboard types.Dashboard
		var v map[string]interface{}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to create dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &v)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to decode create dashboard response, Error: %s", err.Error())
		}

		// check for server error

		utils.CheckMesgErr(v, "createAndPublishDashboard")

		// Create publish payload
		payloadPublish := map[string]interface{}{
			"is_draft": false, "name": name, "slug": v["id"],
		}
		payloadBytes, jsonMarshallError = json.Marshal(payloadPublish)
		if jsonMarshallError != nil {
			return types.Dashboard{}, jsonMarshallError
		}
		args := []string{fmt.Sprintf("%v", v["id"])}

		// Publish dashboard
		uri = GetUrlForResource(ResourceDashboardsGet, args...)
		resp, err = client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))

		if err != nil {
			return types.Dashboard{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard2, Http response error while publishing dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		// Decode publish response
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to publish dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &dashboard)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to decode publish dashboard response, Error: %s", err.Error())
		}
		return dashboard, nil
	}
}

func CreateAndPublishDashboardSpec(dashboardSpecJson string) (string, error) {
	var dashboardSpec types.DashboardSpec
	var dashboard types.Dashboard
	var err error
	var responseSpec types.DashboardSpec
	var response []byte
	err = json.Unmarshal([]byte(dashboardSpecJson), &dashboardSpec)

	if err != nil {
		return "", fmt.Errorf("Error: %s", err.Error())
	}

	if dashboardSpec.Dashboard.Name != "" {
		existingDashboard := getDashboardByName(dashboardSpec.Dashboard.Name)
		if existingDashboard != nil {
			return "", fmt.Errorf("Dashboard with name \"%s\" already exists", dashboardSpec.Dashboard.Name)
		}
		dashboard, err = createAndPublishDashboard(dashboardSpec.Dashboard.Name)
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
		responseSpec.Dashboard = dashboard
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
	} else {
		return "", fmt.Errorf("Error: %s", "Dashboard name is missing in spec")
	}

	if len(dashboardSpec.Datasources) > 0 {
		responseSpec.Datasources = make(map[string]types.Datasource)
		for _, datasource := range dashboardSpec.Datasources {
			if existingDatasource, err := getDataSourceByName(datasource.Name); existingDatasource.Id != 0 {
				datasource.Id = existingDatasource.Id
				responseSpec.Datasources[strconv.Itoa(existingDatasource.Id)] = existingDatasource
			} else if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			} else {
				datasource, err = createDataSource(datasource)
				if err != nil {
					return "", fmt.Errorf("Error: %s", err.Error())
				}
				responseSpec.Datasources[strconv.Itoa(datasource.Id)] = datasource
			}
		}
	}
	if len(dashboardSpec.Widgets) > 0 {
		responseSpec.Widgets = []types.Widget{}
		for _, widget := range dashboardSpec.Widgets {
			if widget.Visualization.Type != "" {
				query := widget.Visualization.Query
				query = getQueryByName(query.Name)

				if query.Id == 0 {
					tempQuery, _ := json.Marshal(widget.Visualization.Query)
					queryPayload := types.CreateQueryPayload{}
					err = json.Unmarshal(tempQuery, &queryPayload)
					if err != nil {
						return "", err
					}
					query, err = createQuery(queryPayload)
					if err != nil {
						return "", err
					}
				}
				if query.IsDraft {
					publishArgs := []string{fmt.Sprintf("%d", query.Id), fmt.Sprintf("%d", query.Version)}
					publishQuery(publishArgs)
				}
				visualization, err := createVisualization(widget.Visualization, query.Id)
				if err != nil {
					return "", fmt.Errorf("Error: %s", err.Error())
				}
				newWidget, err := createWidget(widget, visualization.Id, responseSpec.Dashboard.Id)
				visualization.Query = query
				newWidget.Visualization = visualization
				responseSpec.Widgets = append(responseSpec.Widgets, newWidget)
			}
		}

		if len(dashboardSpec.Alerts) > 0 {
			alertsJson, _ := json.Marshal(dashboardSpec.Alerts)
			alertPayload := []types.CreateAlertPayload{}
			err := json.Unmarshal(alertsJson, &alertPayload)
			if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			}
			payload, err := json.Marshal(alertPayload)
			alertResponse, err := CreateAlert(string(payload))
			if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			}
			alertList := []types.Alert{}
			json.Unmarshal([]byte(alertResponse), &alertList)
			responseSpec.Alerts = alertList
		}
	}

	response, _ = json.MarshalIndent(&responseSpec, "", " ")

	return string(response), nil
}

func getDashboardByName(name string) map[string]interface{} {
	if v, err := GetDashboards(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		dashboards := v["results"].([]interface{})
		for _, dash := range dashboards {
			dashboard := dash.(map[string]interface{})
			if dashboard["name"] == name {
				return dashboard
			}
		}
	}
	return nil
}

func GetDashboards() (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceDashboardsAll)
	client := getHttpClient()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}

	req = utils.AddNetTrace(req)

	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch dashboards, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode dashboards, Error: %s", err.Error())
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error with get dashboards, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch dashboards, Error: %s", err.Error())
	}
}

func listDashboards() {
	if v, err := GetDashboards(); err == nil {
		count := (int)(v["count"].(float64))
		dashboards := v["results"].([]interface{})
		fmt.Println("(", count, ") dashboards found")
		if count > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Slug", "Id"})
			for _, d := range dashboards {
				dash := d.(map[string]interface{})
				slug := dash["slug"].(string)
				name := dash["name"].(string)
				id := (int)(dash["id"].(float64))
				table.Append([]string{name, slug, fmt.Sprintf("%d", id)})
			}

			table.Render()
		}
	} else {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}
}

func GetLogEvents(numDays int) error {
	response, err := ExecutePrometheusQuery(fmt.Sprintf("round(sum(increase(logiq_message_count[%dh])))", numDays))
	if err != nil {
		return err
	}
	if data, ok := response["data"]; ok {
		if data, ok := data.(map[string]interface{}); ok {
			if result, ok := data["result"]; ok {
				if result, ok := result.([]interface{}); ok {
					for _, v := range result {
						if v, ok := v.(map[string]interface{}); ok {
							for k, vv := range v {
								if k == "value" {
									if vv, ok := vv.([]interface{}); ok {
										fmt.Printf("Total Log Events for %d days in %s: %s\n", numDays, viper.GetString(utils.KeyCluster), vv[1])
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func ExecutePrometheusQuery(query string) (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourcePrometheusProxy)
	client := getHttpClient()
	payload := fmt.Sprintf(`{"query":"%s","type":"query"}`, query)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Unable to create widget", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}
	if resp, err := client.Do(req); err == nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to execute Query %s", err.Error())
		}
		respDict := map[string]interface{}{}
		if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
			return nil, fmt.Errorf("unable to decode response")
		}
		return respDict, nil
	} else {
		fmt.Println("ExecutePrometheusQuery err=<", err, ">")
		return nil, err
	}
}

func NewGetLogEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log-events",
		Example: "logiqctl get log-events 7",
		Aliases: []string{"t"},
		Short:   "Get total log events for the duration in days",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			var days int = 1
			if len(args) == 1 {
				d, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse input, [%v], expects an integer\n", args[0])
					os.Exit(-1)
				}
				days = int(d)
			}
			GetLogEvents(days)
		},
	}
	return cmd
}
