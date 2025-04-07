package cleanup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type stageError string

func (e stageError) Error() string {
	return string(e)
}

const errStageMappingError = stageError("Failed to map CleanupStage to message")

type model struct {
	spinner  spinner.Model
	stage    cleanUpStage
	quitting bool
	err      error
}

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func InitialState() model {
	spinnerInit := spinner.New()
	spinnerInit.Spinner = spinner.Dot
	spinnerInit.Style = spinnerStyle

	return model{
		spinner: spinnerInit,
		stage:   FetchPrune,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, startStage(m.stage))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case cleanUpStage:
		m.stage = msg

		if msg == Done {
			return m, tea.Quit
		}

		return m, startStage(m.stage)

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nSomething went wrong: %v\n\n", m.err)
	}

	builder := strings.Builder{}

	stageMsg, mappingErr := mapStageToMessage(m.stage)
	if mappingErr != nil {
		return mappingErr.Error()
	}

	fmt.Fprintf(&builder, "\n  %s %s... ", m.spinner.View(), stageMsg)
	fmt.Fprint(&builder, helpStyle.Render("(q to quit)"))

	if m.quitting {
		builder.WriteString("\n")
		return builder.String()
	}

	return builder.String()
}

func mapStageToMessage(stage cleanUpStage) (string, error) {
	var message string

	switch stage {
	case FetchPrune:
		message = "Fetching remote branches"
	case Find:
		message = "Finding local branches not on remote"
	case Delete:
		message = "Deleting local branches not on remote"
	case Done:
		message = "Done"
	default:
		return "", errStageMappingError
	}

	return message, nil
}
