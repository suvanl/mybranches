//go:build !windows

package git

import (
	"os/exec"
	"strings"
)

func DeleteBranchesNotOnRemote() (int, error) {
	branches, err := getGoneBranches()
	if err != nil {
		return -1, err
	}

	// todo: Delete all the gone branches

	return len(branches), nil
}

func getGoneBranches() ([]string, error) {
	// Introduces a dependency on sh, which may symlink to different shells depending on the platform this is running on.
	// This is simpler than programmatically piping the commands between each other; let's delegate that work for now.
	out, err := exec.Command("sh", "-c", "git branch -vv | grep 'gone]' | awk '{print $1}'").CombinedOutput()
	if err != nil {
		return nil, err
	}

	outStr := strings.TrimSpace(string(out[:]))
	if outStr == "" {
		return []string{}, nil
	}

	branchNames := strings.Split(outStr, "\n")
	return branchNames, nil
}
