package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ApicaSystem/apicactl/defines"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ApicaSystem/apicactl/types"
	"github.com/ApicaSystem/apicactl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewListDatasourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasource",
		Example: "apicactl get datasource|ds <datasource-id>",
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
		Example: "apicactl get datasource all",
		Short:   "List all the available datasources",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDataSources()
		},
	})

	return cmd
}

func printDataSource(args []string) {
	if v, vErr := GetDatasource(args...); v != nil {
		s, _ := json.MarshalIndent(*v, "", "    ")
		fmt.Println(string(s))
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func GetDatasource(args ...string) (*types.Datasource, error) {
	uri := utils.GetUrlForResource(defines.ResourceDatasource, args...)
	client := utils.GetApiClient()
	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	if err == nil {
		defer resp.Body.Close()
		var v = types.Datasource{}
		respBody, err := client.GetResponseString(resp)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(respBody, &v)
		if err != nil {
			return nil, fmt.Errorf("unable to get datasource, %s", err.Error())
		}
		return &v, nil
	} else {
		return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
	}
}

func createDataSource(datasource types.Datasource) (types.Datasource, error) {
	uri := utils.GetUrlForResource(defines.ResourceDatasourceAll)
	client := utils.GetApiClient()

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
	uri := utils.GetUrlForResource(defines.ResourceDatasourceAll)
	client := utils.GetApiClient()
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
