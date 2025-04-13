package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// git fetch --prune
func FetchPrune() error {
	_, err := exec.Command("git", "fetch", "--prune").CombinedOutput()
	return err
}

// Returns branch names found by the `git branch --list <pattern> --format %(refname:short)` command.
func FindBranches(pattern string) []string {
	globPattern := fmt.Sprintf("%s*", pattern)
	out, err := exec.Command("git", "branch", "--list", globPattern, "--format", "%(refname:short)").CombinedOutput()
	if err != nil {
		log.Fatalf("Error finding branches: %v\n", err)
	}

	fromBytes := string(out[:])
	branches := strings.Split(fromBytes, "\n")

	// Last element will be an empty string, let's just drop it here
	return branches[:len(branches)-1]
}

func GetCurrentBranchName() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatalf("Error getting current branch: %v\n", err)
	}

	fromBytes := string(out[:])
	return strings.Split(fromBytes, "\n")[0]
}

// Returns the output of the `git switch` command
func SwitchBranch(branchName string) string {
	out, err := exec.Command("git", "switch", branchName).CombinedOutput()
	if err != nil {
		log.Fatalf("Error switching branch: %v\n", err)
	}

	return string(out[:])
}

// Returns the output of the `git branch -D <branch>` command
func DeleteBranch(branchName string) string {
	out, err := exec.Command("git", "branch", "-D", branchName).CombinedOutput()
	if err != nil {
		log.Fatalf("Error deleting branch: %v\n", err)
	}

	return string(out[:])
}
