package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	semVer    string
	gitCommit string
	buildDate string
)

// KubernetesCluster - Holds the kubectl config
type KubernetesCluster struct {
	APIVersion     string           `json:"apiVersion,omitempty"`
	Clusters       []ClustersObject `json:"clusters,omitempty"`
	Contexts       []ContextsObject `json:"contexts,omitempty"`
	CurrentContext string           `json:"current-context"`
	Kind           string           `json:"kind"`
	// Preferences    struct {
	// } `json:"preferences"`
	Users []UsersObject `json:"users,omitempty"`
}

type ClustersObject struct {
	Cluster ClusterObject `json:"cluster,omitempty"`
	Name    string        `json:"name"`
}

type ClusterObject struct {
	CertificateAuthorityData string `json:"certificate-authority-data"`
	Server                   string `json:"server"`
}

type ContextsObject struct {
	Context ContextObject `json:"context,omitempty"`
	Name    string        `json:"name"`
}

type ContextObject struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

type UsersObject struct {
	Name string     `json:"name"`
	User UserObject `json:"user,omitempty"`
}

type UserObject struct {
	AuthProvider          *AuthProviderObject `json:"auth-provider,omitempty"`
	ClientCertificateData string              `json:"client-certificate-data,omitempty"`
	ClientKeyData         string              `json:"client-key-data,omitempty"`
	Exec                  *ExecObject         `json:"exec,omitempty"`
	Token                 string              `json:"token,omitempty"`
}

type AuthProviderObject struct {
	Config *ConfigObject `json:"config,omitempty"`
	Name   string        `json:"name,omitempty"`
}

type ConfigObject struct {
	AccessToken string `json:"access-token,omitempty"`
	CmdArgs     string `json:"cmd-args,omitempty"`
	CmdPath     string `json:"cmd-path,omitempty"`
	Expiry      string `json:"expiry,omitempty"`
	ExpiryKey   string `json:"expiry-key,omitempty"`
	TokenKey    string `json:"token-key,omitempty"`
}

type ExecObject struct {
	APIVersion string      `json:"apiVersion,omitempty"`
	Args       []string    `json:"args,omitempty"`
	Command    string      `json:"command,omitempty"`
	Env        []EnvObject `json:"env,omitempty"`
}

type EnvObject struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gke-alias",
	Args:  cobra.MinimumNArgs(1),
	Short: `Set a short cluster name for GCP/GKE Clusters`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func getFirstConfig(kubeConfig string) string {

	if strings.Contains(kubeConfig, ":") {
		return strings.Split(kubeConfig, ":")[0]
	}
	return kubeConfig
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gke-alias/config.yml)")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		if _, err := os.Stat(cfgFile); err != nil {
			if os.IsNotExist(err) {
				createRestrictedConfigFile(cfgFile)
			}
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		directory := fmt.Sprintf("%s/%s", home, ".gke-alias")
		if _, err := os.Stat(directory); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(directory, os.ModePerm)
			}
		}
		if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", home, ".gke-alias/config.yml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			logrus.Info("The ~/.gke-alias path is a file and not a directory, please remove the .gke-alias file.")
			os.Exit(1)
		}
	}

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// couldn't read the config file.
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				logrus.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			if runtime.GOOS != "windows" {
				mode := int(0600)
				if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
					logrus.Info("Chmod for config file failed, please set the mode to 0600.")
				}
			}
		}
	}
}
