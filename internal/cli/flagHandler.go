package cli

import (
	"flag"
	"log/slog"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
)

func FlagHandler(pg *postgre.ApiAdapter) {
	name := flag.String("name", "default", "name for feed")
	url := flag.String("url", "url doesn't exist", "url")
	flag.Parse()
	err := app.AddFeedsService(pg, *name, *url)
	if err != nil {
		slog.Error("failed to add feed", "err", err)
	} else {
		slog.Info("feed added successfully", "name", *name, "url", *url)
	}
}
