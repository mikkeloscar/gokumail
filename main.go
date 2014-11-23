package main

import (
	"flag"

	"github.com/op/go-logging"
)

var (
	// Log global logger
	Log    = logging.MustGetLogger("logger")
	format = "%{level:.4s} %{message}"
)

// Conf global config
var Conf *ServerConfig

func main() {
	// config path
	var config string

	flag.StringVar(&config, "c", "/etc/gokumail.conf", "Config path")
	flag.Parse()

	// read config
	Conf = MustReadServerConfig(config)

	// setup logger
	logging.SetLevel(logging.INFO, "logger")
	logging.SetFormatter(logging.MustStringFormatter(format))

	// Run webinterface
	go RunWebInterface(Conf.HTTP.Port)

	// pop3 server
	POP3Server(Conf.POP.Port, Conf.POP.TLS)
}
