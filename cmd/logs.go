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
package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/logiqctl/services"

	"github.com/spf13/cobra"
)

var logsExample = `
Print logs for logiq-flash ingest server
# logiqctl logs logiq-flash

Print logs in json format
# logiqctl -o=json logs logiq-flash

Print logs for logiq-flash ingest server filtered by process logiq-flash-2
In case of Kubernetes deployment a Stateful Set is an application, and each pods in it is a process
The --process (-p) flag lets you view logs for the individual pod
# logiqctl logs -p=logiq-flash-2 logiq-flash

Runs an interactive prompt to let user choose filters
# logiqctl logs i

Search logs for the given text
# logiqctl logs search "your search term"   

If the flag --follow (-f) is specified the logs will be streamed till it over. 

`

var logsLong = `Logs expect a namespace and application to be available to return results.
Set the default namespace using 'logiqctl set-context' command or pass as '-n=NAMESPACE' flag
Application name needs to be passed as an argument to the command. 
If the user is unsure of the application name, they can run an interactive prompt the would help them to choose filters.  See examples below. 

Search command searches at namespace level, flags -p is ignored. 

Global flag '--time-format' is not applicable for this command.
Global flag '--output' only supports json format for this command.`

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Example: logsExample,
	Aliases: []string{"log"},
	Short:   "Print the logs for an application or process",
	Long:    logsLong,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(cmd.UsageString())
			return
		}
		if utils.FlagProcId != "" {
			proc, err := services.GetProcessByApplicationAndProc(args[0], utils.FlagProcId)
			handleError(err)
			query(args[0], "", proc.ProcID, proc.LastSeen)
			return
		}
		app, err := services.GetApplicationByName(args[0])
		handleError(err)
		query(args[0], "", "", app.LastSeen)
	},
}

func query(appName, searchTerm, procId string, lastSeen int64) {
	queryId, err := services.Query(appName, searchTerm, procId, lastSeen)
	handleError(err)
	if queryId != "" {
		for {
			prompt := promptui.Prompt{
				Label:     "View More ",
				IsConfirm: true,
			}

			result, _ := prompt.Run()
			if strings.ToLower(result) == "y" {
				hasMoreData, err := services.GetDataNext(queryId)
				handleError(err)
				if !hasMoreData {
					break
				}
			} else {
				break
			}
		}
	}
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   `Runs an interactive prompt to let user select application and filters`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := services.RunSelectApplicationForNamespacePrompt()
		handleError(err)
		proc, err := services.RunSelectProcessesForNamespaceAndAppPrompt(app.Name)
		handleError(err)
		fmt.Printf("You could also run this directly `logiqctl logs -p=%s %s`\n", proc.ProcID, app.Name)
		fmt.Printf("Fetching logs for %s (namespace), %s (application) and %s (process)\n\n", utils.GetDefaultNamespace(), app.Name, proc.ProcID)
		query(app.Name, "", proc.ProcID, proc.LastSeen)
	},
}

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Example: ``,
	Short:   `Search for given test in logs`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(cmd.Usage())
			return
		}
		query("", args[0], "", -1)
	},
}

func init() {
	logsCmd.PersistentFlags().StringVarP(&utils.FlagLogsSince, "since", "s", "1h", `Only return logs newer than a relative duration like 2m, 3h, or 2h30m.
This is in relative to the last seen log time for a specified application or processes.`)
	logsCmd.PersistentFlags().Uint32Var(&utils.FlagLogsPageSize, "page-size", 30, `Number of log entries to return in one page`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagLogsFollow, "follow", "f", false, `Specify if the logs should be streamed.`)
	logsCmd.Flags().StringVarP(&utils.FlagProcId, "process", "p", "", `Filter logs by  proc id`)
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(interactiveCmd)
	logsCmd.AddCommand(searchCmd)
}