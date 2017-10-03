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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/reports"
	"github.com/ion-channel/ionic/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println("Run the analysis from the . file")
		key := viper.GetString("key")
		api := viper.GetString("api")
		cli, err := ionic.New(key, api)
		if err != nil {
			log.Fatalf("Failed to create Ion Channel Client: %v", err.Error())
		}
		project := viper.GetString("project")
		team := viper.GetString("team")
		branch := getBranch()
		analysisStatus, err := cli.AnalyzeProject(project, team, branch)
		if err != nil {
			log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
		}
		id := analysisStatus.ID

		if viper.GetString("coverage") != "" {
			coverage, err := loadCoverage(viper.GetString("coverage"))
			if err != nil {
				log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
			}

			fmt.Println("Adding external coverage scan data")

			scan := scanner.ExternalScan{}
			scan.Coverage = coverage
			analysisStatus, err = cli.AddScanResult(id, team, project, "accepted", "coverage", scan)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
			}
		}

		if viper.GetString("vulnerabilities") != "" {
			scan, err := loadVulnerabilities(viper.GetString("vulnerabilities"))
			if err != nil {
				log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
			}

			fmt.Println("Adding external vulnerabilities scan data")

			// scan := scanner.ExternalScan{}
			// scan.Vulnerabilities = vulnerabilities
			analysisStatus, err = cli.AddScanResult(id, team, project, "accepted", "vulnerability", *scan)
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
		report, err := cli.GetAnalysisReport(id, team, project)
		if err != nil {
			log.Fatalf("Analysis Report request failed for %s (%s): %v", project, id, err.Error())
		}

		os.Exit(printReport(report))
	},
}

func getBranch() string {
	branch := os.Getenv("GIT_BRANCH")
	if branch != "" {
		fmt.Println("Using branch from environment variable", branch)
		return branch
	}

	branch = os.Getenv("TRAVIS_BRANCH")
	if branch != "" {
		fmt.Println("Using branch from travis-ci", branch)
		return branch
	}

	//TODO: get it from git directly?
	return ""
}

func printReport(report *reports.AnalysisReport) int {
	for _, scanSummary := range report.ScanSummaries {
		fmt.Print(scanSummary.Summary, "...Rule Type: ")
		fmt.Print(scanSummary.Type, "...")
		if scanSummary.Passed {
			fmt.Print("passed")
		} else {
			fmt.Print("not passed")
		}

		fmt.Println("...Risk: ", scanSummary.Risk)
	}

	if !report.Passed {
		fmt.Println("Analysis failed on a rule")
		return 1
	}

	fmt.Println("Analysis passed all rules")
	return 0
}

func loadVulnerabilities(path string) (*scanner.ExternalScan, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Println("Reading coverage value from", path)

		raw, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("Could not open vulnerabilities file %v", err.Error())
		}

		var scan = scanner.ExternalScan{}
		err = json.Unmarshal(raw, &scan)
		if err != nil {
			return nil, fmt.Errorf("Could not parse vulnerabilities file %v", err.Error())
		}

		fmt.Println("Found and loaded vulnerabilities file")
		return &scan, nil
	}
	return nil, fmt.Errorf("File does not exist %s", path)
}

func loadCoverage(path string) (*scanner.ExternalCoverage, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Println("Reading coverage value from", path)
		var value float64
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return nil, fmt.Errorf("Could not open coverage file %v", err.Error())
		}

		_, err = fmt.Fscanln(f, &value)
		if err != nil {
			return nil, fmt.Errorf("Could read coverage from coverage file %v", err.Error())
		}
		fmt.Println("Found coverage", value)
		return &scanner.ExternalCoverage{Value: value}, nil
	}
	return nil, fmt.Errorf("File does not exist %s", path)
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
