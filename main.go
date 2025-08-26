package main

import (
	"fmt"
	"log/slog"
	"os"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
	"rsshub/internal/cli"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		printHelp()
		return
	}
	connStr := "host=localhost port=5432 user=rssuser password=rsspass dbname=rsshub sslmode=disable"

	pgAdapter, err := postgre.NewApiAdapter(connStr)
	if err != nil {
		slog.Error("Postgres connection error", "err", err)
		os.Exit(1)
	}
	defer pgAdapter.Close()
	service := app.NewService(pgAdapter)
	cli.FlagHandler(pgAdapter, service)
}

func printHelp() {
	fmt.Println(`
Usage:
  rsshub COMMAND [OPTIONS]

Available Commands:
  add             Add a new RSS feed
  set-interval    Set RSS fetch interval (in seconds)
  set-workers     Set number of workers for feed processing
  list            List all added RSS feeds
  delete          Delete a feed by its ID or name
  articles        Show latest fetched articles
  fetch           Start the background fetcher (periodic worker pool)
  --help          Show this help message

Example Usages:
  ./rsshub add --name "TechCrunch" --url "https://techcrunch.com/feed/"
  ./rsshub set-interval 60
  ./rsshub set-workers 5
  ./rsshub list
  ./rsshub delete --id 1
  ./rsshub articles
  ./rsshub fetch
`)
}
