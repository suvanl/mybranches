package git

import "os/exec"

// git fetch --prune
func FetchPrune() error {
	// TODO: add prune back!
	_, err := exec.Command("git", "fetch").CombinedOutput()
	return err
}
