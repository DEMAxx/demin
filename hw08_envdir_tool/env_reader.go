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

var (
	ErrInvalidDir = errors.New("invalid directory")
	ErrEmptyDir   = errors.New("empty directory")
)

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
		go processFile(dir, file, osEnv, &wg)
	}

	wg.Wait()
	return osEnv, nil
}

func processFile(dir string, file os.DirEntry, osEnv Environment, wg *sync.WaitGroup) {
	defer wg.Done()
	var envValue EnvValue

	osFile, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		fmt.Println("Couldn't open file", err)
		return
	}
	defer closeFile(osFile)

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
		processLine(scanner, &envValue)
		mu.Lock()
		osEnv[file.Name()] = envValue
		mu.Unlock()
		break
	}

	if err := scanner.Err(); err != nil {
		log.Println("Scanner error:", err)
	}
}

func closeFile(osFile *os.File) {
	err := osFile.Close()
	if err != nil {
		fmt.Println("Couldn't close file")
	}
}

func processLine(scanner *bufio.Scanner, envValue *EnvValue) {
	nullByte := make([]byte, 1)
	nullByte[0] = 0
	text := string(bytes.ReplaceAll([]byte(scanner.Text()), nullByte, []byte("\n")))

	if len(strings.TrimSpace(scanner.Text())) > 0 {
		envValue.NeedRemove = false
	} else {
		envValue.NeedRemove = true
	}

	envValue.Value = text
}
