package cleanup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suvanl/mybranches/shared/git"
)

type errMsg struct{ err error }

func (e errMsg) Error() string {
	return e.err.Error()
}

type cleanUpStage int

const (
	FetchPrune cleanUpStage = iota
	Delete
	Done
)

func startStage(stage cleanUpStage, dryRun bool) tea.Cmd {
	switch stage {
	case FetchPrune:
		return func() tea.Msg {
			err := git.FetchPrune()
			if err != nil {
				return errMsg{err}
			}
			return getNextStage(stage)
		}

	case Delete:
		return func() tea.Msg {
			_, err := git.DeleteBranchesNotOnRemote(dryRun)
			if err != nil {
				return errMsg{err}
			}

			return getNextStage(stage)
		}
	}

	return nil
}

func getNextStage(currentStage cleanUpStage) cleanUpStage {
	var nextStage cleanUpStage

	switch currentStage {
	case FetchPrune:
		nextStage = Delete
	case Delete:
		nextStage = Done
	case Done:
		nextStage = Done
	}

	return cleanUpStage(nextStage)
}
