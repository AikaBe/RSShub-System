package postgre

import (
	"log/slog"
	"time"
)

func (a *ApiAdapter) AddFeed(name, url string) error {
	_, err := a.db.Exec(`insert into feeds (name,url,created_at,updated_at) values ($1,$2,$3,$4)`, name, url, time.Now(), time.Now())
	if err != nil {
		slog.Warn("cannot add data to DB")
		return err
	}
	slog.Info("data Added to the DB!")
	return nil
}

func (pg *ApiAdapter) FeedExists(name, url string) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM feeds WHERE name = $1 OR url = $2
	)`
	var exists bool
	err := pg.db.QueryRow(query, name, url).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
