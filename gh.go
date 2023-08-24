package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"os/exec"

	"gopkg.in/yaml.v2"
)

type GhWorkflow struct{ Name, Status, Id string }

type GhWorkflowInput struct {
	Name,
	Description,
	Default,
	Type string
	Required bool
	Options  []string
}

type GhWorkflowInputsYaml struct {
	Name string
	On   struct {
		WorkflowDispatch struct {
			Inputs yaml.MapSlice
		} `yaml:"workflow_dispatch"`
	}
}

const (
	GhWorkflowInputTypeString      = "string"
	GhWorkflowInputTypeChoice      = "choice"
	GhWorkflowInputTypeBoolean     = "boolean"
	GhWorkflowInputTypeEnvironment = "environment"
)

// GetWorkflows returns a list of active workflows.
func getWorkflows() ([]GhWorkflow, error) {
	// if include disabled, add -a flag
	cmd := exec.Command("gh", "workflow", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	workflows, err := parseWorkflows(out)
	if err != nil {
		return workflows, err
	} else if len(workflows) == 0 {
		return nil, errors.New("No workflows found")
	}

	return workflows, nil
}

// parseWorkflows parses the output from gh workflow list.
func parseWorkflows(src []byte) ([]GhWorkflow, error) {
	var workflows []GhWorkflow
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		scc := bufio.NewScanner(bytes.NewReader(sc.Bytes()))
		scc.Split(bufio.ScanWords)

		var words []string
		for scc.Scan() {
			words = append(words, scc.Text())
		}

		if len(words) != 3 {
			return workflows, errors.New("Error parsing workflow")
		}

		workflows = append(workflows, GhWorkflow{
			Name:   words[0],
			Status: words[1],
			Id:     words[2],
		})
	}

	return workflows, nil
}

// GetWorkflowInputs returns inputs for a workflow.
func (g *GhWorkflow) GetWorkflowInputs(runLinter bool) ([]GhWorkflowInput, error) {
	cmd := exec.Command("gh", "workflow", "view", g.Id, "-y")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	if runLinter {
		lErrors, err := lint(out)
		if err != nil {
			return nil, err
		}
		if lErrors != nil {
			for _, e := range lErrors {
				fmt.Printf("%d:%d: %s\n", e.Line, e.Column, e.Message)
			}

			return nil, errors.New("Workflow file is invalid")
		}
	}

	w, err := parseWorkflowInputs(out)

	return w, nil
}

// parseWorkflowInputs parses the output from gh workflow view.
func parseWorkflowInputs(src []byte) ([]GhWorkflowInput, error) {
	var w []GhWorkflowInput
	r := GhWorkflowInputsYaml{}

	if err := yaml.Unmarshal(src, &r); err != nil {
		return w, err
	}

	inputs := r.On.WorkflowDispatch.Inputs

	if len(inputs) == 0 {
		return w, fmt.Errorf("No inputs found for workflow %s", r.Name)
	}

	for _, v := range inputs {
		name := v.Key.(string)

		var (
			required     bool
			description  string
			defaultValue string
			typeValue    string = GhWorkflowInputTypeString
			options      []string
		)

		if p, ok := v.Value.(yaml.MapSlice); ok {
			for _, vv := range p {
				switch vv.Key.(string) {
				case "required":
					required = vv.Value.(bool)
				case "description":
					description = vv.Value.(string)
				case "default":
					defaultValue, ok = vv.Value.(string)
					if !ok {
						defaultValue = strconv.FormatBool(vv.Value.(bool))
					}
				case "type":
					typeValue = vv.Value.(string)
				case "options":
					for _, o := range vv.Value.([]interface{}) {
						options = append(options, o.(string))
					}
				}
			}
		}

		w = append(w, GhWorkflowInput{
			Name:        name,
			Required:    required,
			Description: description,
			Default:     defaultValue,
			Type:        typeValue,
			Options:     options,
		})
	}

	return w, nil
}

// Run runs a workflow.
func (w *GhWorkflow) Run(inputResult *InputResult) error {
	branch := inputResult.branch
	fieldArgs := inputResult.workflowInputs

	args := []string{"workflow", "run", w.Id, "-r", branch}
	for _, m := range fieldArgs {
		args = append(args, "-f", m.Key+"="+m.Value)
	}
	cmd := exec.Command("gh", args...)
	return cmd.Run()
}
