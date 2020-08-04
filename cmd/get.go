/*
Copyright © 2020 Christopher J. Maahs <cmaahs@gmail.com>

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
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// getCmd represents the show command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the Current Context from KUBECONFIG",
	Long: `EXAMPLE 1
	#> gke-alias get | jq .

	{"clusterAlias": "nonprod-gke-dev1", "clusterName": "gke_nonprod-gke_us-east1_nonprod-gke-dev1"}

EXAMPLE 2

	#> gke-alias get -r

	nonprod-gke-dev1	gke_nonprod-gke_us-east1_nonprod-gke-dev1

EXAMPLE 3

	#> gke-alias get -a

	{"clusterAlias": "nonprod-gke-dev1"}

EXAMPLE 4

	#> gke-alias get -n

	{"clusterName": "gke_nonprod-gke_us-east1_nonprod-gke-dev1"}

EXAMPLE 5

	#> gke-alias get -a -r

	nonprod-gke-dev1

EXAMPLE 6

	#> gke-alias get -n -r

	gke_nonprod-gke_us-east1_nonprod-gke-dev1
`,
	Run: func(cmd *cobra.Command, args []string) {

		onlyAlias, _ := cmd.Flags().GetBool("alias")
		onlyName, _ := cmd.Flags().GetBool("name")
		rawOutput, _ := cmd.Flags().GetBool("raw")
		err := getCurrentContext(onlyAlias, onlyName, rawOutput)
		if err != nil {
			logrus.WithError(err).Error("Error getting the current context")
		}

	},
}

func getCurrentContext(alias bool, name bool, raw bool) error {

	var configData []byte
	var configFile string
	var currentCluster string
	var currentContext string
	var err error

	if kubeConfig := os.Getenv("KUBECONFIG"); kubeConfig != "" {
		configFile = getFirstConfig(kubeConfig)
		configData, err = ioutil.ReadFile(configFile)
		if err != nil {
			logrus.Error("Could not read KUBECONFIG environment for file")
			return err
		}
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logrus.Error("Could not fetch HOME directory")
			return err
		}

		configFile = fmt.Sprintf("%s/%s", home, ".kube/config")
		if _, err := os.Stat(configFile); err != nil {
			if os.IsNotExist(err) {
				logrus.Error("Could not find default ~/.kube/config file")
				return err
			}
		}
		configData, err = ioutil.ReadFile(configFile)
		if err != nil {
			logrus.Error("Could not read default ~/.kube/config file")
			return err
		}
	}

	kcfg := KubernetesCluster{}
	err = yaml.Unmarshal(configData, &kcfg)
	if err != nil {
		logrus.Error("Could not unmarshal YAML data")
		return err
	}

	currentContext = kcfg.CurrentContext

	for idx, ctx := range kcfg.Contexts {
		if ctx.Name == currentContext {
			currentCluster = kcfg.Contexts[idx].Context.Cluster
			break
		}
	}

	if alias && name {
		if raw {
			fmt.Println(fmt.Sprintf("%s\t%s", currentContext, currentCluster))
		} else {
			fmt.Println(fmt.Sprintf("{\"clusterAlias\": \"%s\", \"clusterName\": \"%s\"}", currentContext, currentCluster))
		}
	} else {
		if alias {
			// alias only
			if raw {
				fmt.Println(fmt.Sprintf("%s", currentContext))
			} else {
				fmt.Println(fmt.Sprintf("{\"clusterAlias\": \"%s\"}", currentContext))
			}
		} else {
			// must be name only
			if name {
				if raw {
					fmt.Println(fmt.Sprintf("%s", currentCluster))
				} else {
					fmt.Println(fmt.Sprintf("{\"clusterName\": \"%s\"}", currentCluster))
				}
			} else {
				// none was selected, show all
				if raw {
					fmt.Println(fmt.Sprintf("%s\t%s", currentContext, currentCluster))
				} else {
					fmt.Println(fmt.Sprintf("{\"clusterAlias\": \"%s\", \"clusterName\": \"%s\"}", currentContext, currentCluster))
				}
			}
		}
	}

	return nil
}

func init() {

	getCmd.Flags().BoolP("alias", "a", false, "Show the alias value")
	getCmd.Flags().BoolP("name", "n", false, "Show the cluster full name")
	getCmd.Flags().BoolP("raw", "r", false, "Show values in RAW format")

	rootCmd.AddCommand(getCmd)

}
