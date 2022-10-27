package types

type QueryOptions struct {
	Parameters []map[string]interface{} `json:"parameters"`
}

type QuerySchedule struct {
	DayOfWeek string `json:"day_of_week"`
	Interval  int    `json:"interval"`
	Time      string `json:"time"`
	Until     string `json:"until"`
}

type Query struct {
	Id                int      `json:"id,omitempty"`
	LatestQueryDataId int      `json:"latest_query_data_id,omitempty"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Query             string   `json:"query"`
	IsArchived        bool     `json:"is_archived"`
	IsDraft           bool     `json:"is_draft"`
	DataSourceId      int      `json:"data_source_id,omitempty"`
	Version           int      `json:"version"`
	Tags              []string `json:"tags"`
	QueryOptions      `json:"options"`
	*QuerySchedule    `json:"schedule"`
}

type CreateQueryPayload struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	DatasourceId   int      `json:"data_source_id"`
	Query          string   `json:"query"`
	Tags           []string `json:"tags"`
	QueryOptions   `json:"options"`
	*QuerySchedule `json:"schedule"`
}

type JobStatus int

const (
	READY      JobStatus = 1 << iota
	PROCESSING JobStatus = iota + 1
	FINISHED   JobStatus = iota + 1
	ERROR      JobStatus = iota + 1
)

type QueryResultColumn struct {
	FriendlyName string `json:"friendly_name"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type QueryResultData struct {
	Columns   []*QueryResultColumn `json:"columns"`
	Rows      []*map[string]string `json:"rows"`
	Processed bool                 `json:"-"`
}

type QueryResult struct {
	Id           int               `json:"id,omitempty"`
	DataSourceId int               `json:"data_source_id"`
	Query        string            `json:"query"`
	Data         QueryResultData   `json:"data,omitempty"`
	Runtime      float64           `json:"runtime,omitempty"`
	RetrievedAt  string            `json:"retrieved_at,omitempty"`
	Parameters   map[string]string `json:"parameters"`
	MaxAge       int               `json:"max_age"`
	QueryHash    string            `json:"query_hash"`
}

type JobDetails struct {
	Error         string    `json:"error"`
	Id            string    `json:"id"`
	QueryResultId int       `json:"query_result_id,omitempty"`
	Status        JobStatus `json:"status"`
}
