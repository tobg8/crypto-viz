package common

import "time"

type APIResponse struct {
	Results []ArticleAPI `json:"results"`
}

type ArticleAPI struct {
	ID          int           `json:"id"`
	Kind        string        `json:"kind"`
	SourceAPI   SourceAPI     `json:"source"`
	Title       string        `json:"title"`
	URL         string        `json:"url"`
	Slug        string        `json:"slug"`
	Currencies  []CurrencyAPI `json:"currencies"`
	PublishedAt time.Time     `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
}

type SourceAPI struct {
	Title  string `json:"title"`
	Region string `json:"region"`
	Domain string `json:"domain"`
}

type CurrencyAPI struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	URL   string `json:"url"`
}

type ArticleEvent struct {
	ID          int             `json:"id"`
	Kind        string          `json:"kind,omitempty"`
	Source      string          `json:"source,omitempty"`
	Title       string          `json:"title"`
	URL         string          `json:"url"`
	Currencies  []CurrencyEvent `json:"currencies"`
	PublishedAt time.Time       `json:"published_at"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CurrencyEvent struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}
