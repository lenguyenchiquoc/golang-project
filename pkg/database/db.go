package database

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
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
	seedGenres(db)
	log.Println("✓ Initialize successfully!")
	return db
}

func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id            TEXT PRIMARY KEY,
		username      TEXT UNIQUE NOT NULL,
		email         TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS manga (
		id             TEXT PRIMARY KEY,
		title          TEXT NOT NULL,
		author         TEXT,
		status         TEXT,
		total_chapters INTEGER,
		description    TEXT,
		cover_url      TEXT,
		average_rating REAL    DEFAULT 0,
		rating_count   INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS genres (
		id   TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS manga_genres (	
		manga_id TEXT,
		genre_id TEXT,
		PRIMARY KEY (manga_id, genre_id),
		FOREIGN KEY (manga_id) REFERENCES manga(id),
		FOREIGN KEY (genre_id) REFERENCES genres(id)
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		user_id         TEXT,
		manga_id        TEXT,
		current_chapter INTEGER DEFAULT 0,
		status          TEXT DEFAULT 'reading',
		updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id),
		FOREIGN KEY (user_id)  REFERENCES users(id),
		FOREIGN KEY (manga_id) REFERENCES manga(id)
	);

	CREATE TABLE IF NOT EXISTS user_ratings (
		user_id    TEXT,
		manga_id   TEXT,
		rating     INTEGER CHECK(rating >= 1 AND rating <= 10),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id),
		FOREIGN KEY (user_id)  REFERENCES users(id),
		FOREIGN KEY (manga_id) REFERENCES manga(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Can not create table:", err)
	}
}

func seedGenres(db *sql.DB) {
	genres := []string{
		"action", "adventure", "comedy", "drama", "fantasy",
		"horror", "mystery", "romance", "sci-fi", "slice of life",
		"supernatural", "psychological", "thriller", "sports",
		"mecha", "historical", "isekai", "magic", "school",
		"shounen", "seinen", "shoujo", "josei", "ecchi",
		"harem", "martial arts", "music", "game", "tragedy",
	}

	for _, name := range genres {
		_, err := db.Exec(`
			INSERT INTO genres (id, name)
			VALUES (?, ?)
			ON CONFLICT(name) DO NOTHING
		`, uuid.New().String(), name)

		if err != nil {
			log.Println("seed genre error:", err)
		}
	}

	log.Println("✓ Seed genres done")
}