package services

import (
	"fmt"
	"log"

	"github.com/logiqai/logiq-box/api/v1/query"
	"github.com/logiqai/logiq-box/cfg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(config *cfg.Config, err error) {
	errorStatus := status.Convert(err)
	switch errorStatus.Code() {
	case codes.Unavailable:
		log.Printf("Could not connect, Are you sure you can connect to [%s]?", config.Cluster)
		log.Println("Please verify your configuration...")
		fmt.Println("++++++++++++++++++++++")
		fmt.Print(config.String())
		fmt.Println("++++++++++++++++++++++")
		break
	case codes.Unauthenticated:
		log.Println("Invalid Token, You may need to Login")
		break
	default:
		log.Printf("Error Code: [%d, %s]", errorStatus.Code(), errorStatus.Code().String())
		break
	}
}

func printSyslogHeader() {
	fmt.Println("timestamp|severity_string|hostname|source_ip|proc_id|app_name|facility_string|message")
}

func printSyslogMessage(logMap map[string]interface{}) {
	fmt.Printf("%-33s|%-6s|%s|%s|%-5s|%s|%s|%s\n",
		logMap["timestamp"],
		logMap["severity_string"],
		logMap["hostname"],
		logMap["source_ip"],
		logMap["proc_id"],
		logMap["app_name"],
		logMap["facility_string"],
		logMap["message"],
	)
}

func printSyslogMessageForType(log *query.SysLogMessage) {
	fmt.Printf("%-28s|%-6s|%s|%s|%-5s|%s|%s\n",
		log.Timestamp,
		log.SeverityString,
		log.Hostname,
		log.ProcID,
		log.AppName,
		log.FacilityString,
		log.Message,
	)
}
