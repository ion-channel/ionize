package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	appName        = "IONCI"
	appDescription = "CI/CD logic wrapper around ion-connect"
)

var (
	buildTime  string
	appVersion string
)

type analysisStatus struct {
	AccountID string    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	ProjectID string    `json:"project_id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	app := cli.NewApp()
	app.Name = strings.ToLower(appName)
	app.Version = appVersion
	app.Usage = appDescription

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "account-id, account",
			Usage:  "Account ID for accessing Ion Channel",
			EnvVar: fmt.Sprintf("%v_%v", appName, "ACCOUNT_ID"),
		},
	}

	app.EnableBashCompletion = true
	app.Commands = append(app.Commands, analysisCommands()...)

	app.Run(os.Args)
}

func analysisCommands() []cli.Command {
	var projectID, buildNum string
	analysisCmd := []cli.Command{
		{
			Name:  "analysis",
			Usage: "perform an analysis of a project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "project-id, project",
					Usage:       "Project ID for the analysis",
					Destination: &projectID,
				},
				cli.StringFlag{
					Name:        "buildnum, build-number",
					Usage:       "Build number to inspect",
					Destination: &buildNum,
				},
			},
			Action: func(c *cli.Context) error {
				accountID := c.GlobalString("account_id")
				if projectID == "" {
					return fmt.Errorf("Missing Flag: project_id")
				}
				if buildNum == "" {
					return fmt.Errorf("Missing Flag: buildnumber")
				}

				analysis, err := StartProjectAnalysis(accountID, projectID, buildNum)
				if err != nil {
					return err
				}

				status := make(chan int, 1)
				go func() {
					for {
						<-time.Tick(1 * time.Second)
						analysis, err := GetAnalysisStatus(analysis)
						if err != nil {
							status <- 1
							break
						}

						if analysis.Status == "finished" {
							status <- 0
							break
						}
					}

					close(status)
				}()

				select {
				case s := <-status:
					if s == 0 {
						res, err := GetAnalysis(analysis)
						if err != nil {
							return err
						}

						fmt.Println(string(res))
					} else {
						return fmt.Errorf("Issue getting analysis")
					}
				case <-time.After(time.Second * 20):
					return fmt.Errorf("Timed out waiting for analysis")
				}

				return nil
			},
		},
	}

	return analysisCmd
}

func StartProjectAnalysis(accountID, projectID, buildNumber string) (*analysisStatus, error) {
	cmdName := "ion-connect"
	cmdArgs := []string{"scanner", "analyze-project", "--account-id", accountID, "--project-id", projectID, buildNumber}

	o, err := exec.Command(cmdName, cmdArgs...).Output()

	a := &analysisStatus{}
	err = json.Unmarshal(o, &a)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse analysis response: %v", err.Error())
	}

	return a, err
}

func GetAnalysisStatus(analysis *analysisStatus) (*analysisStatus, error) {
	cmdName := "ion-connect"
	cmdArgs := []string{"scanner", "get-analysis-status", "--account-id", analysis.AccountID, "--project-id", analysis.ProjectID, analysis.ID}

	o, err := exec.Command(cmdName, cmdArgs...).Output()
	a := &analysisStatus{}
	err = json.Unmarshal(o, &a)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse analysis response: %v", err.Error())
	}

	return a, err
}

func GetAnalysis(analysis *analysisStatus) ([]byte, error) {
	cmdName := "ion-connect"
	cmdArgs := []string{"analysis", "get-analysis", "--account-id", analysis.AccountID, "--project-id", analysis.ProjectID, analysis.ID}
	return exec.Command(cmdName, cmdArgs...).Output()
}
