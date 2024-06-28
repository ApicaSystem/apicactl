package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ApicaSystem/apicactl/defines"
	"github.com/ApicaSystem/apicactl/utils"

	"github.com/ApicaSystem/apicactl/types"
)

func CreateWidgetGroup(widget types.Widget, dashboardId int) (*types.Widget, error) {
	uri := utils.GetUrlForResource(defines.ResourceWidgetAll)
	client := utils.GetApiClient()
	payload := map[string]interface{}{
		"text":         widget.Text,
		"type":         widget.Type,
		"options":      widget.Options,
		"width":        widget.Width,
		"dashboard_id": dashboardId,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))
	respString, err := client.GetResponseString(resp)
	if err != nil {
		return nil, err
	}
	var group types.Widget
	json.Unmarshal(respString, &group)
	return &group, nil
}

func CreateWidget(widget types.Widget, visualizationId int, dashboardId int) (types.Widget, error) {
	uri := utils.GetUrlForResource(defines.ResourceWidgetAll)
	client := utils.GetApiClient()
	vSpec := map[string]interface{}{}

	vSpec["visualization_id"] = visualizationId
	vSpec["dashboard_id"] = dashboardId
	vSpec["options"] = widget.Options
	vSpec["width"] = widget.Width
	vSpec["text"] = widget.Text

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		return types.Widget{}, jsonMarshallError
	} else {
		resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))
		if err == nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return types.Widget{}, fmt.Errorf("Unable to read added widget to dashboard response, Error: %s", err.Error())
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				return types.Widget{}, fmt.Errorf("Unable to create widget: %s", err.Error())
			}
			response := types.Widget{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &response); errUnmarshall != nil {
				return types.Widget{}, fmt.Errorf("Unable to decode added widget to dashboard response")
			}

			// utils.CheckMesgErr(response, "creageWidget")

			return response, nil
		} else {
			fmt.Println("createWidget err=<", err, ">")
			return types.Widget{}, err
		}
	}
}
