package app

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"rsshub/config"
	"rsshub/internal/adapter/postgre"
	"rsshub/internal/domain/model"
)

type Service struct {
	db              *postgre.ApiAdapter
	jobs            chan config.Jobs
	mu              sync.Mutex
	ticker          *time.Ticker
	workerCancels   map[int]context.CancelFunc
	currentInterval time.Duration
}

func NewService(db *postgre.ApiAdapter) *Service {
	return &Service{
		db:            db,
		jobs:          make(chan config.Jobs, 100),
		workerCancels: make(map[int]context.CancelFunc),
	}
}

func (s *Service) Start(ctx context.Context) error {
	locked, err := s.db.TryLock()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("background process is already running")
	}
	defer s.db.Unlock()

	interval, err := time.ParseDuration(s.db.GetInterval())
	if err != nil {
		return err
	}
	s.currentInterval = interval
	s.ticker = time.NewTicker(interval)

	workers := s.db.GetWorkers()
	for i := 1; i <= workers; i++ {
		wctx, cancel := context.WithCancel(ctx)
		s.workerCancels[i] = cancel
		go s.worker(i, wctx)
	}
	slog.Info("Background fetch started", "interval", interval, "workers", workers)

	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down background fetcher")
			s.ticker.Stop()
			return nil
		case t := <-s.ticker.C:
			newIntervalStr := s.db.GetInterval()
			newInterval, err := time.ParseDuration(newIntervalStr)
			if err == nil && newInterval != s.currentInterval {
				s.mu.Lock()
				s.ticker.Stop()
				s.ticker = time.NewTicker(newInterval)
				old := s.currentInterval
				s.currentInterval = newInterval
				s.mu.Unlock()
				slog.Info("Interval changed via DB", "from", old, "to", newInterval)
			}

			slog.Info("tick", "time", t)

			feeds, err := s.db.GetOldestFeeds()
			if err != nil {
				return err
			}

			var feedIDs []int
			for _, f := range feeds {
				feedIDs = append(feedIDs, f.Id)
				s.jobs <- config.Jobs{Id: f.Id, Name: f.Name, URL: f.Url}
			}

			slog.Info("Fetched feeds", "ids", feedIDs)
		}
	}
}

func (s *Service) SetInterval(newInterval string) error {
	if newInterval == "" {
		newInterval = "3m"
	}
	_, err := time.ParseDuration(newInterval)
	if err != nil {
		return err
	}

	old := s.db.GetInterval()
	s.db.SetInterval(newInterval)
	slog.Info("Interval updated in DB", "from", old, "to", newInterval)
	return nil
}

func (s *Service) SetWorkers(newWorkers int, parentCtx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if newWorkers < 1 || newWorkers > 10 {
		return errors.New("worker count must be between 1 and 10")
	}

	old := s.db.GetWorkers()
	if newWorkers > old {
		for w := old + 1; w <= newWorkers; w++ {
			ctx, cancel := context.WithCancel(parentCtx)
			s.workerCancels[w] = cancel
			go s.worker(w, ctx)
		}
	} else if newWorkers < old {
		for w := old; w > newWorkers; w-- {
			if cancel, ok := s.workerCancels[w]; ok {
				cancel()
				delete(s.workerCancels, w)
			}
		}
	}
	s.db.SetWorkers(newWorkers)
	slog.Info("Worker count updated", "from", old, "to", newWorkers)
	return nil
}

func (s *Service) worker(id int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-s.jobs:
			resp, err := http.Get(job.URL)
			if err != nil {
				slog.Error("http get error:", err)
				continue
			}
			data, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				slog.Error("read error:", err)
				continue
			}

			var rss model.RSSFeed
			if err := xml.Unmarshal(data, &rss); err != nil {
				slog.Error("xml error:", err)
				continue
			}

			for _, item := range rss.Channel.Item {
				s.mu.Lock()

				exists, err := s.db.ArticleExists(item.Title)
				if err != nil {
					slog.Error("db exists check error:", err)
					s.mu.Unlock()
					continue
				}

				if exists {
					slog.Info("article already exists, skipping", "title", item.Title)
					s.mu.Unlock()
					continue
				}

				if err := s.db.AddArticle(item, job.Id); err != nil {
					slog.Error("db insert error:", err)
				}

				s.mu.Unlock()
			}

		}
	}
}
