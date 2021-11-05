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

func createWidget(vFromDashSpec map[string]interface{}, vId interface{}, dId interface{}) (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceWidgetAll)
	client := getHttpClient()
	vSpec := map[string]interface{}{}

	vSpec["visualization_id"] = vId
	vSpec["dashboard_id"] = dId
	vSpec["options"] = vFromDashSpec["options"]
	vSpec["width"] = vFromDashSpec["width"]
	vSpec["text"] = vFromDashSpec["text"]

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		fmt.Println("json error=", jsonMarshallError)
		return nil, jsonMarshallError
	} else {
		req, err := http.NewRequest("POST",uri,bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Unable to create widget", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.AuthToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
		}
		if resp, err := client.Do(req); err == nil {
			jsonStr, _ := json.MarshalIndent(vSpec, "", "    ")
			fmt.Printf("Successfully added widget to dashboard : %s\n", jsonStr)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read added widget to dashboard response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode added widget to dashboard response")
			}

			utils.CheckMesgErr(respDict, "creageWidget")

			return respDict, nil
		} else {
			fmt.Println("createWidget err=<", err, ">")
			return nil, err
		}
	}
}

