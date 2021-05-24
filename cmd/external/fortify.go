package external

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ion-channel/ionic"
	"github.com/ion-channel/ionic/scanner"
	"github.com/ion-channel/ionize/dropbox"
)

//ParseFortify a Fortify FPR file at the path provided
func ParseFortify(path string) (*Fortify, error) {
	dir, err := unzip(path)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip, %s: %v", path, err)
	}

	fvdlFile := filepath.Join(dir, "audit.fvdl")
	b, err := ioutil.ReadFile(fvdlFile) // just pass the file name
	if err != nil {
		return nil, err
	}

	fvdl := FVDL{}

	err = xml.Unmarshal(b, &fvdl)
	if err != nil {
		return nil, err
	}

	rando, err := dropbox.Randomizer()
	if err != nil {
		return nil, err
	}

	erl, err := dropbox.ParseURL(path, rando)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(fmt.Sprintf("{\"fpr\": \"%v\"}", erl))
	ex := scanner.ExternalScan{}
	ex.Vulnerability = &scanner.ExternalVulnerability{}

	for _, v := range fvdl.Vulnerabilities.Vulnerability {
		confidence, _ := strconv.ParseFloat(v.InstanceInfo.Confidence, 64)
		impact, _ := strconv.ParseFloat(fvdl.Group(v.ClassInfo.ClassID, Impact), 64)
		probability, _ := strconv.ParseFloat(fvdl.Group(v.ClassInfo.ClassID, Probability), 64)
		accuracy, _ := strconv.ParseFloat(fvdl.Group(v.ClassInfo.ClassID, Accuracy), 64)

		likelihood := (accuracy * confidence * probability) / 25
		if impact >= 2.5 && likelihood >= 2.5 {
			ex.Vulnerability.Critcal++
		}

		if impact >= 2.5 && likelihood < 2.5 {
			ex.Vulnerability.High++
		}

		if impact < 2.5 && likelihood >= 2.5 {
			ex.Vulnerability.Medium++
		}

		if impact < 2.5 && likelihood < 2.5 {
			ex.Vulnerability.Low++
		}
	}

	ex.Source = scanner.Source{
		Name: "Fortify",
	}
	ex.Raw = &raw

	return &Fortify{
		FVDL:  &fvdl,
		Value: &ex,
	}, nil
}

//Fortify struct container for encapsulating external vulnerability scan data
type Fortify struct {
	FVDL  *FVDL
	Value *scanner.ExternalScan
}

//Save sends the external vulnerability scan data to ion channel for persistance
func (f *Fortify) Save(aID *AnalysisID, cli *ionic.IonClient) (*scanner.AnalysisStatus, error) {
	fmt.Println("Adding external fortify scan data")

	analysisStatus, err := cli.AddScanResult(aID.ID, aID.TeamID, aID.ProjectID, "accepted", "vulnerability", aID.APIKey, *f.Value)
	if err != nil {
		return nil, fmt.Errorf("Analysis vulnerabilities save failed: %v", err.Error())
	}
	return analysisStatus, nil
}

func unzip(path string) (string, error) {
	dest, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return "", err
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fpath, fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {
			// Make File
			err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			if err != nil {
				return fpath, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fpath, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return fpath, err
			}
		}
	}

	return dest, nil
}
