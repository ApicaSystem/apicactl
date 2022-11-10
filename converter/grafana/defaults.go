package grafana

import "github.com/logiqai/logiqctl/types"

var (
	widthUnits  float64 = 75
	heightUnits float64 = 37
	yUnits      float64 = 38
)

var panelToVisualizationMap = map[string]types.ChartType{
	"stat":       types.COUNTER,
	"barchart":   types.CHART,
	"graph":      types.CHART,
	"gauge":      types.GAUGE,
	"heatmap":    types.CHART,
	"piechart":   types.CHART,
	"table":      types.TABLE,
	"singlestat": types.COUNTER,
	"bargauge":   types.GAUGE,
	"table-old":  types.TABLE,
	"timeseries": types.CHART,
}

var defaultIntervals = map[string]interface{}{
	"locals":      []string{},
	"name":        "interval",
	"title":       "Interval",
	"type":        "enum",
	"value":       "5m",
	"enumOptions": " 5m\n15m\n1h\n6h\n12h\n24h\n2d\n7d\n30d",
}
