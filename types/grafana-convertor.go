package types

type ConverterGrafanaDashboard struct {
	Panels *[]ConverterGrafanaPanel `json:"panels"`
	Rows   *[]ConverterGrafanaRow   `json:"rows"`
}

type ConverterGrafanaRow struct {
	Panels []ConverterGrafanaPanel `json:"panels"`
}

type ConverterGrafanaPanel struct {
	Title   string                   `json:"title"`
	Targets []ConverterGrafanaTarget `json:"targets"`
	IsLines *bool                    `json:"lines"`
	Panels  *[]ConverterGrafanaPanel `json:"panels"`
}

type ConverterGrafanaTarget struct {
	QueryExpr *string `json:"expr"`
	Legend    *string `json:"legendFormat"`
}
