/*
Copyright © 2020 Splice Machine
Author: Christopher Maahs <cmaahs@splicemachine.com>

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

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Express the 'version' of gke-alias.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("{\"VersionInfo\": {")

		fmt.Println(fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"}}", semVer, gitCommit, buildDate))

		fmt.Println("}")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}