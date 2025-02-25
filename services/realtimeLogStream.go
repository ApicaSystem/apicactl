/*
Copyright © 2024 apica.io <support@apica.io>

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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ApicaSystem/apicactl/grpc_utils"

	"github.com/ApicaSystem/apicactl/utils"

	"github.com/ApicaSystem/apicactl/api/v1/realtimeLogStream"
	"github.com/ApicaSystem/apicactl/loglerpart"
	"github.com/logiqai/easymap"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	matchAppMap       = map[string]bool{}
	matchNamespaceMap = map[string]bool{}
	matchLabelsMap    = map[string]string{}
	matchProcessMap   = map[string]bool{}

	tailLabels  = false
	tailDefault = false
)

func isMatch(logMap map[string]interface{}) bool {
	match := false

	log.Debugln("isMatch:", tailLabels)
	log.Debugln("isMatch:", logMap["app_name"], logMap["proc_id"])
	log.Debugln("")

	if tailLabels {
		if mRaw, ok := logMap["message_raw"]; ok {
			var v easymap.EasyMap
			if jErr := json.Unmarshal(([]byte)(mRaw.(string)), &v); jErr == nil {
				labelsLookup := v.Get("labels")
				if labels, nsOk := labelsLookup["kubernetes.labels"]; nsOk {
					for k, v := range labels.(map[string]interface{}) {
						if matchValue, found := matchLabelsMap[k]; found {
							if matchValue == v.(string) {
								log.Debugln("Matched label k:", k, " v:", v)
								match = true
							}
						} else {
							return false
						}
					}
				} else {
					return false
				}
			}
		} else {
			return false
		}
	} else {
		match = true
	}

	return match
}

func setupMatchAttributeMaps(matches []string, m map[string]bool) {
	for _, v := range matches {
		m[v] = true
	}
}

func setupMatchAttributeValueMaps(matches []string, sep string, m map[string]string) {
	log.Debugln(matches)
	for _, v := range matches {
		sp := strings.Split(v, sep)
		if len(sp) != 2 {
			log.Fatal("Labels matches must include key and value with a : or = separator")
		}
		m[sp[0]] = sp[1]
	}
}

func Tail(appName, procId string, tL []string) error {
	namespace := utils.GetDefaultNamespace()

	tailLabels = len(tL) > 0
	output := utils.FlagOut
	log.Debugln("A:", appName, "N:", namespace, "L:", tailLabels, "P:", procId)
	setupMatchAttributeValueMaps(tL, ":=", matchLabelsMap)

	log.Debugln(matchNamespaceMap, matchLabelsMap, matchProcessMap, matchAppMap)

	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := realtimeLogStream.NewLogStreamerServiceClient(conn)
	subName := namespace
	if appName != "" {
		subName = fmt.Sprintf("%s:%s", subName, appName)
	}
	if procId != "" {
		subName = fmt.Sprintf("%s:%s", subName, procId)
	}
	//log.Debugf("====> %s \n", subName)
	sub := &realtimeLogStream.Subscription{
		Applications: []string{subName},
	}
	stream, err := client.StreamLog(grpc_utils.GetGrpcContext(), sub)
	if err != nil {
		return err
	}
	if output == OUTPUT_COLUMNS {
		printSyslogHeader()
	}
	var f *os.File
	var writeToFile bool
	var linecnt = 0
	if utils.FlagFile != "" {
		once.Do(func() {
			writeToFile = true
			if _, err := os.Stat(utils.FlagFile); err == nil {
				utils.HandleError2(err, fmt.Sprintf("Output file %s already exists, cannot override", utils.FlagFile))
				//fmt.Printf("Err> Outfile file %s already exists, cannot override, exit\n", utils.FlagFile)
				//os.Exit(1)
			}

			if fTmp, err := os.OpenFile(utils.FlagFile, os.O_CREATE|os.O_WRONLY, 0600); err != nil {
				utils.HandleError2(err, fmt.Sprintf("Unable to write to file %s", utils.FlagFile))
				//fmt.Printf("1 Unable to write to file: %s \n", err.Error())
				//os.Exit(1)
			} else {
				f = fTmp
				fmt.Printf("Logstream Writing output to %s\n", utils.FlagFile)
			}
		})
		defer f.Close()
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		logMap := make(map[string]interface{})

		err = json.Unmarshal([]byte(response.Log), &logMap)
		if err != nil {
			log.Print("Cannot read payload, this should not happen!")
		}
		if isMatch(logMap) {

			logMap["mypp"] = "NonePat"

			if utils.FlagEnablePsmod {
				//if pp=="NoPat" {

				loglerpart.IncLogLineCount()

				msg := logMap["message"].(string)
				PS := loglerpart.ProcessLogCmd(msg)
				pp := loglerpart.PsCheckAndReturnTag(PS, msg)
				logMap["mypp"] = pp
			}

			if writeToFile {
				linecnt += 1
				line := fmt.Sprintf("%s %s %s %s %s %s %s",
					logMap["timestamp"],
					logMap["mypp"],
					logMap["severity_string"],
					logMap["namespace"],
					logMap["app_name"],
					logMap["proc_id"],
					logMap["message"],
				)
				if strings.HasSuffix(line, "\n") {
					line = strings.ReplaceAll(line, "\n", "")
				}
				line = fmt.Sprintf("%s\n", line)

				if _, err := f.WriteString(line); err != nil {
					fmt.Printf("Cannot write file: %s\n", err.Error())
					return nil
					// os.Exit(1)
				}
				if linecnt > utils.FlagMaxLogLines {
					fmt.Printf("Max log line %s reached. Control file size using -m\n", utils.FlagMaxLogLines)

					/*
						if stat, err := os.Stat(utils.FlagFile); err == nil {
							if stat.Size() > int64(utils.FlagMaxFileSize*1048576) {
								fmt.Printf("Max file size reached. Control file size using -m\n")
								_ = stream.CloseSend()
								return nil
								// os.Exit(1)
							}
						}
					*/
				}
			} else {
				printSyslogMessage(logMap, output)
			}
		}
	}
	return nil
}

//map[app_name:CRON facility:10 tag:CRON Id:bf4ef224-dc36-450b-a1d2-5909002f2aaf facility_string:security/authorization hostname:work severity:6 severity_string:info timestamp:2019-06-22T14:43:01.814112+05:30 uuid4:0ceacea4-600f-433d-ace4-7722401583fb message:pam_unix(cron:session): session opened for user tito by (uid=0) priority:86 proc_id:24280 source_ip:127.0.0.1]
