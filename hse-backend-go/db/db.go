package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN has not been set in .env")
	}

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database (ensure migrations have been run and MySQL is running): %v", err)
	}

	DB = conn
	log.Println("Connected to the MySQL database")
}
