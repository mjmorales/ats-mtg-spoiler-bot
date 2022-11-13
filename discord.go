package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Image struct {
	URL string `json:"url"`
}

type Embed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Color string `json:"color"`
	Image Image  `json:"image"`
}

type DiscordWebhookRequest struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

// sends a discord webhook request to the webhook url for new spoilers
func NotifyNewSpoilers(config Config, cards []Card) {
	var embeds []Embed

	for _, card := range cards {
		embed := Embed{
			Title: card.BaseUrl,
			URL:   card.BaseUrl,
			Color: "16773997",
			Image: Image{
				URL: card.ImageURL,
			},
		}
		embeds = append(embeds, embed)
	}

	contentString := fmt.Sprintf("Found **%d** new spoilers :biting_lip:", len(cards))
	newSpoiler := DiscordWebhookRequest{
		Content: contentString,
		Embeds:  embeds,
	}

	// post the request to the discord webhook url
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(newSpoiler)
	req, err := http.NewRequest("POST", config.DiscordWebhookURL, payloadBuf)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Posting to Webhook URL:", resp.Status)
}
