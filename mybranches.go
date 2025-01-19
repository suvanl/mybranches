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
	branches    []string
	cursorIndex int
	selected    map[int]struct{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			_, ok := m.selected[m.cursorIndex]
			if ok {
				// Remove the entry from the map, it should no longer be selected
				delete(m.selected, m.cursorIndex)
			} else {
				m.selected[m.cursorIndex] = struct{}{}
			}
		}
	}

	return m, nil
}

// The UI is just a string that gets updated by the Update method.
//
// There's no need to implement redrawing logic - bubbletea takes care of redrawing for us.
func (m model) View() string {
	s := fmt.Sprintf("Branches containing '%s'\n\n", getUsernamePattern())

	for i, branch := range m.branches {
		// Render cursor if the current item is selected
		cursor := " "
		if m.cursorIndex == i {
			cursor = ">"
		}

		// Render checkmark (x) for selected choice
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Render row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, branch)
	}

	s += "\nPress q to quit\n"

	return s
}

func main() {
	program := tea.NewProgram(initialState())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}

func initialState() model {
	return model{
		branches: findBranches(getUsernamePattern()),
		selected: make(map[int]struct{}),
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

	strFromBytes := string(out[:])
	branches := strings.Split(strFromBytes, "\n")

	// Last element will be an empty string, let's just drop it here
	return branches[:len(branches)-1]
}
