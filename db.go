//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=petuser password=petpass dbname=petclinic sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL!")
}
