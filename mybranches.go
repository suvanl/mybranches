package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	branches       []string
	cursorIndex    int
	selectedBranch string
}

const deselectedIndicator = "( )"
const selectedIndicator = "(*)"

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
	fmt.Fprintf(&builder, "Branches containing '%s'\n\n", getUsernamePattern())

	for i := range m.branches {
		if m.cursorIndex == i {
			builder.WriteString(selectedIndicator + " ")
		} else {
			builder.WriteString(deselectedIndicator + " ")
		}
		builder.WriteString(m.branches[i])
		if getCurrentBranchName() == m.branches[i] {
			builder.WriteString(" (current)")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("\n(press q to quit)\n")

	return builder.String()
}

func initialState(branches []string) model {
	return model{
		branches: branches,
	}
}

func getUsernamePattern() string {
	user, err := user.Current()

	if err != nil {
		defaultUsername := "user"
		fmt.Printf("Failed to determine your username. Defaulting to %s.", defaultUsername)
		return defaultUsername
	}

	// The usage in mind is that branches will be named "<name>/mybranchname", but there's a chance
	// another char, such as `-` may be commonplace too, so we'll return the username without any
	// trailing chars for now.
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
	pattern := getUsernamePattern()
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
