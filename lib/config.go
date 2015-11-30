package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	siteName         string
	author           string
	description      string
	port             int
	recentPostsCount int
	slogan           string
	baseURL          string
	*github
}

type github struct {
	repository string
	branch     string
}

func getConfig(configFile string) *config {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("Error opening the configuration file: %s", err)
	}

	decoder := json.NewDecoder(file)
	conf := &config{}

	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Printf("Error reading the configuration file: %s", err)
	}

	return conf
}
