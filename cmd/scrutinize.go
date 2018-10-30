package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/pagination"
	"github.com/ion-channel/ionic/projects"
	"github.com/ion-channel/ionize/dropbox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ScrutinizeCmd represents the doAnalysis command
var scrutinizeCmd = &cobra.Command{
	Use:   "scrutinize url name version",
	Short: "Perform an analysis on a url and wait for report",
	Long: `Perform an analysis on a url and wait for report. For example:

ionize scrutinize url name version

Will read the configuration from the $PWD/.ionize.yaml file and begin an analysis.
`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run the analysis from the . file")
		key := viper.GetString("key")
		api := viper.GetString("api")
		team := viper.GetString("team")
		cli, err := ionic.New(api)
		if err != nil {
			log.Fatalf("Failed to create Ion Channel Client: %v", err.Error())
		}

		rando, err := dropbox.Randomizer()
		if err != nil {
			log.Fatalf("Failed to parse url: %v\n", err.Error())
		}

		url, err := dropbox.ParseURL(args[0], rando)
		if err != nil {
			log.Fatalf("Failed to parse url: %v\n", err.Error())
		}

		name := args[1]
		version := args[2]
		ty := "artifact"

		rulesets, err := cli.GetRuleSets(team, key, nil)
		if err != nil || rulesets == nil || len(rulesets) == 0 {
			log.Fatalf("Failed to retrieve rulesets for team, make sure a valid ruleset exists\n")
		}

		project := &projects.Project{
			Name:      &name,
			Branch:    &version,
			Source:    &url,
			Type:      &ty,
			POCEmail:  "",
			POCName:   "",
			TeamID:    &team,
			Active:    true,
			RulesetID: &rulesets[0].ID,
		}
		project, err = cli.CreateProject(project, team, key)
		if err != nil {
			projects, er := cli.GetProjects(team, key, pagination.AllItems)
			if er != nil {
				log.Fatalf("Failed to receive projects: %v", err.Error())
			}
			for i, p := range projects {
				if *p.Source == url && *p.Branch == version {
					project = &projects[i]
					break
				}
			}
			if project == nil {
				log.Fatalf("Failed to create project: %v", err.Error())
			}
		} else {
			_, err = cli.AddAlias(*project.ID, team, name, version, key)
			if err != nil {
				log.Fatalf("Failed to add alias to project, analysis depth will be reduced: %v", err.Error())
			}
			fmt.Printf("Created alias %s for %v (%v) %v", name, project.ID, project.TeamID, project.Aliases)
		}
		fmt.Printf("Created project: %v (%v) for %s\n", project.ID, project.TeamID, url)

		analysisStatus, err := cli.AnalyzeProject(*project.ID, team, version, key)
		if err != nil {
			log.Fatalf("Analysis request failed for %v: %v", project.ID, err.Error())
		}
		id := analysisStatus.ID

		fmt.Printf("Waiting for analysis (%s) to finish", id)
		for analysisStatus.Status == "accepted" {
			fmt.Print(".")
			time.Sleep(10 * time.Second)
			analysisStatus, err = cli.GetAnalysisStatus(id, team, *project.ID, key)
			if err != nil {
				log.Fatalf("Analysis Status request failed for %v: %v", project.Name, err.Error())
			}
		}
		fmt.Printf("%s\n", analysisStatus.Status)

		fmt.Println("Checking status of scans")
		report, err := cli.GetAnalysisReport(id, team, *project.ID, key)
		if err != nil {
			log.Fatalf("Analysis Report request failed for %v (%s): %v", project.Name, id, err.Error())
		}

		os.Exit(printReport(report))
	},
}
