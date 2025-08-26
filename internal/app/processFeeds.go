package app

import (
	"context"
	"encoding/xml"
	"errors"
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
	jobs            = make(chan config.Jobs, 100)
	mu              sync.Mutex
	ticker          *time.Ticker
	workerCancels   = make(map[int]context.CancelFunc)
	currentInterval time.Duration
)

func Start(ctx context.Context, pg *postgre.ApiAdapter) error {
	locked, err := pg.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("Background process is already running")
	}
	defer pg.Unlock()

	interval, err := time.ParseDuration(pg.GetInterval())
	if err != nil {
		return err
	}
	currentInterval = interval
	ticker = time.NewTicker(interval)

	workers := pg.GetWorkers()
	for i := 1; i <= workers; i++ {
		wctx, cancel := context.WithCancel(ctx)
		workerCancels[i] = cancel
		go Workers(i, pg, wctx)
	}
	slog.Info("The background process for fetching feeds has started ", "interval", interval, "workers", workers)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down background fetcher...")
			if ticker != nil {
				ticker.Stop()
			}
			return nil
		case t := <-ticker.C:
			newIntervalStr := pg.GetInterval()
			newInterval, err := time.ParseDuration(newIntervalStr)
			if err == nil && newInterval != currentInterval {
				mu.Lock()
				ticker.Stop()
				ticker = time.NewTicker(newInterval)
				old := currentInterval
				currentInterval = newInterval
				mu.Unlock()
				slog.Info("Interval changed via DB", "from", old, "to", newInterval)
			}

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

func SetInterval(newInterval string, pg *postgre.ApiAdapter) error {
	if newInterval == "" {
		newInterval = "3m"
	}
	_, err := time.ParseDuration(newInterval)
	if err != nil {
		return err
	}

	old := pg.GetInterval()
	pg.SetInterval(newInterval)
	slog.Info("Interval of fetching feeds saved to DB", "from", old, "to", newInterval)
	return nil
}

func SetWorkers(newWorkers int, pg *postgre.ApiAdapter, parentCnxt context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if newWorkers < 1 || newWorkers > 10 {
		slog.Error("worker count must be between 1 and 10")
		return errors.New("worker count must be between 1 and 10")
	}
	workers := pg.GetWorkers()
	old := workers

	if newWorkers > workers {
		for w := workers + 1; w <= newWorkers; w++ {
			ctx, cancel := context.WithCancel(parentCnxt)
			workerCancels[w] = cancel
			go Workers(w, pg, ctx)
		}
	}

	if newWorkers < workers {
		for w := workers; w > newWorkers; w-- {
			if cancel, ok := workerCancels[w]; ok {
				cancel()
				delete(workerCancels, w)
			}
		}
	}

	pg.SetWorkers(newWorkers)
	slog.Info("Number of workers changed", "from", old, "to", newWorkers)
	return nil
}
