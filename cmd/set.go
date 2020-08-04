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

// setCmd represents the show command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Get a list of repositories for your Organization",
	Long: `EXAMPLE 1
	#> gke-alias set -a nonprod-gke-dev1

EXAMPLE 2

	#> gke-alias set --alias nonprod-gke-dev1

`,
	Run: func(cmd *cobra.Command, args []string) {
		alias, _ := cmd.Flags().GetString("alias")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if alias == "lias" {
			logrus.Fatal("When using the flag 'alias', plese use two dashes '--alias', otherwise use the shortcut '-a'")
		}
		err := setAlias(alias)
		if err != nil {
			logrus.WithError(err).Error("Error setting the alias for the current-context")
		}
		if verbose {
			logrus.Info(fmt.Sprintf("Successfully set current context to alias %s", alias))
		}

	},
}

func setAlias(newAlias string) error {

	var configData []byte
	var configFile string
	var currentCluster string
	var currentContext string
	var err error

	if kubeConfig := os.Getenv("KUBECONFIG"); kubeConfig != "" {
		configFile = getFirstConfig(kubeConfig)
		configData, err = ioutil.ReadFile(configFile)
		if err != nil {
			logrus.Error("Could not read KUBECONFIG environment for config file")
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
			kcfg.Contexts[idx].Name = newAlias
			currentCluster = kcfg.Contexts[idx].Context.Cluster
			break
		}
	}
	kcfg.CurrentContext = newAlias

	yamlCfg, yerr := yaml.Marshal(&kcfg)
	if yerr != nil {
		logrus.Error("Error marshalling YAML data")
		return yerr
	}

	werr := ioutil.WriteFile(configFile, yamlCfg, 0644)
	if werr != nil {
		logrus.Error("Error writing to the config file", configFile)
		return werr
	}

	fmt.Println(fmt.Sprintf("{\"clusterAlias\": \"%s\", \"clusterName\": \"%s\"}", newAlias, currentCluster))

	return nil
}

func init() {

	getCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	setCmd.Flags().StringP("alias", "a", "", "Set the Alias for the current-context")
	setCmd.MarkFlagRequired("alias")
	rootCmd.AddCommand(setCmd)

}
