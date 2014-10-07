package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yml"
)

type Config struct {
	KU_work_mail   string
	From_whitelist []string
	To_whitelist   []string
	Whitelist      []string
	Blacklist      []string
}

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
	config.Whitelist = config.To_whitelist // TODO merge of from and to whitelist
	config.Blacklist = append(config.Blacklist, config.KU_work_mail)

	return &config, nil
}
