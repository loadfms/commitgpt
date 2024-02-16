package services

import (
	"fmt"
	"os"
	"os/user"

	"github.com/loadfms/commitgpt/models"
	"github.com/pelletier/go-toml/v2"
)

type FileConfig struct {
	ApiKey struct {
		Key string `toml:"key"`
	} `toml:"apikey"`
	Prompt struct {
		Custom string `toml:"custom"`
	} `toml:"prompt"`
}

type ConfigService struct {
	ApiKey       string
	CustomPrompt string
}

func NewConfigService() (svc *ConfigService) {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Config file not found")
	}

	return &ConfigService{ApiKey: cfg.ApiKey.Key, CustomPrompt: cfg.Prompt.Custom}
}

func SaveConfigFile(filePath string, cfg FileConfig) error {
	// Create or open the config file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Error opening/creating file: %v", err)
	}
	defer file.Close()

	// Marshal the Config struct into TOML format
	tomlData, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error marshaling TOML data: %v", err)
	}

	// Write the TOML data to the file
	_, err = file.Write(tomlData)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

func loadConfig() (result FileConfig, err error) {
	currentUser, err := user.Current()
	if err != nil {
		return result, fmt.Errorf("Error getting current user")
	}

	file, err := os.Open(currentUser.HomeDir + models.CONFIG_FOLDER + models.FILENAME)
	if err != nil {
		return result, fmt.Errorf("Error opening TOML file: %v", err)
	}
	defer file.Close()

	// Unmarshal the TOML content into a struct
	if err := toml.NewDecoder(file).Decode(&result); err != nil {
		return result, fmt.Errorf("Error parsing TOML file: %v", err)
	}

	return result, nil
}
