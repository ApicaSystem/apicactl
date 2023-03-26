package types

import (
	"fmt"
	"time"
)

type Mapping struct {
	Application string `json:"application"`
	Createts    int    `json:"createts"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Namespace   string `json:"namespace"`
	Updatets    int    `json:"updatets"`
}

func (m Mapping) GetTableData() map[string]string {
	result := make(map[string]string)

	createdAt := time.Unix(int64(m.Createts), 0)
	updatedAt := time.Unix(int64(m.Updatets), 0)

	result["ID"] = fmt.Sprintf("%d", m.ID)
	result["Name"] = m.Name

	result["Application"] = m.Application
	if m.Application == "-1" {
		result["Application"] = "Default"
	}

	result["Namespace"] = m.Namespace
	if m.Namespace == "-1" {
		result["Namespace"] = "Default"
	}

	result["Created At"] = createdAt.Format("02/01/2006 15:04")
	result["Updated At"] = updatedAt.Format("02/01/2006 15:04")

	return result
}

func (m Mapping) GetColumns() []string {
	return []string{
		"ID", "Name", "Namespace", "Application", "Created At", "Updated At",
	}
}