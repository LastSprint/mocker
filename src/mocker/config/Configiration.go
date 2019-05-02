package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Config contains configuration for server instance
type Config struct {
	MocksRootDir string `json:"mocksRootDir"`
	Port         int    `json:"Port"`
}

// LoadConfig load config by filepath
func LoadConfig(filePath string) (Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var conf = Config{}

	err = json.Unmarshal(data, &conf)

	log.Println(conf)

	if err != nil {
		return Config{}, err
	}
	var abs string
	abs, err = filepath.Abs(conf.MocksRootDir)

	if err != nil {
		return Config{}, err
	}

	conf.MocksRootDir = abs
	return conf, nil
}
