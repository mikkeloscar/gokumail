package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type popState int

const (
	stateUnauthorized popState = iota
	stateTransaction
	stateUpdate
)

const (
	eol = "\r\n"
)

// POP3Server spawn a simple pop3 server which acts as a proxy to KUmail
func POP3Server(port int) {
	tcpPort := fmt.Sprintf(":%d", port)

	ln, err := net.Listen("tcp", tcpPort)
	if err != nil {
		Log.Error("listen error: " + err.Error())
	} else {
		Log.Info("POP3 server listening on port: %d", port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			Log.Error("accept error: " + err.Error())
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	kumailClient := new(KUmail)
	config, err := ReadConfig()
	if err != nil {
		Log.Error("config error: " + err.Error())
		return
	}

	var (
		state = stateUnauthorized
	)

	reader := bufio.NewReader(conn)

	writeClient(conn, "+OK simple KUmail POP3 -> IMAP proxy")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			Log.Error(err.Error())
			return
		}

		Log.Debug("-> %s", line)

		// Parse command
		cmd, args := readCommand(line)

		if cmd == "USER" && state == stateUnauthorized {
			// accept username and wait for PASS command
			username, _ := getSafeArgs(args, 0)
			kumailClient.User = username
			writeClient(conn, "+OK user accepted")
		} else if cmd == "PASS" && state == stateUnauthorized {
			pass, _ := getSafeArgs(args, 0)
			kumailClient.Pass = pass
			if kumailClient.Init(config) {
				writeClient(conn, "+OK pass accepted")
				state = stateTransaction
			} else {
				writeClient(conn, "-ERR Username or password incorrect!")
			}
		} else if cmd == "CAPA" && state == stateTransaction {
			writeClient(conn, "+OK Capability list follows")
			writeClient(conn, "UIDL")
			writeClient(conn, "USER")
			writeClient(conn, ".")
		} else if cmd == "UIDL" && state == stateTransaction {
			list, err := kumailClient.UIDL()
			if err != nil {
				kumailClient.Close()
				Log.Error(err.Error())
				writeClient(conn, "-ERR unable to perform UIDL")
				return
			}
			writeClient(conn, "+OK")
			// TODO buffer this stuff and send as a single command?
			for _, msg := range list {
				writeClient(conn, "%s %d", msg.ID, msg.UID)
			}
			writeClient(conn, ".")
		} else if cmd == "LIST" && state == stateTransaction {
			list, total, err := kumailClient.ListAll()
			if err != nil {
				kumailClient.Close()
				Log.Error(err.Error())
				writeClient(conn, "-ERR unable to perform LIST")
				return
			}
			writeClient(conn, "+OK %d messages (%d octets)", len(list), total)
			// TODO buffer this stuff and send as a single command?
			for _, msg := range list {
				writeClient(conn, "%s %d", msg.ID, msg.size)
			}
			writeClient(conn, ".")
		} else if cmd == "RETR" && state == stateTransaction {
			id, _ := getSafeArgs(args, 0)
			msg, octets, err := kumailClient.GetMessage(id)
			if err != nil {
				kumailClient.Close()
				Log.Error(err.Error())
				writeClient(conn, "-ERR no such message")
				return
			}

			writeClient(conn, "+OK %d octets", octets)
			// send message
			fmt.Fprintf(conn, msg+eol)
			writeClient(conn, ".")
		} else if cmd == "DELE" && state == stateTransaction {
			writeClient(conn, "-ERR you are not allowed to delete messages on this server")
		} else if cmd == "QUIT" {
			// take down IMAP connection
			kumailClient.Close()
			writeClient(conn, "+OK Bye bye!")
			return
		} else {
			writeClient(conn, "-ERR invalid command")
		}
	}
}

// read commands send by client
func readCommand(line string) (string, []string) {
	line = strings.Trim(line, "\r \n")
	cmd := strings.Split(line, " ")
	return cmd[0], cmd[1:]
}

func getSafeArgs(args []string, n int) (string, error) {
	if n < len(args) {
		return args[n], nil
	}
	return "", errors.New("out of range")
}

// write message to client and print the message in the server log
func writeClient(conn net.Conn, msg string, args ...interface{}) {
	fmt.Fprintf(conn, msg+eol, args...)
	Log.Debug("<- "+msg+eol, args...)
}
