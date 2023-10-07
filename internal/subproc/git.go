package subproc

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

// getBranchName returns the current branch name.
func GetBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// getRemoteBranches returns the list of remote branches.
func GetRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r")
	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	branches := parseBranchesResult(&out)

	return branches, nil
}

// parseBranchesResult parses the output of `git branch -r` and returns the list of branches.
func parseBranchesResult(out *[]byte) []string {
	sc := bufio.NewScanner(bytes.NewReader(*out))
	branches := []string{}
	for sc.Scan() {
		// to exclude results such as
		// origin/HEAD -> origin/main
		scc := bufio.NewScanner(bytes.NewReader(sc.Bytes()))
		scc.Split(bufio.ScanWords)
		var words []string

		for scc.Scan() {
			words = append(words, scc.Text())
		}

		if len(words) == 1 {
			branches = append(branches, words[0])
		}
	}

	// Remove the "origin/" prefix
	formattedBranches := []string{}
	for _, b := range branches {
		formattedBranches = append(formattedBranches, strings.Split(b, "/")[1])
	}

	return formattedBranches
}
