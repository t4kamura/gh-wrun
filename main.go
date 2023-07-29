package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

const version = "0.0.0"

func main() {
	// TODO:ブランチ、選択方式に変更

	branch, err := GetBranchName()
	if err != nil {
		log.Fatal(err)
	}
	answer, err := AskConfirm("Running git branch: "+branch, true)
	if err != nil {
		log.Fatal(err)
	}
	rBranches, err := GetRemoteBranches()
	if err != nil {
		log.Fatal(err)
	}

	if !answer {
		answer, err := AskChoices("Which branch would you like to run a warflow on?", rBranches)
		if err != nil {
			log.Fatal(err)
		}
		branch = answer
	} else {
		// check if the current branch is in the remote
		isRemote := false
		for _, b := range rBranches {
			if b == branch {
				isRemote = true
				break
			}
		}
		if !isRemote {
			log.Fatal("The current branch is not in the remote")
		}
	}

	workflows, err := GetWorkflows()
	if err != nil {
		log.Fatal(err)
	}

	if len(workflows) == 0 {
		log.Fatal("No active workflows found")
	}

	selectedWorkflow, err := SelectWorkflowByUser(workflows)
	if err != nil {
		log.Fatal(err)
	}

	workflowInputs, err := GetWorkflowInputs(selectedWorkflow.Id)
	if err != nil {
		log.Fatal(err)
	}

	selectedWorkflowFile := filepath.Base(selectedWorkflow.Name)
	fieldArgs, err := workflowInputs.AskToUser()

	// generate table
	tableData := [][]string{
		{"Targets", "Git branch", branch},
		{"Targets", "Workflow", selectedWorkflowFile},
	}
	for k, v := range fieldArgs {
		tableData = append(tableData, []string{"Inputs", k, v})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(tableData)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render()

	// confirm
	answer, err = AskConfirm("Would you like to run a warflow with the above information?", true)
	if err != nil {
		log.Fatal(err)
	}

	if !answer {
		log.Fatal("Canceled")
	}

	if err := selectedWorkflow.Run(branch, fieldArgs); err != nil {
		log.Fatal(err)
	}
}
