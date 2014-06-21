package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SiteName         string
	Author           string
	Description      string
	Port             int
	RecentPostsCount int
	Slogan           string
	BaseURL          string
	*Github
}

type Github struct {
	Repository string
	Branch     string
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
