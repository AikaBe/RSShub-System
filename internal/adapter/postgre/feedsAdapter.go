package postgre

import (
	"errors"
	"log/slog"
	"time"

	"rsshub/internal/domain/model"
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
	query := `SELECT EXISTS (        SELECT 1 FROM feeds WHERE name = $1 OR url = $2
    )`
	var exists bool
	err := pg.db.QueryRow(query, name, url).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (a *ApiAdapter) GetFeeds(limit int) ([]model.Feed, error) {
	rows, err := a.db.Query(`
	select id, name, url
	from feeds order by created_at DESC 
	limit $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feeds []model.Feed
	for rows.Next() {
		var id int
		var name, url string
		if err := rows.Scan(&id, &name, &url); err != nil {
			return nil, err
		}
		feeds = append(feeds, model.Feed{
			Id:   id,
			Name: name,
			Url:  url,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (a *ApiAdapter) DeleteFeedByName(name string) error {
	var id int
	err := a.db.QueryRow(`select id from feeds where name = $1`, name).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("feed not found")
		}
		return err
	}
	_, err = a.db.Exec(`DELETE FROM feeds WHERE name = $1`, name)
	if err != nil {
		return err
	}

	slog.Info("Feed deleted from database", "name", name)
	return nil
}

func (a *ApiAdapter) GetArticlesByFeedName(feedName string, limit int) ([]model.Article, error) {
	rows, err := a.db.Query(`
		SELECT a.id, a.feed_id, a.title, a.link, a.published_at, a.description, a.created_at, a.updated_at
		FROM articles a
		JOIN feeds f ON a.feed_id = f.id
		WHERE f.name = $1
		ORDER BY a.published_at DESC
		LIMIT $2
	`, feedName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var art model.Article
		if err := rows.Scan(
			&art.ID, &art.FeedID, &art.Title, &art.Link,
			&art.PublishedAt, &art.Description, &art.CreatedAt, &art.UpdatedAt,
		); err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, rows.Err()
}
