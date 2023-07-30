package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

const version = "0.0.0"

func main() {
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		fmt.Printf("ghrun version %s\n", version)
		os.Exit(0)
	}

	branch, err := AskBranch()
	if err != nil {
		log.Fatal(err)
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
	if err != nil {
		log.Fatal(err)
	}

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
	answer, err := AskConfirm("Run this?", true)
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
