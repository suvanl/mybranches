package main

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"
)

func getUsernamePattern() string {
	user, err := user.Current()

	if err != nil {
		defaultUsername := "user"
		fmt.Printf("Failed to determine your username. Defaulting to %s.", defaultUsername)
		return defaultUsername
	}

	// Here we go with Windows being funny
	if runtime.GOOS == "windows" {
		return withoutDomain(user.Username)
	}

	// Different branch naming conventions exist, but all usually start with the author's name.
	// The character after this often differs (":", "/", "-" are commonly used), so we won't include it in default pattern.
	// If needed, it can be included in the value provided for the `--pattern` flag.
	return user.Username
}

// This should only affect Windows, where the username format is DOMAIN\username
func withoutDomain(withPossibleDomain string) string {
	split := strings.Split(withPossibleDomain, "\\")
	if len(split) == 2 {
		return split[1]
	}

	return withPossibleDomain
}
