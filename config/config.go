package config

import "os"

var (
	DbHost = os.Getenv("POSTGRES_HOST")
	DbPort = os.Getenv("POSTGRES_PORT")
	DbUser = os.Getenv("POSTGRES_USER")
	DbPass = os.Getenv("POSTGRES_PASSWORD")
	DbName = os.Getenv("POSTGRES_DBNAME")

	Interval = os.Getenv("CLI_APP_TIMER_INTERVAL")
	Workers  = os.Getenv("CLI_APP_WORKERS_COUNT")
)

type Jobs struct {
	Name string
	URL  string
}
