package scrapper

import (
	"fmt"

	"github.com/gocolly/colly"
)

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
}

type NewCurrency struct {
	Name        string
	Acronym     string
	Cours       string
	Chaine      string
	Variation1h string
	Variation1d string
	Volume24d   string
	FDV         string
	LastAdded   string
}

// ScrapCurrencies scraps currency information from "coingecko.com"
func ScrapCurrencies(url string) ([]Currency, error) {
	c := colly.NewCollector()
	c.UserAgent = "Go scraping"

	var currencies []Currency
	var countCurrency int

	// Create a channel to capture errors
	errorCh := make(chan error, 1)

	// Handle errors asynchronously
	c.OnError(func(r *colly.Response, err error) {
		errorCh <- err
	})

	// Create a Currency Object from each found rows
	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			currency := Currency{
				ID:            el.ChildText("td:nth-child(2)"),
				Name:          el.ChildText("td:nth-child(3) span:first-child"),
				Acronym:       el.ChildText("td:nth-child(3) span:nth-child(2)"),
				Variation1h:   el.ChildText("td:nth-child(4)"),
				Variation1d:   el.ChildText("td:nth-child(5)"),
				Variation1w:   el.ChildText("td:nth-child(6)"),
				Volume:        el.ChildText("td:nth-child(7)"),
				MarketCapital: el.ChildText("td:nth-child(8)"),
			}

			currencies = append(currencies, currency)
		})

		countCurrency = e.DOM.Find("tr").Length()
	})

	c.Visit(url)

	// Check for errors asynchronously
	select {
	case err := <-errorCh:
		return nil, fmt.Errorf("error while scraping: %v", err)
	default:
	}

	if len(currencies) != countCurrency {
		return nil, fmt.Errorf("expected %d currencies but got %d", countCurrency, len(currencies))
	}

	return currencies, nil
}

func ScrapNewCurrencies(url string) ([]NewCurrency, error) {
	c := colly.NewCollector()
	c.UserAgent = "Go scrapping"

	var currencies []NewCurrency
	var countCurrency int

	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			currency := NewCurrency{
				Name:        el.ChildText("td:nth-child(2)"),
				Acronym:     el.ChildText("td:nth-child(3) span:first-child"),
				Cours:       el.ChildText("td:nth-child(3) span:nth-child(2)"),
				Chaine:      el.ChildText("td:nth-child(4)"),
				Variation1h: el.ChildText("td:nth-child(5)"),
				Variation1d: el.ChildText("td:nth-child(6)"),
				Volume24d:   el.ChildText("td:nth-child(7)"),
				FDV:         el.ChildText("td:nth-child(8)"),
				LastAdded:   el.ChildText("td:nth-child(10)"),
			}

			currencies = append(currencies, currency)
		})

		countCurrency = e.DOM.Find("tr").Length()
	})

	c.Visit(url)

	if len(currencies) != countCurrency {
		return nil, fmt.Errorf("expected %d currencies but got %d", countCurrency, len(currencies))
	}

	return currencies, nil
}
