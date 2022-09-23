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

	if err != nil {
		return []types.Resource{}, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []types.Resource{}, err
	}
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
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Error: Alert does not exist")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return nil, fmt.Errorf("Error: %s", errorResponse["message"])
	}
	var result types.Resource
	alert := types.Alert{}
	json.Unmarshal(responseData, &alert)
	alert.FormatAlert(utils.FlagTimeFormat)
	result = alert
	return result, nil
}

func createAlert(alert types.CreateAlertPayload) (types.Alert, error) {
	client := ApiClient{}
	uri := GetUrlForResource(ResourceAlertsAll)

	payload, err := json.Marshal(alert)
	if err != nil {
		return types.Alert{}, fmt.Errorf("%s", err.Error())
	}
	resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBufferString(string(payload)))
	if err != nil {
		return types.Alert{}, fmt.Errorf("%s", err.Error())
	}
	defer resp.Body.Close()
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Alert{}, err
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
	newAlert.FormatAlert(utils.FlagTimeFormat)
	return newAlert, nil
}

func isAlertList(payload string) bool {
	trimmedPayload := bytes.TrimLeft([]byte(payload), " \t\r\n")
	if len(trimmedPayload) > 0 && trimmedPayload[0] == '[' {
		return true
	}
	return false
}

func CreateAlert(jsonPayload string) (string, error) {
	var alertList []types.CreateAlertPayload
	if isAlertList(jsonPayload) {
		err := json.Unmarshal([]byte(jsonPayload), &alertList)
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
	} else {
		var alert types.CreateAlertPayload
		err := json.Unmarshal([]byte(jsonPayload), &alert)
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
		alertList = append(alertList, alert)
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
