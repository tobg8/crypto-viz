package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly/v2"
)

type Currency struct {
	ID            string
	Name          string
	Acronym       string
	Cours         string
	variation1h   string
	variation1d   string
	variation1w   string
	volume        string
	marketCapital string
}

func main() {
	my_scheduler := gocron.NewScheduler(time.UTC)
	my_scheduler.Every(3).Seconds().Do(scrapCurrencies)
	my_scheduler.StartBlocking()
}

func scrapCurrencies() {
	c := colly.NewCollector()
	c.UserAgent = "Go scrapping"

	var currencies []Currency
	timer := 0

	// Will be executed on every request
	c.OnRequest(func(r *colly.Request) {
		timer = timer + 1
		fmt.Println(timer)
	})

	c.OnHTML(".coingecko-table .coin-table tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			currency := Currency{
				ID:            el.ChildText("td:nth-child(2)"),
				Name:          el.ChildText("td:nth-child(3) span:first-child"),
				Acronym:       el.ChildText("td:nth-child(3) span:nth-child(2)"),
				variation1h:   el.ChildText("td:nth-child(4)"),
				variation1d:   el.ChildText("td:nth-child(5)"),
				variation1w:   el.ChildText("td:nth-child(6)"),
				volume:        el.ChildText("td:nth-child(7)"),
				marketCapital: el.ChildText("td:nth-child(8)"),
			}

			currencies = append(currencies, currency)
		})

		log.Printf("CURRENCIES: %v", currencies)
	})

	// Will be executed on every response
	// c.OnResponse(func(r *colly.Response) {

	// 	fmt.Println("-----------------------------")
	// 	fmt.Println("RESPONSE")

	// 	fmt.Println(r.StatusCode)
	// 	for key, value := range *r.Headers {
	// 		fmt.Printf("%s: %s\n", key, value)
	// 	}
	// })

	c.Visit("https://www.coingecko.com/fr")
}
