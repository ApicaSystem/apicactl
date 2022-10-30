package types

type VisualizationScale struct {
	Width  float64
	Height float64
	Y      float64
}

func GetVisualizationScale() VisualizationScale {
	return VisualizationScale{
		Width:  312,
		Height: 50,
		Y:      50,
	}
}

type ChartType string

var (
	COUNTER ChartType = "COUNTER"
	CHART   ChartType = "CHART"
	GAUGE   ChartType = "GAUGE"
	TABLE   ChartType = "TABLE"
)
