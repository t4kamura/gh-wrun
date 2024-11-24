package input

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/t4kamura/gh-wrun/internal/interactive"
	"github.com/t4kamura/gh-wrun/internal/subproc"
	"github.com/t4kamura/gh-wrun/internal/table"
)

type InputResult struct {
	Branch         string
	Workflow       subproc.GhWorkflow
	WorkflowInputs []struct{ Key, Value string }
	IsRun          bool
}

// NewInputResult asks the user to all the required inputs to run a workflow.
// The answers are stored in InputResult receiver.
func NewInputResult(branchAuto bool) (*InputResult, error) {
	r := &InputResult{}

	if err := r.askBranch(branchAuto); err != nil {
		return r, err
	} else if err := r.askWorkflow(); err != nil {
		return r, err
	} else if err := r.askWorkflowInputs(); err != nil {
		return r, err
	} else if err := r.askRunWithRenderTable(); err != nil {
		return r, err
	}

	return r, nil
}

// AskBranch asks the user to select a branch.
// The answer is stored in InputResult receiver.
// If the auto flag is true, automatically set the current branch
func (r *InputResult) askBranch(auto bool) error {
	currentBranch, err := subproc.GetBranchName()
	if err != nil {
		return err
	}

	if auto {
		r.Branch = currentBranch
		return nil
	}

	rBranches, err := subproc.GetRemoteBranches()
	if err != nil {
		return err
	}

	if len(rBranches) == 0 {
		return errors.New("No remote branches found")
	}

	if len(rBranches) == 1 && rBranches[0] == currentBranch {
		answer := interactive.AskConfirm("Run on this branch: " + currentBranch)
		if !answer {
			return errors.New("No other executable branch found")
		}
		r.Branch = currentBranch
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

	answer, err := interactive.AskChoices("Select a branch", rBranches, currentBranch)
	if err != nil {
		return err
	}

	r.Branch = answer

	return nil
}

// selectWorkflow asks the user to select a workflow.
// If there is only one workflow, it ask ok or cancel.
// The answer is stored in InputResult receiver.
func (r *InputResult) askWorkflow() error {
	var selectedWorkflow subproc.GhWorkflow
	workflows, err := subproc.GetWorkflows()
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
		ok := interactive.AskConfirm(fmt.Sprintf("Do you want to run [%s]", workflowNames[0]))
		if !ok {
			return errors.New("Canceled")
		}

		r.Workflow = workflows[0]
		return nil
	}

	selectedWorkflowName, err := interactive.AskChoices("Select the workflow you wish to run", workflowNames, workflowNames[0])
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

	r.Workflow = selectedWorkflow
	return nil
}

// askWorkflowInputs asks workflow inputs to user.
func (r *InputResult) askWorkflowInputs() error {
	var err error
	answers := []struct{ Key, Value string }{}

	if r.Workflow.Name == "" {
		return errors.New("No workflow found. Need to run AskWorkflow() before AskWorkflowInputs()")
	}

	w, err := r.Workflow.GetWorkflowInputs()
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
		case subproc.GhWorkflowInputTypeChoice:
			answer, err = interactive.AskChoices(message, v.Options, v.Options[0])
		case subproc.GhWorkflowInputTypeBoolean:
			var ok bool
			d, _ := strconv.ParseBool(v.Default)
			ok, err = interactive.AskBool(message, d)
			answer = strconv.FormatBool(ok)
		case subproc.GhWorkflowInputTypeEnvironment:
			envs, err := subproc.GetEnvironments()
			if err != nil {
				return err
			}
			if len(envs) == 0 {
				return fmt.Errorf("no environments exist")
			}
			answer, err = interactive.AskChoices(message, envs, envs[0])
		default:
			answer, err = interactive.AskInput(message, v.Default)
		}

		if err != nil {
			return err
		}

		answers = append(answers, struct{ Key, Value string }{
			Key:   v.Name,
			Value: answer,
		})
	}

	r.WorkflowInputs = answers
	return nil
}

// AskRun asks the user to confirm the execution.
// Render the table and ask if it is ok to run.
// The answer is stored in InputResult receiver.
func (r *InputResult) askRunWithRenderTable() error {
	table.Render(r.genTableData())
	answer := interactive.AskConfirm("Run this?")

	r.IsRun = answer
	return nil
}

// genTableData generates table data from InputResult receiver.
// It is used to render the table.
func (r *InputResult) genTableData() [][]string {
	selectedWorkflowFile := filepath.Base(r.Workflow.Name)
	tableData := [][]string{
		{"Targets", "Git branch", r.Branch},
		{"Targets", "Workflow", selectedWorkflowFile},
	}
	for _, m := range r.WorkflowInputs {
		tableData = append(tableData, []string{"Inputs", m.Key, m.Value})
	}

	return tableData
}
