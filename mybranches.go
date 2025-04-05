package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	pattern := flag.String("pattern", getUsernamePattern(), "Custom pattern to use. Defaults to your username.")
	flag.Parse()

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

	uiModel, ok := m.(model)
	selectedBranch := uiModel.selectedBranch
	deletionRequested := uiModel.deletionContext.shouldDelete

	if !ok {
		log.Fatal("m is not of type model")
		return
	}

	if selectedBranch != "" {
		switchOut := switchBranch(selectedBranch)
		fmt.Printf("\n---\n\n%s\n", switchOut)
		return
	}

	if deletionRequested {
		deleteOut := deleteBranch(uiModel.deletionContext.branchName)
		fmt.Printf("\n---\n\n%s\n", deleteOut)
		return
	}
}
