package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get config file directory: %v", err)
	}

	fileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %v", err)
	}

	var config Config
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return config, nil
}

func (config *Config) SetUser(user string) error {
	config.CurrentUserName = user
	configMarshal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config struct: %v", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get config file directory: %v", err)
	}

	err = os.WriteFile(configFilePath, configMarshal, 0600)
	if err != nil {
		return fmt.Errorf("failed to write to config file: %v", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	return userHomeDir + "/" + configFileName, nil
}
