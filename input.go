package main

import (
	"errors"
	"strconv"
)

type InputResultWorkflowInput struct{ Key, Value string }

type InputResult struct {
	branch         string
	workflow       GhWorkflow
	workflowInputs []InputResultWorkflowInput
	isRun          bool
}

// NewInputResult asks the user to all the required inputs to run a workflow.
// The answers are stored in InputResult receiver.
func NewInputResult() (*InputResult, error) {
	r := &InputResult{}

	if err := r.askBranch(); err != nil {
		return r, err
	} else if err := r.askWorkflow(); err != nil {
		return r, err
	} else if err := r.askWorkflowInputs(); err != nil {
		return r, err
	} else if err := r.askRun(); err != nil {
		return r, err
	}

	return r, nil
}

// AskBranch asks the user to select a branch.
// The answer is stored in InputResult receiver.
func (r *InputResult) askBranch() error {
	currentBranch, err := getBranchName()
	if err != nil {
		return err
	}

	rBranches, err := getRemoteBranches()
	if err != nil {
		return err
	}

	if len(rBranches) == 0 {
		return errors.New("No remote branches found")
	}

	if len(rBranches) == 1 && rBranches[0] == currentBranch {
		answer, err := AskConfirm("Run on this branch: "+currentBranch, true)
		if err != nil {
			return err
		}

		if !answer {
			return errors.New("No other executable branch found")
		}
		r.branch = currentBranch
		return nil
	}

	// order rBranches so that currentBranch is at the top
	for i, b := range rBranches {
		if b == currentBranch {
			// remove currentBranch from rBranches
			rBranches = rBranches[:i+copy(rBranches[i:], rBranches[i+1:])]

			// append currentBranch to the top
			rBranches = append([]string{currentBranch}, rBranches...)
		}
	}

	answer, err := AskChoices("Select a branch", rBranches, currentBranch)
	if err != nil {
		return err
	}

	r.branch = answer

	return nil
}

// selectWorkflow asks the user to select a workflow.
// If there is only one workflow, it ask ok or cancel.
// The answer is stored in InputResult receiver.
func (r *InputResult) askWorkflow() error {
	var selectedWorkflow GhWorkflow
	workflows, err := getWorkflows()
	if err != nil {
		return err
	}

	if len(workflows) == 0 {
		return errors.New("No active workflows found")
	}

	workflowNames := []string{}
	for _, workflow := range workflows {
		workflowNames = append(workflowNames, workflow.Name)
	}

	if len(workflowNames) == 1 {
		ok, err := AskConfirm("Can I proceed with the following file? \n"+workflowNames[0], true)

		if err != nil {
			return err
		}

		if !ok {
			return errors.New("Canceled")
		}

		r.workflow = workflows[0]
		return nil
	}

	selectedWorkflowName, err := AskChoices("Select the workflow you wish to run", workflowNames, workflowNames[0])
	if err != nil {
		return err
	}

	for _, w := range workflows {
		if w.Name == selectedWorkflowName {
			selectedWorkflow = w
		}
	}

	if selectedWorkflow.Name == "" {
		return errors.New("No workflow found")
	}

	r.workflow = selectedWorkflow
	return nil
}

// askWorkflowInputs asks workflow inputs to user.
func (r *InputResult) askWorkflowInputs() error {
	var err error
	answers := []InputResultWorkflowInput{}

	if r.workflow.Name == "" {
		return errors.New("No workflow found. Need to run AskWorkflow() before AskWorkflowInputs()")
	}

	w, err := r.workflow.GetWorkflowInputs()
	if err != nil {
		return err
	}

	for _, v := range w {
		message := v.Description
		if message == "" {
			message = v.Name
		}

		var answer string
		switch v.Type {
		case "choice":
			answer, err = AskChoices(message, v.Options, v.Options[0])
		case "bool":
			var ok bool
			d, _ := strconv.ParseBool(v.Default)
			ok, err = AskConfirm(message, d)
			answer = strconv.FormatBool(ok)
		default:
			answer, err = AskInput(message, v.Default)
		}

		if err != nil {
			return err
		}

		answers = append(answers, InputResultWorkflowInput{
			Key:   v.Name,
			Value: answer,
		})
	}

	r.workflowInputs = answers
	return nil
}

// AskRun asks the user to confirm the execution.
// The answer is stored in InputResult receiver.
func (r *InputResult) askRun() error {
	renderTable(*r)
	answer, err := AskConfirm("Run this?", true)
	if err != nil {
		return err
	}

	r.isRun = answer
	return nil
}
