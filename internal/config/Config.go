package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	DBURL       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := homeDir + "/" + configFileName
	return path, nil
}

func write(cfg Config) error {
	jsonPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(jsonPath, cfgBytes, 0644)

}

func Read() (Config, error) {

	jsonPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(jsonPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUser = username
	return write(*cfg)
}
