package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDb() error {
	log.Println("setting up database connection..")
	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		log.Fatalf("no database host provided, exiting..")
	}

	name := os.Getenv("DATABASE_NAME")
	if name == "" {
		log.Fatalf("no database name provided, exiting..")
	}

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		log.Fatalf("no database user provided, exiting..")
	}

	pass := os.Getenv("DATABASE_PASSWORD")
	if pass == "" {
		log.Fatalf("no database password provided, exiting..")
	}

	log.Printf("connecting to database: host - %s, name - %s, user - %s\n", host, name, user)

	var err error
	connStr := fmt.Sprintf("connect_timeout=10 dbname=%s host=%s user=%s password=%s", name, host, user, pass)
	Db, err = sql.Open("postgres", connStr)
	return err
}
