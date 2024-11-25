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
	Metric       string            `json:"metric"`
}

type Row struct {
	Panels []GrafanaPanel `json:"panels"`
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
	IsLines     *bool                  `json:"lines"`
	Panels      *[]GrafanaPanel        `json:"panels"`
}

func (p *GrafanaPanel) UnmarshalJSON(data []byte) error {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	if ds, found := temp["datasource"]; found && ds != nil {
		if _, ok := ds.(string); !ok && ds != nil {
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
	Inputs     []map[string]string `json:"__inputs"`
	Requires   []map[string]string `json:"__requires"`
	Panels     []GrafanaPanel
	Rows       *[]Row `json:"rows"`
	Templating map[string][]GrafanaTemplate
	Title      string `json:"title"`
}

func (d *GrafanaDashboard) parseGrafanaRows(data map[string]interface{}) {
	var parsed []GrafanaPanel
	var panels []interface{}
	if _, ok := data["panels"]; ok {
		panels = data["panels"].([]interface{})
	}
	if _, ok := data["rows"]; ok {
		for _, i := range data["rows"].([]interface{}) {
			r := i.(map[string]interface{})
			r["type"] = "row"
			panels = append(panels, r)
		}
	}
	for _, i := range panels {
		r := i.(map[string]interface{})
		if r["type"].(string) == "row" {
			row := GrafanaPanel{}
			row.Title = r["title"].(string)
			row.Type = "row"
			parsed = append(parsed, row)
			if embeddedPanels, found := r["panels"]; found && embeddedPanels != nil {
				for _, j := range embeddedPanels.([]interface{}) {
					embeddedPanel := j.(map[string]interface{})
					jsonStr, _ := json.Marshal(embeddedPanel)
					p := GrafanaPanel{}
					json.Unmarshal(jsonStr, &p)
					parsed = append(parsed, p)
				}
			}
		} else {
			jsonStr, _ := json.Marshal(r)
			p := GrafanaPanel{}
			json.Unmarshal(jsonStr, &p)
			parsed = append(parsed, p)
		}
	}
	d.Panels = parsed
}

func (d *GrafanaDashboard) UnmarshalJSON(data []byte) error {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   d,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	decoder.Decode(temp)
	d.parseGrafanaRows(temp)
	templateStr, _ := json.Marshal(temp["templating"])
	json.Unmarshal(templateStr, &d.Templating)
	return nil
}

const (
	GrafanaLineChart string = "graph"
	GrafanaBarChart  string = "barchart"
	GrafanaCounter   string = "counter"
	GrafanaHeatMap   string = "heatmap"
	GrafanaPieChart  string = "piechart"
	GrafanaBarGauge  string = "bargauge"
)
