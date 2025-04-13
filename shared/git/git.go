package git

import "os/exec"

// git fetch --prune
func FetchPrune() error {
	_, err := exec.Command("git", "fetch", "--prune").CombinedOutput()
	return err
}
