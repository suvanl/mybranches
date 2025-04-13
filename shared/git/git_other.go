//go:build !windows

package git

import (
	"os/exec"
	"strings"
)

// DeleteBranchesNotOnRemote finds 'gone' branches and optionally deletes them if dryRun is false.
//
// Returns number of branches deleted and error (if any)
func DeleteBranchesNotOnRemote(dryRun bool) (int, error) {
	branches, err := getGoneBranches(dryRun)
	if err != nil {
		return -1, err
	}

	return len(branches), nil
}

func getGoneBranches(dryRun bool) ([]string, error) {
	// Introduces a dependency on sh, which may symlink to different shells depending on the platform this is running on.
	// This is simpler than programmatically piping the commands between each other; let's delegate that work for now.
	// If it starts to lead to inconsistent behaviour across platforms, might be a good idea to revisit this.

	cmd := "git branch -vv | grep 'gone]' | awk '{print $1}'"
	if !dryRun {
		cmd += " | xargs git branch -D"
	}

	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}

	result := strings.TrimSpace(string(out[:]))
	if result == "" {
		return []string{}, nil
	}

	branchNames := strings.Split(result, "\n")
	return branchNames, nil
}
