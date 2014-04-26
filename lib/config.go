package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SiteName         string
	Port             string
	BaseURL          string
	RecentPostsCount int
	Slogan           string
}

func GetConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Errorf("Error opening the configuration file: %s", err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		fmt.Errorf("Error reading the configuration file: %s", err)
	}

	return *config
}
