package main

import (
	proto "database-service/genproto/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func CloseDb() {
	if db != nil {
		db.Close()
	}
}

func GetMessages(limit int32) (*proto.Messages, error) {
	sqlStr := "SELECT * FROM messages ORDER BY time DESC"
	if limit > 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", limit)
	}

	if db == nil {
		msg := "unable to get all messages, database connection has not been initialized"
		log.Println(msg)
		return nil, errors.New(msg)
	}

	rows, err := db.Query(sqlStr)
	if err != nil {
		log.Printf("error occurred querying for all messages with limit: %d - %s\n", limit, err)
		return nil, err
	}
	defer rows.Close()

	msgs := make([]*proto.Message, 0)
	for rows.Next() {
		msg := &proto.Message{}
		if err := rows.Scan(&msg.Name, &msg.Message, &msg.Time); err != nil {
			log.Printf("error occurred scanning one of the rows from query response - %s\n", err)
			return &proto.Messages{Messages: msgs}, err
		}
		msgs = append(msgs, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error occurred scanning the rows from the query response - %s\n", err)
		return &proto.Messages{Messages: msgs}, err
	}
	return &proto.Messages{Messages: msgs}, nil
}

func PostMessage(msg *proto.Message) error {
	if db == nil {
		msg := "unable to post message, database connection has not been initialized"
		log.Println(msg)
		return errors.New(msg)
	}

	_, err := db.Exec(fmt.Sprintf("INSERT INTO messages(message, name, time) VALUES ('%s', '%s', %d)", msg.Message, msg.Name, msg.Time))
	if err != nil {
		log.Printf("error occurred attempting to insert message - message: %s, name: %s, time: %d - %s\n", msg.Message, msg.Name, msg.Time, err)
		return err
	}
	return nil
}

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

	log.Printf("connecting to database with - name: %s, user: %s\n", name, user)

	var err error
	connStr := fmt.Sprintf("connect_timeout=10 dbname=%s host=%s user=%s password=%s sslmode=disable", name, host, user, pass)
	db, err = sql.Open("postgres", connStr)
	return err
}
