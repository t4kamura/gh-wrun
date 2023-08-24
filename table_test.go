package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRenderTable(t *testing.T) {
	input := InputResult{
		branch: "main",
		workflow: GhWorkflow{
			Name:   ".github/workflows/test.yml",
			Status: "active",
			Id:     "123456789",
		},
		workflowInputs: []InputResultWorkflowInput{
			{
				Key: "env", Value: "test",
			},
			{
				Key: "message", Value: "test",
			},
			{
				Key: "server", Value: "app",
			},
		},
	}

	want := "+---------+------------+----------+\n" +
		"| Targets | Git branch | main     |\n" +
		"+         +------------+----------+\n" +
		"|         | Workflow   | test.yml |\n" +
		"+---------+------------+----------+\n" +
		"| Inputs  | env        | test     |\n" +
		"+         +------------+          +\n" +
		"|         | message    |          |\n" +
		"+         +------------+----------+\n" +
		"|         | server     | app      |\n" +
		"+---------+------------+----------+\n"

	orgStdout := os.Stdout
	pr, pw, _ := os.Pipe()
	os.Stdout = pw

	renderTable(input)

	pw.Close()
	os.Stdout = orgStdout

	buf := bytes.Buffer{}
	io.Copy(&buf, pr)
	output := buf.String()

	if output != want {
		t.Errorf("Expected is %s but got %s\n", want, output)
	}
}
