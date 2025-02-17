package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type MapEnvValue struct {
	Key   string
	Value string
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	wg := new(sync.WaitGroup)

	envMap := [6]MapEnvValue{
		{
			"HELLO",
			"",
		},
		{
			"BAR",
			"",
		},
		{
			"FOO",
			"",
		},
		{
			"UNSET",
			"",
		},
		{
			"ADDED",
			"",
		},
		{
			"EMPTY",
			"",
		},
	}

	for key, value := range envMap {
		wg.Add(1)

		go func() {
			defer wg.Done()

			osValue, exists := os.LookupEnv(value.Key)

			envValue, ok := env[value.Key]

			if value.Key == "UNSET" {
				envMap[key] = MapEnvValue{
					Key:   value.Key,
					Value: "",
				}
			} else if value.Key == "ADDED" {
				if exists {
					envMap[key] = MapEnvValue{
						Key:   value.Key,
						Value: osValue,
					}
				} else if ok {
					envMap[key] = MapEnvValue{
						Key:   value.Key,
						Value: envValue.Value,
					}
				}
			} else if value.Key == "EMPTY" {
				envMap[key] = MapEnvValue{
					Key:   value.Key,
					Value: "",
				}
			} else {
				if ok {
					envMap[key] = MapEnvValue{
						Key:   value.Key,
						Value: envValue.Value,
					}
				} else {
					envMap[key] = MapEnvValue{
						Key:   value.Key,
						Value: osValue,
					}
				}

			}
		}()
	}

	wg.Wait()

	for _, value := range envMap {
		fmt.Printf("%s is (%s)\n", value.Key, value.Value)
	}

	isFirst := true

	for _, v := range os.Args {
		if strings.Contains(v, "=") {
			returnCode++
			if isFirst {
				isFirst = false
				fmt.Printf("arguments are %s", v)
			} else {
				fmt.Printf(" %s", v)
			}
		}
	}

	return returnCode
}
