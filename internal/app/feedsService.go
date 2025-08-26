package app

import (
	"errors"
	"log/slog"
	"rsshub/internal/adapter/postgre"
)

func AddFeedsService(pg *postgre.ApiAdapter, name, url string) error {
	if name == "" || url == "" {
		return errors.New("Name and Url cannot be empty!")
	}
	exists, err := pg.FeedExists(name, url)
	if err != nil {
		return err
	}
	if exists {
		slog.Error("Feed with this name or url already exists!")
		return errors.New("Feed with this name or url already exists!")
	}

	err = pg.AddFeed(name, url)
	if err != nil {
		return err
	}
	return nil
}
