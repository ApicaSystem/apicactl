package grafana

import (
	"bytes"
	"fmt"
	"github.com/logiqai/logiqctl/defaults"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var mutexLock = sync.RWMutex{}
var queryCache = map[string]types.QueryResultData{}

type grafanaWorker struct {
	grafanaDashboard types.GrafanaDashboard
	inputMap         map[string]map[string]string
	panelIndex       int
	dashboardTitle   string
}

func (w grafanaWorker) parseLegend(legendFormat string) string {
	for startPos := strings.Index(legendFormat, "{{"); startPos != -1; startPos = strings.Index(legendFormat, "{{") {
		endPos := strings.Index(legendFormat, "}}")
		if startPos == 0 {
			legendFormat = legendFormat[endPos+2:]
		} else {
			legendFormat = legendFormat[:startPos-1] + legendFormat[endPos+2:]
		}
	}
	var words []string
	var buffer bytes.Buffer
	for _, ch := range legendFormat {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			buffer.WriteRune(ch)
		} else {
			words = append(words, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		words = append(words, buffer.String())
		buffer.Reset()
	}
	return strings.Trim(strings.Join(words, " "), " ")
}

func convertPanelToWidget(grafanaDashboard types.GrafanaDashboard, inputMap map[string]map[string]string, dashboardTitle string, datasourceId string, targetIdx int, widgetIdx int, panelIndex int, dataChannel *chan visualizationChannelData, errChannel *chan error) {
	result := types.Widget{}
	w := &grafanaWorker{
		grafanaDashboard: grafanaDashboard,
		inputMap:         inputMap,
		panelIndex:       panelIndex,
		dashboardTitle:   dashboardTitle,
	}
	result.Width = 1
	legendFormat := w.grafanaDashboard.Panels[panelIndex].Targets[targetIdx].LegendFormat
	title := w.grafanaDashboard.Panels[panelIndex].Title
	if legendFormat != "" {
		legendFormat = w.parseLegend(legendFormat)
		if legendFormat != "" {
			title += " - " + strings.Trim(legendFormat, " ")
		}
	}
	visualization, err := w.convertVisualization(datasourceId, targetIdx, title)
	if err != nil {
		*errChannel <- err
		return
	}
	if err != nil {
		*errChannel <- err
		return
	}
	result.Visualization = visualization
	*dataChannel <- visualizationChannelData{
		visualization: visualization,
		position:      widgetIdx,
		panelIndex:    panelIndex,
	}
}

func (w grafanaWorker) convertVisualization(datasourceId string, targetIdx int, title string) (*types.Visualization, error) {
	var payload = types.Visualization{}
	panel := w.grafanaDashboard.Panels[w.panelIndex]
	payload.Name = title
	if panel.Format == "bytes" {
		payload.Name += " (MB)"
	} else if panel.Format == "s" {
		payload.Name += " (Min)"
	}
	payload.Type = panelToVisualizationMap[panel.Type]
	newQuery, err := w.convertQuery(panel, datasourceId, targetIdx, title)
	payload.Options = w.getVisualizationOptions(payload.Type, targetIdx, newQuery)
	if err != nil {
		return nil, err
	}
	result, err := ui.CreateVisualization(&payload, newQuery.Id)
	if err != nil {
		return nil, fmt.Errorf("error creating visualization, %s", err.Error())
	}
	result.Query = &newQuery
	return result, nil
}

func colInLegend(legend string, columns []*types.QueryResultColumn) string {
	t := strings.Split(legend, " ")
	for i := len(t) - 1; i >= 0; i-- {
		for _, col := range columns {
			if t[i] == fmt.Sprintf("{{%s}}", col.Name) {
				return col.Name
			}
		}
	}
	return ""
}

func (w grafanaWorker) getChartOptions(grafanaPanel types.GrafanaPanel, targetIdx int, query types.Query) map[string]interface{} {
	options := map[string]interface{}{
		"stacking":   "none",
		"showLegend": true,
		"xMode":      "-",
	}
	if grafanaPanel.Type == types.GrafanaBarChart {
		options["globalSeriesType"] = "column"
	} else if grafanaPanel.Type == types.GrafanaHeatMap {
		options["globalSeriesType"] = "heatmap"
	} else if grafanaPanel.Type == types.GrafanaPieChart {
		options["globalSeriesType"] = "pie"
	} else {
		options["globalSeriesType"] = "line"
	}
	grafanaLegend, _ := grafanaPanel.Options["legend"].(map[string]interface{})
	showLegend := true
	if grafanaLegend != nil {
		if legend, ok := grafanaLegend["showLegend"]; ok {
			showLegend = legend.(bool)
		}
	}
	options["showLegend"] = showLegend
	xMode, ok := grafanaPanel.XAxis["mode"]
	if ok {
		if xMode == "time" {
			options["xMode"] = "datetime"
		}
	}
	stacking, ok := grafanaPanel.Options["stacking"].(string)
	if ok && stacking != "none" {
		options["stacking"] = "stack"
	} else {
		options["stacking"] = nil
	}
	options["columnMapping"] = map[string]string{
		"value":     "y",
		"timestamp": "x",
	}
	legend := grafanaPanel.Targets[targetIdx].LegendFormat
	if strings.Contains(legend, "{{") {
		queryResult := w.executeQuery(query.Query, fmt.Sprintf("%d", query.DataSourceId), &query.Parameters)
		if queryResult != nil {
			seriesOptions := map[string]map[string]interface{}{}
			valueOptions := map[string]interface{}{}
			if colName := colInLegend(legend, queryResult.Columns); colName != "" {
				if options["globalSeriesType"] != "pie" {
					options["columnMapping"].(map[string]string)[colName] = "series"
				} else {
					for k, v := range options["columnMapping"].(map[string]string) {
						if v == "x" {
							delete(options["columnMapping"].(map[string]string), k)
							break
						}
					}
					options["columnMapping"].(map[string]string)[colName] = "x"
				}
				for idx, row := range queryResult.Rows {
					if options["globalSeriesType"] == "pie" {
						valueOptions[(*row)[colName]] = map[string]string{}
					} else {
						series := map[string]interface{}{
							"type":   options["globalSeriesType"],
							"yAxis":  0,
							"zIndex": idx,
							"index":  0,
						}
						seriesOptions[(*row)[colName]] = series
					}
				}
			}
			if len(seriesOptions) > 0 {
				options["seriesOptions"] = seriesOptions
			}
			if len(valueOptions) > 0 {
				options["valuesOptions"] = valueOptions
			}
		}
	}
	result := defaults.GetVisualizationOptions(options, types.CHART)
	return result
}

func (w grafanaWorker) getGaugeOptions(grafanaPanel types.GrafanaPanel) map[string]interface{} {
	options := map[string]interface{}{}
	grafanaFieldConfig := grafanaPanel.FieldConfig.Defaults
	options["unit"] = grafanaFieldConfig["unit"]
	options["title"] = grafanaPanel.Title
	stepRanges := []map[string]interface{}{}
	steps := ((grafanaFieldConfig["thresholds"].(map[string]interface{}))["steps"]).([]interface{})
	min := grafanaFieldConfig["min"].(float64)
	for i := 1; i < len(steps); i++ {
		t := steps[i]
		s := map[string]interface{}{}
		step := t.(map[string]interface{})
		prev := steps[i-1].(map[string]interface{})
		if step["value"] != nil {
			s["color"] = prev["color"]
			s["min"] = fmt.Sprintf("%.0f", min)
			s["max"] = fmt.Sprintf("%.0f", step["value"])
			min = step["value"].(float64) + 1.0
			stepRanges = append(stepRanges, s)
		}
	}
	stepRanges = append(stepRanges, map[string]interface{}{
		"min":   fmt.Sprintf("%.0f", min),
		"max":   fmt.Sprintf("%.0f", grafanaFieldConfig["max"]),
		"color": steps[len(steps)-1].(map[string]interface{})["color"],
	})
	if len(stepRanges) > 0 {
		options["stepRanges"] = stepRanges
	}
	return defaults.GetVisualizationOptions(options, types.GAUGE)
}

func (w grafanaWorker) getTableOptions(query types.Query) map[string]interface{} {
	tableOptions := []map[string]interface{}{}
	queryResult := w.executeQuery(query.Query, fmt.Sprintf("%d", query.DataSourceId), &query.Parameters)
	if queryResult == nil {
		return map[string]interface{}{}
	}
	for _, col := range queryResult.Columns {
		option := map[string]interface{}{}
		option["name"] = col.Name
		option["type"] = col.Type
		option["title"] = col.FriendlyName
		option["displayAs"] = col.Type
		if col.Type == "datetime" {
			option["dateTimeFormat"] = "DD/MM/YY HH:mm"
		}
		tableOptions = append(tableOptions, option)
	}
	options := map[string]interface{}{
		"columns": tableOptions,
	}
	result := defaults.GetVisualizationOptions(options, types.TABLE)
	return result
}

func (w grafanaWorker) getCounterOptions(grafanaPanel types.GrafanaPanel) map[string]interface{} {
	options := map[string]interface{}{
		"counterLabel": grafanaPanel.Title,
	}
	result := defaults.GetVisualizationOptions(options, types.COUNTER)
	return result
}

func (w grafanaWorker) getVisualizationOptions(visualizationType types.ChartType, targetIdx int, query types.Query) map[string]interface{} {
	grafanaPanel := w.grafanaDashboard.Panels[w.panelIndex]
	if visualizationType == types.GAUGE {
		return w.getGaugeOptions(grafanaPanel)
	} else if visualizationType == types.TABLE {
		return w.getTableOptions(query)
	} else if visualizationType == types.COUNTER {
		return w.getCounterOptions(grafanaPanel)
	} else {
		return w.getChartOptions(grafanaPanel, targetIdx, query)
	}
}

func (w grafanaWorker) convertQuery(grafanaPanel types.GrafanaPanel, datasourceId string, targetIdx int, title string) (types.Query, error) {
	queryPayload := types.CreateQueryPayload{}
	templateMappings := w.getTemplateMappings(w.grafanaDashboard.Templating, grafanaPanel.Type)
	queryExp, queryMappings := w.parseGrafanaQuery(grafanaPanel.Targets[targetIdx].Query, templateMappings)
	queryExp = w.formatQuery(queryExp, grafanaPanel.Format)
	if grafanaPanel.Type != types.GrafanaCounter {
		queryExp += "&duration={{duration}}&step={{step}}"
		queryMappings = append(queryMappings, templateMappings[0], templateMappings[1])
	}
	queryPayload.Query = queryExp
	queryPayload.Name = title + " " + "query"
	queryPayload.DatasourceId, _ = strconv.Atoi(datasourceId)
	queryPayload.QueryOptions = types.QueryOptions{
		Parameters: queryMappings,
	}
	queryPayload.Tags = []string{
		w.dashboardTitle,
	}
	query, err := ui.CreateQuery(queryPayload)
	if err != nil {
		return types.Query{}, fmt.Errorf("error creating query, %s", err.Error())
	}
	publishArgs := []string{fmt.Sprintf("%d", query.Id), fmt.Sprintf("%d", query.Version)}
	_, err = ui.PublishQuery(publishArgs)
	if err != nil {
		fmt.Printf("error publishing query %d: %s\n", query.Id, err.Error())
	} else {
		query.IsDraft = false
	}
	return query, nil
}

func (w grafanaWorker) getTemplateMappings(templates map[string][]types.GrafanaTemplate, chartType string) []map[string]interface{} {
	var queryMappings []map[string]interface{}
	if chartType != types.GrafanaCounter {
		queryMappings = append(queryMappings, map[string]interface{}{
			"locals": []string{},
			"name":   "duration",
			"title":  "Duration",
			"type":   "text",
			"value":  "1h",
		},
			map[string]interface{}{
				"locals": []string{},
				"name":   "step",
				"title":  "Step",
				"type":   "text",
				"value":  "1m",
			})
	}
	for _, templateList := range templates {
		for _, t := range templateList {
			if t.Label == "" {
				t.Label = t.Name
			}
			m := map[string]interface{}{
				"locals": []string{},
				"name":   t.Name,
				"title":  strings.ToUpper(t.Label[0:1]) + t.Label[1:],
				"type":   "text",
				"value":  "",
			}
			if t.Type == "query" {
				var query string
				qType := reflect.TypeOf(t.Query)
				if qType.Kind() == reflect.String {
					query = t.Query.(string)
				} else {
					query = t.Query.(map[string]interface{})["query"].(string)
				}
				enumOptions, err := w.getEnumOptions(query, t.Datasource, t.Name, queryMappings)
				if err != nil {
					fmt.Println(err.Error())
					m["enumOptions"] = ""
				}
				m["type"] = "enum"
				m["value"] = ""
				m["enumOptions"] = ""
				if len(enumOptions) > 0 {
					m["enumOptions"] = strings.Join(enumOptions, "\n")
					m["value"] = enumOptions[0]
				}
			} else if t.Type == "interval" || t.Type == "custom" {
				intervals := t.Query.(string)
				m["enumOptions"] = strings.Replace(intervals, ",", "\n", len(intervals))
				m["type"] = "enum"
				for _, option := range t.Options {
					if option["selected"].(bool) {
						m["value"] = option["value"].(string)
						break
					}
				}
			} else if t.Type == "textbox" {
				m["value"] = t.Query.(string)
			}
			queryMappings = append(queryMappings, m)
		}
	}

	return queryMappings
}

func (w grafanaWorker) formatQuery(query, format string) string {
	if format == "bytes" {
		query += " / (1024 ^ 2)"
		query = fmt.Sprintf("(%s)", query)
	} else if format == "s" {
		query += " / 60"
		query = fmt.Sprintf("(%s)", query)
	}
	return query
}

func (w grafanaWorker) removeVariableFromQuery(query string, mapping map[string]interface{}) string {
	for {
		varIndex := strings.Index(query, "$"+mapping["name"].(string))
		if varIndex == -1 {
			break
		}
		beforeVar := 0
		for i := varIndex - 1; i >= 0; i-- {
			if query[i] == '{' || query[i] == '[' {
				beforeVar = i
				break
			}
			if query[i] == '|' || query[i] == ',' {
				beforeVar = i - 1
				break
			}
		}
		exp := fmt.Sprintf("$%s", mapping["name"])
		varEndIndex := varIndex + len(exp)
		for {
			if query[varEndIndex] == ',' {
				varEndIndex++
				break
			}
			if query[varEndIndex] == '}' {
				break
			}
			if query[varEndIndex] == ']' {
				break
			}
			varEndIndex++
		}
		query = query[:beforeVar+1] + query[varEndIndex:]
	}
	return query
}

func (w grafanaWorker) parseGrafanaQuery(query string, queryParameters []map[string]interface{}) (string, []map[string]interface{}) {
	var usedVars []map[string]interface{}
	var intervalParam map[string]interface{}
	for _, v := range queryParameters {
		var mapping map[string]interface{}
		mapping = v
		if mapping["value"] == "" {
			query = w.removeVariableFromQuery(query, mapping)
		} else {
			exp := fmt.Sprintf("{{%s}}", mapping["name"])
			varIndex := strings.Index(query, "$"+mapping["name"].(string))
			if varIndex != -1 {
				usedVars = append(usedVars, mapping)
			}
			query = strings.Replace(query, "$"+mapping["name"].(string), exp, -1)
		}
		if mapping["name"] == "interval" {
			intervalParam = mapping
		}
	}
	if strings.Index(query, "$__interval") != -1 {
		query = strings.Replace(query, "$__interval", "{{interval}}", -1)
		usedVars = append(usedVars, intervalParam)
	}
	return query, usedVars
}

func (w grafanaWorker) getEnumOptions(query string, dsVar string, name string, templateMappings []map[string]interface{}) ([]string, error) {
	key := "label_values("
	if strings.Contains(query, "query_result(") {
		key = "query_result("
	}
	query = strings.Replace(query, key, "", -1)
	query = query[:len(query)-1]
	q := strings.Split(query, "}")
	if len(q) == 1 {
		q = append(q, name)
	}
	q[0] += "}"
	if key == "label_values(" {
		if len(q) > 1 && q[1][0] == ',' {
			q[1] = q[1][1:]
		} else {
			q = strings.Split(query, ",")
		}
	}
	if len(q) > 1 {
		q[1] = strings.Trim(q[1], " ")
	}
	for i := 0; i < len(dsVar); i++ {
		if (dsVar[i] >= 'A' && dsVar[i] <= 'Z') || (dsVar[i] >= 'a' && dsVar[i] <= 'z') {
			dsVar = dsVar[i:]
			break
		}
	}
	if dsVar[len(dsVar)-1] == '}' {
		dsVar = dsVar[:len(dsVar)-1]
	}
	datasourceId, _ := w.inputMap["datasource"][dsVar]
	var queryMappings []map[string]interface{}
	query, queryMappings = w.parseGrafanaQuery(q[0], templateMappings)
	queryResult := w.executeQuery(query, datasourceId, &queryMappings)
	var enumOptions []string
	if queryResult == nil {
		return []string{}, nil
	}
	for _, row := range queryResult.Rows {
		var val string
		if len(q) > 1 {
			val = (*row)[strings.Trim(q[1], " ")]
		} else {
			val = (*row)[name]
		}
		if len(val) > 0 {
			enumOptions = append(enumOptions, val)
		}
	}
	return enumOptions, nil
}

func (w grafanaWorker) executeQuery(query string, datasourceId string, parameters *[]map[string]interface{}) *types.QueryResultData {
	mutexLock.RLock()
	if result, ok := queryCache[query]; ok && result.Processed {
		mutexLock.RUnlock()
		return &result
	}
	mutexLock.RUnlock()
	ds, _ := strconv.Atoi(datasourceId)
	result, err := ui.ExecuteRawQueries(query, ds, parameters)
	if err != nil {
		mutexLock.Lock()
		queryCache[query] = types.QueryResultData{
			Processed: true,
		}
		mutexLock.Unlock()
		return nil
	}
	mutexLock.Lock()
	result.Processed = true
	queryCache[query] = *result
	mutexLock.Unlock()
	return result
}
