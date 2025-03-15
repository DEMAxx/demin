package env_reader

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

type EnvValue struct {
	Value      string
	NeedRemove bool
}

var mu sync.Mutex

func ReadDir(dir string) (Environment, error) {
	log.Println("ReadDir")
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

		osFile, err := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))

		if err != nil {
			fmt.Println("Couldn't open file", err)
			continue
		}

		go processFile(osFile, osEnv, &wg)
	}

	wg.Wait()

	fmt.Println("osEnv: ", osEnv)

	return osEnv, nil
}

func ReadFile(file string) (Environment, error) {
	osEnv := make(Environment)
	wg := sync.WaitGroup{}

	wg.Add(1)

	osFile, err := os.OpenFile(file, os.O_RDONLY, 0644)

	if err != nil {
		return nil, err
	}

	processFile(osFile, osEnv, &wg)

	wg.Wait()

	return osEnv, nil
}

func processFile(osFile *os.File, osEnv Environment, wg *sync.WaitGroup) {
	defer wg.Done()
	defer closeFile(osFile)

	scanner := bufio.NewScanner(osFile)

	for scanner.Scan() {
		processLine(scanner, osEnv)
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

func processLine(scanner *bufio.Scanner, osEnv Environment) {
	line := scanner.Text()
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	nullByte := make([]byte, 1)
	nullByte[0] = 0
	value = string(bytes.ReplaceAll([]byte(value), nullByte, []byte("\n")))

	envValue := EnvValue{
		Value:      value,
		NeedRemove: len(value) == 0,
	}

	mu.Lock()
	osEnv[key] = envValue
	mu.Unlock()
}
