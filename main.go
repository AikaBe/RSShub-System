package main

import (
	"fmt"
	"os"
	"rsshub/config"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/cli"

	_ "github.com/lib/pq"
)

func main() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName,
	)

	pgAdapter, err := postgre.NewApiAdapter(connStr)
	if err != nil {
		// slog.Error("Postgres connection error", "err", err)
		os.Exit(1)
	}
	defer pgAdapter.Close()
	cli.FlagHandler(pgAdapter)
}
