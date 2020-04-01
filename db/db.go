package db

import (
	"sync"

	"github.com/logiqai/logiqbox/cfg"
)

var once sync.Once

var LATEST = []byte("LATEST_QID")

func UpdateLatestQID(qId string) {
	instance := cfg.GetConfigDB()
	instance.Put(LATEST, []byte(qId))
	instance.Close()
}

func GetLastQuery() (string, error) {
	instance := cfg.GetConfigDB()
	defer instance.Close()
	val, err := instance.Get(LATEST)
	if err != nil {
		return "", nil
	}
	return string(val), nil
}

func HandleGetDataStatus(status string) {
	if "COMPLETE" == status {
		//remove the last qid from cache
		instance := cfg.GetConfigDB()
		defer instance.Close()
		instance.Delete(LATEST)
	}
}
