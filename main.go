package rsshub

import (
	"log/slog"
	"os"
	"rsshub/internal/adapter/postgre"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=rsshub-db password=rsspass dbname=rsshub sslmode=disable"
	pgAdapter, err := postgre.NewApiAdapter(connStr)
	if err != nil {
		slog.Error("Postgres connection error", "err", err)
		os.Exit(1)
	}
	defer pgAdapter.Close()
}
