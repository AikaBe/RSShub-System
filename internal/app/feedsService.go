package app

import (
	"errors"
	"rsshub/internal/adapter/postgre"
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
