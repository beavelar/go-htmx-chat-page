package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb(host string, name string, user string, pass string) error {
    log.Printf("initializing db connection - host: %s, name: %s, user: %s\n", host, name, user)

    var err error
    connStr := fmt.Sprintf("connect_timeout=10 dbname=%s host=%s user=%s password=%s", name, host, user, pass)
	Db, err = sql.Open("postgres", connStr)
    return err
}
