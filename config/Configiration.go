package config

import (
	"os"
	"strconv"
)

const (
	MocksRootDirEnvKey string = "MOCKER_MOCKS_ROOT_DIR"
	ServerPortEnvKey   string = "MOCKER_SERVER_PORT"
	LogPathEnvKey      string = "MOCKER_LOG_PATH"
)

// Config contains configuration for server instance
type Config struct {
	MocksRootDir string
	Port         int
	LogsPath     string
}

// EnvOrCurrent Returns environment variable with key name value or def value
func EnvOrCurrent(key string, def string) string {

	env := os.Getenv(key)

	if len(env) == 0 {
		return def
	}

	return env
}

// LoadConfig load config by filepath
func LoadConfig() (Config, error) {
	mocksRootDir := EnvOrCurrent(MocksRootDirEnvKey, "sandbox/mocs")
	port, err := strconv.Atoi(EnvOrCurrent(ServerPortEnvKey, "1111"))
	logsPath := EnvOrCurrent(LogPathEnvKey, "sandbox/logs.json")

	if err != nil {
		return Config{}, err
	}

	return Config{MocksRootDir: mocksRootDir, Port: port, LogsPath: logsPath}, nil
}
