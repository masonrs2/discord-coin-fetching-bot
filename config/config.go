package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token string
	BotPrefix string
}

var (
	Token string
	BotPrefix string
	config *Config
)

func ConfigureBot() error {

	file, err := os.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Error opening config.json file")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error parsing config.json file")
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}