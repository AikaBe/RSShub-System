package app

import (
	"errors"
	"log/slog"

	"rsshub/internal/adapter/postgre"
	"rsshub/internal/domain/model"
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

func ListFeedsService(pg *postgre.ApiAdapter, num int) ([]model.Feed, error) {
	if num <= 0 {
		return nil, errors.New("number must be greater")
	}
	feeds, err := pg.GetFeeds(num)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func DeleteFeedService(pg *postgre.ApiAdapter, name string) error {
	if name == "" {
		return errors.New("Can not be enpty")
	}
	err := pg.DeleteFeedByName(name)
	if err != nil {
		return err
	}
	return nil
}

func GetArticlesService(pg *postgre.ApiAdapter, feedName string, limit int) ([]model.Article, error) {
	if feedName == "" {
		return nil, errors.New("feed name cannot be empty")
	}
	if limit <= 0 {
		limit = 3
	}

	return pg.GetArticlesByFeedName(feedName, limit)
}
