package main

import (
	"fmt"
	"log/slog"
	"os"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/cli"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DBNAME")

	//	interval := os.Getenv("CLI_APP_TIMER_INTERVAL")
	//	workers := os.Getenv("CLI_APP_WORKERS_COUNT")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	pgAdapter, err := postgre.NewApiAdapter(connStr)
	if err != nil {
		slog.Error("Postgres connection error", "err", err)
		os.Exit(1)
	}
	defer pgAdapter.Close()
	cli.FlagHandler(pgAdapter)
}
