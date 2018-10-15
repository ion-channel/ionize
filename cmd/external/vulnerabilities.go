package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/scanner"
)

//ParseVulnerabilities - given a path to a file containing Ion channel
//formatted data will parse the file and return a struct representation
func ParseVulnerabilities(path string) (*Vulnerabilities, error) {
	eScan, err := loadVulnerabilities(path)
	if err != nil {
		return nil, fmt.Errorf("Analysis request failed: %v", err.Error())
	}

	return &Vulnerabilities{
		Value: eScan,
	}, nil
}

//Vulnerabilities struct representation of external vulnerability scan data
type Vulnerabilities struct {
	Value *scanner.ExternalScan
}

//Save sends the external vulnerability scan data to ion channel for persistance
func (c *Vulnerabilities) Save(aID *AnalysisID, cli *ionic.IonClient) (*scanner.AnalysisStatus, error) {
	fmt.Println("Adding external coverage scan data")

	analysisStatus, err := cli.AddScanResult(aID.ID, aID.TeamID, aID.ProjectID, "accepted", "vulnerability", aID.APIKey, *c.Value)
	if err != nil {
		return nil, fmt.Errorf("Analysis vulnerabilities save failed: %v", err.Error())
	}
	return analysisStatus, nil
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
