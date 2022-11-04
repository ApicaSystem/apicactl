package types

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

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

func (p *GrafanaPanel) UnmarshalJSON(data []byte) error {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	if ds, found := temp["datasource"]; found {
		if _, ok := ds.(string); !ok {
			datasource := ds.(map[string]interface{})
			temp["datasource"] = datasource["uid"].(string)
		}
	}
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   p,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	decoder.Decode(temp)
	return nil
}

func (p *GrafanaPanel) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
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

func (t *GrafanaTemplate) UnmarshalJSON(data []byte) error {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	if ds, found := temp["datasource"]; found && ds != nil {
		if _, ok := temp["datasource"].(string); !ok {
			datasource := temp["datasource"].(map[string]interface{})
			temp["datasource"] = datasource["uid"].(string)
		}
	}
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   t,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	decoder.Decode(temp)
	return nil
}

func (t *GrafanaTemplate) MarshalJson() ([]byte, error) {
	return json.Marshal(t)
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
