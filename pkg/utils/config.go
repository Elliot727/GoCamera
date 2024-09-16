package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"GoCamera/internal/types"

	"gopkg.in/yaml.v2"
)

func LoadConfig() (*types.Config, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	configPath := filepath.Join(basePath, "..", "..", "cmd", "config.yaml")

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Config file not found, using default settings.")
			return &types.Config{
				Server:   false,
				Transfer: true,
				Organise: false,
				Source:   "",
				Dest:     "",
				Port:     "8080",
			}, nil
		}
		log.Printf("Error opening config file: %v\n", err)
		return nil, err
	}
	defer file.Close()

	var config types.Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Printf("Error decoding config file: %v\n", err)
		return nil, err
	}

	return &config, nil
}
