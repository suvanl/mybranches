package main

import "os/exec"

type DarwinClipboard struct{}

func (c DarwinClipboard) Copy(text string) error {
	echo := exec.Command("echo", text)
	pbcopy := exec.Command("pbcopy")

	pipe, pipeErr := echo.StdoutPipe()
	if pipeErr != nil {
		return pipeErr
	}

	defer pipe.Close()

	pbcopy.Stdin = pipe

	startErr := echo.Start()
	if startErr != nil {
		return startErr
	}

	pbcopy.Output()

	waitErr := echo.Wait()
	if waitErr != nil {
		return waitErr
	}

	return nil
}
