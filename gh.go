package main

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"

	"os/exec"

	"gopkg.in/yaml.v3"
)

type GitHubWorkflow struct {
	Name   string
	Status string
	Id     string
}

// Run runs a workflow.
func (w *GitHubWorkflow) Run(branch string, fieldArgs map[string]string) error {
	args := []string{"workflow", "run", w.Id, "-r", branch}
	for k, v := range fieldArgs {
		args = append(args, "-f", k+"="+v)
	}
	cmd := exec.Command("gh", args...)
	return cmd.Run()
}

// SelectWorkflow returns a workflow selected by user.
// If there is only one workflow, it ask ok or cancel.
func SelectWorkflowByUser(workflows []GitHubWorkflow) (GitHubWorkflow, error) {
	var selectedWorkflow GitHubWorkflow

	workflowNames := []string{}
	for _, workflow := range workflows {
		workflowNames = append(workflowNames, workflow.Name)
	}

	if len(workflowNames) == 1 {
		ok, err := AskConfirm("Can I proceed with the following file? \n"+workflowNames[0], true)

		if err != nil {
			return selectedWorkflow, err
		}

		if !ok {
			return selectedWorkflow, errors.New("Canceled")
		}

		return workflows[0], nil
	}

	selectedWorkflowName, err := AskChoices("Select the workflow you wish to run", workflowNames, workflowNames[0])
	if err != nil {
		return selectedWorkflow, err
	}

	for _, w := range workflows {
		if w.Name == selectedWorkflowName {
			selectedWorkflow = w
		}
	}

	if selectedWorkflow.Name == "" {
		return selectedWorkflow, errors.New("No workflow found")
	}
	return selectedWorkflow, nil
}

// GetWorkflows returns a list of active workflows.
func GetWorkflows() ([]GitHubWorkflow, error) {
	var workflows []GitHubWorkflow

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

		workflows = append(workflows, GitHubWorkflow{
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

type GitHubWorkflowInputs map[string]struct {
	Required    bool
	Description string
	Default     string
	Type        string
	Options     []string
}

type RespGitHubWorkflowInputs struct {
	Name string
	On   struct {
		WorkflowDispatch struct {
			Inputs GitHubWorkflowInputs
		} `yaml:"workflow_dispatch"`
	}
}

// AskToUser asks inputs to user.
func (w *GitHubWorkflowInputs) AskToUser() (map[string]string, error) {
	var err error
	answers := make(map[string]string)

	for k, v := range *w {
		message := v.Description
		if message == "" {
			message = k
		}

		var answer string
		switch v.Type {
		case "choice":
			answer, err = AskChoices(message, v.Options, v.Options[0])
		case "bool":
			var ok bool
			println(v.Default)
			// var defaultInput bool
			d, _ := strconv.ParseBool(v.Default)
			ok, err = AskConfirm(message, d)
			answer = strconv.FormatBool(ok)
		default:
			answer, err = AskInput(message, v.Default)
		}

		if err != nil {
			return nil, err
		}

		answers[k] = answer
	}

	return answers, nil
}

// GetWorkflowInputs returns inputs for a workflow.
func GetWorkflowInputs(workflowId string) (GitHubWorkflowInputs, error) {
	var workflowInputs GitHubWorkflowInputs
	cmd := exec.Command("gh", "workflow", "view", workflowId, "-y")
	out, err := cmd.Output()

	if err != nil {
		return workflowInputs, err
	}

	w := RespGitHubWorkflowInputs{}
	err = yaml.Unmarshal(out, &w)
	if err != nil {
		return workflowInputs, err
	}

	return w.On.WorkflowDispatch.Inputs, nil
}
