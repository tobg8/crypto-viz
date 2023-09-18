package common

// Currency represents a currency model
type Currency struct {
	ID            string
	Name          string
	Acronym       string
	Cours         string
	Variation1h   string
	Variation1d   string
	Variation1w   string
	Volume        string
	MarketCapital string
	Chaine        string
	Volume24d     string
	FDV           string
	LastAdded     string
}

// CurrencyEvent represents a currency kafka model event
type CurrencyEvent struct {
	ID            int
	Name          string
	Acronym       string
	Cours         float64
	Variation1h   float64
	Variation1d   float64
	Variation1w   float64
	Volume        float64
	MarketCapital float64
	Chaine        string
	Volume24d     float64
	FDV           float64
	LastAdded     string
}
