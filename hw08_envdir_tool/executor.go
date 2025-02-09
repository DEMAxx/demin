package main

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for _, v := range cmd {
		println("cmd", v)
	}

	for _, v := range env {
		println("env", v.Value)
	}

	return
}
