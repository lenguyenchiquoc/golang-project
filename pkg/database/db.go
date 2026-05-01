package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Can not open", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Can not connect", err)
	}

	createTables(db)
	log.Println("✓ Initialize successfully!")
	return db
}

func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id           TEXT PRIMARY KEY,
		username     TEXT UNIQUE NOT NULL,
		email        TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS manga (
		id             TEXT PRIMARY KEY,
		title          TEXT NOT NULL,
		author         TEXT,
		genres         TEXT,
		status         TEXT,
		total_chapters INTEGER,
		description    TEXT,
		cover_url 	   TEXT
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		user_id         TEXT,
		manga_id        TEXT,
		current_chapter INTEGER DEFAULT 0,
		status          TEXT DEFAULT 'reading',
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (manga_id) REFERENCES manga(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Can not create table", err)
	}
}
