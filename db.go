package main

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const table = "user_settings"

var usernameRe = regexp.MustCompile(`^[b-df-hj-np-tv-xz]{3}\d{3}$`)

// Settings user_settings
type Settings struct {
	User          string
	Workmail      string
	FromWhitelist []string
	ToWhitelist   []string
	Blacklist     []string
}

// Whitelist a combined list of FromWhitelist and ToWhitelist
func (s *Settings) Whitelist() []string {
	return s.ToWhitelist
}

func connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@/%s", Conf.DB.User, Conf.DB.Pass, Conf.DB.DBname)
	return sql.Open("mysql", dsn)
}

// GetSettings get settings for user
func GetSettings(user string) (*Settings, error) {
	if !usernameRe.MatchString(user) {
		return nil, errors.New("invalid ku-username format")
	}

	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var from string
	var to string
	var blacklist string

	stmt := fmt.Sprintf("SELECT username, workmail, fromwhitelist, towhitelist, blacklist FROM %s WHERE username=?", table)

	row := db.QueryRow(stmt, user)
	s := new(Settings)
	err = row.Scan(&s.User, &s.Workmail, &from, &to, &blacklist)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	s.FromWhitelist = splitWithoutEmpty(from, ";")
	s.ToWhitelist = splitWithoutEmpty(to, ";")
	s.Blacklist = splitWithoutEmpty(blacklist, ";")

	return s, err
}

// Create new user_settings entry in DB
func (s *Settings) Create() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt := fmt.Sprintf("INSERT INTO %s (username, workmail, fromwhitelist, towhitelist, blacklist) VALUES (?, ?, ?, ?, ?)", table)
	_, err = db.Exec(
		stmt,
		s.User,
		s.Workmail,
		joinWithoutEmpty(s.FromWhitelist, ";"),
		joinWithoutEmpty(s.ToWhitelist, ";"),
		joinWithoutEmpty(s.Blacklist, ";"))

	return err
}

// Update user settings
func (s *Settings) Update() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt := fmt.Sprintf("UPDATE %s SET workmail=?, fromwhitelist=?, towhitelist=?, blacklist=? WHERE username=?", table)

	_, err = db.Exec(
		stmt,
		s.Workmail,
		joinWithoutEmpty(s.FromWhitelist, ";"),
		joinWithoutEmpty(s.ToWhitelist, ";"),
		joinWithoutEmpty(s.Blacklist, ";"),
		s.User)

	return err
}

// join elements of a into a single string seperated by sep. Ignore empty
// strings in a
func joinWithoutEmpty(a []string, sep string) string {
	joined := ""
	for _, s := range a {
		if s != "" {
			if joined != "" {
				joined += sep
			}
			joined += s
		}
	}
	return joined
}

// split string at sep into []string not containing any strings of length 0
func splitWithoutEmpty(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, sep)
}
