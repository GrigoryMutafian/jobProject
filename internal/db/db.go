package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "123456789"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "subs"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable options='-c application_name=subs-app'",
		host, user, password, dbname, port)

	var err error
	log.Printf("connecting to host=%s dbname=%s port=%s user=%s sslmode=disable", host, dbname, port, user)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("checking connection error: %v", err)
	}

	fmt.Println("database connection successful")
	return nil
}
