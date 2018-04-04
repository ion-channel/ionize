package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/pagination"
	"github.com/ion-channel/ionic/projects"
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

		url, err := parseURL(args[0])
		if err != nil {
			log.Fatalf("Failed to parse url: %v\n", err.Error())
		}

		name := args[1]
		version := args[2]

		rulesets, err := cli.GetRuleSets(team, key, nil)
		if err != nil || rulesets == nil || len(rulesets) == 0 {
			log.Fatalf("Failed to retrieve rulesets for team, make sure a valid ruleset exists\n")
		}

		project := &projects.Project{
			Name:      name,
			Branch:    version,
			Source:    url,
			Type:      "artifact",
			POCEmail:  "",
			POCName:   "",
			TeamID:    team,
			Active:    true,
			RulesetID: rulesets[0].ID,
		}
		project, err = cli.CreateProject(project, team, key)
		if err != nil {
			projects, er := cli.GetProjects(team, key, pagination.AllItems)
			if er != nil {
				log.Fatalf("Failed to receive projects: %v", err.Error())
			}
			for i, p := range projects {
				if p.Source == url && p.Branch == version {
					project = &projects[i]
					break
				}
			}
			if project == nil {
				log.Fatalf("Failed to create project: %v", err.Error())
			}
		} else {
			_, err = cli.AddAlias(project.ID, team, name, version, key)
			if err != nil {
				log.Fatalf("Failed to add alias to project, analysis depth will be reduced: %v", err.Error())
			}
			fmt.Printf("Created alias %s for %s (%s) %v", name, project.ID, project.TeamID, project.Aliases)
		}
		fmt.Printf("Created project: %s (%s) for %s\n", project.ID, project.TeamID, url)

		analysisStatus, err := cli.AnalyzeProject(project.ID, team, version, key)
		if err != nil {
			log.Fatalf("Analysis request failed for %v: %v", project.ID, err.Error())
		}
		id := analysisStatus.ID

		fmt.Printf("Waiting for analysis (%s) to finish", id)
		for analysisStatus.Status == "accepted" {
			fmt.Print(".")
			time.Sleep(10 * time.Second)
			analysisStatus, err = cli.GetAnalysisStatus(id, team, project.ID, key)
			if err != nil {
				log.Fatalf("Analysis Status request failed for %s: %v", project.Name, err.Error())
			}
		}
		fmt.Printf("%s\n", analysisStatus.Status)

		fmt.Println("Checking status of scans")
		report, err := cli.GetAnalysisReport(id, team, project.ID, key)
		if err != nil {
			log.Fatalf("Analysis Report request failed for %s (%s): %v", project.Name, id, err.Error())
		}

		os.Exit(printReport(report))
	},
}

func parseURL(input string) (string, error) {
	url, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %v", err.Error())
	}

	if url.Scheme == "file" {
		// get reader
		f, err := os.Open(url.Hostname() + url.EscapedPath())
		if err != nil {
			return "", fmt.Errorf("failed to read file for url (%s): %v", url.String(), err.Error())
		}
		// upload to dropbox.ionchannel.io
		mySession, _ := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		uploader := s3manager.NewUploader(mySession)
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("dropbox.ionchannel.io"),
			Key:    aws.String("ionize/" + f.Name()),
			Body:   f,
		})

		if err != nil {
			return "", fmt.Errorf("failed to write file to Ion Channel: %v", err.Error())
		}
		// get url to upload
		return result.Location, nil
		// return s3 upload
	}

	return url.String(), nil
}
