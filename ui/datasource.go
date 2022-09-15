package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewListDatasourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasource",
		Example: "logiqctl get datasource|ds <datasource-id>",
		Aliases: []string{"ds"},
		Short:   "Get a datasource",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing datasource id")
				os.Exit(-1)
			}
			printDataSource(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get datasource all",
		Short:   "List all the available datasources",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDataSources()
		},
	})

	return cmd
}

func printDataSource(args []string) {
	if v, vErr := getDatasource(args); v != nil {
		s, _ := json.MarshalIndent(*v, "", "    ")
		fmt.Println(string(s))
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func getDatasource(args []string) (*map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceDatasource, args...)
	client := getHttpClient()
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Unable to get datasource ", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode datasources, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
	}
}

func createDataSource(datasource types.Datasource) (types.Datasource, error) {
	uri := GetUrlForResource(ResourceDatasourceAll)
	client := utils.ApiClient{}

	if payloadBytes, jsonMarshallError := json.Marshal(datasource); jsonMarshallError != nil {
		return types.Datasource{}, jsonMarshallError
	} else {
		resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))
		if err == nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return types.Datasource{}, fmt.Errorf("Unable to read create datasource response, Error: %s", err.Error())
			}
			response := types.Datasource{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &response); errUnmarshall != nil {
				return types.Datasource{}, fmt.Errorf("%s", errUnmarshall.Error())
			}

			return response, nil
		} else {
			return types.Datasource{}, err
		}
	}
}

func getDataSourceByName(name string) (types.Datasource, error) {
	if v, err := getDatasources(); err != nil {
		return types.Datasource{}, err
	} else {
		for _, ds := range v {
			if ds.Name == name {
				return ds, nil
			}
		}
	}
	return types.Datasource{}, nil
}

func listDataSources() {
	if v, err := getDatasources(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		count := len(v)
		fmt.Println("(", count, ") datasources found")
		if count > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Type", "Id"})
			for _, datasource := range v {
				name := datasource.Name
				typ := datasource.DatasourceType
				id := (int)(datasource.Id)
				table.Append([]string{name, typ,
					fmt.Sprintf("%d", id)})
			}
			table.Render()
		}
	}
}

func getDatasources() ([]types.Datasource, error) {
	uri := GetUrlForResource(ResourceDatasourceAll)
	client := utils.ApiClient{}
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)

	if err == nil {
		defer resp.Body.Close()
		var v = []types.Datasource{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode datasources, Error: %s", err.Error())
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
	}
}
