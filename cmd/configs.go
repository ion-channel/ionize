package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(configsCmd)
}

var configsCmd = &cobra.Command{
	Use:   "configs",
	Short: "Print the configs in use.",
	Long:  `Print out the configs and their values that have been loaded into ionize.`,
	Run:   runConfigsCmd,
}

func runConfigsCmd(cmd *cobra.Command, args []string) {
	fmt.Fprintf(output, "Config File: %v\n", viper.ConfigFileUsed())
	fmt.Fprintf(output, "Secret Key: %v\n", viper.GetString("key"))
	fmt.Fprintf(output, "API: %v\n", viper.GetString("api"))
}
