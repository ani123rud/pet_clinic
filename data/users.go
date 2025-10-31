package data

import (
	"database/sql"
)

type UserRow struct {
	ID           int
	Email        string
	PasswordHash string
}

func EmailExists(db *sql.DB, email string) (bool, error) {
	var c int
	if err := db.QueryRow("SELECT COUNT(1) FROM users WHERE email=$1", email).Scan(&c); err != nil {
		return false, err
	}
	return c > 0, nil
}

func CreateUser(db *sql.DB, email, passwordHash string) (UserRow, error) {
	var u UserRow
	err := db.QueryRow("INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id, email, password_hash", email, passwordHash).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

func FindUserByEmail(db *sql.DB, email string) (UserRow, error) {
	var u UserRow
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}
