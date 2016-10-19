package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const configFilename = "$HOME/.gomuche/config.json"

// Config represents gomuche config file
type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NewConfig returns a new Config.
func NewConfig(clientID, clientSecret string) *Config {
	return &Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// NewConfigFromFile read config from file and returns it.
func NewConfigFromFile() *Config {
	filename := os.ExpandEnv(configFilename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	config := new(Config)

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

// SaveConfig saves config to file.
func SaveConfig(conf *Config) {
	bytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	filename := os.ExpandEnv(configFilename)
	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(filename, bytes, 0755)
	if err != nil {
		log.Fatalln(err)
	}
}
