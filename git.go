package main

import (
	"bufio"
	"bytes"
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
