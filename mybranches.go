package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	branches       []string
	cursorIndex    int
	selectedBranch string
}

const deselectedIndicator = "( )"
const selectedIndicator = "(*)"

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
		formatHelpSection("q", "quit"),
	}

	return "\n" + strings.Join(sections, helpStyleSecondary.Render(" • ")) + "\n"
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
	return user.Username
}

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

func main() {
	patternFlag := flag.String("pattern", getUsernamePattern(), "Custom pattern to use. Defaults to your username.")
	flag.Parse()

	pattern := *patternFlag

	if strings.TrimSpace(pattern) == "" {
		fmt.Println("pattern flag cannot be blank")
		return
	}

	branches := findBranches(pattern)
	if len(branches) == 0 {
		fmt.Printf("Couldn't find any branches containing '%s'\n", pattern)
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
