package external

import (
	"fmt"
	"os"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/scanner"
	"github.com/spf13/viper"
)

func ParseCoverage(contents string) (*Coverage, error) {
	coverage, err := loadCoverage(viper.GetString("coverage"))
	if err != nil {
		return nil, fmt.Errorf("Analysis request failed: %v", err.Error())
	}

	return &Coverage{
		Value: coverage,
	}, nil
}

type Coverage struct {
	Value *scanner.ExternalCoverage
}

func (c *Coverage) Save(aID *AnalysisID, cli *ionic.IonClient) (*scanner.AnalysisStatus, error) {
	fmt.Println("Adding external coverage scan data")

	scan := scanner.ExternalScan{}
	scan.Coverage = c.Value
	analysisStatus, err := cli.AddScanResult(aID.ID, aID.TeamID, aID.ProjectID, "accepted", "coverage", aID.APIKey, scan)
	if err != nil {
		return nil, fmt.Errorf("Analysis coverage save failed: %v", err.Error())
	}
	return analysisStatus, nil
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
