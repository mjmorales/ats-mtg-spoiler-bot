# ats-mtg-spoiler-bot

This bot scrapes mythic spoiler for new card images and posts to discord.
Images are stored in a Postgres Database so user's aren't alerted of dupes.

## CLI Flags
```
mmorales@manny-desktop:/mnt/d/Development/ats-mtg-spoiler-bot$ go run . --help
Usage of /tmp/go-build2260391231/b001/exe/ats-mtg-spoiler-bot:
  -database string
        database name (default "ats_mtg_spoiler_bot")
  -database-host string
        database host (default "0.0.0.0")
  -database-password string
        database password (default "password")
  -database-port string
        database port (default "7000")
  -database-sslmode string
        database sslmode (default "disable")
  -database-user string
        database user (default "postgres")
  -discord-webhook-url string
        discord webhook (default "https://discord.com/api/webhooks/1234567890/abcdefghijklmnopqrstuvwxyz")
  -interval string
        time interval to run scraper (default "1h")
```
