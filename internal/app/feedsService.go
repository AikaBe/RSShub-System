package app

import (
	"errors"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/domain/model"
)

func AddFeedsService(pg *postgre.ApiAdapter, name, url string) error {
	if name == "" || url == "" {
		return errors.New("Name and Url cannot be empty!")
	}
	err := pg.AddFeed(name, url)
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
