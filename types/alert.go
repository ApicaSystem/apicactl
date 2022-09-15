package types

import (
	"fmt"
	"reflect"
	"strconv"
)

type AlertOption struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  int    `json:"value"`
}

type Alert struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	AlertOption   `json:"options"`
	State         string `json:"state"`
	LastTriggered string `json:"last_triggered_at"`
	Rearm         int    `json:"rearm"`
	Query         `json:"query"`
	QueryId       int `json:"query_id"`
}

func (a Alert) GetColumns() []string {
	return []string{
		"ID", "Name", "State", "Last Triggered", "Query", "Command",
	}
}

func (a Alert) GetTableData() map[string]string {
	result := make(map[string]string)
	columns := a.GetColumns()
	outputFields := []outputField{
		{
			label:     columns[0],
			attribute: "Id",
		},
		{
			label:     columns[1],
			attribute: "Name",
		},
		{
			label:     columns[2],
			attribute: "State",
		},
		{
			label:     columns[3],
			attribute: "LastTriggered",
		},
		{
			label:     columns[4],
			attribute: "Query",
		},
		{
			label:     columns[5],
			attribute: "AlertOption",
		},
	}
	r := reflect.ValueOf(a)
	alertReference := reflect.Indirect(r)

	for _, output := range outputFields {
		switch output.attribute {
		case "AlertOption":
			result[output.label] = fmt.Sprintf("%s %s %d", a.AlertOption.Column, a.AlertOption.Op, a.AlertOption.Value)
		case "Query":
			result[output.label] = a.Query.Query
		default:
			value := ""
			field := alertReference.FieldByName(output.attribute)
			if field.Kind().String() == "int" {
				value = strconv.Itoa(int(field.Int()))
			} else {
				value = field.String()
			}
			result[output.label] = value
		}
	}

	return result
}

func (a *Alert) FormatAlert(time_format string) {
	a.LastTriggered = FormatTime(a.LastTriggered, time_format)
}

type CreateAlertPayload struct {
	Name        string `json:"name"`
	AlertOption `json:"options"`
	Rearm       int `json:"rearm"`
	QueryId     int `json:"query_id"`
}
