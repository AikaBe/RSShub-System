package postgre

import (
	"log/slog"
	"time"

	"rsshub/internal/domain/model"
)

func (a *ApiAdapter) GetOldestFeeds() ([]model.Feed, error) {
	workers := a.GetWorkers()
	slog.Info("workers count ", workers)
	rows, err := a.db.Query(`
		SELECT id, name, url 
		FROM feeds 
		ORDER BY updated_at asc
		LIMIT $1
	`, workers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []model.Feed
	var id int
	for rows.Next() {
		var name, url string
		if err := rows.Scan(&id, &name, &url); err != nil {
			return nil, err
		}
		feeds = append(feeds, model.Feed{
			Id:   id,
			Name: name,
			Url:  url,
		})
		a.changeUpdate(id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (a *ApiAdapter) changeUpdate(id int) {
	_, err := a.db.Exec(`update feeds set updated_at = NOW() where id = $1`, id)
	if err != nil {
		slog.Warn("cannot update the field updated_at from feeds table")
		return
	}
}

func (a *ApiAdapter) AddArticle(item model.RSSItem, feedID int) error {
	pubAt := parsePubDate(item.PubDate)

	_, err := a.db.Exec(`
		insert into articles (feed_id, title, link, published_at, description, created_at, updated_at)
		values ($1, $2, $3, $4, $5, NOW(), NOW())
	`, feedID, item.Title, item.Link, pubAt, item.Description)

	return err
}

func (a *ApiAdapter) SetWorkers(workersCount int) {
	_, err := a.db.Exec(`update settings set workers = $1 `, workersCount)
	if err != nil {
		slog.Warn("cannot update workers count!")
		return
	}
}

func (a *ApiAdapter) GetWorkers() int {
	var workers int
	err := a.db.QueryRow(`SELECT workers FROM settings LIMIT 1`).Scan(&workers)
	if err != nil {
		slog.Warn("cannot get workers!", "err", err)
		return 0
	}

	return workers
}

func (a *ApiAdapter) SetInterval(newInterval string) {
	_, err := a.db.Exec(`update settings set interval = $1 `, newInterval)
	if err != nil {
		slog.Warn("cannot update interval!")
		return
	}
}

func (a *ApiAdapter) GetInterval() string {
	var interval string
	err := a.db.QueryRow(`SELECT interval FROM settings LIMIT 1`).Scan(&interval)
	if err != nil {
		slog.Warn("cannot get workers!", "err", err)
		return ""
	}

	return interval
}

func (a *ApiAdapter) TryLock() (bool, error) {
	var ok bool
	err := a.db.QueryRow("SELECT pg_try_advisory_lock(12345)").Scan(&ok)
	return ok, err
}

func (a *ApiAdapter) Unlock() error {
	_, err := a.db.Exec("SELECT pg_advisory_unlock(12345)")
	return err
}

func parsePubDate(s string) *time.Time {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
