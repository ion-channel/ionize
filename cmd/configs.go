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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Config File: %v\n", viper.ConfigFileUsed())
		fmt.Printf("Secret Key: %v\n", viper.GetString("key"))
		fmt.Printf("API: %v\n", viper.GetString("api"))
	},
}
