package config

import (
	"log"
	"os"

	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/envreader"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	Server ServerConf
	DB     DBConf
}

type LoggerConf struct {
	Level  string
	Output string
}

type ServerConf struct {
	Host     string
	Port     string
	GrpcPort string
}

type DBConf struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func NewConfig(fileOrDir string) Config {
	if fileOrDir == "" {
		fileOrDir = "../../configs"
	}

	fileInfo, err := os.Stat(fileOrDir)
	if err != nil {
		log.Fatalf("Error stating file or directory: %v", err)
	}

	var env envreader.Environment

	if fileInfo.IsDir() {
		env, err = envreader.ReadDir(fileOrDir)
	} else {
		env, err = envreader.ReadFile(fileOrDir)
	}

	if err != nil {
		log.Fatalf("Error reading environment: %v", err)
	}

	return Config{
		Logger: LoggerConf{
			Level:  env["level"].Value,
			Output: env["output"].Value,
		},
		Server: ServerConf{
			Host:     env["host"].Value,
			Port:     env["port"].Value,
			GrpcPort: env["grpc_port"].Value,
		},
		DB: DBConf{
			User:     env["user"].Value,
			Password: env["password"].Value,
			Host:     env["db_host"].Value,
			Port:     env["db_port"].Value,
			Name:     env["name"].Value,
		},
	}
}
