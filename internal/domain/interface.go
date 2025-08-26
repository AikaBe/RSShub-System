package domain

import (
	"context"
)

type FeedManager interface {
	Start(ctx context.Context) error
	SetInterval(newInterval string) error
	SetWorkers(newWorkers int, parentCtx context.Context) error
}
