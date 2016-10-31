// ion-cli.go
//
// Copyright (C) 2016 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Jeffail/gabs"
)

const (
	app            = "ION-CI"
	appDescription = "CI/CD logic wrapper around ion-connect"
)

var (
	buildTime  string
	appVersion string
)

func main() {

	argsWithoutProg := os.Args[1:]
	//fmt.Println(argsWithoutProg)
	err, analysisid := AnalyzeProject(argsWithoutProg[0], argsWithoutProg[1], argsWithoutProg[2])
	if err {
		fmt.Println("The build passed : ", false)
		os.Exit(1)
	}
	var status string
	for ok := true; ok; ok = (status == "accepted" || err == true) {
		time.Sleep(500 * time.Millisecond)
		err, status = GetAnalysisStatus(analysisid)
		if err {
			fmt.Println("The build passed : ", false)
			os.Exit(1)
		}
	}
	err, passfail := GetAnalysis(analysisid)
	if err {
		fmt.Println("The build passed : ", false)
		os.Exit(1)
	}
	fmt.Println("The build passed : ", passfail)
}

func AnalyzeProject(projectid string, accountid string, buildnumber string) (bool, string) {
	cmdName := "ion-connect"
	cmdArgs := []string{"scanner", "analyze-project", "--account-id", accountid, "--project-id", projectid, buildnumber}
	out, err2 := exec.Command(cmdName, cmdArgs...).Output()
	if err2 != nil {
		return true, ""
	}
	jsonParsed2, err2 := gabs.ParseJSON(out)
	if str, ok := jsonParsed2.Path("id").Data().(string); ok {
		return false, str
	}
	return true, ""
}

func GetAnalysisStatus(analysisid string) (bool, string) {
	cmdName := "ion-connect"
	cmdArgs := []string{"scanner", "get-analysis-status", "--account-id", "account_id", "--project-id", "7b9a0a87-fbe6-40c1-aa37-89acc6e5c191", analysisid}
	out, err2 := exec.Command(cmdName, cmdArgs...).Output()
	if err2 != nil {
		return true, ""
	}
	jsonParsed2, err2 := gabs.ParseJSON(out)
	if str, ok := jsonParsed2.Path("status").Data().(string); ok {
		return false, str
	}
	return true, ""
}

func GetAnalysis(analysisid string) (bool, bool) {
	cmdName := "ion-connect"
	cmdArgs := []string{"analysis", "get-analysis", "--account-id", "account_id", "--project-id", "7b9a0a87-fbe6-40c1-aa37-89acc6e5c191", analysisid}
	out, err2 := exec.Command(cmdName, cmdArgs...).Output()
	if err2 != nil {
		return true, false
	}
	jsonParsed2, err2 := gabs.ParseJSON(out)
	return false, jsonParsed2.Path("passed").Data().(bool)
}
