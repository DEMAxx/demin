package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var err error
	if len(cmd) == 0 {
		fmt.Println("Error: no command provided")
		return 1
	}

	safeCmd, err := sanitizeInput(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}

	command := exec.Command(safeCmd[0], safeCmd[1:]...) //nolint:gosec
	command.Env = os.Environ()

	for name, value := range env {
		if value.NeedRemove {
			command.Env = removeEnv(command.Env, name)
		} else {
			cleanValue := strings.ReplaceAll(value.Value, "\x00", "")
			if strings.Trim(name, " ") == "" {
				command.Env = append(command.Env, cleanValue)
			} else {
				command.Env = append(command.Env, name+"="+cleanValue)
			}
		}
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return 1
	}

	err = command.Wait()
	if err != nil {
		fmt.Println("Error waiting for command:", err)
		return 1
	}

	return command.ProcessState.ExitCode()
}

func sanitizeInput(input []string) ([]string, error) {
	for _, arg := range input {
		if strings.ContainsAny(arg, `;&|<>`) {
			return nil, errors.New("invalid input")
		}
	}
	return input, nil
}

func removeEnv(env []string, name string) []string {
	result := make([]string, 0, len(env))
	prefix := name + "="
	for _, e := range env {
		if len(e) > len(prefix) && e[:len(prefix)] != prefix {
			result = append(result, e)
		}
	}
	return result
}
