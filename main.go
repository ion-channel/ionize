package main

import (
	"fmt"
	"os"

	"github.com/ion-channel/ionize/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
