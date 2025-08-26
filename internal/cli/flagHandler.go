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
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		num := listCmd.Int("num", 5, "num of feeds")
		listCmd.Parse(os.Args[2:])
		feeds, err := app.ListFeedsService(pg, *num)
		if err != nil {
			slog.Error("failed to list", "err", err)
		} else {
			slog.Info("Listening available RSS feeds:")
			for i, feed := range feeds {
				fmt.Printf("%d. Name: %s\n   URL: %s\n", i+1, feed.Name, feed.Url)
			}
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		name := deleteCmd.String("name", "", "name of the feed to delete")

		deleteCmd.Parse(os.Args[2:])

		if *name == "" {
			slog.Error("Feed name cannot be empty")
			os.Exit(1)
		}

		err := app.DeleteFeedService(pg, *name)
		if err != nil {
			slog.Error("failed to delete feed", "err", err)
		} else {
			slog.Info("Feed deleted successfully", "name", *name)
		}

	case "fetch":
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		if err := app.Start(ctx, pg); err != nil {
			slog.Error("Wrong background work! ", err)
		}

	case "set-interval":
		interval := os.Args[2]
		err := app.SetInterval(interval, pg)
		if err != nil {
			slog.Error("Cannot change the interval", err)
		}

	case "set-workers":
		strWorkers := os.Args[2]
		workers, err := strconv.Atoi(strWorkers)
		if err != nil {
			slog.Error("Invalid workers count", err)
		}
		err = app.SetWorkers(workers, pg, context.Background())
		if err != nil {
			slog.Error("Cannot change the workers", err)
		}

	default:
		fmt.Println("unknown command:", "cmd", os.Args[1])
	}
}
