package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

var cards []Card

func AppendCard(e *colly.HTMLElement) {
	url := e.Request.URL.String()
	imageUrl := strings.Replace(url, ".html", ".jpg", 1)
	cards = append(cards, Card{
		BaseUrl:  url,
		ImageURL: imageUrl,
	})
}

func ScrapeSpoilers(dbConfig DBConfig) {
	log.Println("Scraping spoilers...")

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
	defer db.Close()

	spoilerBase := "https://mythicspoiler.com/"
	spoilerURL := spoilerBase + "newspoilers.html"

	c := colly.NewCollector(
		colly.AllowedDomains("mythicspoiler.com"),
	)

	// find cards on the newspoilers page
	c.OnHTML(".grid-card a:nth-child(1)", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "cards") {
			fullLink := fmt.Sprintf("%s%s", spoilerBase, link)
			fullLink = strings.ReplaceAll(fullLink, "\n", "")
			if CheckCardExists(db, fullLink) {
				log.Printf("Card already exists: %s", link)
				return
			}
			e.Request.Visit(fullLink)
		}
	})

	// find the actual card image
	c.OnHTML("td", func(e *colly.HTMLElement) {
		width := e.Attr("width")
		if width == "535" || width == "530" {
			AppendCard(e)
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(spoilerURL)
	for _, card := range cards {
		log.Printf("Inserting card: %s", card.BaseUrl)
		InsertCard(db, card)
	}

	log.Println("Finished scraping spoilers.")
}
