package main

import (
	"github.com/op/go-logging"
)

const (
	pop3Port = 1100
)

// Log global logger
var Log = logging.MustGetLogger("logger")

var format = "%{level:.4s} %{message}"

func main() {
	// setup logger
	logging.SetLevel(logging.INFO, "logger")
	logging.SetFormatter(logging.MustStringFormatter(format))

	// pop3 server
	POP3Server(pop3Port)
}
