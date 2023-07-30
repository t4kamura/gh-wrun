package main

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// getBranchName returns the current branch name.
func getBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// getRemoteBranches returns the list of remote branches.
func getRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r")
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	sc := bufio.NewScanner(bytes.NewReader(out))
	branches := []string{}
	for sc.Scan() {
		branches = append(branches, sc.Text())
	}

	// Remove the "origin/" prefix
	formattedBranches := []string{}
	for _, b := range branches {
		formattedBranches = append(formattedBranches, strings.Split(b, "/")[1])
	}

	return formattedBranches, nil
}

// AskBranch asks the user to select a branch.
func AskBranch() (string, error) {
	currentBranch, err := getBranchName()
	if err != nil {
		return "", err
	}

	rBranches, err := getRemoteBranches()
	if err != nil {
		return "", err
	}

	if len(rBranches) == 0 {
		return "", errors.New("No remote branches found")
	}

	if len(rBranches) == 1 && rBranches[0] == currentBranch {
		answer, err := AskConfirm("Run on this branch: "+currentBranch, true)
		if err != nil {
			return "", err
		}

		if !answer {
			return "", errors.New("No other executable branch found")
		}
		return currentBranch, nil
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
		return "", err
	}

	return answer, nil
}
