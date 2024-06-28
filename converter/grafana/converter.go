package grafana

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/ApicaSystem/apicactl/types"
	"github.com/ApicaSystem/apicactl/ui"
)

type emptyCols struct {
	col      float64
	maxSizeY float64
}

type nextCoords struct {
	row float64
	col float64
}

type Converter struct {
	grafanaDashboard *types.GrafanaDashboard
	widgetScale      map[string]int
	inputMap         map[string]map[string]string
	filledRows       map[float64]emptyCols
	nextCoords
}

type visualizationChannelData struct {
	position      int
	visualization *types.Visualization
	panelIndex    int
}

func (c *Converter) createWidgets(widgets []types.Widget, dashboardId int, channelData map[int]visualizationChannelData) ([]types.Widget, error) {
	var result []types.Widget
	for idx, widget := range widgets {
		if widget.Type == "group" {
			var row float64 = 0
			var pos map[string]interface{}
			pos = map[string]interface{}{
				"autoHeight": false,
				"col":        0,
				"maxSizeX":   6,
				"maxSizeY":   1000,
				"minSizeX":   1,
				"minSizeY":   6,
				"row":        row,
				"sizeX":      6,
				"sizeY":      1,
			}
			if idx != 0 {
				prevWidgetPos := result[idx-1].Options["position"].(map[string]interface{})
				row = prevWidgetPos["row"].(float64) + prevWidgetPos["sizeY"].(float64)
				pos["row"] = row
			}
			widget.Options = map[string]interface{}{
				"position": pos,
			}
			w, err := ui.CreateWidgetGroup(widget, dashboardId)
			if err != nil {
				return []types.Widget{}, fmt.Errorf("error: %s", err.Error())
			}

			c.filledRows = map[float64]emptyCols{
				row + 1: {
					maxSizeY: 0,
					col:      0,
				},
			}
			result = append(result, *w)

			continue
		} else {
			visualization := widget.Visualization
			if visualization == nil {
				continue
			}
			if visualization.Query == nil {
				continue
			}
			parameterMappings := c.getParameterMappings(visualization.Query.Parameters)
			widget.Options = map[string]interface{}{
				"position":          c.adjustPanelPosition(channelData[idx].panelIndex),
				"parameterMappings": parameterMappings,
			}
		}
		w, err := ui.CreateWidget(widget, widget.Visualization.Id, dashboardId)
		if err != nil {
			return []types.Widget{}, fmt.Errorf("error: %s", err.Error())
		}
		w.Visualization = widget.Visualization
		result = append(result, w)
	}
	return result, nil
}

func (c *Converter) getParameterMappings(queryParams []map[string]interface{}) map[string]types.WidgetParameterMappings {
	mappings := map[string]types.WidgetParameterMappings{}
	for _, param := range queryParams {
		if param["title"] == "" {
			param["title"] = param["name"]
		}
		mappings[param["name"].(string)] = types.WidgetParameterMappings{
			Name:  param["title"].(string),
			Type:  "dashboard-level",
			MapTo: param["name"].(string),
		}
	}
	return mappings
}

func (c *Converter) getWidgetCoords(sizeX float64, sizeY float64) (float64, float64) {
	var nextRow float64 = 0
	var row float64 = -1
	var col float64 = 0
	for r, space := range c.filledRows {
		if space.col+sizeX <= 6 {
			row = r
			col = space.col
			space.col += sizeX
			if sizeY > space.maxSizeY {
				space.maxSizeY = sizeY
			}
			c.filledRows[r] = space
			break
		}
		next := r + c.filledRows[r].maxSizeY
		if nextRow < next {
			nextRow = next
		}
		if space.col >= 6 {
			delete(c.filledRows, r)
		}
	}
	if row == -1 {
		row = nextRow
		c.filledRows[row] = emptyCols{
			maxSizeY: sizeY,
			col:      sizeX,
		}
	}
	return row, col
}

func (c *Converter) adjustPanelPosition(panelIndex int) map[string]interface{} {
	widgetPos := map[string]interface{}{
		"maxSizeX":   6,
		"maxSizeY":   1000,
		"minSizeX":   1,
		"minSizeY":   6,
		"autoHeight": false,
	}
	grafanaPanel := c.grafanaDashboard.Panels[panelIndex]
	widgetScale := types.GetVisualizationScale()
	widgetPos["sizeX"] = math.Ceil(float64(grafanaPanel.GridPos["w"]) * widthUnits / widgetScale.Width)
	widgetPos["sizeY"] = math.Ceil((float64(grafanaPanel.GridPos["h"])*heightUnits + 35.0) / widgetScale.Height)
	sizeX := widgetPos["sizeX"].(float64)
	sizeY := widgetPos["sizeY"].(float64)
	if sizeY < 8 {
		sizeY = 8
	}
	if sizeX <= 1 {
		if grafanaPanel.Type != "stat" && grafanaPanel.Type != "singlestat" {
			sizeX = 2
		} else {
			sizeX = 1
		}
	}
	row, col := c.getWidgetCoords(sizeX, sizeY)
	widgetPos["row"] = row
	widgetPos["col"] = col
	widgetPos["sizeY"] = sizeY
	widgetPos["sizeX"] = sizeX
	return widgetPos
}

func (c *Converter) CreateAndPublishDashboard(dashboardName string) (*types.DashboardSpec, error) {
	result := types.DashboardSpec{}
	tags := []string{
		"Grafana Import",
	}
	dashboard, err := ui.CreateAndPublishDashboard(dashboardName, tags)
	if err != nil {
		return nil, fmt.Errorf("error creating dashboard: %s", err.Error())
	}
	result.Dashboard = dashboard
	result.Datasources = map[string]types.Datasource{}
	var datasourceId string
	for _, value := range c.inputMap["datasource"] {
		if value == "" {
			return nil, fmt.Errorf("error: Your template is missing with input for datasource. Datasource Id is required to import this dashboard. Please update your template with datasource input")
		}
		datasource, err := ui.GetDatasource(value)
		if err != nil {
			return nil, fmt.Errorf("error: %s", err.Error())
		}
		result.Datasources[value] = *datasource
		datasourceId = value
	}
	result.Widgets = []types.Widget{}
	channelData := map[int]visualizationChannelData{}
	dataChannel := make(chan visualizationChannelData)
	errChannel := make(chan error)
	numOfRoutines := 0
	widgetsCount := 0
	for idx, _ := range c.grafanaDashboard.Panels {
		panel := c.grafanaDashboard.Panels[idx]
		if c.grafanaDashboard.Panels[idx].Type == "row" {
			group := types.Widget{
				Width: 1,
				Type:  "group",
				Text:  panel.Title,
			}
			result.Widgets = append(result.Widgets, group)
			widgetsCount++
			continue
		}
		for targetIdx := 0; targetIdx < len(c.grafanaDashboard.Panels[idx].Targets); targetIdx++ {
			result.Widgets = append(result.Widgets, types.Widget{Width: 1})
			numOfRoutines++
			go convertPanelToWidget(*c.grafanaDashboard, c.inputMap, dashboard.Name, datasourceId, targetIdx, widgetsCount, idx, &dataChannel, &errChannel)
			widgetsCount++
		}
	}
	for i := 0; i < numOfRoutines; i++ {
		select {
		case res := <-dataChannel:
			result.Widgets[res.position].Visualization = res.visualization
			channelData[res.position] = res
		case e := <-errChannel:
			{
				err = e
				break
			}
		}
	}
	close(dataChannel)
	close(errChannel)
	if err != nil {
		return nil, err
	}
	result.Widgets, err = c.createWidgets(result.Widgets, dashboard.Id, channelData)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}
	return &result, nil
}

func (c *Converter) ParseInput(inputData *[]byte) (*map[string]map[string]string, error) {
	result := map[string]map[string]string{}
	if inputData == nil {
		if len(c.grafanaDashboard.Inputs) == 0 {
			for _, t := range c.grafanaDashboard.Templating["list"] {
				if t.Type == "datasource" {
					ds := t.Query.(string)
					input := map[string]string{
						"type":       "datasource",
						"name":       fmt.Sprintf("%s %s", ds, t.Name),
						"label":      t.Label,
						"pluginId":   ds,
						"pluginName": ds,
					}
					c.grafanaDashboard.Inputs = append(c.grafanaDashboard.Inputs, input)
				}
			}
		}
		for _, input := range c.grafanaDashboard.Inputs {
			var reference string
			fmt.Printf("Enter %s %s id for %s: ", input["label"], input["type"], input["name"])
			_, err := fmt.Scanf("%s", &reference)
			if err != nil {
				return nil, fmt.Errorf("error parsing input, %s", err.Error())
			}
			_, ok := result[input["type"]]
			if !ok {
				result[input["type"]] = map[string]string{}
			}
			result[input["type"]][input["name"]] = reference
		}
	} else {
		err := json.Unmarshal(*inputData, &result)
		if err != nil {
			return nil, fmt.Errorf("error reading inputs, %s", err.Error())
		}
	}
	c.inputMap = result
	return &result, nil
}

func NewGrafanaConverter(grafanaDashboard *types.GrafanaDashboard) Converter {
	return Converter{
		grafanaDashboard: grafanaDashboard,
		filledRows: map[float64]emptyCols{
			0: {
				maxSizeY: 0,
				col:      0,
			},
		},
	}
}
