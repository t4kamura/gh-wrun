package main

import (
	"bufio"
	"bytes"
	"errors"

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

// TODO: test

// GetWorkflows returns a list of active workflows.
func getWorkflows() ([]GhWorkflow, error) {
	var workflows []GhWorkflow

	// if include disabled, add -a flag
	cmd := exec.Command("gh", "workflow", "list")
	out, err := cmd.Output()

	if err != nil {
		return workflows, err
	}

	sc := bufio.NewScanner(bytes.NewReader(out))
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

	if len(workflows) == 0 {
		return workflows, errors.New("No workflows found")
	}

	return workflows, nil
}

// TODO: test

// GetWorkflowInputs returns inputs for a workflow.
func (g *GhWorkflow) GetWorkflowInputs() ([]GhWorkflowInput, error) {
	var w []GhWorkflowInput
	cmd := exec.Command("gh", "workflow", "view", g.Id, "-y")
	out, err := cmd.Output()

	if err != nil {
		return w, err
	}

	r := struct {
		Name string
		On   struct {
			WorkflowDispatch struct {
				Inputs yaml.MapSlice
			} `yaml:"workflow_dispatch"`
		}
	}{}

	err = yaml.Unmarshal(out, &r)
	if err != nil {
		return w, err
	}

	for _, v := range r.On.WorkflowDispatch.Inputs {
		name := v.Key.(string)

		var (
			required     bool
			description  string
			defaultValue string
			typeValue    string
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
					defaultValue = vv.Value.(string)
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
