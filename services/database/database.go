package main

import (
	proto "database-service/genproto/database"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool

func CloseDb() {
	if pool != nil {
		pool.Close()
	}
}

func GetMessages(limit int32) (*proto.Messages, error) {
	if err := checkPool(); err != nil {
		return nil, err
	}

	sqlStr := "SELECT * FROM messages ORDER BY time DESC"
	if limit > 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := pool.Query(sqlStr)
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
	if err := checkPool(); err != nil {
		return err
	}

	if _, err := pool.Exec(fmt.Sprintf("INSERT INTO messages(message, name, time) VALUES ('%s', '%s', %d)", msg.Message, msg.Name, msg.Time)); err != nil {
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

	connStr := fmt.Sprintf("connect_timeout=10 dbname=%s host=%s user=%s password=%s sslmode=disable", name, host, user, pass)
	config, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		log.Fatalf("error occurred parsing database connection string, exiting..")
	}

	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: config})
	return err
}

func checkPool() error {
	if pool == nil {
		log.Println("database connection has not been initialized, database operation will not be conducted")
		return errors.New("database connection has not been initialized, database operation will not be conducted")
	}
	return nil
}
