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
	"github.com/ion-channel/ionize/cmd/external"
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
		cli, err := ionic.New(api)
		if err != nil {
			log.Fatalf("Failed to create Ion Channel Client: %v", err.Error())
		}
		project := viper.GetString("project")
		team := viper.GetString("team")
		branch := getBranch()
		analysisStatus, err := cli.AnalyzeProject(project, team, branch, key)
		if err != nil {
			log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
		}
		id := analysisStatus.ID
		aID := &external.AnalysisID{
			ID:        id,
			TeamID:    team,
			ProjectID: project,
			APIKey:    key,
		}

		if viper.IsSet("coverage") {
			coverage, err := external.ParseCoverage(viper.GetString("coverage"))
			if err != nil {
				log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
			}

			fmt.Println("Adding external coverage scan data")

			analysisStatus, err = coverage.Save(aID, cli)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
			}
		}

		if viper.IsSet("vulnerabilities") {
			files := make([]string, 1)
			if viper.GetString("vulnerabilities") != "" {
				files[0] = viper.GetString("vulnerabilities")
			} else {
				files = viper.GetStringSlice("vulnerabilities")
			}
			for _, file := range files {
				scan, err := loadVulnerabilities(file)
				if err != nil {
					log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
				}

				fmt.Println("Adding external vulnerabilities scan data")
				analysisStatus, err = cli.AddScanResult(id, team, project, "accepted", "vulnerability", key, *scan)
				if err != nil {
					log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
				}
			}
		}

		fmt.Print("Waiting for analysis to finish")
		for analysisStatus.Status == "accepted" {
			fmt.Print(".")
			time.Sleep(10 * time.Second)
			analysisStatus, err = cli.GetAnalysisStatus(id, team, project, key)
			if err != nil {
				log.Fatalf("Analysis Status request failed for %s: %v", project, err.Error())
			}
		}
		fmt.Printf("%s\n", analysisStatus.Status)

		fmt.Println("Checking status of scans")
		report, err := cli.GetAnalysisReport(id, team, project, key)
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
