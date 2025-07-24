package config

import (
	"encoding/json"
	"log"
	"os"
)

var config Config

type Config struct {
	CustomDataKeys map[string]string `json:"customDataKeys"`
	CLIArgs
}

type CLIArgs struct {
	InputDir  string
	OutputDir string
}

func ParseConfig(filePath string, cliArgs CLIArgs) (Config, error) {
	if filePath == "" {
		config = Config{}
		return config, nil
	}

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(byteValue, &config); err != nil {
		log.Fatal(err)
	}

	config.CLIArgs = cliArgs
	return config, nil
}

func Get() Config {
	return config
}
