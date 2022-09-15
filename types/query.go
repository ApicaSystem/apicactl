package types

type QueryOptions struct {
	Parameters []interface{} `json:"parameters"`
}

type Query struct {
	Id                int      `json:"id"`
	LatestQueryDataId int      `json:"latest_query_data_id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Query             string   `json:"query"`
	IsArchived        bool     `json:"is_archived"`
	IsDraft           bool     `json:"is_draft"`
	DataSourceId      int      `json:"data_source_id"`
	Version           int      `json:"version"`
	Tags              []string `json:"tags"`
	QueryOptions      `json:"options"`
}

type CreateQueryPayload struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DatasourceId int    `json:"data_source_id"`
	Query        string `json:"query"`
	QueryOptions `json:"options"`
}
