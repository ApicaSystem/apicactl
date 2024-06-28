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

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ApicaSystem/apicactl/grpc_utils"
	"github.com/ApicaSystem/apicactl/services"
	"github.com/ApicaSystem/apicactl/types"
	"github.com/ApicaSystem/apicactl/ui"
	"github.com/ApicaSystem/apicactl/utils"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <resource_name>",
	Short: "Display one or more of your Apica Ascent resources",
	Long:  `Prints a table that displays the most important information about the Apica Ascent resource you specify.`,
	Example: `
List all applications for the selected context
apicactl get applications

List all applications for all available contexts
apicactl get applications all

List all dashboards
apicactl get dashboards all

Get dashboard
apicactl get dashboard dashboard-slug

List events for the available namespace
apicactl get events

List events for all available namespaces
apicactl get events all

List or export eventrules
apicactl get eventrules

List all namespaces
apicactl get namespaces

List all processes
apicactl get processes

List all queries
apicactl get queries all

List all forwarder mappings
apicactl get mappings
		
List all forwarders
apicactl get forwards, apicactl get forwards all

List config of a specific forwarder
apicactl get forwards <id,name> --output <json,yaml>

List all alerts
apicactl get alert all

Get httpingestkey
apicactl get httpingestkey

Get query
apicactl get query query-slug

Get alert
apicactl get alert alert-slug
`,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(NewListNameSpaceCommand())
	getCmd.AddCommand(NewListApplicationsCommand())
	getCmd.AddCommand(NewListProcessesCommand())
	getCmd.AddCommand(NewListEventsCommand())
	getCmd.AddCommand(NewGetEventRulesCommand())
	getCmd.AddCommand(NewListDashboardsCommand())
	getCmd.AddCommand(ui.NewListQueriesCommand())
	getCmd.AddCommand(ui.NewListDatasourcesCommand())
	getCmd.AddCommand(getForwardsCommand())
	getCmd.AddCommand(getMappersCommand())
	getCmd.AddCommand(getHttpingestkeyCommand())
	getCmd.AddCommand(ui.NewGetLogEvents())
	getCmd.AddCommand(getAlertsCommand())
}

func getMappersCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "mappers",
		Example: "apicactl get mappers",
		Aliases: []string{"mappers"},
		Short:   "Get logflow log mappers",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			mappings, err := ui.GetMappings()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			forwarders, err := ui.GetForwarders()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			result := make([]types.Resource, 0)
			for _, m := range mappings {
				for _, f := range forwarders {
					if m.ID == f.Id {
						m.Name = f.Name
						break
					}
				}

				result = append(result, m)
			}

			utils.PrintResult(result, true)
		},
	}
}

func getForwardsCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "forwards",
		Example: "apicactl get forwards id",
		Aliases: []string{"forwards"},
		Short:   "Get logflow forwards by ID",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			forwards, err := ui.GetForwarders()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Just checking if there's and ID or Name specified
			if len(os.Args) < 4 || os.Args[3] == "all" {
				res := make([]types.Resource, 0)
				for _, f := range forwards {
					res = append(res, f)
				}

				utils.PrintResult(res, true)
				return
			}

			slug := os.Args[3]

			for _, f := range forwards {
				if strconv.Itoa(f.Id) == slug || f.Name == slug {
					if utils.FlagOut == "table" {
						fmt.Println("Please specify output as json or yaml")
						os.Exit(1)
					}

					utils.PrintResponse(f)
					return
				}
			}

			fmt.Println("Invalid Forwarder Id/Name")
		},
	}
}

func getHttpingestkeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "httpingestkey",
		Example: "apicactl get httpingestkey",
		Aliases: []string{"ingestkey"},
		Short:   "Get httpingestkey",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if u, cookieJar, err := grpc_utils.GetCookies(); err != nil {
				fmt.Println("Error getting httpingestkey: ", err.Error())
			} else {
				tokFound := false
				for _, c := range cookieJar.Cookies(u) {
					if "x-api-key" == c.Name {
						tokFound = true
						fmt.Println(c.Value)
					}
				}
				if !tokFound {
					fmt.Println("Error getting the httpingestkey")
				}
			}
		},
	}
	return cmd
}

func NewListNameSpaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespaces",
		Example: "apicactl get namespaces|ns|n",
		Aliases: []string{"n", "ns"},
		Short:   "List all available namespaces",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.ListNamespaces()
		},
	}

	return cmd
}

func NewListApplicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "applications",
		Example: "apicactl get applications|apps|a",
		Aliases: []string{"a", "apps"},
		Short:   "List all available applications within the default namespace",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetApplicationsV2(false)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "apicactl get applications all",
		Short:   "List all available applications across namespaces",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetApplicationsV2(true)
		},
	})
	return cmd
}

func NewListEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "events",
		Example: `
List last 30 events
- apicactl get events|e

List events by application 
- apicactl get events -a=sshd

`,
		Aliases: []string{"ev"},
		Short:   "List all available events for the namespace",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEvents(application, process)
		},
	}
	cmd.Flags().StringVarP(&application, "application", "a", "", `Filter events by application name`)
	cmd.Flags().StringVarP(&process, "process", "p", "", `Filter events by application and process name`)
	return cmd
}

func NewListProcessesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "processes",
		Example: "apicactl get processes|proc|p",
		Aliases: []string{"p", "proc"},
		Short:   "List all available processes. This command runs an interactive prompt that lets you choose an application from a list of available applications.",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.ListProcesses()
		},
	}
}

func NewGetEventRulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eventrules",
		Example: "apicactl get eventrule",
		Aliases: []string{"eventrule", "er"},
		Short:   "List event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			help := `Usage:
apicactl get eventrules all
apicactl get eventrules all -w filename.json
apicactl get eventrules groups
apicactl get eventrules groups -g=group1,group2,...
apicactl get eventrules groups -g=group1,group2,... -w filename.json
`
			fmt.Print(help)
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "apicactl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEventRules(args, nil)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "groups",
		Example: "apicactl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEventRuleGroups(args)
		},
	})
	cmd.Flags().StringVarP(&utils.EventRuleGroupsFlag, "groups", "g", "", "list of groups separated by comma")
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "w", "", "Path to file to be written to")
	return cmd
}

func getAlertsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alert",
		Example: "apicactl get alert <alert-id>",
		Short:   "Get available alert",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing alert id")
				os.Exit(-1)
			}
			var alert types.Resource
			var err error
			alert, err = ui.GetAlert(args[0])

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(-1)
			}
			utils.PrintResult([]types.Resource{alert}, false)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "apicactl get alert all",
		Short:   "List all available alerts",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			var alertList []types.Resource
			var err error
			alertList, err = ui.ListAlerts()

			if err != nil {
				fmt.Println(err.Error())
			}
			utils.PrintResult(alertList, true)
		},
	})

	return cmd
}

func NewListDashboardsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "apicactl get dashboard|d <dashboard-slug>",
		Aliases: []string{"d"},
		Short:   "Get a dashboard",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing dashboard slug")
				os.Exit(-1)
			}
			response, err := ui.GetDashboard(args)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(-1)
			}
			fmt.Println(response)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "apicactl get dashboard all",
		Short:   "List all the available dashboards",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			ui.ListDashboards()
		},
	})

	return cmd
}
