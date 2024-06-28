package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ApicaSystem/apicactl/defines"
	"github.com/ApicaSystem/apicactl/types"
	"github.com/ApicaSystem/apicactl/utils"
)

type MappingsResponse struct {
	Mappings []types.Mapping `json:"mappings"`
}

func GetMappings() ([]types.Mapping, error) {
	client := utils.GetApiClient()
	uri := utils.GetUrlForResource(defines.ResourceMappingsAll)
	res, err := client.MakeApiCall(http.MethodGet, uri, nil)
	if err != nil {
		return []types.Mapping{}, err
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []types.Mapping{}, err
	}

	if res.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return nil, fmt.Errorf("Error: %s", errorResponse["message"])
	}

	var mappings MappingsResponse
	json.Unmarshal(responseData, &mappings)

	return mappings.Mappings, nil
}
