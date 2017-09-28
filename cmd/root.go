// Copyright © 2017 Ion Channel dev@ionchannel.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.ionize.yaml)")

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

	viper.BindEnv("key", KEY)
	viper.BindEnv("api", API)

	viper.SetDefault("api", "https://api.ionchannel.io")
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
