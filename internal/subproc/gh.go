package subproc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"os/exec"

	"gopkg.in/yaml.v2"
)

type GhWorkflow struct {
	Id     json.Number `json:"id"`
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Status string      `json:"state"`
}

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

// GetGhVersion returns the version of gh.
func GetGhVersion() (string, error) {
	cmd := exec.Command("gh", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var words []string
	sc := bufio.NewScanner(bytes.NewReader(out))
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		words = append(words, sc.Text())
	}
	if len(words) < 3 {
		return "", errors.New("Error parsing gh version")
	}
	return string(words[2]), nil
}

// GetWorkflows returns a list of active workflows.
func GetWorkflows() ([]GhWorkflow, error) {
	// if include disabled, add -a flag
	cmd := exec.Command("gh", "workflow", "list", "--json", "id,name,path,state")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var workflows []GhWorkflow
	if err := json.Unmarshal(out, &workflows); err != nil {
		return nil, err
	}

	if len(workflows) == 0 {
		return nil, errors.New("No workflows found")
	}

	return workflows, nil
}

// GetWorkflowInputs returns inputs for a workflow.
func (g *GhWorkflow) GetWorkflowInputs() ([]GhWorkflowInput, error) {
	cmd := exec.Command("gh", "workflow", "view", string(g.Id), "-y")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	w, err := parseWorkflowInputs(out)
	if err != nil {
		return nil, err
	}

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

	// this is blank inputs case
	if len(inputs) == 0 {
		return []GhWorkflowInput{}, nil
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
func (w *GhWorkflow) Run(branch string, fieldArgs []struct{ Key, Value string }) error {
	args := []string{"workflow", "run", string(w.Id), "-r", branch}
	for _, m := range fieldArgs {
		args = append(args, "-f", m.Key+"="+m.Value)
	}
	cmd := exec.Command("gh", args...)
	return cmd.Run()
}

type GhApiGetEnvironmentsResult struct {
	TotalCount   int                                     `json:"total_count"`
	Environments []GhApiGetEnvironmentsResultEnvironment `json:"environments"`
}

type GhApiGetEnvironmentsFailedResult struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Status           string `json:"status"`
}

type GhApiGetEnvironmentsResultEnvironment struct {
	Name string `json:"name"`
}

// GetEnvironments returns current repository GitHub Environments.
func GetEnvironments() ([]string, error) {
	repoWithOwner, err := getCurrentRepositoryWithOwner()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/repos/%s/environments", repoWithOwner)

	cmd := exec.Command("gh", "api", "-H", "Accept: application/vnd.github+json", "-H", "X-GitHub-Api-Version: 2022-11-28", endpoint)
	out, err := cmd.Output()
	if err != nil {
		if len(out) > 0 {
			var failedRes GhApiGetEnvironmentsFailedResult
			if err := json.Unmarshal(out, &failedRes); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("failed to get environments, message: %s, status: %s, doc url: %s", failedRes.Message, failedRes.Status, failedRes.DocumentationURL)
		}
		return nil, err
	}

	var environmentsResult GhApiGetEnvironmentsResult
	if err := json.Unmarshal(out, &environmentsResult); err != nil {
		return nil, err
	}

	environments := make([]string, 0, environmentsResult.TotalCount)
	for _, e := range environmentsResult.Environments {
		environments = append(environments, e.Name)
	}
	return environments, nil
}

type GhRepoViewForNameWithOwnerResult struct {
	Name string `json:"nameWithOwner"`
}

// getCurrentRepositoryWithOwner returns current repository name with owner
// e.g. "t4kamura/gh-wrun"
func getCurrentRepositoryWithOwner() (string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "nameWithOwner")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var res GhRepoViewForNameWithOwnerResult
	if err := json.Unmarshal(out, &res); err != nil {
		return "", err
	}

	if res.Name == "" {
		return "", fmt.Errorf("not found repository")
	}

	return res.Name, err
}
