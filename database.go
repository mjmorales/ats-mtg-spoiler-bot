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

func CheckCardExists(db *sql.DB, url string) bool {
	var count int

	row := db.QueryRow("SELECT COUNT(*) FROM scraped_cards WHERE url = $1", url)
	err := row.Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count > 0
}

func InsertCard(db *sql.DB, card Card) {
	_, err := db.Exec("INSERT INTO scraped_cards (url, image_url) VALUES ($1, $2)", card.BaseUrl, card.ImageURL)
	if err != nil {
		log.Fatal(err)
	}
}