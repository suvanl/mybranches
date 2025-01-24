package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	branches       []string
	cursorIndex    int
	selectedBranch string
}

var (
	selectedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	currentStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	helpStylePrimary   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	helpStyleSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.branches) == 0 {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		// Standard quit keys
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursorIndex > 0 {
				m.cursorIndex--
			}

		case "down", "j":
			if m.cursorIndex < len(m.branches)-1 {
				m.cursorIndex++
			}

		case "c":
			err := handleCopy(m.branches[m.cursorIndex])
			if err != nil {
				if err == ErrClipboardNotSupported {
					// Do nothing
				}

				log.Fatal(err)
				return m, tea.Quit
			}

		case "enter", " ": // spacebar is represented by space char
			m.selectedBranch = m.branches[m.cursorIndex]
			return m, tea.Quit
		}
	}

	return m, nil
}

// The UI is just a string that gets updated by the Update method.
//
// There's no need to implement redrawing logic - bubbletea takes care of redrawing for us.
func (m model) View() string {
	const deselectedIndicator string = "( )"
	const selectedIndicator string = "(*)"

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "Branches containing '%s'\n\n", selectedStyle.Render(getUsernamePattern()))

	for i := range m.branches {
		if m.cursorIndex == i {
			builder.WriteString(selectedIndicator + " ")
		} else {
			builder.WriteString(deselectedIndicator + " ")
		}
		builder.WriteString(m.branches[i])
		if getCurrentBranchName() == m.branches[i] {
			builder.WriteString(currentStyle.Render(" (current)"))
		}

		builder.WriteString("\n")
	}

	builder.WriteString(buildHelpFooter())

	return builder.String()
}

func initialState(branches []string) model {
	return model{
		branches: branches,
	}
}

func formatHelpSection(key string, value string) string {
	return helpStylePrimary.Render(key) + " " + helpStyleSecondary.Render(value)
}

func buildHelpFooter() string {
	sections := []string{
		formatHelpSection("↑/k", "up"),
		formatHelpSection("↓/j", "down"),
		formatHelpSection("c", "copy"),
		formatHelpSection("q", "quit"),
	}

	return "\n" + strings.Join(sections, helpStyleSecondary.Render(" • ")) + "\n"
}

func handleCopy(text string) error {
	clipboard := getPlatformClipboard()
	if clipboard == nil {
		return ErrClipboardNotSupported
	}

	return clipboard.Copy(text)
}
