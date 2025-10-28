package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	connStr := "host=localhost port=5432 user=petuser password=petpass dbname=petclinic sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	fmt.Println("Connected to PostgreSQL!")
	return nil
}
