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

func ListAlerts(client utils.Api) ([]types.Resource, error) {
	uri := GetUrlForResource(ResourceAlertsAll)
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	defer resp.Body.Close()

	if err != nil {
		return []types.Resource{}, err
	}
	responseData, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: %s", responseData)
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

func GetAlert(client utils.Api, id string) (types.Resource, error) {
	uri := GetUrlForResource(ResourceAlert, id)
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Error: Alert does not exist")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: Unable to fetch alert")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
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

func CreateAlert(client utils.Api, jsonPayload string) (string, error) {
	uri := GetUrlForResource(ResourceAlertsAll)
	resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBufferString(jsonPayload))
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("Error: %s", err.Error())
	}
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error: %s", string(respString))
	}
	var alert types.Alert
	json.Unmarshal(respString, &alert)
	jsonResponse, err := json.MarshalIndent(alert, "", " ")
	if err != nil {
		return "", fmt.Errorf("Alert Created. Unable to print json response")
	}
	return string(jsonResponse), nil
}
