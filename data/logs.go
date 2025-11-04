package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// InitLogsTable creates the logs table if it doesn't exist
func InitLogsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			level VARCHAR(10) NOT NULL,
			message TEXT NOT NULL,
			file VARCHAR(255) NOT NULL,
			function VARCHAR(100) NOT NULL,
			user_id INTEGER NULL,
			user_email VARCHAR(255) NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
		CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
		CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
		CREATE INDEX IF NOT EXISTS idx_logs_user_email ON logs(user_email);
	`)
	if err != nil {
		return fmt.Errorf("failed to create logs table: %w", err)
	}
	log.Println("Logs table initialized successfully")
	return nil
}

type LogEntry struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file"`
	Function  string    `json:"function"`
	UserID    *int      `json:"user_id,omitempty"`
	UserEmail *string   `json:"user_email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// SaveLogWithUser saves a log entry with optional user_id and user_email (NULL when not valid)
func SaveLogWithUser(db *sql.DB, level, message, file, function string, userID sql.NullInt64, userEmail sql.NullString) error {
	_, err := db.Exec(
		`INSERT INTO logs (level, message, file, function, user_id, user_email) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		level, message, file, function, userID, userEmail,
	)
	return err
}

// SaveLog keeps backward compatibility without user_id
func SaveLog(db *sql.DB, level, message, file, function string) error {
	return SaveLogWithUser(db, level, message, file, function, sql.NullInt64{}, sql.NullString{})
}

// GetLogs retrieves logs with optional filters
func GetLogs(db *sql.DB, level string, limit, offset int) ([]LogEntry, error) {
	query := `SELECT id, level, message, file, function, user_id, user_email, created_at 
	          FROM logs 
	          WHERE ($1 = '' OR level = $1)
	          ORDER BY created_at DESC
	          LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, level, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		var uid sql.NullInt64
		var uemail sql.NullString
		err := rows.Scan(
			&log.ID,
			&log.Level,
			&log.Message,
			&log.File,
			&log.Function,
			&uid,
			&uemail,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if uid.Valid {
			u := int(uid.Int64)
			log.UserID = &u
		}
		if uemail.Valid {
			s := uemail.String
			log.UserEmail = &s
		}
		logs = append(logs, log)
	}

	return logs, nil
}
