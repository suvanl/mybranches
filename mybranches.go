package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	pattern := flag.String("pattern", getUsernamePattern(), "Custom pattern to use. Defaults to your username.")
	flag.Parse()

	if strings.TrimSpace(*pattern) == "" {
		fmt.Println("pattern flag cannot be blank")
		return
	}

	branches := findBranches(*pattern)
	if len(branches) == 0 {
		fmt.Printf("Couldn't find any branches containing '%s'\n", *pattern)
		return
	}

	program := tea.NewProgram(initialState(branches))
	m, err := program.Run()

	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok && m.selectedBranch != "" {
		switchOut := switchBranch(m.selectedBranch)
		fmt.Printf("\n---\n\n%s\n", switchOut)
	}
}
