package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}

	cmds := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	envs := make([]string, 0)
	for k, v := range env {
		if _, ok := os.LookupEnv(k); ok {
			if err := os.Unsetenv(k); err != nil {
				return 1
			}
		}
		if !v.NeedRemove {
			s := fmt.Sprintf("%s=%s", k, v.Value)
			envs = append(envs, s)
		}
	}

	cmds.Env = append(os.Environ(), envs...)
	cmds.Stdin = os.Stdin
	cmds.Stdout = os.Stdout
	cmds.Stderr = os.Stderr
	if err := cmds.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cmds.ProcessState.ExitCode()
}
