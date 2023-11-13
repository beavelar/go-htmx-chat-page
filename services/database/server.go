package main

import (
	"log"
	"os"
)

func main() {
    dbHost := os.Getenv("DATABASE_HOST")
    if dbHost == "" {
        log.Println("no database host provided, exiting..")
        os.Exit(1)
    }

    dbName := os.Getenv("DATABASE_NAME")
    if dbName == "" {
        log.Println("no database name provided, exiting..")
        os.Exit(1)
    }

    dbUser := os.Getenv("DATABASE_USER")
    if dbUser == "" {
        log.Println("no database user provided, exiting..")
        os.Exit(1)
    }

    dbPass := os.Getenv("DATABASE_PASSWORD")
    if dbPass == "" {
        log.Println("no database password provided, exiting..")
        os.Exit(1)
    }

    err := InitDb(dbHost, dbName, dbUser, dbPass)
    if err != nil {
        log.Printf("failed to create database session, exiting - %s\n", err)
        os.Exit(1)
    }
    defer Db.Close()
}
