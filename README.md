RSSHub

RSSHub is a CLI application for fetching, storing, and displaying articles from RSS feeds. It uses a background worker pool to periodically fetch RSS data and store it in PostgreSQL.

 Features

 CLI-based interface

 Add, list, and delete RSS feeds

 Periodically fetch articles using ticker and workers

 Dynamically change fetch interval and worker count

 Graceful shutdown of ticker and workers

 Postgres-based article storage

 Project Requirements

 Code must be formatted using gofumpt

 No third-party dependencies (except PostgreSQL driver)

 Must compile without errors using:

go build -o rsshub .


 Must run without data races using:

go run -race main.go


 Dynamic interval & worker resizing via CLI (without restart)

 Clear error messages and proper exit codes on failure

 Docker Setup

The project includes docker-compose.yml to run:

PostgreSQL

RSSHub CLI application

docker-compose up --build

 Configuration

Set in .env or directly in Docker Compose:

# CLI App
CLI_APP_TIMER_INTERVAL=3m
CLI_APP_WORKERS_COUNT=3

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=rssuser
POSTGRES_PASSWORD=rsspass
POSTGRES_DBNAME=rsshub

 CLI Commands
Start Background Fetcher
./rsshub fetch


Starts a ticker + worker pool. You will see:

The background process for fetching feeds has started (interval = 3 minutes, workers = 3)


Only one instance is allowed at a time.

Set Interval (Live)
./rsshub set-interval 2m


Changes RSS fetch interval without restarting the app:

Interval of fetching feeds changed from 3 minutes to 2 minutes

Set Workers (Live)
./rsshub set-workers 5


Resizes the background worker pool dynamically:

Number of workers changed from 3 to 5

Add RSS Feed
./rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"

List Feeds
./rsshub list --num 5


Shows the 5 most recently added feeds.

Delete Feed
./rsshub delete --name "tech-crunch"

Show Articles
./rsshub articles --feed-name "tech-crunch" --num 5

Help
./rsshub --help

 RSS Structure Example
<rss>
  <channel>
    <title>RSS Feed</title>
    <link>https://example.com</link>
    <item>
      <title>Post 1</title>
      <link>https://example.com/1</link>
      <pubDate>Mon, 06 Sep 2021 12:00:00 GMT</pubDate>
      <description>Summary of post 1</description>
    </item>
  </channel>
</rss>


Parsed into:

type RSSFeed struct {
	Channel struct {
		Title string    `xml:"title"`
		Link  string    `xml:"link"`
		Item  []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Feed struct {
	Id   int
	Name string
	Url  string
}

type Article struct {

	ID          int
	FeedID      int
	Title       string
	Link        string
	Description string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	
}

 Background Aggregator

Default interval: 3m

Default workers: 3

Fetches N oldest feeds from DB every tick

Sends jobs to workers via a buffered channel

Workers fetch and parse RSS, then store articles in Postgres

Interval and worker count can change while running

 Safety Rules
Problem	Solution
 Data Race	Use sync.Mutex or atomic
 Goroutine Leaks	Use context.Context for cancellation
 Duplicate Tickers	Always Stop() old ticker before creating a new one
 Closing Channel Twice	Only one goroutine closes the channel
 ticker.Reset() Panic	Don’t call Reset() on a stopped ticker
 Deadlock on jobs	Make sure workers are always reading from jobs channel

 Migrations

Directory: migrations/

-- для самих каналов
create table feeds (
    id serial primary key,
    name text not null unique,
    url text not null unique,
    created_at timestamp not null ,
    updated_at timestamp not null
);

-- для самих статей
create table articles (
    id serial primary key,
    feed_id int references feeds(id) on delete cascade,
    title text not null,
    link text not null,
    published_at timestamp,
    description text not null ,
    created_at timestamp,
    updated_at timestamp
);

create table settings (
                          id serial primary key,
                          interval text not null default '3m',
                          workers int not null default 3
);
insert into settings (interval, workers) values ('3m', 3);

 Example Workflow

Terminal 1:

./rsshub fetch
# → Background process for fetching feeds has started


Terminal 2:

./rsshub set-interval 1m
./rsshub set-workers 5


To stop:

Press Ctrl+C

OR send a signal (e.g., kill)

Logs:

Graceful shutdown: aggregator stopped

 Warnings

Do not spam RSS feed servers

Print logs per request

Be ready to interrupt (Ctrl+C) if too many requests occur

Interface Overview
type Aggregator interface {
  Start(ctx context.Context) error
  Stop() error
  SetInterval(d time.Duration)
  Resize(workers int) error
}

 Summary

This project meets all constraints:

 Safe concurrency

 Dynamic control over workers & interval

 Works with PostgreSQL

🧼Clean and idiomatic Go code