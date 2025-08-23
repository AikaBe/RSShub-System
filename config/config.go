package config

var (
	Interval = "3m"
	Workers  = 3
)

type Jobs struct {
	Id   int
	Name string
	URL  string
}
