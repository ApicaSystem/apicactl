package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/logiqai/logiqctl/types"
)

func createWidget(widget types.Widget, visualizationId int, dashboardId int) (types.Widget, error) {
	uri := GetUrlForResource(ResourceWidgetAll)
	client := ApiClient{}
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
