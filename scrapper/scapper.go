package scrapper

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/tobg8/crypto-viz/common"
)

// ScrapCurrencies scraps currency information from "coingecko.com"
func ScrapCurrencies(url string) ([]common.CurrencyEvent, error) {
	c := colly.NewCollector()
	c.UserAgent = "Go scraping"

	var currencies []common.Currency
	var countCurrency int

	// Create a Currency Object from each found rows
	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			currency := common.Currency{
				ID:            el.ChildText("td:nth-child(2)"),
				Name:          el.ChildText("td:nth-child(3) span:first-child"),
				Acronym:       el.ChildText("td:nth-child(3) span:nth-child(2)"),
				Cours:         el.ChildText("td:nth-child(4)"),
				Variation1h:   el.ChildText("td:nth-child(5)"),
				Variation1d:   el.ChildText("td:nth-child(6)"),
				Variation1w:   el.ChildText("td:nth-child(7)"),
				Volume:        el.ChildText("td:nth-child(8)"),
				MarketCapital: el.ChildText("td:nth-child(9)"),
			}
			currencies = append(currencies, currency)
		})

		countCurrency = e.DOM.Find("tr").Length()
	})

	c.Visit(url)

	if len(currencies) != countCurrency {
		return nil, fmt.Errorf("expected %d currencies but got %d", countCurrency, len(currencies))
	}

	currenciesEvents, err := common.CurrenciesToCurrencyEvents(currencies)
	if err != nil {
		return nil, fmt.Errorf("could not convert Currency to CurrencyEvent: %w", err)
	}
	return currenciesEvents, nil
}

// ScrapNewCurrenciesscraps currency information from newly trending crypto currencies.
func ScrapNewCurrencies(url string) ([]common.CurrencyEvent, error) {
	c := colly.NewCollector()
	c.UserAgent = "Go scrapping"

	var currencies []common.Currency
	var countCurrency int

	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			currency := common.Currency{
				Name:        el.ChildText("td:nth-child(3) span:first-child"),
				Acronym:     el.ChildText("td:nth-child(3) span:nth-child(2)"),
				Cours:       el.ChildText("td:nth-child(4)"),
				Chaine:      el.ChildText("td:nth-child(5)"),
				Variation1h: el.ChildText("td:nth-child(6)"),
				Variation1d: el.ChildText("td:nth-child(7)"),
				Volume24d:   el.ChildText("td:nth-child(8)"),
				FDV:         el.ChildText("td:nth-child(9)"),
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

	currenciesEvents, err := common.CurrenciesToCurrencyEvents(currencies)
	if err != nil {
		return nil, fmt.Errorf("could not convert newCurrency to CurrencyEvent: %w", err)
	}
	return currenciesEvents, nil
}
