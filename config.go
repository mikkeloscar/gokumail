package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yml"
)

// Config defines a struct repesenting the config for the app
type Config struct {
	KUWorkMail    string   `yaml:"ku_work_mail"`
	FromWhitelist []string `yaml:"from_whitelist"`
	ToWhitelist   []string `yaml:"to_whitelist"`
	Whitelist     []string `yaml:"whitelist"`
	Blacklist     []string `yaml:"blacklist"`
}

// ReadConfig reads a config.yml file and returns a pointer to a Config struct
func ReadConfig() (*Config, error) {
	config := Config{}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// update lists
	config.Whitelist = config.ToWhitelist // TODO merge of from and to whitelist
	config.Blacklist = append(config.Blacklist, config.KUWorkMail)

	return &config, nil
}
