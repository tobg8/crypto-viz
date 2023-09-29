package common

import "time"

type NewsEvent struct {
	ID          string    `json:"id"`
	RssURL      string    `json:"rss_url"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	ImageURL    string    `json:"image_url"`
	PubDate     time.Time `json:"publication_date"`
	Author      string    `json:"author"`
	Categories  []string  `json:"categories"`
	Description string    `json:"description"`
}
