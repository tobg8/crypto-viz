package scrapper

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
	"github.com/tobg8/crypto-viz/common"
)

// ScrapCurrencies scraps currency information from "coingecko.com"
func ScrapCurrencies(url string) ([]common.CurrencyEvent, error) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.75 Safari/537.36"
	// c.AllowedDomains = []string{"www.coingecko.com"}
	// c.IgnoreRobotsTxt = false

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "\nError:", err, r.StatusCode)
	})

	var currencies []common.Currency
	var countCurrency int
	// Create a Currency Object from each found rows
	c.OnHTML(".cmc-table tbody", func(e *colly.HTMLElement) {

		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			log.Print(e)
			name := el.ChildText("td:nth-child(3)")
			if el.ChildText("td:nth-child(3) span:nth-child(2)") != "" {
				name = el.ChildText("td:nth-child(3) span:nth-child(2)")
			}

			currency := common.Currency{
				ID:            el.ChildText("td:nth-child(2)"),
				Name:          name,
				Cours:         el.ChildText("td:nth-child(4)"),
				Variation1h:   el.ChildText("td:nth-child(5)"),
				Variation1d:   el.ChildText("td:nth-child(6)"),
				Variation1w:   el.ChildText("td:nth-child(7)"),
				MarketCapital: el.ChildText("td:nth-child(8)"),
				Volume:        el.ChildText("td:nth-child(9)"),
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

// // ScrapNewCurrenciesscraps currency information from newly trending crypto currencies.
// func ScrapNewCurrencies(url string) ([]common.CurrencyEvent, error) {
// 	c := colly.NewCollector()
// 	c.UserAgent = "Go scrapping"

// 	var currencies []common.Currency
// 	var countCurrency int

// 	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
// 		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
// 			currency := common.Currency{
// 				Name:        el.ChildText("td:nth-child(3) span:first-child"),
// 				Acronym:     el.ChildText("td:nth-child(3) span:nth-child(2)"),
// 				Cours:       el.ChildText("td:nth-child(4)"),
// 				Chaine:      el.ChildText("td:nth-child(5)"),
// 				Variation1h: el.ChildText("td:nth-child(6)"),
// 				Variation1d: el.ChildText("td:nth-child(7)"),
// 				Volume24d:   el.ChildText("td:nth-child(8)"),
// 				FDV:         el.ChildText("td:nth-child(9)"),
// 				LastAdded:   el.ChildText("td:nth-child(10)"),
// 			}

// 			currencies = append(currencies, currency)
// 		})

// 		countCurrency = e.DOM.Find("tr").Length()
// 	})

// 	c.Visit(url)

// 	if len(currencies) != countCurrency {
// 		return nil, fmt.Errorf("expected %d currencies but got %d", countCurrency, len(currencies))
// 	}

// 	currenciesEvents, err := common.CurrenciesToCurrencyEvents(currencies)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not convert newCurrency to CurrencyEvent: %w", err)
// 	}
// 	return currenciesEvents, nil
// }
