// Package cmd
//Copyright © 2020 Logiq.ai <cli@logiq.ai>
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//

package cmd

import (
	"fmt"
	"github.com/logiqai/logiqctl/api/v1/applications"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/logiqctl/services"

	"github.com/spf13/cobra"
)

var logsExample = `
Print logs for the LOGIQ ingest server
- logiqctl logs -a <application_name>

Print logs in JSON format
- logiqctl -o=json logs -a <application_name>

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for logiq-flash ingest server filtered by process logiq-flash-2
The --process (-p) flag lets you view logs for the individual pod
- logiqctl logs -p=<proc_id> -a <application_name>

Runs an interactive prompt that lets you choose filters
- logiqctl logs interactive|i

Search logs for specific keywords or terms
- logiqctl logs -a <application_name> search <searchterm>
- logiqctl logs -a <application_name> -p <proc_id> search <searchterm>

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

- stream logs contains log pattern-signature (PS).
- Example:  % logiqctl config set-context <namespace>
            % logiqctl logs -a <proc_id> -s 10s -f 
            % logiqctl logs -a <application_name> -p -s 10s -f
            % logiqctl logs -a <application_name> -s 10s -w outputfile.txt
  (You might want to pipe above dump into file for later cross-reference)
- after done logs streaming, two files will be created.
  notice that these files are reset for every logs query session.
  * ps_stat.out: compute byte and log counts and percentage for each pattern signature 

`

var logsLong = `
The 'logs' command is used to view historical logs. This command expects a namespace and an application to be available to return results. You can set the default namespace using the 'logiqctl set-context' command or pass the namespace as '-n=NAMESPACE' flag. The application name also needs to be passed as an argument to the command. You can also use the 'interactive' command to choose from the list of available applications and processes.   

**Note:**
- The global flag '--time-format' is not applicable for this command.
- The global flag '--output' only supports JSON format for this command.`

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Example: logsExample,
	Aliases: []string{"log"},
	Short:   "View logs for the given namespace and application",
	Long:    logsLong,
	Run: func(cmd *cobra.Command, args []string) {
		hasApp := false
		hasProc := false
		if utils.FlagProcId != "" {
			hasProc = true
		}
		if utils.FlagAppName != "" {
			hasApp = true
		}
		if hasApp && hasProc {
			proc, err := services.GetProcessByApplicationAndProc(utils.FlagAppName, utils.FlagProcId)
			handleError(err)
			services.DoQuery(utils.FlagAppName, "", proc.ProcID, proc.LastSeen)
			return
		} else if hasApp {
			app, err := services.GetApplicationByName(utils.FlagAppName)
			handleError(err)
			services.DoQuery(utils.FlagAppName, "", "", app.LastSeen)
		} else {
			fmt.Println(cmd.UsageString())
			return
		}
	},
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   `Runs an interactive prompt to display logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := services.RunSelectApplicationForNamespacePrompt(false)
		handleError(err)
		proc, err := services.RunSelectProcessesForNamespaceAndAppPrompt(app.Name, false)
		handleError(err)
		fmt.Printf("You could also run this directly `logiqctl logs -p=%s %s`\n", proc.ProcID, app.Name)
		fmt.Printf("Fetching logs for %s (namespace), %s (application) and %s (process)\n\n", utils.GetDefaultNamespace(), app.Name, proc.ProcID)
		services.DoQuery(app.Name, "", proc.ProcID, proc.LastSeen)
	},
}

var searchCmd = &cobra.Command{
	Use:     "search [searchterm]",
	Aliases: []string{"s"},
	Example: "logiqctl -a <application_name> -p <proc_id> logs search <somestring>\nlogiqctl -a <application_name> logs search <somestring>",
	Short:   `Search logs for specific keywords or terms.`,
	Long:    `Search for specific keywords or terms in logs within a namespace, app, proc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Args : ", args)
		if len(args) != 1 {
			fmt.Println(cmd.Usage())
			return
		}
		var app *applications.ApplicationV2 = nil
		hasApp := false
		hasProc := false
		if utils.FlagAppName == "" {
			a, err := services.RunSelectApplicationForNamespacePrompt(false)
			handleError(err)
			app = a
		} else {
			a, err := services.GetApplicationByName(utils.FlagAppName)
			handleError(err)
			app = a
		}
		hasApp = true

		if utils.FlagProcId != "" {
			hasProc = true
		}
		if hasApp && hasProc {
			proc, err := services.GetProcessByApplicationAndProc(utils.FlagAppName, utils.FlagProcId)
			handleError(err)
			services.DoQuery(app.Name, args[0], proc.ProcID, proc.LastSeen)
			return
		} else if hasApp {
			services.DoQuery(app.Name, args[0], "", app.LastSeen)
		} else {
			fmt.Println(cmd.UsageString())
			return
		}
	},
}

func init() {
	logsCmd.PersistentFlags().StringVarP(&utils.FlagLogsSince, "since", "s", "1h", `Only return logs newer than a relative duration. This is in relative to the last
seen log time for a specified application or processes within the namespace.
A duration string is a possibly signed sequence of decimal numbers, each with optional
fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h"`)
	logsCmd.PersistentFlags().Uint32Var(&utils.FlagLogsPageSize, "page-size", 30, `Number of log entries to return in one page`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagLogsFollow, "follow", "f", false, `Specify if the logs should be streamed.`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagProcId, "process", "p", "", `Filter logs by  proc id`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagAppName,"application","a","",`Filter logs by application`)
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(interactiveCmd)
	logsCmd.AddCommand(searchCmd)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagFile, "write-to-file", "w", "", "Path to file")
	logsCmd.PersistentFlags().IntVarP(&utils.FlagMaxFileSize, "max-file-size", "m", 10, "Max output file size")
}
