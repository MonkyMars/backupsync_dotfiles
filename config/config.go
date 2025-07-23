// Package config contains the logic for implementing configuration into the tool.
package config

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed config.json
var configData []byte

type Config struct {
	DriveName            string
	SourceDirectory      string `json:"source_directory"`
	DestinationDirectory string `json:"destination_directory"`
}

func ParseConfig() (Config, error) {
	var config Config
	err := json.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, err
	}

	config.DriveName = getDriveName(config.DestinationDirectory)

	return config, nil
}

func getDriveName(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "mnt" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
