package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

func CreateCard(e *colly.HTMLElement) Card {
	url := e.Request.URL.String()
	imageUrl := strings.Replace(url, ".html", ".jpg", 1)
	return Card{
		BaseUrl:  url,
		ImageURL: imageUrl,
	}
}

func ScrapeSpoilers(config Config, db *sql.DB) {
	log.Println("Scraping spoilers...")

	InitializeDB(db)

	spoilerBase := "https://mythicspoiler.com/"
	spoilerURL := spoilerBase + "newspoilers.html"

	c := colly.NewCollector(
		colly.AllowedDomains("mythicspoiler.com"),
	)

	cards := []Card{}

	// find cards on the newspoilers page
	c.OnHTML(".grid-card a:nth-child(1)", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "cards") {
			fullLink := fmt.Sprintf("%s%s", spoilerBase, link)
			fullLink = strings.ReplaceAll(fullLink, "\n", "")
			if !CheckCardExists(db, fullLink) {
				// log.Printf("Card already exists: %s", link)
				e.Request.Visit(fullLink)
			}
		}
	})

	// find the actual card image
	c.OnHTML("td", func(e *colly.HTMLElement) {
		width := e.Attr("width")
		if width == "535" || width == "530" {
			cards = append(cards, CreateCard(e))
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

	// insert new cards into the database
	for _, card := range cards {
		log.Printf("Inserting card: %s", card.BaseUrl)
		InsertCard(db, card)
	}

	// if there are new cards, notify discord
	if len(cards) > 0 {
		NotifyNewSpoilers(config, cards)
	} else {
		log.Println("No new cards found")
	}

	log.Println("Finished scraping spoilers.")
}
