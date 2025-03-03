package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var ErrInvalidDir = errors.New("invalid directory")
var ErrEmptyDir = errors.New("empty directory")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var mu sync.Mutex

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	osDir, err := os.ReadDir(dir)

	if err != nil {
		return nil, ErrInvalidDir
	}

	if len(osDir) == 0 {
		return nil, ErrEmptyDir
	}

	wg := sync.WaitGroup{}
	osEnv := make(Environment)

	for _, file := range osDir {
		wg.Add(1)

		go func() {
			defer wg.Done()
			var envValue EnvValue

			osFile, err := os.Open(dir + "/" + file.Name())

			if err != nil {
				fmt.Println("Couldn't open file", err)
				return
			}

			defer func(osFile *os.File) {
				err := osFile.Close()
				if err != nil {
					fmt.Println("Couldn't close file")
				}
			}(osFile)

			scanner := bufio.NewScanner(osFile)

			if file.Name() == "UNSET" {
				envValue.NeedRemove = true
				envValue.Value = ""
				mu.Lock()
				osEnv[file.Name()] = envValue
				mu.Unlock()
				return
			}

			for scanner.Scan() {

				if file.Name() == "!FOO" {
					newScanBytes := make([]byte, len(scanner.Text()))

					for _, vb := range scanner.Bytes() {
						if vb == 0 {
							newScanBytes = append(newScanBytes, 32)
							newScanBytes = append(newScanBytes, 32)
							newScanBytes = append(newScanBytes, 32)
						} else {
							newScanBytes = append(newScanBytes, vb)
						}
					}
				}

				nullByte := make([]byte, 1)

				nullByte[0] = 0

				text := string(bytes.Replace([]byte(scanner.Text()), nullByte, []byte("\n"), -1))

				if len(strings.TrimSpace(scanner.Text())) > 0 {
					envValue.NeedRemove = false
				} else {
					envValue.NeedRemove = true
				}

				envValue.Value = text

				mu.Lock()
				osEnv[file.Name()] = envValue
				mu.Unlock()
				return
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	wg.Wait()

	return osEnv, nil
}
