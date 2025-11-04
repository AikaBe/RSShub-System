#RSSHub

RSSHub is a CLI application that aggregates RSS feeds, fetches articles periodically, and stores them in PostgreSQL. It allows users to manage feeds, view the latest articles, and adjust the fetch interval and worker count dynamically.

##🚀 Live Demo

Since RSSHub is a CLI application, you can test it locally by running it with Docker Compose:

docker-compose up --build

##🛠️ Technologies Used

Backend: Go (CLI Application)

Database: PostgreSQL

Deployment: Docker & Docker Compose

##✨ Features

CLI-based interface for managing RSS feeds

Add, list, and delete RSS feeds

Periodic article fetching with a configurable interval

Worker pool for concurrent feed processing

Dynamic interval and worker resizing without restarting

Graceful shutdown of background processes

Safe concurrency (no data races)

##📦 Installation

Clone the repository:

git clone https://github.com/yourusername/rsshub.git
cd rsshub


Create a .env file with the following:

# CLI App
CLI_APP_TIMER_INTERVAL=3m
CLI_APP_WORKERS_COUNT=3

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=changem
POSTGRES_DBNAME=rsshub


##Run the project with Docker Compose:

docker-compose up --build


Alternatively, build locally:

go build -o rsshub .

##🎯 Usage

Start background fetching:

./rsshub fetch


Add a new RSS feed:

./rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"


List feeds:

./rsshub list --num 5


Set fetch interval dynamically:

./rsshub set-interval 2m


Resize workers dynamically:

./rsshub set-workers 5


View latest articles:

./rsshub articles --feed-name "tech-crunch" --num 5

##🏗️ Architecture

Aggregator Interface: Handles start, stop, interval changes, and worker resizing.

Worker Pool: Fetches and parses RSS feeds concurrently.

Ticker: Runs periodic feed aggregation.

PostgreSQL Storage: Stores feeds metadata and articles.

CLI Commands: Provides interface to control feeds, workers, and intervals.

##📸 Screenshots

Since this is a CLI application, output examples in terminal:

###Start Fetching

$ ./rsshub fetch
The background process for fetching feeds has started (interval = 3 minutes, workers = 3)


List Feeds

$ ./rsshub list --num 3
1. Name: tech-crunch | URL: https://techcrunch.com/feed/ | Added: 2025-06-10 15:34
2. Name: hacker-news | URL: https://news.ycombinator.com/rss | Added: 2025-06-10 15:37
3. Name: bbc-world | URL: http://feeds.bbci.co.uk/news/world/rss.xml | Added: 2025-06-11 09:15

##🔮 Future Improvements

Web interface for easier feed management

RSS feed categorization and search

Notifications for new articles

Multi-user support with role-based access