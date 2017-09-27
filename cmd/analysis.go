// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

// AnalyzeCmd represents the doAnalysis command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform an analysis and wait for report",
	Long: `Perform an analysis and wait for report. For example:

ionize analyze

Will read the configuration from the $PWD/.ionize.yaml file and begin an analysis.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Run the analysis from the . file")
		key := viper.GetString("key")
		api := viper.GetString("api")
		cli, err := ionic.New(key, api)
		if err != nil {
			log.Fatalf("Failed to create Ion Channel Client: %v", err.Error())
		}
		project := viper.GetString("project")
		team := viper.GetString("team")
		analysisStatus, err := cli.AnalyzeProject(team, project)
		if err != nil {
			log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
		}
		id := analysisStatus.ID

		value := viper.GetFloat64("coverage")
		if value != 0.0 {
			log.Println("Adding external coverage scan data")
			coverage := scanner.ExternalCoverage{value}
			scan := scanner.ExternalScan{}
			scan.Coverage = &coverage
			analysisStatus, err = cli.AddScanResult(id, team, project, "accepted", "coverage", scan)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
			}
		}

		fmt.Print("Waiting for analysis to finish")
		for analysisStatus.Status == "accepted" {
			fmt.Print(".")
			time.Sleep(10 * time.Second)
			analysisStatus, err = cli.GetAnalysisStatus(id, team, project)
			if err != nil {
				log.Fatalf("Analysis Status request failed for %s: %v", project, err.Error())
			}
		}
		fmt.Printf("%s\n", analysisStatus.Status)

		fmt.Println("Checking status of scans")
		report, err := cli.GetReport(id, team, project)
		if err != nil {
			log.Fatalf("Analysis Report request failed for %s (%s): %v", project, id, err.Error())
		}

		if !report.Passed {
			fmt.Println("Analysis failed on a rule")
		} else {
			fmt.Println("Analysis passed all rules")
		}

		for _, scanSummary := range report.ScanSummaries {
			scanData := make(map[string]interface{})

			err := json.Unmarshal(scanSummary, &scanData)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s (%s): %v", project, id, err.Error())
			}

			fmt.Print(scanData["summary"], "...Rule Type: ")
			fmt.Print(scanData["type"], "...")
			if scanData["passed"].(bool) {
				fmt.Print("passed")
			} else {
				fmt.Print("not passed")
			}

			fmt.Println("...Risk: ", scanData["risk"])
		}

	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doAnalysisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doAnalysisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
