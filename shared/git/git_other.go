//go:build !windows

package git

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func DeleteBranchesNotOnRemote() (int, error) {
	branches, err := getGoneBranches()
	if err != nil {
		return -1, err
	}

	// todo: Delete all the gone branches

	return len(branches), nil
}

func getGoneBranches() ([]string, error) {
	branch := exec.Command("git", "branch", "-vv")
	grep := exec.Command("grep", "'gone]'")
	awk := exec.Command("awk", "'{print $1}'")

	// Set up pipes between commands

	branchOut, err := branch.StdoutPipe()
	if err != nil {
		return nil, err
	}

	grepIn, err := grep.StdinPipe()
	if err != nil {
		return nil, err
	}

	grepOut, err := grep.StdoutPipe()
	if err != nil {
		return nil, err
	}

	awkIn, err := awk.StdinPipe()
	if err != nil {
		return nil, err
	}

	// Start commands

	if err := branch.Start(); err != nil {
		return nil, err
	}

	if err := grep.Start(); err != nil {
		return nil, err
	}

	if err := awk.Start(); err != nil {
		return nil, err
	}

	// Join pipes

	go func() {
		defer grepIn.Close()
		if _, err := bufio.NewReader(branchOut).WriteTo(grepIn); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer awkIn.Close()
		if _, err := bufio.NewReader(grepOut).WriteTo(awkIn); err != nil {
			log.Fatalln(err)
		}
	}()

	// Finish commands + collect output

	var awkOut bytes.Buffer
	awk.Stdout = &awkOut

	if err := branch.Wait(); err != nil {
		return nil, err
	}

	// ! fixme
	if err := grep.Wait(); err != nil {
		return nil, err
	}

	if err := awk.Wait(); err != nil {
		return nil, err
	}

	output := strings.TrimSpace(awkOut.String())
	branchNames := strings.Split(output, "\n")

	return branchNames, nil
}
