package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

func createVisualization(vFromDashSpec map[string]interface{}, qId interface{}) (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceVisualizationAll)
	client := getHttpClient()
	vSpec := map[string]interface{}{}
	vSpec["name"] = vFromDashSpec["name"]
	vSpec["options"] = vFromDashSpec["options"]
	vSpec["description"] = vFromDashSpec["description"]
	vSpec["type"] = vFromDashSpec["type"]
	vSpec["query_id"] = qId

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		req, err := http.NewRequest("POST",uri,bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Unable to create visualization ", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.KeyUiToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", ))
		}
		if resp, err := client.Do(req); err == nil {
			jsonStr, _ := json.MarshalIndent(vSpec, "", "    ")
			fmt.Printf("Successfully created visualization : %s", jsonStr)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read create visualization response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode create visualization response")
			}

			return respDict, nil
		} else {
			return nil, err
		}
	}
}
