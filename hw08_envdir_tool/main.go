package main

func main() {
	env, err := ReadDir("./testdata/env")

	if err != nil {
		panic(err)
	}

	cmd := make([]string, len(env))

	runCommand := RunCmd(cmd, env)

	println("runCommand", runCommand)
}
