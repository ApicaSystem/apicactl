package types

type Dashboard struct {
	Id   int      `json:"id,omitempty"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type DashboardSpec struct {
	Dashboard   Dashboard             `json:"dashboard"`
	Datasources map[string]Datasource `json:"datasources"`
	Widgets     []Widget              `json:"widgets,omitempty"`
	Alerts      []Alert               `json:"alerts,omitempty"`
}

type WidgetOptions struct {
	IsHidden          bool                         `json:"isHidden"`
	ParameterMappings map[string]map[string]string `json:"parameterMappings"`
	Position          map[string]interface{}       `json:"position"`
}

type Visualization struct {
	Id          int                    `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Type        ChartType              `json:"type"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
	*Query      `json:"query"`
}

type VisualizationPayload struct {
	Name    string                 `json:"name"`
	Type    ChartType              `json:"type"`
	Options map[string]interface{} `json:"options"`
	Query   `json:"query"`
}

type Widget struct {
	Id             int                    `json:"id,omitempty"`
	Text           string                 `json:"text"`
	Width          int                    `json:"width"`
	Options        map[string]interface{} `json:"options"`
	*Visualization `json:"visualization"`
}

type WidgetParameterMappings struct {
	MapTo string `json:"mapTo"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
