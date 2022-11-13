package main

import (
	"database/sql"
	"fmt"
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
	var id int

	row := db.QueryRow("SELECT id FROM scraped_cards WHERE url = $1 LIMIT 1", url)
	err := row.Scan(&id)

	// db.QueryRow returns ErrNoRows if there are no rows
	// in the result set so we can use that to check if
	// the card already exists
	return err == nil
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

// connect to the database
func ConnectDB(config Config) *sql.DB {
	dbConfig := config.DB
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.DatabaseUser,
		dbConfig.DatabasePassword,
		dbConfig.DatabaseHost,
		dbConfig.DatabasePort,
		dbConfig.Database,
		dbConfig.DatabaseSslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
