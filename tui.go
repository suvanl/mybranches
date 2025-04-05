package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI state model
type model struct {
	branches       []string
	cursorIndex    int
	selectedBranch string

	deletionContext struct {
		branchName   string
		shouldDelete bool
	}
}

var (
	// Style for reiterating the selected item
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	// Style for the '(current)' badge next to the name of the branch the user is currently on
	currentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	// Style for primary items in the help footer
	helpStylePrimary = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Style for secondary items in the help footer
	helpStyleSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.branches) == 0 {
		return m, tea.Quit
	}

	isInDeleteMode := m.deletionContext.branchName != ""
	if isInDeleteMode {
		return m.handleDeleteBranchViewUpdate(msg)
	}

	return m.handleMainViewUpdate(msg)
}

// The UI is just a string that gets updated by the Update method.
//
// There's no need to implement redrawing logic - bubbletea takes care of redrawing for us.
func (m model) View() string {
	if m.deletionContext.branchName != "" {
		return m.deleteBranchView()
	}

	return m.mainView()
}

func (m model) mainView() string {
	const deselectedIndicator string = "( )"
	const selectedIndicator string = "(*)"

	builder := strings.Builder{}

	// pattern flag defaults to username, so this should always have a value
	pattern := flag.Lookup("pattern").Value.String()

	fmt.Fprintf(&builder, "Branches containing '%s'\n\n", selectedStyle.Render(pattern))

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

func (m model) deleteBranchView() string {
	builder := strings.Builder{}

	if getCurrentBranchName() == m.deletionContext.branchName {
		fmt.Fprint(&builder, "\n:( Can't delete the branch you're currently on. Switch to a different branch first.\n\n")

		footerSections := []string{
			formatHelpSection("n", "go back"),
			formatHelpSection("q", "quit"),
		}
		builder.WriteString(strings.Join(footerSections, helpStyleSecondary.Render(" • ")) + "\n")

		return builder.String()
	}

	fmt.Fprintf(&builder, "\nDelete '%s'?\n\n", selectedStyle.Render(m.deletionContext.branchName))

	builder.WriteString(buildDeleteHelpFooter())

	return builder.String()
}

func (m model) handleMainViewUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case "d":
			m.deletionContext.branchName = m.branches[m.cursorIndex]

		case "enter", " ": // spacebar is represented by space char
			m.selectedBranch = m.branches[m.cursorIndex]
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) handleDeleteBranchViewUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "n":
			m.deletionContext.branchName = ""

		case "y":
			// Prevent sending the deletion request to git if it's guaranteed to never succeed.
			// Specifically, if the branch we're currently on is the branch we're trying to delete.
			isDeletable := getCurrentBranchName() != m.deletionContext.branchName
			if isDeletable {
				m.deletionContext.shouldDelete = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
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
		formatHelpSection("d", "delete"),
		formatHelpSection("q", "quit"),
	}

	os := runtime.GOOS
	if getPlatformClipboard(os) != nil {
		// force "copy" to be before "delete"
		sections = slices.Insert(sections, len(sections)-2, formatHelpSection("c", "copy"))
	}

	return "\n" + strings.Join(sections, helpStyleSecondary.Render(" • ")) + "\n"
}

func buildDeleteHelpFooter() string {
	sections := []string{
		formatHelpSection("y", "yes"),
		formatHelpSection("n", "no"),
		formatHelpSection("q", "quit"),
	}

	return strings.Join(sections, helpStyleSecondary.Render(" • ")) + "\n"
}

func handleCopy(text string) error {
	os := runtime.GOOS
	clipboard := getPlatformClipboard(os)
	if clipboard == nil {
		return ErrClipboardNotSupported
	}

	return clipboard.Copy(text)
}
