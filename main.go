package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
)

const (
	AppVersion = "0.0.1"
	BaseUrl    = "https://circleci.com/api/v1.1"
)

var (
	argVersion  = flag.Bool("version", false, "Print version")
	argUsername = flag.String("username", os.Getenv("CIRCLECI_USERNAME"), "CircleCI User Name")
	// argProject  = flag.String("project", os.Getenv("CIRCLECI_PROJECT"), "Project(Repository) Name")
	argVcs = flag.String("vcs", "github", "VCS Type [github|bitbacket]")
)

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Project struct {
	Branches struct {
		Master struct {
			RecentBuilds []struct {
				Outcome     string `json:"outcome"`
				Status      string `json:"status"`
				BuildNum    int    `json:"build_num"`
				VcsRevision string `json:"vcs_revision"`
				// PushedAt      time.Time `json:"pushed_at"`
				IsWorkflowJob bool `json:"is_workflow_job"`
				Is20Job       bool `json:"is_2_0_job"`
				// AddedAt       time.Time `json:"added_at"`
			} `json:"recent_builds"`
		} `json:"master"`
	} `json:"branches"`
	VcsURL string `json:"vcs_url"`
}

func displyTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)

	for _, value := range data {
		table.Append(value)
	}
	table.Render()
}

func describeProjects(apiEndpoint string) (allVariables [][]string) {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").Get(apiEndpoint)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// fmt.Println(resp)
	jsonBytes := ([]byte)(resp.Body())
	var ps []Project
	if err := json.Unmarshal(jsonBytes, &ps); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var allProject [][]string
	for _, p := range ps {
		var latestBuildNo, latestBuildStatus string
		if len(p.Branches.Master.RecentBuilds) > 0 {
			latestBuildNo = strconv.Itoa(p.Branches.Master.RecentBuilds[0].BuildNum)
			latestBuildStatus = p.Branches.Master.RecentBuilds[0].Status
		} else {
			latestBuildNo = "N/A"
			latestBuildStatus = "N/A"
		}
		splitedVcsUrl := strings.Split(p.VcsURL, "/")
		pjName := splitedVcsUrl[len(splitedVcsUrl)-1]
		project := []string{pjName, latestBuildNo, latestBuildStatus}
		allProject = append(allProject, project)
	}

	// fmt.Println(allProject)
	return allProject
}

func describeVariables(apiEndpoint string) (allVariables [][]string) {
	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").Get(apiEndpoint)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonBytes := ([]byte)(resp.Body())
	var vs []Variable
	if err := json.Unmarshal(jsonBytes, &vs); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var allVariable [][]string
	for _, v := range vs {
		variable := []string{v.Name, v.Value}
		allVariable = append(allVariable, variable)
	}

	return allVariable
}

func addVariable(apiEndpoint string, varName string, varValue string) {
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(Variable{Name: varName, Value: varValue}).Post(apiEndpoint)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return
}

func delVariable(apiEndpoint string) {
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json").Delete(apiEndpoint)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return
}

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	if os.Getenv("CIRCLECI_TOKEN") == "" {
		fmt.Println("Please supply the CircleCI Personal Token using `CIRCLECI_TOKEN` environment variable.")
		os.Exit(1)
	}

	// subcommand `variables`
	varCommand := flag.NewFlagSet("variables", flag.ExitOnError)
	varProject := varCommand.String("project", os.Getenv("CIRCLECI_PROJECT"), "Project(Repository) Name")
	varActAddFlag := varCommand.Bool("add", false, "Add variable")
	varActDelFlag := varCommand.Bool("del", false, "Delete variable")
	varNameFlag := varCommand.String("name", "", "Variable Name")
	varValueFlag := varCommand.String("value", "", "Variable Value")

	// subcommand `projects`
	prjCommand := flag.NewFlagSet("projects", flag.ExitOnError)

	if len(os.Args) == 1 {
		fmt.Println("usage: wakka <command> [<args>]")
		fmt.Println(" variables   Manage Project Environment Variables.")
		return
	}

	switch os.Args[1] {
	case "variables":
		varCommand.Parse(os.Args[2:])
	case "projects":
		prjCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if varCommand.Parsed() {
		apiEndpoint := fmt.Sprintf("%s/project/%s/%s/%s/envvar?circle-token=%s",
			BaseUrl, *argVcs, *argUsername, *varProject, os.Getenv("CIRCLECI_TOKEN"))

		if *varActAddFlag || *varActDelFlag {
			if *varNameFlag == "" {
				fmt.Println("Please supply the variable name using -name option.")
				return
			}
			apiEndpoint := fmt.Sprintf("%s/project/%s/%s/%s/envvar/%s?circle-token=%s",
				BaseUrl, *argVcs, *argUsername, *varProject, *varNameFlag, os.Getenv("CIRCLECI_TOKEN"))
			delVariable(apiEndpoint)
		}

		if *varActAddFlag {
			if *varValueFlag == "" {
				fmt.Println("Please supply the variable value using -value option.")
				return
			}
			addVariable(apiEndpoint, *varNameFlag, *varValueFlag)
		}

		allVariable := describeVariables(apiEndpoint)
		header := []string{"Variable Name", "Variable Value"}
		displyTable(header, allVariable)
	}

	if prjCommand.Parsed() {
		apiEndpoint := fmt.Sprintf("%s/projects?circle-token=%s", BaseUrl, os.Getenv("CIRCLECI_TOKEN"))
		allProject := describeProjects(apiEndpoint)
		header := []string{"Project Name", "Master Latest Build No", "Master Latest Build Status"}
		displyTable(header, allProject)
	}

}
