package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"fired-calendar/config"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create tables
	createTables()
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		recovery_phrase TEXT NOT NULL,
		fired_date DATE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		is_deleted BOOLEAN DEFAULT FALSE
	);`

	calendarTable := `
	CREATE TABLE IF NOT EXISTS calendar_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		date DATE NOT NULL,
		checked BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		UNIQUE(user_id, date)
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(calendarTable)
	if err != nil {
		log.Fatal("Failed to create calendar_entries table:", err)
	}
}
