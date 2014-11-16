package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	configFile = "config.yml"
)

// ServerConfig defining configuration for pop, imap
type ServerConfig struct {
	POP  pop
	IMAP imapClient
	DB   db
	HTTP httpClient
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

type httpClient struct {
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
