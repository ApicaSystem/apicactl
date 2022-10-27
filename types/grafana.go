package types

type FieldConfig struct {
	Defaults  map[string]interface{}   `json:"defaults"`
	Overrides []map[string]interface{} `json:"overrides"`
}

type Target struct {
	Datasource   map[string]string `json:"datasource"`
	Query        string            `json:"expr"`
	Interval     string            `json:"interval"`
	LegendFormat string            `json:"legendFormat"`
}

type GrafanaPanel struct {
	FieldConfig `json:"fieldConfig"`
	GridPos     map[string]int         `json:"gridPos"`
	Targets     []Target               `json:"targets"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Options     map[string]interface{} `json:"options"`
	XAxis       map[string]interface{} `json:"xaxis"`
	Format      string                 `json:"format"`
	Datasource  string                 `json:"datasource"`
}

type GrafanaTemplate struct {
	Datasource string                   `json:"datasource"`
	Label      string                   `json:"label"`
	Name       string                   `json:"name"`
	Query      interface{}              `json:"query"`
	Refresh    int                      `json:"refresh"`
	Type       string                   `json:"type"`
	Options    []map[string]interface{} `json:"options,omitempty"`
}

type GrafanaRow struct {
	Title  string         `json:"title"`
	Panels []GrafanaPanel `json:"panels"`
}

type GrafanaDashboard struct {
	Inputs     []map[string]string          `json:"__inputs"`
	Requires   []map[string]string          `json:"__requires"`
	Panels     []GrafanaPanel               `json:"panels"`
	Templating map[string][]GrafanaTemplate `json:"templating"`
	PanelRows  []GrafanaRow                 `json:"rows"`
	Title      string                       `json:"title"`
}

const (
	GrafanaLineChart string = "graph"
	GrafanaBarChart  string = "barchart"
	GrafanaCounter   string = "counter"
	GrafanaHeatMap   string = "heatmap"
	GrafanaPieChart  string = "piechart"
	GrafanaBarGauge  string = "bargauge"
)
