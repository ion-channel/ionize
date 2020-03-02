package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	key    = "IONCHANNEL_SECRET_KEY"
	api    = "IONCHANNEL_ENDPOINT_URL"
	bucket = "IONCHANNEL_DROP_BUCKET"
)

var (
	output  io.Writer
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ionize",
	Short: "Provides a macro level interface for performing supply chain analysis of a project",
	Long: `ionize is a CLI tool that allows for rich interaction with the Ion Channel API to
perform supply chain analysis for a project.
`,
}

func init() {
	output = os.Stdout

	cobra.OnInitialize(initDefaults, initEnvs, initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.ionize.yaml)")
}

func initDefaults() {
	viper.SetDefault("api", "https://api.ionchannel.io")
	viper.SetDefault("bucket", "dropbox.ionchannel.io")
}

func initEnvs() {
	viper.BindEnv("key", key)
	viper.BindEnv("api", api)
	viper.BindEnv("bucket", bucket)
}

func initConfig() {
	viper.SetConfigType("yaml")

	viper.SetConfigFile(".ionize.yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed reading config: %v\n", err.Error())
	}
}

func init() {
	RootCmd.AddCommand(scrutinizeCmd)
}
