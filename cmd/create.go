package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ApicaSystem/apicactl/api/v1/eventRules"
	"github.com/ApicaSystem/apicactl/converter"
	"github.com/ApicaSystem/apicactl/services"
	"github.com/ApicaSystem/apicactl/ui"
	"github.com/ApicaSystem/apicactl/utils"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <resource_name>",
	Short: "Create a Apica Ascent resource",
	Long:  `This command helps you create Apica Ascent resources such as dashboards and event rules from a resource specification.`,
	Example: `
Create a dashboard
apicactl create dashboard -f <path to dashboard_spec_file.json>

Create eventrules
apicactl create eventrules -f <path to eventrules_file.json>

Create alerts
apicactl create alert -f <path to alert.json>
`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(NewDashboardCreateCommand())
	createCmd.AddCommand(NewCreateEventRulesCommand())
	createCmd.AddCommand(CreateAlertCommand())
}

func NewDashboardCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "apicactl create dashboard|d -f <path to dashboard spec> [-s <source of the template>] [-d <dashboard name>]",
		Aliases: []string{"d"},
		Short:   "Create a dashboard",
		Long: `

The crowd-sourced dashboards available in https://github.com/logiqai/logiqhub and https://grafana.com/grafana/dashboards can be downloaded and applied to any clusters. 
One can also export dashboards created using "apicactl get dashboard" command and apply on different clusters.

Dashboards from Grafana are imported by passing "grafana" as a value to source (-s) flag.
`,
		PreRun: utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.FlagFile == "" {
				fmt.Println("Missing dashboard spec file")
				os.Exit(-1)
			} else {
				fmt.Println("Dashboard spec file :", utils.FlagFile)
				fileBytes, err := ioutil.ReadFile(utils.FlagFile)
				if err != nil {
					fmt.Println("Unable to read file ", utils.FlagFile)
					os.Exit(-1)
				}
				var dashboardSpec string
				if utils.FlagDashboardSource == "grafana" {
					dashboardSpec, err = converter.ConvertToApicaDashboard(string(fileBytes), utils.FlagDashboardSource,
						utils.FlagDashboardName, nil)
				} else {
					dashboardSpec, err = ui.CreateAndPublishDashboardSpec(string(fileBytes))
				}
				if err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println(dashboardSpec)
			}
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	cmd.Flags().StringVarP(&utils.FlagDashboardSource, "source", "s", "", "Source of the dashboard template")
	cmd.Flags().StringVarP(&utils.FlagDashboardName, "name", "d", "", "dasboard title")
	return cmd
}

func NewCreateEventRulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eventrules",
		Example: "apicactl create eventrules -f <path to event rules file>",
		Aliases: []string{"eventrule", "er"},
		Short:   "Create an event rule",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.FlagFile == "" {
				fmt.Println("Missing event rules file")
				return
			} else {
				fmt.Println("Event rules file :", utils.FlagFile)
				if fileBytes, err := ioutil.ReadFile(utils.FlagFile); err != nil {
					fmt.Println("Unable to read file", utils.FlagFile)
					return
				} else {
					var rules []eventRules.EventRule
					if err = json.Unmarshal(fileBytes, &rules); err != nil {
						fmt.Println("Unable to decode event rules from ", utils.FlagFile)
					} else {
						services.CreateEventRules(rules)
					}
				}
			}
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}

func CreateAlertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "alert",
		Example: "apicactl create alert -f <path to alert file>",
		Short:   "Create an alert for query",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.FlagFile == "" {
				fmt.Println("Missing alert file")
				os.Exit(-1)
			}
			fileBytes, err := ioutil.ReadFile(utils.FlagFile)
			if err != nil {
				fmt.Println("Unable to read json value from ", utils.FlagFile)
				os.Exit(-1)
			}
			message, err := ui.CreateAlert(string(fileBytes))

			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(message)
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}
