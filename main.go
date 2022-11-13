package main

import (
	"flag"
	"log"
	"time"

	gocron "github.com/go-co-op/gocron"
)

var interval string
var database string
var databaseUser string
var databasePassword string
var databaseHost string
var databasePort string
var databaseSslmode string

func init() {
	flag.StringVar(&interval, "interval", "1h", "time interval to run scraper")
	flag.StringVar(&database, "database", "ats_mtg_spoiler_bot", "database name")
	flag.StringVar(&databaseUser, "database-user", "postgres", "database user")
	flag.StringVar(&databasePassword, "database-password", "password", "database password")
	flag.StringVar(&databaseHost, "database-host", "0.0.0.0", "database host")
	flag.StringVar(&databasePort, "database-port", "7000", "database port")
	flag.StringVar(&databaseSslmode, "database-sslmode", "disable", "database sslmode")
	flag.Parse()
}

func main() {
	log.Println("Starting Scheduler with interval:", interval)

	dbConfig := DBConfig{
		Database:         database,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseHost:     databaseHost,
		DatabasePort:     databasePort,
		DatabaseSslmode:  databaseSslmode,
	}

	s := gocron.NewScheduler(time.UTC)
	s.Every(interval).SingletonMode().Do(func() { ScrapeSpoilers(dbConfig) })
	s.StartBlocking()
}
