package types

type DatasourceOptions struct {
	Url string `json:"url"`
}

type Datasource struct {
	Id             int               `json:"id"`
	Name           string            `json:"name"`
	DatasourceType string            `json:"type"`
	Options        DatasourceOptions `json:"options"`
}
