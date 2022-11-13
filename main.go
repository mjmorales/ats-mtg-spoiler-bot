package main

import (
	"flag"
	"log"
	"time"

	gocron "github.com/go-co-op/gocron"
)

var interval string

func init() {
	flag.StringVar(&interval, "interval", "1h", "help message for flagname")
	flag.Parse()
}

func main() {
	log.Println("Starting Scheduler with interval:", interval)
	s := gocron.NewScheduler(time.UTC)
	s.Every(interval).SingletonMode().Do(ScrapeSpoilers)
	s.StartBlocking()
}
