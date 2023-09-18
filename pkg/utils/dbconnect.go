package utils

import (
	"database/sql"
	"log"
)

func ConnectToPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db, err
}
