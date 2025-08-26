package cli

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"rsshub/internal/adapter/postgre"
	"rsshub/internal/app"
	"rsshub/internal/domain"
)

type Handler struct {
	manager domain.FeedManager
}

func NewHandler(manager domain.FeedManager) *Handler {
	return &Handler{manager: manager}
}

func FlagHandler(pg *postgre.ApiAdapter, manager domain.FeedManager) {
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
		num := listCmd.Int("num", 3, "num of feeds")
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
		if err := manager.Start(ctx); err != nil {
			slog.Error("Failed to start fetcher", "err", err)
		}

	case "set-interval":
		interval := os.Args[2]
		if err := manager.SetInterval(interval); err != nil {
			slog.Error("Cannot change interval", "err", err)
		}
	case "set-workers":
		strWorkers := os.Args[2]
		workers, err := strconv.Atoi(strWorkers)
		if err != nil {
			slog.Error("Invalid workers count", "err", err)
			return
		}
		if err := manager.SetWorkers(workers, context.Background()); err != nil {
			slog.Error("Cannot change workers", "err", err)
		}
	case "articles":
		articlesCmd := flag.NewFlagSet("articles", flag.ExitOnError)
		feedName := articlesCmd.String("feed-name", "", "name of the feed")
		num := articlesCmd.Int("num", 3, "number of articles to show")
		articlesCmd.Parse(os.Args[2:])

		if *feedName == "" {
			slog.Error("Feed name is required")
			os.Exit(1)
		}

		articles, err := app.GetArticlesService(pg, *feedName, *num)
		if err != nil {
			slog.Error("failed to get articles", "err", err)
			os.Exit(1)
		}

		fmt.Printf("Feed: %s\n\n", *feedName)
		for i, a := range articles {
			published := "unknown"
			if a.PublishedAt != nil {
				published = a.PublishedAt.Format("2006-01-02")
			}
			fmt.Printf("%d. [%s] %s\n   %s\n\n", i+1, published, a.Title, a.Link)
		}

	default:
		fmt.Println("unknown command:", "cmd", os.Args[1])
	}
}
