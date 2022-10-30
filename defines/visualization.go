package defines

import "github.com/logiqai/logiqctl/types"

func getTableOptions(options map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"itemsPerPage":   25,
		"autoHeight":     true,
		"defaultRows":    14,
		"defaultColumns": 3,
		"minColumns":     2,
	}
	order := 100000
	columns := options["columns"].([]map[string]interface{})
	var tableColumns []map[string]interface{}
	for _, colOption := range columns {
		option := map[string]interface{}{
			"booleanValues":      []string{"false", "true"},
			"imageUrlTemplate":   "{{ @ }}",
			"imageTitleTemplate": "{{ @ }}",
			"imageWidth":         "",
			"imageHeight":        "",
			"linkUrlTemplate":    "{{ @ }}",
			"linkTextTemplate":   "{{ @ }}",
			"linkTitleTemplate":  "{{ @ }}",
			"linkOpenInNewTab":   true,
			"visible":            true,
			"allowSearch":        false,
			"alignContent":       "left",
			"allowHTML":          true,
			"highlightLinks":     false,
			"order":              order,
		}
		for k, v := range colOption {
			option[k] = v
		}
		tableColumns = append(tableColumns, option)
		order++
	}
	result["columns"] = tableColumns
	return result
}

func getCounterOptions(options map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"counterColName":  "value",
		"rowNumber":       1,
		"targetRowNumber": nil,
		"stringDecimal":   0,
		"stringDecChar":   ".",
		"stringThouSep":   ",",
		"defaultColumns":  2,
		"defaultRows":     5,
		"bgColor":         nil,
		"textColor":       nil,
		"targetColName":   nil,
		"countRow":        false,
	}
	for k, v := range options {
		result[k] = v
	}
	return result
}

func getGaugeOptions(options map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"ranges":     []string{},
		"title":      "",
		"stepRanges": []map[string]interface{}{},
		"color":      "blue",
		"unit":       "percent",
		"colname":    "value",
	}
	for k, v := range options {
		result[k] = v
	}
	return result
}

func getChartOptions(options map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	result["globalSeriesType"] = options["globalSeriesType"]
	result["yAxis"] = []interface{}{
		map[string]string{
			"type": "linear",
		},
		map[string]interface{}{
			"type":     "linear",
			"opposite": true,
		},
	}
	result["xAxis"] = map[string]interface{}{
		"type": options["xMode"],
		"labels": map[string]bool{
			"enabled": true,
		},
	}
	result["error_y"] = map[string]interface{}{
		"type":    "data",
		"visible": true,
	}
	result["series"] = map[string]interface{}{
		"stacking": options["stacking"],
		"error_y":  result["error_y"],
	}
	result["seriesOptions"] = map[string]interface{}{
		"value": map[string]interface{}{
			"type":   options["globalSeriesType"],
			"yAxis":  0,
			"zIndex": 0,
			"index":  0,
		},
	}
	seriesOptions, ok := options["seriesOptions"]
	if ok {
		result["seriesOptions"] = seriesOptions
	}
	result["columnMapping"] = options["columnMapping"]
	result["legend"] = map[string]bool{
		"enabled": options["showLegend"].(bool),
	}
	valuesOptions, ok := options["valuesOptions"]
	if ok {
		result["valuesOptions"] = valuesOptions
	}
	result["numberFormat"] = "0,0[.]00000"
	result["percentFormat"] = "0[.]00%"
	result["textFormat"] = ""
	result["defaultColumns"] = 3
	result["defaultRows"] = 8
	result["minColumns"] = 1
	result["minRows"] = 5
	result["customCode"] = "// Available variables are x, ys, element, and Plotly\n// Type console.log(x, ys); for more info about x and ys\n// To plot your graph call Plotly.plot(element, ...)\n// Plotly examples and docs: https://plot.ly/javascript/"
	result["showDataLabels"] = false
	result["dateTimeFormat"] = "DD/MM/YY HH:mm"
	return result
}

func GetVisualizationOptions(options map[string]interface{}, chartType types.ChartType) map[string]interface{} {
	if chartType == types.TABLE {
		return getTableOptions(options)
	} else if chartType == types.COUNTER {
		return getCounterOptions(options)
	} else if chartType == types.GAUGE {
		return getGaugeOptions(options)
	} else if chartType == types.CHART {
		return getChartOptions(options)
	}
	return map[string]interface{}{}
}
