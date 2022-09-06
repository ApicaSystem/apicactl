package types

import (
	"strconv"
	"time"
)

type Resource interface {
	GetTableData() map[string]string
	GetColumns() []string
}

type outputField struct {
	label     string
	attribute string
}

func FormatTime(t string, timeFormat string) string {
	if t == "" {
		return ""
	}
	parsedTime, _ := time.Parse(time.RFC3339, t)
	if timeFormat == "epoch" {
		return strconv.Itoa(int(parsedTime.Unix()))
	}
	return t
}
