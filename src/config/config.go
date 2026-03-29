package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ScreenWidth   int32  `json:"screenWidth"`
	ScreenHeight  int32  `json:"screenHeight"`
	LogicalWidth  int32  `json:"logicalWidth"`
	LogicalHeight int32  `json:"logicalHeight"`
	FPSLimit      int32  `json:"fpsLimit"`
	Title         string `json:"title"`
	Debug         bool   `json:"debug"`
}

func LoadConfig() (*Config, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve executable path: %w", err)
	}

	filePath := filepath.Join(filepath.Dir(exe), "..", "config.json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err = json.Unmarshal(jsonData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func DefaultConfig() *Config {
	return &Config{
		ScreenWidth:   1280,
		ScreenHeight:  720,
		LogicalWidth:  1280,
		LogicalHeight: 720,
		FPSLimit:      60,
		Title:         "Coldline Miami",
		Debug:         false,
	}
}
