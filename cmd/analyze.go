package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/reports"
	"github.com/ion-channel/ionize/cmd/external"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	async = false
)

func init() {
	RootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().BoolVarP(&async, "async", "a", false, "run the command asynchronously without waiting for completion")
}

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
			analysisStatus, err = coverage.Save(aID, cli)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
			}
		}

		if viper.IsSet("vulnerabilities") {
			files := viper.GetStringSlice("vulnerabilities")

			for _, file := range files {
				vulns, err := external.ParseVulnerabilities(file)
				if err != nil {
					log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
				}
				analysisStatus, err = vulns.Save(aID, cli)
				if err != nil {
					log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
				}
			}
		}

		if viper.IsSet("fortify") {
			fortify, err := external.ParseFortify(viper.GetString("fortify"))
			if err != nil {
				log.Fatalf("Analysis request failed for %s: %v", project, err.Error())
			}
			analysisStatus, err = fortify.Save(aID, cli)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s: %v", project, err.Error())
			}
		}

		if !async {
			fmt.Print("Waiting for analysis to finish")
			for !analysisStatus.Done() {
				fmt.Print(".")
				time.Sleep(10 * time.Second)
				analysisStatus, err = cli.GetAnalysisStatus(id, team, project, key)
				if err != nil {
					log.Fatalf("Analysis Status request failed for %s: %v", project, err.Error())
				}
			}
			fmt.Printf("%s\n", analysisStatus.Status)
			if analysisStatus.Status == "errored" {
				log.Fatalf("Analysis error occured: %v", analysisStatus.Message)
			}

			fmt.Println("Checking status of scans")
			report, err := cli.GetAnalysisReport(id, team, project, key)
			if err != nil {
				log.Fatalf("Analysis Report request failed for %s (%s): %v", project, id, err.Error())
			}

			os.Exit(printReport(report))
		}
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
	for _, scanSummary := range report.Report.RulesetEvaluation.RuleEvaluationSummary.Ruleresults {
		fmt.Print(scanSummary.Summary, "...Rule Type: ")
		fmt.Print(scanSummary.Type, "...")
		if scanSummary.Passed {
			fmt.Print("passed")
		} else {
			fmt.Print("not passed")
		}

		fmt.Println("...Risk: ", scanSummary.Risk)
	}

	if !report.Report.RulesetEvaluation.RuleEvaluationSummary.Passed {
		fmt.Println("Analysis failed on a rule")
		return 1
	}

	fmt.Println("Analysis passed all rules")
	return 0
}
