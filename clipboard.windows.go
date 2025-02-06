package main

import (
	"fmt"
	"os/exec"
)

type WindowsClipboard struct{}

func (c WindowsClipboard) Copy(text string) error {
	// Possibly another bit of Windows "fun". If we run `exec.Command("echo", text)`,
	// we get the error 'exec: "echo": executable file not found in %PATH%'.
	// Hopefully there's a better solution to this, but running it within a
	// PowerShell 7.x (pwsh) seems to fix things. Obviously, this won't work if the
	// user doesn't have pwsh installed (e.g. they may only have classic powershell).

	pwsh := exec.Command("pwsh", "-nologo", "-noprofile")
	pwshIn, pwshInErr := pwsh.StdinPipe()
	if pwshInErr != nil {
		return pwshInErr
	}

	go func() {
		defer pwshIn.Close()
		fmt.Fprintf(pwshIn, "echo '%s' | clip\r\n", text)
	}()

	_, pwshOutErr := pwsh.CombinedOutput()
	if pwshOutErr != nil {
		return pwshOutErr
	}

	return nil
}
