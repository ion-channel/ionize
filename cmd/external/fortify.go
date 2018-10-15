package external

import (
	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/scanner"
)

//ParseFortify a Fortify FPR file at the path provided
func ParseFortify(path string) (*Fortify, error) {
	return nil, nil
}

//Fortify struct container for encapsalating external vulnerability scan data
type Fortify struct {
	Value *scanner.ExternalScan
}

//Save sends the external vulnerability scan data to ion channel for persistance
func (c *Fortify) Save(aID *AnalysisID, cli *ionic.IonClient) (*scanner.AnalysisStatus, error) {
	return nil, nil
}
