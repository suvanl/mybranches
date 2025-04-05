package cleanup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cleanupStage int
type stageError string

const errStageMappingError = stageError("Failed to map CleanupStage to message")

func (e stageError) Error() string {
	return string(e)
}

const (
	FetchPrune cleanupStage = iota
	Find
	Delete
)

type model struct {
	spinner  spinner.Model
	stage    cleanupStage
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
	return m.spinner.Tick
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

	case error:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
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

func mapStageToMessage(stage cleanupStage) (string, error) {
	var message string

	switch stage {
	case FetchPrune:
		message = "Fetching remote branches"
	case Find:
		message = "Finding local branches not on remote"
	case Delete:
		message = "Deleting local branches not on remote"
	default:
		return "", errStageMappingError
	}

	return message, nil
}
