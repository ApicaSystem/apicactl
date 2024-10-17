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
	"path"

	"github.com/ApicaSystem/apicactl/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Update this before publishing the release!!!
var currentReleaseVersion = "3.0.0"

var rootCmd = &cobra.Command{
	Short:   "apicactl - CLI for Apica Ascent",
	Use:     "apicactl [flags] [options]",
	Version: currentReleaseVersion,
	Long: `
Apica Ascent comes with an inbuilt command-line toolkit that lets you interact with the Apica Ascent Observability platform without logging into the UI. Using apicactl, you can:
- Stream logs in real-time
- Query historical application logs
- Search within logs across namespaces
- Query and view events across your Apica Ascent stack
- View and create event rules
- Create and manage dashboards
- Query and view all your resources on Apica Ascent such as applications, dashboards, namespaces, processes, and queries
- Manage Apica Ascent licenses
- Log pattern signature extraction and reporting (max 50,000 log-lines)
- Create alerts by passing the JSON file which has alert specifications.
- Create alerts when creating the dashboard using create dashboard command
- View alerts
- Grafana dashboard import - Import dashboards from grafana public repository

Find more information docs.apica.io, please contact support@apica.io.
`,
}

func Execute() {
	//doc.GenMarkdownTree(rootCmd, "./docs")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&utils.FlagOut, "output", "o", "table", `Output format. One of: table|json|yaml. 
JSON output is not indented, use '| jq' for advanced JSON operations`)
	rootCmd.PersistentFlags().StringVarP(&utils.FlagTimeFormat, "time-format", "t", "relative", `Time formatting options. One of: relative|epoch|RFC3339. 
This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds.`)
	rootCmd.PersistentFlags().StringVarP(&utils.FlagNamespace, "namespace", "n", "", "Override the default context set by `apicactl set-context' command")
	rootCmd.PersistentFlags().StringVarP(&utils.FlagCluster, "cluster", "c", "", "Override the default cluster set by `apicactl set-cluster' command")
	rootCmd.PersistentFlags().BoolVarP(&utils.FlagNetTrace, "nettrace", "", false, "Network trace enable")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.DisableAutoGenTag = true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir := path.Join(home, ".apicactl")
	exists, err := exists(configDir)
	cfgFile := path.Join(configDir, "apicactl.toml")
	if err != nil {
		fmt.Println(fmt.Errorf("Cannot create config: %s ", err.Error()))
		return
	}
	if !exists {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("Cannot create config: %s ", err.Error()))
			return
		}
		viper.SetConfigFile(cfgFile)
		viper.Set("apicactl", currentReleaseVersion)
		viper.Set(utils.LineBreaksKey, false)
		viper.WriteConfig()
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot Use Config file:", viper.ConfigFileUsed())
	}
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
