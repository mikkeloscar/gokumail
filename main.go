package main

import "github.com/op/go-logging"

var (
	// Log global logger
	Log    = logging.MustGetLogger("logger")
	format = "%{level:.4s} %{message}"
)

// Conf global config
var Conf = MustReadServerConfig("gokumail.conf")

func main() {
	// setup logger
	logging.SetLevel(logging.INFO, "logger")
	logging.SetFormatter(logging.MustStringFormatter(format))

	// pop3 server
	POP3Server(Conf.POP.Port)
}
