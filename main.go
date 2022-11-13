package main

import (
	"flag"
	"log"
	"time"

	gocron "github.com/go-co-op/gocron"
)

type Config struct {
	DB                DBConfig
	DiscordWebhookURL string
}

var interval string
var database string
var databaseUser string
var databasePassword string
var databaseHost string
var databasePort string
var databaseSslmode string
var discordWebhookURL string

func init() {
	flag.StringVar(&interval, "interval", "1h", "time interval to run scraper")
	flag.StringVar(&database, "database", "ats_mtg_spoiler_bot", "database name")
	flag.StringVar(&databaseUser, "database-user", "postgres", "database user")
	flag.StringVar(&databasePassword, "database-password", "password", "database password")
	flag.StringVar(&databaseHost, "database-host", "0.0.0.0", "database host")
	flag.StringVar(&databasePort, "database-port", "7000", "database port")
	flag.StringVar(&databaseSslmode, "database-sslmode", "disable", "database sslmode")
	flag.StringVar(&discordWebhookURL, "discord-webhook-url", "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz", "discord webhook")
	flag.Parse()
}

func main() {
	log.Println("Starting Scheduler with interval:", interval)

	config := Config{
		DiscordWebhookURL: discordWebhookURL,
		DB: DBConfig{
			Database:         database,
			DatabaseUser:     databaseUser,
			DatabasePassword: databasePassword,
			DatabaseHost:     databaseHost,
			DatabasePort:     databasePort,
			DatabaseSslmode:  databaseSslmode,
		},
	}

	db := ConnectDB(config)
	defer db.Close()

	s := gocron.NewScheduler(time.UTC)
	s.Every(interval).SingletonMode().Do(func() { ScrapeSpoilers(config, db) })
	s.StartBlocking()
}
