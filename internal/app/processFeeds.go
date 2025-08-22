package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"rsshub/config"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/domain/model"
	"time"
)

var jobs = make(chan config.Jobs, 100)

func Start(ctx context.Context, pg *postgre.ApiAdapter) error {
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down ticker...")
			return nil
		case t := <-ticker.C:
			fmt.Println("Tick at:", t)

			feeds, err := pg.GetOldestFeeds()
			if err != nil {
				return err
			}

			for _, f := range feeds {
				jobs <- config.Jobs{Name: f.Name, URL: f.Url}
			}
		}
	}
}

// workers
func Workers(n int) {
	for w := 1; w <= n; w++ {
		go func(id int) {
			for job := range jobs {
				resp, err := http.Get(job.URL)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					return
				}
				var rss model.RSSFeed
				if err := xml.Unmarshal(data, &rss); err != nil {
					return
				}
				for _, item := range rss.Channel.Item {
					err := api.AddArticle(item, job.FeedID)
					if err != nil {
						fmt.Println("db insert err:", err)
					}
				}
			}
		}(w)
	}
}
