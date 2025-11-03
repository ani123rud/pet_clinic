package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	user := getenvDefault("DB_USER", "petuser")
	password := getenvDefault("DB_PASSWORD", "petpass")
	dbname := getenvDefault("DB_NAME", "petclinic")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func getenvDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
