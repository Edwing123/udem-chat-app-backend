package main

import (
	"database/sql"
	"fmt"
	"net/url"
)

func NewSQLServerDatabase(details ConnectionDetails) (*sql.DB, error) {
	query := url.Values{}
	query.Add("app name", "Nameless")
	query.Add("database", "Nameless")

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(details.User, details.Password),
		Host:     fmt.Sprintf("%s:%d", details.Host, details.Port),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
