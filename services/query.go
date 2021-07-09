/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"fmt"
	"github.com/araddon/dateparse"
	"os"
	"strings"
	"time"
	"errors"

	"github.com/logiqai/logiqctl/grpc_utils"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"

	"github.com/logiqai/logiqctl/loglerpart"
)

const (
	OUTPUT_COLUMNS = "column"
	OUTPUT_RAW     = "raw"
	OUTPUT_JSON    = "json"
)

var st time.Time
var et time.Time

// dateparse is a robust date/time parser
// here's is what it supports
// However, only yyyy-mm-dd hh:MM:ss.sss is used
//
/*
    "May 8, 2009 5:57:51 PM",
    "Mon Jan  2 15:04:05 2006",
    "Mon Jan  2 15:04:05 MST 2006",
    "Mon Jan 02 15:04:05 -0700 2006",
    "Monday, 02-Jan-06 15:04:05 MST",
    "Mon, 02 Jan 2006 15:04:05 MST",
    "Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
    "Mon, 02 Jan 2006 15:04:05 -0700",
    "Thu, 4 Jan 2018 17:53:36 +0000",
    "Mon Aug 10 15:44:11 UTC+0100 2015",
    "Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
    "12 Feb 2006, 19:17",
    "12 Feb 2006 19:17",
    "03 February 2013",
    "2013-Feb-03",
    //   mm/dd/yy
    "3/31/2014",
    "03/31/2014",
    "08/21/71",
    "8/1/71",
    "4/8/2014 22:05",
    "04/08/2014 22:05",
    "4/8/14 22:05",
    "04/2/2014 03:00:51",
    "8/8/1965 12:00:00 AM",
    "8/8/1965 01:00:01 PM",
    "8/8/1965 01:00 PM",
    "8/8/1965 1:00 PM",
    "8/8/1965 12:00 AM",
    "4/02/2014 03:00:51",
    "03/19/2012 10:11:59",
    "03/19/2012 10:11:59.3186369",
    // yyyy/mm/dd
    "2014/3/31",
    "2014/03/31",
    "2014/4/8 22:05",
    "2014/04/08 22:05",
    "2014/04/2 03:00:51",
    "2014/4/02 03:00:51",
    "2012/03/19 10:11:59",
    "2012/03/19 10:11:59.3186369",
    // Chinese
    "2014年04月08日",
    //   yyyy-mm-ddThh
    "2006-01-02T15:04:05+0000",
    "2009-08-12T22:15:09-07:00",
    "2009-08-12T22:15:09",
    "2009-08-12T22:15:09Z",
    //   yyyy-mm-dd hh:mm:ss
    "2014-04-26 17:24:37.3186369",
    "2012-08-03 18:31:59.257000000",
    "2014-04-26 17:24:37.123",
    "2013-04-01 22:43",
    "2013-04-01 22:43:22",
    "2014-12-16 06:20:00 UTC",
    "2014-12-16 06:20:00 GMT",
    "2014-04-26 05:24:37 PM",
    "2014-04-26 13:13:43 +0800",
    "2014-04-26 13:13:44 +09:00",
    "2012-08-03 18:31:59.257000000 +0000 UTC",
    "2015-09-30 18:48:56.35272715 +0000 UTC",
    "2015-02-18 00:12:00 +0000 GMT",
    "2015-02-18 00:12:00 +0000 UTC",
    "2017-07-19 03:21:51+00:00",
    "2014-04-26",
    "2014-04",
    "2014",
    "2014-05-11 08:20:13,787",
    // mm.dd.yy
    "3.31.2014",
    "03.31.2014",
    "08.21.71",
    //  yyyymmdd and similar
    "20140601",
    // unix seconds, ms
    "1332151919",
    "1384216367189",
*/

func timeFormat(t time.Time) string {
	if !utils.FlagSubSecond {
		ts := t.Format(time.RFC3339Nano)
		// fmt.Println("ts=", ts)
		return (ts)
	} else {
		return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%06dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		)
	}
}

func setTimeRange(lastSeen int64) {

	//var st time.Time
	//var et time.Time

	var err error

	if utils.FlagBegTime!="" && utils.FlagEndTime!="" {

		st, err = dateparse.ParseLocal(utils.FlagBegTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
			//fmt.Println("ERR> beg time error string ", utils.FlagBegTime)
			//os.Exit(1)
		}
		et, err = dateparse.ParseLocal(utils.FlagEndTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
			//fmt.Println("ERR> end time error string ", utils.FlagEndTime)
			//os.Exit(1)
		}

	} else if utils.FlagLogsSince!="" && utils.FlagEndTime!="" {

		et, err = dateparse.ParseLocal(utils.FlagEndTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("end time (-e) error string <%s>", utils.FlagEndTime))
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		st = et.Add(-1*d)

	} else if utils.FlagLogsSince!="" && utils.FlagBegTime!="" {

		st, err = dateparse.ParseLocal(utils.FlagBegTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		et = st.Add(d)

	} else if utils.FlagLogsSince!="" {

		if lastSeen > 0 {
			et = utils.GetStartTime(lastSeen)
		} else {
			et = time.Now().UTC()
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		st = et.Add(-1*d)

	} else {
		// default search
		err := errors.New("Need to set search period or search time range (-s -b -e)")
		utils.HandleError(err)
		// utils.HandleError2(err, fmt.Sprintf("Need to set search period or search time range (-s -b -e)"))
	}


	//in.StartTime = timeFormat(st)
	//in.EndTime = timeFormat(et)

	if et.UnixNano()-st.UnixNano() < 0 {

		// doesn't work
		//in.StartTime = et.Format(time.RFC3339Nano)
		//in.EndTime = st.Format(time.RFC3339Nano)
		in.StartTime = timeFormat(st)
		in.EndTime = timeFormat(et)

		// fmt.Printf("ERR> EndTime %s is older than BegTime: %s \n", in.EndTime, in.StartTime)
		// fmt.Printf("ERR> difference is %d \n", et.UnixNano()-st.UnixNano())
		// os.Exit(1)
	}  else {
		in.StartTime = timeFormat(et)
		in.EndTime = timeFormat(st)

	}

	fmt.Println(" StartTime: ", in.EndTime)
	fmt.Println(" EndTime  : ", in.StartTime)

	//fmt.Println("BegTime2:", in.StartTime)
	//fmt.Println("EndTime2:", in.EndTime)
}

var in *query.QueryProperties

func postQuery(applicationName, searchTerm, procId string, lastSeen int64) (string, query.QueryServiceClient, error) {
	//fmt.Println("Enter postQuery2")
	//fmt.Println("searchTerm ", searchTerm)
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return "", nil, err
	}
	client := query.NewQueryServiceClient(conn)

	in = &query.QueryProperties{
		Namespace: utils.GetDefaultNamespace(),
		PageSize:  utils.GetPageSize(),
		QType:     query.QueryType_Fetch,
	}

	if searchTerm!="" {
		setTimeRange(lastSeen)
	}

	if applicationName != "" {
		in.ApplicationNames = []string{applicationName}
		//  process filter is only available for query
		if procId != "" {
			filterValuesMap := make(map[string]*query.FilterValues)
			filterValuesMap["ProcId"] = &query.FilterValues{
				Values: []string{procId},
			}
			in.Filters = filterValuesMap
		}

	}

	if searchTerm != "" {
		in.KeyWord = searchTerm
		in.QType = query.QueryType_Search
	}

	queryResponse, err := client.Query(grpc_utils.GetGrpcContext(), in)
	if err != nil {
		return "", nil, err
	}
	return queryResponse.QueryId, client, nil
}

func DoQuery(appName, searchTerm, procId string, lastSeen int64) {

	//fmt.Println("Enter DoQuery2")
	//fmt.Println("BegTime:", utils.FlagBegTime)
	//fmt.Println("EndTime:", utils.FlagEndTime)

	search := searchTerm != ""

	queryId, client, err := postQuery(appName, searchTerm, procId, lastSeen)
	utils.HandleError(err)
	if queryId != "" {
		var f *os.File
		var writeToFile bool

		if utils.FlagFile != "" {
			once.Do(func() {
				writeToFile = true
				fn, _ := os.Stat(utils.FlagFile)
				if fn!=nil {
					fmt.Printf("Outfile file %s already exists, please remove it before proceed\n", utils.FlagFile)
					os.Exit(-1)
					//utils.HandleError2(err, fmt.Sprintf("Outfile file %s already exists, cannot override", utils.FlagFile))
				}

				if fTmp, err := os.OpenFile(utils.FlagFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
					utils.HandleError2( err, fmt.Sprintf("Err> Unable to write to file: %s \n", err.Error()))
				} else {
					// fmt.Printf("Info> file opened %s\n", utils.FlagFile)
					f = fTmp
					fmt.Printf("Info> Writing output to %s\n", utils.FlagFile)
				}
			})
			defer f.Close()
		}
		for {
			var response *query.GetDataResponse
			var err error
			if search {
				response, err = client.GetDataPrevious(grpc_utils.GetGrpcContext(), &query.GetDataRequest{
					QueryId: queryId,
				})
				if err != nil {
					utils.HandleError(err)
				}
			} else {
				response, err = client.GetDataNext(grpc_utils.GetGrpcContext(), &query.GetDataRequest{
					QueryId: queryId,
				})
				if err != nil {
					utils.HandleError(err)
				}
			}
			if len(response.Data) > 0 {
				for _, entry := range response.Data {

					pp := "NonePat"
					for kk := range entry.StructuredData {
						if entry.StructuredData[kk].Key == "PatternId"	{
							pp = entry.StructuredData[kk].Values[0]
							break
						}
					}
					if utils.FlagEnablePsmod {
						//if pp=="NoPat" {

						loglerpart.IncLogLineCount()

						if pp=="NonePat" {
							msg:=fmt.Sprintf("%s", entry.Message)
							// fmt.Println("msg=", msg)
							PS := loglerpart.ProcessLogCmd(msg)
							pp = loglerpart.PsCheckAndReturnTag(PS, msg)
						}
					}

					if writeToFile {
						line := fmt.Sprintf("%s %s %s %s %s - %s",
							entry.Timestamp,
							pp,
							entry.SeverityString,
							entry.ProcID,
							entry.Message,
						)
						if strings.HasSuffix(line, "\n") {
							line = strings.ReplaceAll(line, "\n", "")
						}
						line = fmt.Sprintf("%s\n", line)
						if _, err := f.WriteString(line); err != nil {
							fmt.Printf("Info> Cannot write file: %s\n", err.Error())
							return
							// os.Exit(1)
						}
						if stat, err := os.Stat(utils.FlagFile); err == nil {
							if stat.Size() > int64(utils.FlagMaxFileSize*1048576) {
								fmt.Printf("Info> Max file size reached. Control file size using -m\n")
								return
								// os.Exit(1)
							}
						}
					} else {
						PrintSyslogMessageForType(entry, "raw")
						time.Sleep(20 * time.Millisecond)
					}
				}
			} else {
				if response.Remaining <= 0 && response.Status == "COMPLETE" {
					return
					//os.Exit(0)
				}
				time.Sleep(2 * time.Second)
			}
		}
	}
}
