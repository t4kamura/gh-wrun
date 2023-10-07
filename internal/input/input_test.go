package input

import (
	"reflect"
	"testing"

	"github.com/t4kamura/gh-wrun/internal/subproc"
)

func TestGenTableData(t *testing.T) {
	input := InputResult{
		Branch:   "main",
		Workflow: subproc.GhWorkflow{Name: "test.yml", Status: "active", Id: "12345678"},
		WorkflowInputs: []struct{ Key, Value string }{
			{Key: "env", Value: "dev"},
			{Key: "message", Value: "test message"},
			{Key: "server", Value: "app"},
		},
		IsRun: true,
	}

	want := [][]string{
		{"Targets", "Git branch", "main"},
		{"Targets", "Workflow", "test.yml"},
		{"Inputs", "env", "dev"},
		{"Inputs", "message", "test message"},
		{"Inputs", "server", "app"},
	}

	result := input.genTableData()

	if !reflect.DeepEqual(result, want) {
		t.Errorf("Expected is %v but got %v\n", want, result)
	}
}
