package common

import "time"

type ArticleResponse struct {
	Results []ArticleAPI `json:"results"`
}

type ArticleAPI struct {
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

type Listing struct {
	ID                           string    `json:"id"`
	Symbol                       string    `json:"symbol"`
	Name                         string    `json:"name"`
	Image                        string    `json:"image"`
	CurrentPrice                 float64   `json:"current_price"`
	MarketCap                    int64     `json:"market_cap"`
	MarketCapRank                float64   `json:"market_cap_rank"`
	FullyDilutedValuation        float64   `json:"fully_diluted_valuation"`
	TotalVolume                  float64   `json:"total_volume"`
	High24H                      float64   `json:"high_24h"`
	Low24H                       float64   `json:"low_24h"`
	PriceChange24H               float64   `json:"price_change_24h"`
	PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64   `json:"circulating_supply"`
	TotalSupply                  float64   `json:"total_supply"`
	MaxSupply                    float64   `json:"max_supply"`
	Ath                          float64   `json:"ath"`
	AthChangePercentage          float64   `json:"ath_change_percentage"`
	AthDate                      time.Time `json:"ath_date"`
	Atl                          float64   `json:"atl"`
	AtlChangePercentage          float64   `json:"atl_change_percentage"`
	AtlDate                      time.Time `json:"atl_date"`
	LastUpdated                  time.Time `json:"last_updated"`
}
