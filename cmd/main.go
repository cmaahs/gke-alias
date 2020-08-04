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
	APIVersion string `json:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data"`
			Server                   string `json:"server"`
		} `json:"cluster"`
		Name string `json:"name"`
	} `json:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
		Name string `json:"name"`
	} `json:"contexts"`
	CurrentContext string `json:"current-context"`
	Kind           string `json:"kind"`
	Preferences    struct {
	} `json:"preferences"`
	Users []struct {
		Name string `json:"name"`
		User struct {
			AuthProvider struct {
				Config struct {
					CmdArgs   string `json:"cmd-args"`
					CmdPath   string `json:"cmd-path"`
					ExpiryKey string `json:"expiry-key"`
					TokenKey  string `json:"token-key"`
				} `json:"config"`
				Name string `json:"name"`
			} `json:"auth-provider"`
		} `json:"user"`
	} `json:"users"`
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
