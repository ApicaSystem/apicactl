package converter

import (
	"encoding/json"
	"fmt"

	"github.com/logiqai/logiqctl/converter/grafana"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
)

func ConvertToLogiqDashboard(dashboardTemplate string, dashboardType string, dashboardName string, input *[]byte) (string, error) {
	if dashboardType == "grafana" {
		var grafanaDashboard types.GrafanaDashboard
		err := json.Unmarshal([]byte(dashboardTemplate), &grafanaDashboard)
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
		if dashboardName == "" {
			dashboardName = grafanaDashboard.Title
		}
		existingDashboard, err := ui.GetDashboardByName(dashboardName)
		if err != nil {
			return "", err
		}
		if existingDashboard != nil {
			return "", fmt.Errorf("error: Dashboard with name '%s' already exist", dashboardName)
		}
		grafanaConverter := grafana.NewGrafanaConverter(&grafanaDashboard)
		_, err = grafanaConverter.ParseInput(input)
		if err != nil {
			return "", fmt.Errorf("error: %s", err.Error())
		}
		dashboardSpec, err := grafanaConverter.CreateAndPublishDashboard(dashboardName)
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		result, err := json.MarshalIndent(dashboardSpec, "", " ")
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}

		return string(result), nil
	}
	return "", nil
}
