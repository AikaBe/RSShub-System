package cli

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
	"syscall"
)

func Start(pg *postgre.ApiAdapter) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx, pg); err != nil {
		slog.Error("Wrong background work! ", err)
	}
	return nil
}
