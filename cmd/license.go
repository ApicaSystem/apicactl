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
limitations under the License..
*/

package cmd

import (
	"github.com/ApicaSystem/apicactl/services"
	"github.com/ApicaSystem/apicactl/utils"
	"github.com/spf13/cobra"
)

var licenseExample = `
Upload your Apica Ascent platform license
  % apicactl license set -f license.jws

View your Apica Ascent license information 
  % apicactl license get 
 
You can obtain a valid license by contacting Apica Ascent at support@apica.io.
This command lets you view your existing Apica Ascent license or apply a new one. 
`
var licenseLong = `
The Apica Ascent Observability platform comes preconfigured with a 30-day trial license. You can obtain a valid license by contacting Apica Ascent at support@apica.io.
This command lets you view your existing Apica Ascent license or apply a new one. 
`
var licenseCmd = &cobra.Command{
	Use:     "license",
	Example: licenseExample,
	Aliases: []string{"licence"},
	Short:   "View or update Apica Ascent license",
	Long:    licenseLong,
}

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(NewGetLicenseCommand())
	licenseCmd.AddCommand(NewSetLicenseCommand())
}

func NewGetLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Example: "apicactl license get",
		Aliases: []string{},
		Short:   "View license information",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetLicense()
		},
	}
	return cmd
}

func NewSetLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Example: "apicactl license set -f <license-file-path>",
		Aliases: []string{},
		Short:   "Configure license for Apica Ascent",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.SetLicense()
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}
