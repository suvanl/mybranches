package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/suvanl/mybranches/cleanup"
	"github.com/suvanl/mybranches/shared/git"
)

func main() {
	patternFlag := flag.String("pattern", getUsernamePattern(), "Custom pattern to use. Defaults to your username.")
	cleanupFlag := flag.Bool("cleanup", false, "Delete all local branches that aren't on remote")

	flag.Parse()

	if *cleanupFlag {
		hasCustomPatternFlag := *patternFlag != getUsernamePattern()
		if hasCustomPatternFlag {
			fmt.Println("  ⚠️ Ignoring --pattern flag because --cleanup flag was set to true")
		}

		// Run cleanup program
		program := tea.NewProgram(cleanup.InitialState())
		_, err := program.Run()

		if err != nil {
			fmt.Printf("Something went wrong: %v", err)
			os.Exit(1)
		}

		return
	}

	branches := git.FindBranches(*patternFlag)
	if len(branches) == 0 {
		fmt.Printf("Couldn't find any branches containing '%s'\n", *patternFlag)
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
		switchOut := git.SwitchBranch(selectedBranch)
		fmt.Printf("\n---\n\n%s\n", switchOut)
		return
	}

	if deletionRequested {
		deleteOut := git.DeleteBranch(uiModel.deletionContext.branchName)
		fmt.Printf("\n---\n\n%s\n", deleteOut)
		return
	}
}
