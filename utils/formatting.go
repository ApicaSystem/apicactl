/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"

	"github.com/dustin/go-humanize"
	"github.com/logiqai/logiqctl/types"
	"github.com/olekukonko/tablewriter"
)

var FlagOut string
var FlagTimeFormat string
var FlagNamespace string
var FlagCluster string
var FlagLogsSince string
var FlagLogsPageSize uint32
var FlagLogsFollow bool
var FlagProcId string
var FlagAppName string
var FlagFile string
var FlagMaxLogLines int
var EventRuleGroupsFlag string
var FlagBegTime string
var FlagEndTime string
var FlagSubSecond bool
var FlagEnablePsmod bool
var FlagEnableSerial bool
var FlagParPeriod int
var FlagParCopies int
var FlagRegex bool
var FlagNetTrace bool
var FlagDashboardName string
var FlagInputMap string
var FlagDashboardSource string

const LineBreaksKey = "lineBreaksAfterEachLogEntry"

func GetTimeAsString(s int64) string {
	t := time.Unix(s, 0)
	switch FlagTimeFormat {
	case "epoch":
		return fmt.Sprintf("%d", s)
	case "RFC3339":
		return t.Format(time.RFC3339)
	default:
		return humanize.Time(t)
	}
}

func PrintResponse(data interface{}) bool {
	if FlagOut == "json" {
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		return true
	} else if FlagOut == "yaml" {
		b, err := yaml.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		return true
	}
	return false
}

func GetStartTime(lastSeen int64) time.Time {
	return time.Unix(lastSeen, 0).UTC()
}

func NeedsLineBreak() bool {
	if val, ok := viper.Get(LineBreaksKey).(bool); ok {
		return val
	}
	return false
}

func GetPageSize() uint32 {
	return FlagLogsPageSize
}

func PrintTable(responseList []types.Resource) {
	table := tablewriter.NewWriter(os.Stdout)
	isHeadersSet := false
	columns := []string{}
	for _, resource := range responseList {
		formattedResponse := resource.GetTableData()
		tableRow := []string{}
		if !isHeadersSet {
			columns = resource.GetColumns()
			table.SetHeader(columns)
			isHeadersSet = true
		}
		for _, col := range columns {
			tableRow = append(tableRow, formattedResponse[col])
		}
		table.Append(tableRow)
	}

	table.Render()
}

func PrintResult(responseList []types.Resource, asList bool) {
	if FlagOut == "table" {
		PrintTable(responseList)
	} else if FlagOut == "json" {
		var jsonOutput []byte
		var jsonErr error
		if asList {
			jsonOutput, jsonErr = json.MarshalIndent(*&responseList, "", " ")
		} else {
			jsonOutput, jsonErr = json.MarshalIndent(*&responseList[0], "", " ")
		}
		if jsonErr != nil {
			fmt.Println("Error: " + jsonErr.Error())
		}
		fmt.Println(string(jsonOutput))
	} else if FlagOut == "yaml" {
		var yamlOutput []byte
		var yamlErr error
		if asList {
			yamlOutput, yamlErr = yaml.Marshal(responseList)
		} else {
			yamlOutput, yamlErr = yaml.Marshal(responseList[0])
		}
		if yamlErr != nil {
			fmt.Println("Error: " + yamlErr.Error())
		} else {
			fmt.Println(string(yamlOutput))
		}
	} else {
		fmt.Println("Error: Unsupported Output Format")
	}
}
