package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yml"
)

// ServerConfig defining configuration for pop, imap
type ServerConfig struct {
	POP  pop
	IMAP imapClient
	DB   db
	HTTP http
}

type pop struct {
	Port int
	TLS  bool
	Cert string
	Key  string
}

type imapClient struct {
	Server      string
	Port        int
	UsernameFmt string `toml:"username_fmt"`
	AddressFmt  string `toml:"address_fmt"`
	Folder      string
}

type db struct {
	DBname string
	User   string
	Pass   string
	Host   string
	Port   int
}

type http struct {
	Port int
}

// MustReadServerConfig from path
func MustReadServerConfig(path string) *ServerConfig {
	config, err := ReadServerConfig(path)
	if err != nil {
		panic("unable to read config: " + err.Error())
	}
	return config
}

// ReadServerConfig from path
func ReadServerConfig(path string) (*ServerConfig, error) {
	var config = ServerConfig{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

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
