package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strings"
)

func findBranches(pattern string) []string {
	globPattern := fmt.Sprintf("%s*", pattern)
	out, err := exec.Command("git", "branch", "--list", globPattern, "--format", "%(refname:short)").CombinedOutput()
	if err != nil {
		log.Fatalf("Error finding branches: %v", err)
	}

	fromBytes := string(out[:])
	branches := strings.Split(fromBytes, "\n")

	// Last element will be an empty string, let's just drop it here
	return branches[:len(branches)-1]
}

func getCurrentBranchName() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatalf("Error getting current branch: %v", err)
	}

	fromBytes := string(out[:])
	return strings.Split(fromBytes, "\n")[0]
}

// Returns the output of the `git switch` command
func switchBranch(branchName string) string {
	out, err := exec.Command("git", "switch", branchName).CombinedOutput()
	if err != nil {
		log.Fatalf("Error switching branch: %v", err)
	}

	return string(out[:])
}

func getUsernamePattern() string {
	user, err := user.Current()

	if err != nil {
		defaultUsername := "user"
		fmt.Printf("Failed to determine your username. Defaulting to %s.", defaultUsername)
		return defaultUsername
	}

	// Different branch naming conventions exist, but all usually start with the author's name.
	// The character after this often differs (":", "/", "-" are commonly used), so we won't include it in default pattern.
	// If needed, it can be included in the value provided for the `--pattern` flag.
	return withoutDomain(user.Username)
}

// This should only affect Windows, where the username format is DOMAIN\username
func withoutDomain(withPossibleDomain string) string {
	split := strings.Split(withPossibleDomain, "\\")
	if len(split) == 2 {
		return split[1]
	}

	return withPossibleDomain
}
