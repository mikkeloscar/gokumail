package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
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
		log.Fatalf("listen error: %v", err)
	} else {
		fmt.Printf("Server listening on port: %d\n", port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("accept error: %v", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	kumailClient := new(KUmail)
	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Error: " + err.Error())
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
			fmt.Println("Error: " + err.Error())
			return
		}

		log.Printf("-> %s", line)

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
			writeClient(conn, "+OK")
			writeClient(conn, ".")
		} else if cmd == "LIST" && state == stateTransaction {
			list, total, err := kumailClient.ListAll()
			if err != nil { // TODO better way of handling these kind of errors
				kumailClient.Close()
				fmt.Printf("error: %s\n", err)
				return
			}
			writeClient(conn, "+OK %d messages (%d octets)", len(list), total)
			// TODO buffer this stuff and send as a single command?
			for _, msg := range list {
				writeClient(conn, "%s %d", msg.ID, msg.size)
			}
			writeClient(conn, ".")
		} else if cmd == "RETR" && state == stateTransaction {
			writeClient(conn, "+OK SEND MSG")
		} else if cmd == "DELE" && state == stateTransaction {
			msgID, _ := getSafeArgs(args, 0)
			writeClient(conn, "+OK message %s deleted (NOT IMPLEMENTED)", msgID)
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
		return args[0], nil
	}
	return "", errors.New("Out of range")
}

// write message to client and print the message in the server log
func writeClient(conn net.Conn, msg string, args ...interface{}) {
	fmt.Fprintf(conn, msg+eol, args...)
	log.Printf("<- "+msg+eol, args...)
}
