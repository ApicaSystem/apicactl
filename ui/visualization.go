package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ApicaSystem/apicactl/defines"
	"github.com/ApicaSystem/apicactl/utils"
	"io/ioutil"
	"net/http"

	"github.com/ApicaSystem/apicactl/types"
)

func CreateVisualization(visualization *types.Visualization, queryId int) (*types.Visualization, error) {
	uri := utils.GetUrlForResource(defines.ResourceVisualizationAll)
	client := utils.GetApiClient()
	vSpec := map[string]interface{}{}
	vSpec["name"] = visualization.Name
	vSpec["options"] = visualization.Options
	vSpec["description"] = visualization.Description
	vSpec["type"] = visualization.Type
	vSpec["query_id"] = queryId

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))
		if err == nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read create visualization response, Error: %s", err.Error())
			}

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				return nil, fmt.Errorf("Unable to create visualization: %s", err.Error())
			}

			respDict := types.Visualization{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode create visualization response")
			}

			// utils.CheckMesgErr(respDict, "createVisualization")

			return &respDict, nil
		} else {
			fmt.Println("err=<", err, ">")
			return nil, err
		}
	}
}
