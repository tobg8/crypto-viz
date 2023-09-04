package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/tobg8/crypto-viz/scrapper"
)

func main() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(3).Seconds().Do(Init)
	scheduler.StartBlocking()

	// TODO PRODUCER
	// Init producer

	// TODO CONSUMER
	// Init Consumer
}

func Init() {
	currencies, err := scrapper.ScrapCurrencies("https://www.coingecko.com/fr")
	if err != nil {
		log.Print(err)
	}
	log.Printf("CURRENCIES: %v", currencies)

	newCurrencies, err := scrapper.ScrapNewCurrencies("https://www.coingecko.com/fr/new-cryptocurrencies")
	if err != nil {
		log.Print(err)
	}
	log.Printf("NEW CURRENCIES: %v", newCurrencies)

	// TODO PRODUCER
	// create and send message with producer

	// TODO CONSUMER
	// read producer message and load them into a pocketbase DB
}
