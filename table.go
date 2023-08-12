package main

import (
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

func renderTable(r InputResult) {
	// generate table
	selectedWorkflowFile := filepath.Base(r.workflow.Name)
	tableData := [][]string{
		{"Targets", "Git branch", *r.branch},
		{"Targets", "Workflow", selectedWorkflowFile},
	}
	for k, v := range *r.workflowInputs {
		tableData = append(tableData, []string{"Inputs", k, v})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(tableData)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render()
}
