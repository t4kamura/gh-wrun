package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

// GetBranchName returns the current branch name.
func GetBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GetRemoteBranches returns the list of remote branches.
// but exclude HEAD.
func GetRemoteBranches() ([]string, error) {
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

	formattedBranches := []string{}
	for _, b := range branches {
		formattedBranches = append(formattedBranches, strings.Split(b, "/")[1])
	}

	return formattedBranches, nil
}
