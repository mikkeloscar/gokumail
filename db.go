package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const table = "user_settings"

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
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt := "SELECT username, workmail, fromwhitelist, towhitelist, blacklist FROM ? WHERE username=?"

	row := db.QueryRow(stmt, table, user)
	s := new(Settings)
	err = row.Scan(&s.User, &s.Workmail, &s.FromWhitelist, &s.ToWhitelist, &s.Blacklist)

	return s, err
}

// Create new user_settings entry in DB
func (s *Settings) Create() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt := "INSERT INTO ? (username, workmail, fromwhitelist, towhitelist, blacklist) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(
		stmt,
		table,
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

	stmt := "UPDATE ? SET workmail=?, fromwhitelist=?, towhitelist=?, blacklist=? WHERE username=?"

	_, err = db.Exec(
		stmt,
		table,
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
	for i, s := range a {
		if s != "" {
			joined += s
			if i > len(a)-1 {
				joined += sep
			}
		}
	}
	return joined
}
