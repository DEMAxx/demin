package main

import (
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/env_reader" // Import the env_reader package
	"log"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Server ServerConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(fileOrDir string) Config {
	if fileOrDir == "" {
		fileOrDir = "../../configs"
	}

	fileInfo, err := os.Stat(fileOrDir)

	if err != nil {
		log.Fatalf("Error stating file or directory: %v", err)
	}

	var env env_reader.Environment

	if fileInfo.IsDir() {
		env, err = env_reader.ReadDir(fileOrDir)
	} else {
		env, err = env_reader.ReadFile(fileOrDir)
	}

	return Config{
		Logger: LoggerConf{
			Level: env["level"].Value,
		},
		Server: ServerConf{
			Host: env["host"].Value,
			Port: env["port"].Value,
		},
	}
}
