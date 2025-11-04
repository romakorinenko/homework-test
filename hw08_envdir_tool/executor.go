package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}

	envs := make([]string, 0)
	for k, v := range env {
		if _, ok := os.LookupEnv(k); ok {
			if err := os.Unsetenv(k); err != nil {
				return 1
			}
		}
		if !v.NeedRemove {
			varString := fmt.Sprintf("%s=%s", k, v.Value)
			envs = append(envs, varString)
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = append(os.Environ(), envs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		return 1
	}

	return 0
}
