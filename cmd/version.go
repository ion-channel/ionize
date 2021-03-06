package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var (
	// Version of Ionize.
	Version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Ionize",
	Long:  `All software has versions. This is mine`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Ionize %v\n", Version)
	},
}
