package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Println("error cmd")
		return 1
	}

	for name, value := range env {
		var err error
		if value.NeedRemove {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, value.Value)
		}

		if err != nil {
			return 1
		}
	}

	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
