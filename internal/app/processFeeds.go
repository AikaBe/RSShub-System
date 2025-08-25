package app

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"log/slog"
	"net/http"
	"rsshub/config"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/domain/model"
	"sync"
	"time"
)

var (
	jobs          = make(chan config.Jobs, 100)
	mu            sync.Mutex
	ticker        *time.Ticker
	workerCancels = make(map[int]context.CancelFunc)
)

func Start(ctx context.Context, pg *postgre.ApiAdapter) error {
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return err
	}
	ticker = time.NewTicker(interval)

	for i := 1; i <= config.Workers; i++ {
		wctx, cancel := context.WithCancel(ctx)
		workerCancels[i] = cancel
		go Workers(i, pg, wctx)
	}

	slog.Info("background fetch started", "interval", config.Interval, "workers", config.Workers)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down background fetcher...")
			if ticker != nil {
				ticker.Stop()
			}
			return nil
		case t := <-ticker.C:
			slog.Info("tick", "time", t)

			feeds, err := pg.GetOldestFeeds()
			if err != nil {
				return err
			}

			for _, f := range feeds {
				slog.Info("enqueue feed", "id", f.Id, "name", f.Name, "url", f.Url)
				jobs <- config.Jobs{Id: f.Id, Name: f.Name, URL: f.Url}
			}
		}
	}
}

// workers
func Workers(id int, pg *postgre.ApiAdapter, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				log.Printf("worker %d jobs channel closed\n", id)
				return
			}

			resp, err := http.Get(job.URL)
			if err != nil {
				log.Println("http get error:", err)
				continue
			}
			data, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Println("read error:", err)
				continue
			}

			var rss model.RSSFeed
			if err := xml.Unmarshal(data, &rss); err != nil {
				log.Println("xml error:", err)
				continue
			}

			for _, item := range rss.Channel.Item {
				mu.Lock()
				if err := pg.AddArticle(item, job.Id); err != nil {
					log.Println("db insert error:", err)
				}
				mu.Unlock()
			}
		}
	}
}

func SetInterval(newInterval string) error {
	d, err := time.ParseDuration(newInterval)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	if ticker != nil {
		ticker.Stop()
	}

	ticker = time.NewTicker(d)
	old := config.Interval
	config.Interval = newInterval
	slog.Info("new interval ", config.Interval)
	slog.Info("Interval of fetching feeds changed", "from", old, "to", newInterval)
	return nil
}

func SetWorkers(newWorkers int, pg *postgre.ApiAdapter, parentCnxt context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	old := config.Workers

	if newWorkers > config.Workers {
		for w := config.Workers + 1; w <= newWorkers; w++ {
			ctx, cancel := context.WithCancel(parentCnxt)
			workerCancels[w] = cancel
			go Workers(w, pg, ctx)
		}
	}

	if newWorkers < config.Workers {
		for w := config.Workers; w > newWorkers; w-- {
			if cancel, ok := workerCancels[w]; ok {
				cancel()
				delete(workerCancels, w)
			}
		}
	}

	config.Workers = newWorkers
	slog.Info("workers count ", config.Workers)
	slog.Info("Number of workers changed", "from", old, "to", newWorkers)
	return nil
}
