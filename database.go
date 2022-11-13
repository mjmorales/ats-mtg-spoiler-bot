package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Card struct {
	BaseUrl  string
	ImageURL string
}

type DBConfig struct {
	Database         string
	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseSslmode  string
}

// checks if the card URL already exists in the database
func CheckCardExists(db *sql.DB, url string) bool {
	var count int

	row := db.QueryRow("SELECT COUNT(*) FROM scraped_cards WHERE url = $1", url)
	err := row.Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count > 0
}

// inserts a card record into the database
func InsertCard(db *sql.DB, card Card) {
	_, err := db.Exec("INSERT INTO scraped_cards (url, image_url) VALUES ($1, $2)", card.BaseUrl, card.ImageURL)
	if err != nil {
		log.Fatal(err)
	}
}

// initialize the database
func InitializeDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS scraped_cards (id SERIAL PRIMARY KEY, url TEXT, image_url TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}
