package cli

import (
	"context"
	"log"
	"rsshub/config"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
	"strconv"
)

func Start(pg *postgre.ApiAdapter) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workers, err := strconv.Atoi(config.Workers)
	if err != nil {
		return err
	}
	app.Workers(workers)

	if err := app.Start(ctx, pg); err != nil {
		log.Fatal(err)
	}
	return nil
}
