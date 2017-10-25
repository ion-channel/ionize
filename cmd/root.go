package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	KEY = "IONCHANNEL_SECRET_KEY"
	API = "IONCHANNEL_ENDPOINT_URL"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ionize",
	Short: "Provides a macro level interface for performing supply chain analysis of a project",
	Long: `ionize is a CLI tool that allows for rich interaction with the Ion Channel API to
perform supply chain analysis for a project.
`,
}

func init() {
	cobra.OnInitialize(initDefaults, initEnvs, initConfig)

	// Here you will define your flags and configuration settings.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.ionize.yaml)")
}

func initDefaults() {
	viper.SetDefault("api", "https://api.ionchannel.io")
}

func initEnvs() {
	viper.BindEnv("key", KEY)
	viper.BindEnv("api", API)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in pwd directory with name ".ionize" (without extension).
		viper.AddConfigPath(dir)
		// viper.AddConfigPath(".")
		viper.SetConfigName(".ionize")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
