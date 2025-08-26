package model

import "time"

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
