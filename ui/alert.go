package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/utils"
)

func ListAlerts() ([]types.Resource, error) {
	client := ApiClient{}
	uri := GetUrlForResource(ResourceAlertsAll)
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	defer resp.Body.Close()

	if err != nil {
		return []types.Resource{}, err
	}
	responseData, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return nil, fmt.Errorf("Error: %s", errorResponse["message"])
	}
	var result []types.Resource = make([]types.Resource, 0)
	alertList := []types.Alert{}

	json.Unmarshal(responseData, &alertList)
	for _, alert := range alertList {
		alert.FormatAlert(utils.FlagTimeFormat)
		result = append(result, alert)
	}
	return result, nil
}

func GetAlert(id string) (types.Resource, error) {
	client := ApiClient{}
	uri := GetUrlForResource(ResourceAlert, id)
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Error: Alert does not exist")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return nil, fmt.Errorf("Error: %s", errorResponse["message"])
	}
	if err != nil {
		return nil, err
	}
	var result types.Resource
	alert := types.Alert{}
	json.Unmarshal(responseData, &alert)
	alert.FormatAlert(utils.FlagTimeFormat)
	result = alert
	return result, nil
}

func createAlert(alert types.Alert) (types.Alert, error) {
	client := ApiClient{}
	uri := GetUrlForResource(ResourceAlertsAll)

	payload, err := json.Marshal(alert)
	if err != nil {
		return types.Alert{}, fmt.Errorf("%s", err.Error())
	}
	resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBufferString(string(payload)))
	defer resp.Body.Close()
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Alert{}, fmt.Errorf("%s", err.Error())
	}
	if resp.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(respString, &errorResponse)
		return types.Alert{}, fmt.Errorf("%s", errorResponse["message"])
	}
	var newAlert types.Alert
	err = json.Unmarshal(respString, &newAlert)
	if err != nil {
		fmt.Println()
		return types.Alert{}, err
	}
	return newAlert, nil
}

func CreateAlert(jsonPayload string) (string, error) {
	var alertList []types.Alert
	err := json.Unmarshal([]byte(jsonPayload), &alertList)
	if err != nil {
		return "", fmt.Errorf("Error: %s", err.Error())
	}
	var response []types.Alert
	for _, alert := range alertList {
		newAlert, err := createAlert(alert)
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}

		response = append(response, newAlert)
	}
	jsonResponse, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return "", fmt.Errorf("Alert Created. Unable to print json response")
	}

	return string(jsonResponse), nil
}
