package types

type DataExplorerTab struct {
	Key         string                  `json:"key"`
	Order       int                     `json:"order"`
	QueriesList []DataExplorerQueryItem `json:"queriesList"`
	Title       string                  `json:"title"`
	Type        string                  `json:"type"`
}

type DataExplorerQueryItem struct {
	Name      string                  `json:"name"`
	Options   DataExplorerQueryOption `json:"options"`
	Query     string                  `json:"query"`
	Schema    string                  `json:"schema"`
	ChartType *string                 `json:"chart_type,omitempty"`
}

type DataExplorerQueryOption struct {
	Description string                `json:"description"`
	Order       int                   `json:"order"`
	Plot        DataExplorerQueryPlot `json:"plot"`
	UpperLimit  string                `json:"upperLimit"`
}

type DataExplorerQueryPlot struct {
	ErrorColumn string   `json:"errorColumn"`
	GroupBy     string   `json:"groupBy"`
	X           string   `json:"x"`
	XLabel      string   `json:"xLabel"`
	Y           []string `json:"y"`
	YLabel      string   `json:"yLabel"`
}

func NewQueryPlot() DataExplorerQueryPlot {
	return DataExplorerQueryPlot{
		ErrorColumn: "",
		GroupBy:     "",
		X:           "Timestamp",
		XLabel:      "Timestamp",
		Y:           []string{"Timestamp"},
		YLabel:      "value",
	}
}
