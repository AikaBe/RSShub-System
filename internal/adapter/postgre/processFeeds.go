package postgre

import (
	"log/slog"
	"rsshub/config"
	"rsshub/internal/domain/model"
	"time"
)

func (a *ApiAdapter) GetOldestFeeds() ([]model.Feed, error) {
	slog.Info("workers count ", config.Workers)
	rows, err := a.db.Query(`
		SELECT id, name, url 
		FROM feeds 
		ORDER BY updated_at asc
		LIMIT $1
	`, config.Workers)
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	a.changeUpdate(id)
	return feeds, nil
}

func (a *ApiAdapter) changeUpdate(id int) {
	_, err := a.db.Exec(`update feeds set updated_at = NOW() where id = $1`, id)
	if err != nil {
		slog.Warn("cannot update the field ipdated_at from feeds table")
		return
	}
}

func (a *ApiAdapter) AddArticle(item model.RSSItem, feedID int) error {
	var pubAt *time.Time
	if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
		pubAt = &t
	}

	_, err := a.db.Exec(`
		insert into articles (feed_id, title, link, published_at, description, created_at, updated_at)
		values ($1, $2, $3, $4, $5, NOW(), NOW())
	`, feedID, item.Title, item.Link, pubAt, item.Description)

	return err
}

func (a *ApiAdapter) ReadArticle() ([]model.Article, error) {
	rows, err := a.db.Query(`
		select id, feed_id, title, link, published_at, description, created_at, updated_at 
		from articles`)
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
