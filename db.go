package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Settings struct {
	User          string
	WorkMail      string
	FromWhitelist []string
	ToWhitelist   []string
	Whitelist     []string
	Blacklist     []string
}

func connect() (*sql.DB, error) {
	return sql.Open("mysql", "")
}

// GetSettings get settings for user
func GetSettings(user string) (*Settings, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM user_settings WHERE username=?", user)
	s := new(Settings)
	err = row.Scan(&s.User, &s.WorkMail, &s.FromWhitelist, &s.ToWhitelist, &s.Blacklist)

	return s, err
}

// func UpdateSettings(user string, s *Settings) {

// }
