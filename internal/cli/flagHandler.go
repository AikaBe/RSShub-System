package cli

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
	"strconv"
	"syscall"
)

func FlagHandler(pg *postgre.ApiAdapter) {
	if len(os.Args) < 2 {
		slog.Error("expected subcommands ! ")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		name := addCmd.String("name", "default", "name for feed")
		url := addCmd.String("url", "", "url for feed")

		addCmd.Parse(os.Args[2:])

		err := app.AddFeedsService(pg, *name, *url)
		if err != nil {
			slog.Error("failed to add feed", "err", err)
		} else {
			slog.Info("feed added successfully", "name", *name, "url", *url)
		}

	case "fetch":
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := app.Start(ctx, pg); err != nil {
			slog.Error("Wrong background work! ", err)
		}

	case "set-interval":
		interval := os.Args[2]
		err := app.SetInterval(interval)
		if err != nil {
			slog.Error("Cannot change the interval", err)
		}

	case "set-workers":
		strWorkers := os.Args[2]
		workers, err := strconv.Atoi(strWorkers)
		if err != nil {
			slog.Error("Cannot change the workers", err)
		}
		err = app.SetWorkers(workers, pg, context.Background())
		if err != nil {
			slog.Error("Cannot change the workers", err)
		}

	default:
		fmt.Println("unknown command:", "cmd", os.Args[1])
	}
}
